package main

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/tachunwu/graviton/pkg/watcher"
)

const ControlPlaneURL = nats.DefaultURL

func main() {
	watcher, err := watcher.NewClusterStatusWatcher(ControlPlaneURL)
	if err != nil {
		log.Fatalln("Failed to initialize the ClusterStatusWatcher:", err)
	}

	errCh := watcher.WatchClusterStatus()

	// Watch for errors in a separate Goroutine.
	go func() {
		if err, ok := <-errCh; ok {
			log.Fatalln("Failed to watch cluster status:", err)
		}
	}()
	select {}
}
