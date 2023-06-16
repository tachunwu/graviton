package server

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/nats-io/nats.go"
	"github.com/tachunwu/graviton/pkg/api"
	"go.uber.org/zap"
)

type LeafServer struct {
	db     *badger.DB
	nc     *nats.Conn
	logger *zap.Logger
}

func NewLeafServer(db *badger.DB, nc *nats.Conn, logger *zap.Logger) *LeafServer {
	return &LeafServer{
		db:     db,
		nc:     nc,
		logger: logger,
	}
}

func (s *LeafServer) Start() {
	handler := &api.LeafHandler{
		NatsConn: s.nc,
		DB:       s.db,
		Logger:   s.logger,
	}

	// Subscribe to subject "$GVTN.KV.*"

	_, err := s.nc.Subscribe("$GVTN.KV.*", handler.HandleKVMsg)
	if err != nil {
		s.logger.Error("Error subscribing to subject", zap.Error(err))
	}
	s.logger.Info("Leaf KV",
		zap.String("KV single operation API", "$GVTN.KV.*"),
	)

	// Subscribe to subject "$GVTN.KV"
	_, err = s.nc.Subscribe("$GVTN.KV", handler.HandleScanMsg)
	if err != nil {
		s.logger.Error("Error subscribing to subject", zap.Error(err))
	}
	s.logger.Info("Leaf KV",
		zap.String("KV multi read API", "$GVTN.KV"),
	)
}
