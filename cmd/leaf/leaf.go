package main

import (
	"log"

	"github.com/dgraph-io/badger/v4"
	"github.com/nats-io/nats.go"
	"github.com/spf13/pflag"
	"github.com/tachunwu/graviton/pkg/server"
	"go.uber.org/zap"
)

func main() {
	// Define command-line flags
	natsURL := pflag.String("nats-url", nats.DefaultURL, "NATS server URL")
	pflag.Parse()

	// Initialize Logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// Initialize BadgerDB
	opts := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opts)
	if err != nil {
		logger.Fatal("Failed to open BadgerDB", zap.Error(err))
	}
	defer db.Close()

	// Initialize NATS connection
	nc, err := nats.Connect(*natsURL)
	if err != nil {
		logger.Fatal("Failed to connect to NATS", zap.Error(err))
	}
	defer nc.Close()

	// Initialize LeafServer
	leafServer := server.NewLeafServer(db, nc, logger)

	// Start LeafServer
	leafServer.Start()

	// Infinite loop, waiting for termination signal
	select {}
}
