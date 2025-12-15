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

// Options holds configuration options that can override file values
type Options struct {
	ConfigFile      string
	Port            int
	ShutdownTimeout time.Duration
}

// ParseFlags parses command-line flags and returns Options
func ParseFlags() *Options {
	opts := &Options{}
	flag.StringVar(&opts.ConfigFile, "config", "config.yml", "path to config file")
	flag.IntVar(&opts.Port, "port", 0, "server port (overrides config file)")
	flag.DurationVar(&opts.ShutdownTimeout, "shutdown-timeout", 0, "shutdown timeout (overrides config file)")
	flag.Parse()
	return opts
}

// Load loads configuration from file with given options
func Load(opts *Options) (*Config, error) {
	// Start with defaults
	cfg := DefaultConfig()

	// Load from file if specified
	if opts.ConfigFile != "" {
		data, err := os.ReadFile(opts.ConfigFile)
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

	// Apply option overrides
	if opts.Port != 0 {
		cfg.Server.Port = opts.Port
	}
	if opts.ShutdownTimeout != 0 {
		cfg.Server.ShutdownTimeout = opts.ShutdownTimeout
	}

	return cfg, nil
}
