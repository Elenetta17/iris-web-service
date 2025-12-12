// internal/config/config.go
package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Port            int           `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

// Load loads configuration from file and command-line flags
// Command-line flags override file values
func Load() (*Config, error) {
	// Define flags
	configFile := flag.String("config", "config.yml", "path to config file")
	port := flag.Int("port", 0, "server port (overrides config file)")
	shutdownTimeout := flag.Duration("shutdown-timeout", 0, "shutdown timeout (overrides config file)")
	flag.Parse()

	// Start with defaults
	cfg := DefaultConfig()

	if *configFile != "" {
		data, err := os.ReadFile(*configFile)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("reading config file: %w", err)
			}
			// File doesn't exist, use defaults
		} else {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("parsing config file: %w", err)
			}
		}
	}

	// Command-line flags override file values
	if *port != 0 {
		cfg.Server.Port = *port
	}
	if *shutdownTimeout != 0 {
		cfg.Server.ShutdownTimeout = *shutdownTimeout
	}

	return cfg, nil
}
