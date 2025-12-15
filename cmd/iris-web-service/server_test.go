package main

import (
	"net/http"
	"os"
	"strings" // Add this import
	"syscall"
	"testing"
	"time"

	"github.com/Elenetta17/iris-web-service/internal/config"
)

func TestRunFunction(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:            8887,
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    5 * time.Second,
			IdleTimeout:     30 * time.Second,
			ShutdownTimeout: 2 * time.Second,
		},
	}

	// Run() in background - it will block until shutdown
	done := make(chan error, 1)
	go func() {
		done <- Run(cfg)
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Verify server is running
	resp, err := http.Get("http://localhost:8887/")
	if err != nil {
		t.Fatalf("server not responding: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Send actual signal to the process to trigger shutdown
	// (Run creates its own signal channel, so we need to send a real signal)
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("failed to find process: %v", err)
	}

	if err := proc.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("failed to send signal: %v", err)
	}

	// Wait for shutdown
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run() returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("server did not shut down in time")
	}
}

func TestRunServerStartupError(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:            8882,
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    5 * time.Second,
			IdleTimeout:     30 * time.Second,
			ShutdownTimeout: 2 * time.Second,
		},
	}

	quit1 := make(chan os.Signal, 1)

	// Start first server
	done1 := make(chan error, 1)
	go func() {
		done1 <- RunWithSignal(cfg, quit1)
	}()

	// Wait for first server to start
	time.Sleep(200 * time.Millisecond)

	// Verify first server is running
	resp, err := http.Get("http://localhost:8882/")
	if err != nil {
		t.Fatalf("first server not responding: %v", err)
	}
	resp.Body.Close()

	// Try to start second server on same port (should fail immediately)
	quit2 := make(chan os.Signal, 1)
	done2 := make(chan error, 1)
	go func() {
		done2 <- RunWithSignal(cfg, quit2)
	}()

	// Should get error from second server
	select {
	case err := <-done2:
		if err == nil {
			t.Error("expected error when port is already in use")
		}
		// Verify it's a server error (not shutdown error)
		if !strings.Contains(err.Error(), "server error") {
			t.Errorf("expected 'server error', got: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("expected immediate error for port conflict")
	}

	// Clean up first server
	quit1 <- os.Interrupt
	<-done1
}

func TestRunShutdownTimeout(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:            8886,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			IdleTimeout:     60 * time.Second,
			ShutdownTimeout: 1 * time.Nanosecond, // Essentially 0 - will timeout immediately
		},
	}

	quit := make(chan os.Signal, 1)

	done := make(chan error, 1)
	go func() {
		done <- RunWithSignal(cfg, quit)
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Create multiple concurrent slow requests to ensure connections are active
	for i := 0; i < 5; i++ {
		go func() {
			client := &http.Client{Timeout: 30 * time.Second}
			client.Get("http://localhost:8886/slow")
		}()
	}

	// Make sure requests have started
	time.Sleep(200 * time.Millisecond)

	// Trigger shutdown
	quit <- os.Interrupt

	// Should get shutdown error
	select {
	case err := <-done:
		if err == nil {
			t.Error("expected error due to shutdown timeout, got nil")
		}
		if err != nil {
			t.Logf("got error: %v", err)
			if !strings.Contains(err.Error(), "shutdown") {
				t.Errorf("expected 'shutdown' in error message, got: %v", err)
			}
		}
	case <-time.After(3 * time.Second):
		t.Error("test timeout - server didn't return error")
	}
}
