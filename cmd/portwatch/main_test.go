package main

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func TestLoadConfig_Default(t *testing.T) {
	cfg, err := loadConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	def := config.DefaultConfig()
	if cfg.Protocol != def.Protocol {
		t.Errorf("protocol: got %q want %q", cfg.Protocol, def.Protocol)
	}
	if cfg.StartPort != def.StartPort {
		t.Errorf("start port: got %d want %d", cfg.StartPort, def.StartPort)
	}
}

func TestLoadConfig_FromFile(t *testing.T) {
	content := `protocol = "tcp"
start_port = 8000
end_port   = 8100
interval   = "10s"
`
	f, err := os.CreateTemp("", "portwatch-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := loadConfig(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StartPort != 8000 {
		t.Errorf("start port: got %d want 8000", cfg.StartPort)
	}
	if cfg.EndPort != 8100 {
		t.Errorf("end port: got %d want 8100", cfg.EndPort)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := loadConfig("/nonexistent/path/portwatch.toml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
