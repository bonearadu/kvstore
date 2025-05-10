package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/bonearadu/kvstore/config"
)

// Server represents an HTTP server
type Server struct {
	httpServer *http.Server
}

// New creates a new server with the given configuration and handler
func New(cfg *config.ServerConfig, handler http.Handler) *Server {
	addr := fmt.Sprintf(":%d", cfg.Port)
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

// Start starts the server in a goroutine
func (s *Server) Start() {
	go func() {
		log.Printf("Starting server on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
