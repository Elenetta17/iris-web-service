package config

import (
	"flag"
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 15*time.Second {
		t.Errorf("expected read timeout 15s, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Server.WriteTimeout != 15*time.Second {
		t.Errorf("expected write timeout 15s, got %v", cfg.Server.WriteTimeout)
	}
	if cfg.Server.IdleTimeout != 60*time.Second {
		t.Errorf("expected idle timeout 60s, got %v", cfg.Server.IdleTimeout)
	}
	if cfg.Server.ShutdownTimeout != 30*time.Second {
		t.Errorf("expected shutdown timeout 30s, got %v", cfg.Server.ShutdownTimeout)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary config file
	content := `server:
  port: 9090
  read_timeout: 20s
  write_timeout: 20s
  idle_timeout: 120s
  shutdown_timeout: 45s
`
	tmpfile, err := os.CreateTemp("", "config-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpfile.Close()

	// Create options pointing to temp file
	opts := &Options{
		ConfigFile: tmpfile.Name(),
	}

	cfg, err := Load(opts)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Server.Port)
	}
	if cfg.Server.ShutdownTimeout != 45*time.Second {
		t.Errorf("expected shutdown timeout 45s, got %v", cfg.Server.ShutdownTimeout)
	}
}

func TestLoadConfigWithMissingFile(t *testing.T) {
	opts := &Options{
		ConfigFile: "nonexistent.yml",
	}

	cfg, err := Load(opts)
	if err != nil {
		t.Fatalf("Load() should not fail with missing file: %v", err)
	}

	// Should fall back to defaults
	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}
}

func TestLoadConfigWithInvalidYAML(t *testing.T) {
	// Create a temporary invalid config file
	content := `server:
  port: not-a-number
`
	tmpfile, err := os.CreateTemp("", "config-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpfile.Close()

	opts := &Options{
		ConfigFile: tmpfile.Name(),
	}

	_, err = Load(opts)
	if err == nil {
		t.Error("Load() should fail with invalid YAML")
	}
}

func TestLoadConfigWithOverrides(t *testing.T) {
	// Create a config file with some values
	content := `server:
  port: 9090
  shutdown_timeout: 45s
`
	tmpfile, err := os.CreateTemp("", "config-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpfile.Close()

	// Override with options
	opts := &Options{
		ConfigFile:      tmpfile.Name(),
		Port:            3000,
		ShutdownTimeout: 60 * time.Second,
	}

	cfg, err := Load(opts)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Options should override file values
	if cfg.Server.Port != 3000 {
		t.Errorf("expected port 3000 (from options), got %d", cfg.Server.Port)
	}
	if cfg.Server.ShutdownTimeout != 60*time.Second {
		t.Errorf("expected shutdown timeout 60s (from options), got %v", cfg.Server.ShutdownTimeout)
	}
}

func TestLoadConfigWithReadError(t *testing.T) {
	// Try to read a directory as a file (causes read error)
	tmpdir, err := os.MkdirTemp("", "config-dir-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	opts := &Options{
		ConfigFile: tmpdir, // Directory, not file
	}

	_, err = Load(opts)
	if err == nil {
		t.Error("Load() should fail when trying to read a directory")
	}
}

func TestParseFlags(t *testing.T) {
	// Save original state
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	}()

	// Reset flag.CommandLine for this test
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Test with flags
	os.Args = []string{"cmd", "-config", "test.yml", "-port", "9090", "-shutdown-timeout", "60s"}

	opts := ParseFlags()

	if opts.ConfigFile != "test.yml" {
		t.Errorf("expected config file 'test.yml', got '%s'", opts.ConfigFile)
	}
	if opts.Port != 9090 {
		t.Errorf("expected port 9090, got %d", opts.Port)
	}
	if opts.ShutdownTimeout != 60*time.Second {
		t.Errorf("expected timeout 60s, got %v", opts.ShutdownTimeout)
	}
}

func TestParseFlagsDefaults(t *testing.T) {
	// Save original state
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	}()

	// Reset flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Test with no flags (defaults)
	os.Args = []string{"cmd"}

	opts := ParseFlags()

	if opts.ConfigFile != "config.yml" {
		t.Errorf("expected default config file 'config.yml', got '%s'", opts.ConfigFile)
	}
	if opts.Port != 0 {
		t.Errorf("expected default port 0, got %d", opts.Port)
	}
	if opts.ShutdownTimeout != 0 {
		t.Errorf("expected default timeout 0, got %v", opts.ShutdownTimeout)
	}
}
