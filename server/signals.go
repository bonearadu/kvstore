package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// WaitForShutdownSignal waits for an interrupt or terminate signal
func WaitForShutdownSignal() os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	return <-stop
}

// GracefulShutdown performs a graceful shutdown of the server
func GracefulShutdown(srv *Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server gracefully stopped")
}
