package main

import (
	"log"

	"github.com/bonearadu/kvstore/api"
	"github.com/bonearadu/kvstore/config"
	"github.com/bonearadu/kvstore/kv_store"
	"github.com/bonearadu/kvstore/server"
)

func main() {
	// Parse configuration
	cfg := config.ParseFlags()

	// Initialize components
	var store kv_store.KeyValueStore
	switch cfg.Mode {
	case config.InMemory:
		store = kv_store.NewInMemoryStore()
		log.Printf("Using in-memory KV store")
	case config.Persistent:
		store = kv_store.NewPersistentStore(cfg.StorePath)
		log.Printf("Using persistent KV store. Store root path: %s", cfg.StorePath)
	case config.PersistentCached:
		store = kv_store.NewPersistentCachedStore(cfg.StorePath, cfg.CacheCapacity)
		log.Printf("Using persistent KV store with caching. Store root path: %s. Cache size: %d",
			cfg.StorePath, cfg.CacheCapacity)
	default:
		log.Panicf("Unknown map store implementation specified: %d", cfg.Mode)
	}

	handler := api.NewHandler(store)
	srv := server.New(cfg, handler)

	// Start server
	srv.Start()

	// Wait for shutdown signal
	server.WaitForShutdownSignal()

	// Perform graceful shutdown
	server.GracefulShutdown(srv)
}
