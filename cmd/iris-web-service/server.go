package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Elenetta17/iris-web-service/internal/config"
	"github.com/Elenetta17/iris-web-service/internal/httpapi"
)

// Run starts the server with the given configuration
func Run(cfg *config.Config) error {
	return RunWithSignal(cfg, nil)
}

// RunWithSignal starts the server and listens for signals on the provided channel
// If quit is nil, it creates a default signal channel
func RunWithSignal(cfg *config.Config, quit chan os.Signal) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", httpapi.FormPage)
	mux.HandleFunc("POST /hello", httpapi.HelloHandler)

	// Add a slow endpoint for testing shutdown behavior
	mux.HandleFunc("GET /slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("done"))
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in background
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Server running at http://localhost:%d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Set up signal handling
	if quit == nil {
		quit = make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	}

	select {
	case <-quit:
		log.Println("Shutting down server...")
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	}

	// Give ongoing requests time to complete
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server stopped")
	return nil
}
