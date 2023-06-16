package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {

	// NATS connection
	nc, _ := nats.Connect(nats.DefaultURL)

	// Create cluster config KV store
	js, _ := nc.JetStream()
	kv, _ := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:      "config",
		Description: "Gravition cluster config",
	})

	// Check status
	status, err := kv.Status()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Cluster status bucket name:", status.Bucket())

	// Handler for cluster config
}
