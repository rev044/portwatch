package runner_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/runner"
)

func defaultCfg() *config.Config {
	cfg := config.DefaultConfig()
	cfg.StartPort = 19900
	cfg.EndPort = 19910
	cfg.Interval = 50 * time.Millisecond
	return cfg
}

func TestNew_ValidOptions(t *testing.T) {
	var buf bytes.Buffer
	_, err := runner.New(runner.Options{
		Config:  defaultCfg(),
		Output:  &buf,
		UseJSON: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_JSONOutput(t *testing.T) {
	var buf bytes.Buffer
	_, err := runner.New(runner.Options{
		Config:  defaultCfg(),
		Output:  &buf,
		UseJSON: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_InvalidProtocol(t *testing.T) {
	cfg := defaultCfg()
	cfg.Protocol = "udp" // scanner only supports tcp/tcp4/tcp6
	var buf bytes.Buffer
	_, err := runner.New(runner.Options{Config: cfg, Output: &buf})
	if err == nil {
		t.Fatal("expected error for invalid protocol")
	}
	if !strings.Contains(err.Error(), "runner: scanner") {
		t.Errorf("error should mention runner/scanner, got: %v", err)
	}
}

func TestStart_CancelStops(t *testing.T) {
	var buf bytes.Buffer
	r, err := runner.New(runner.Options{
		Config:  defaultCfg(),
		Output:  &buf,
		UseJSON: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		r.Start(ctx)
		close(done)
	}()

	select {
	case <-done:
		// runner stopped cleanly after context cancellation
	case <-time.After(500 * time.Millisecond):
		t.Fatal("runner did not stop within deadline")
	}
}
