package watcher

import (
	"log"

	"github.com/nats-io/nats.go"
)

type ClusterStatusWatcher struct {
	nc *nats.Conn
	js nats.JetStreamContext
	kv nats.KeyValue
}

func NewClusterStatusWatcher(controlPlaneURL string) (*ClusterStatusWatcher, error) {
	watcher := &ClusterStatusWatcher{}

	// Initialize nats connection
	nc, err := nats.Connect(controlPlaneURL)
	if err != nil {
		return nil, err
	}
	watcher.nc = nc

	// Initialize JetStream context
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}
	watcher.js = js

	// Get cluster config kv store
	kv, err := js.KeyValue("config")
	if err != nil {
		log.Println("fail to get config kv:", err)
		return nil, err
	}
	watcher.kv = kv

	return watcher, nil
}

func (w *ClusterStatusWatcher) WatchClusterStatus() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		defer w.nc.Drain()
		defer close(errCh)

		watcher, err := w.kv.WatchAll()
		if err != nil {
			log.Println("fail to watch config kv", err)
			errCh <- err
			return
		}

		for {
			c := <-watcher.Updates()
			if c == nil {
				continue
			}
			log.Println("config kv updated:", c.Key(), string(c.Value()), c.Revision())
		}
	}()

	return errCh
}
