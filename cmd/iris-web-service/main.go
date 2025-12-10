package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Elenetta17/iris-web-service/internal/httpapi"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", httpapi.FormPage)
	mux.HandleFunc("POST /hello", httpapi.HelloHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server in background
	go func() {
		log.Println("Server running at http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Wait for interrupt signal (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give ongoing requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
