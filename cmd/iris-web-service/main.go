package main

import (
	"log"

	"github.com/Elenetta17/iris-web-service/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run() error {
	opts := config.ParseFlags()

	cfg, err := config.Load(opts)
	if err != nil {
		return err
	}

	return Run(cfg)
}
