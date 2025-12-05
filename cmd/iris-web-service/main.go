package main

import (
	"log"
	"net/http"

	"github.com/Elenetta17/iris-web-service/internal/httpapi"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", httpapi.FormPage)        // Only "/" exactly
	mux.HandleFunc("POST /hello", httpapi.HelloHandler) // POST only

	log.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
