package main

import (
	"log"
	"net/http"

	"github.com/Elenetta17/iris-web-service/internal/httpapi"
)

func main() {
	http.HandleFunc("/", httpapi.HelloHandler)

	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
