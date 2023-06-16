package api

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type LeafHandler struct {
	NatsConn *nats.Conn
	DB       *badger.DB
	Logger   *zap.Logger
}

type Response struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
