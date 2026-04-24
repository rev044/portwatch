package config_test

import (
	"os"
	"testing"
	"time"

	"portwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", cfg.Protocol)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
	if len(cfg.Ports) == 0 {
		t.Error("expected at least one default port range")
	}
}

func TestLoad_ValidFile(t *testing.T) {
	path := writeTemp(t, `{"protocol":"tcp","interval":10000000000,"format":"json","ports":[{"from":80,"to":443}]}`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "json" {
		t.Errorf("expected format json, got %s", cfg.Format)
	}
	if len(cfg.Ports) != 1 || cfg.Ports[0].From != 80 {
		t.Errorf("unexpected ports: %+v", cfg.Ports)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestValidate_InvalidProtocol(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Protocol = "icmp"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid protocol")
	}
}

func TestValidate_InvalidPortRange(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Ports = []config.PortRange{{From: 500, To: 100}}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for inverted port range")
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Format = "xml"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for unsupported format")
	}
}
