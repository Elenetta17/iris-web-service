package main

import (
	"flag"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestRunFunctionWithInvalidConfig(t *testing.T) {
	// Save original state
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	}()

	// Create invalid config file
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

	// Reset flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-config", tmpfile.Name()}

	// Should return error from config parsing
	err = run()
	if err == nil {
		t.Error("run() should return error with invalid config")
	}
}

func TestRunFunctionConfigLoadError(t *testing.T) {
	// Save original state
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	}()

	// Try to read a directory as config file
	tmpdir, err := os.MkdirTemp("", "config-dir-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	// Reset flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-config", tmpdir}

	// Should return error from config loading
	err = run()
	if err == nil {
		t.Error("run() should return error when config cannot be read")
	}
}

func TestRunFunctionSuccess(t *testing.T) {
	// Save original state
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	}()

	// Create valid config file
	content := `server:
  port: 8889
  read_timeout: 5s
  write_timeout: 5s
  idle_timeout: 30s
  shutdown_timeout: 1s
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

	// Reset flag.CommandLine
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-config", tmpfile.Name()}

	// Run in background
	done := make(chan error, 1)
	go func() {
		done <- run()
	}()

	// Give server time to start
	time.Sleep(200 * time.Millisecond)

	// Verify server is running
	resp, err := http.Get("http://localhost:8889/")
	if err != nil {
		t.Fatalf("server not responding: %v", err)
	}
	resp.Body.Close()

	// Send signal to shut down
	proc, _ := os.FindProcess(os.Getpid())
	proc.Signal(syscall.SIGTERM)

	// Wait for shutdown
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("run() returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("timeout waiting for shutdown")
	}
}
