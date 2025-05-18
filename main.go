package main

import (
	"github.com/bonearadu/kvstore/api"
	"github.com/bonearadu/kvstore/config"
	"github.com/bonearadu/kvstore/kv_store"
	"github.com/bonearadu/kvstore/server"
)

func main() {
	// Parse configuration
	cfg := config.ParseFlags()

	// Initialize components
	store := kv_store.NewInMemoryStore()
	handler := api.NewHandler(store)
	srv := server.New(cfg, handler)

	// Start server
	srv.Start()

	// Wait for shutdown signal
	server.WaitForShutdownSignal()

	// Perform graceful shutdown
	server.GracefulShutdown(srv)
}
