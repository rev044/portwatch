package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	Ports    []PortRange   `json:"ports"`
	Interval time.Duration `json:"interval"`
	Protocol string        `json:"protocol"`
	Output   string        `json:"output"`
	Format   string        `json:"format"`
}

// PortRange defines an inclusive range of ports to monitor.
type PortRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Ports:    []PortRange{{From: 1, To: 1024}},
		Interval: 30 * time.Second,
		Protocol: "tcp",
		Output:   "stdout",
		Format:   "text",
	}
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	cfg := DefaultConfig()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate checks that the Config fields are sensible.
func (c *Config) Validate() error {
	if c.Protocol != "tcp" && c.Protocol != "udp" {
		return fmt.Errorf("config: unsupported protocol %q (want tcp or udp)", c.Protocol)
	}
	if c.Interval <= 0 {
		return fmt.Errorf("config: interval must be positive")
	}
	for _, r := range c.Ports {
		if r.From < 1 || r.To > 65535 || r.From > r.To {
			return fmt.Errorf("config: invalid port range %d-%d", r.From, r.To)
		}
	}
	if c.Format != "text" && c.Format != "json" {
		return fmt.Errorf("config: unsupported format %q (want text or json)", c.Format)
	}
	return nil
}
