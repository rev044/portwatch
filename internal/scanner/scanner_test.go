package scanner_test

import (
	"net"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/scanner"
)

func startTestListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestScan_DetectsOpenPort(t *testing.T) {
	port, cleanup := startTestListener(t)
	defer cleanup()

	s := scanner.New(200 * time.Millisecond)
	results, err := s.Scan("tcp", port, port)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(results))
	}
	if results[0].Port != port {
		t.Errorf("expected port %d, got %d", port, results[0].Port)
	}
}

func TestScan_InvalidProtocol(t *testing.T) {
	s := scanner.New(200 * time.Millisecond)
	_, err := s.Scan("udp", 80, 80)
	if err == nil {
		t.Error("expected error for unsupported protocol, got nil")
	}
}

func TestScan_InvalidPortRange(t *testing.T) {
	s := scanner.New(200 * time.Millisecond)
	_, err := s.Scan("tcp", 9000, 8000)
	if err == nil {
		t.Error("expected error for invalid port range, got nil")
	}
}

func TestPortState_String(t *testing.T) {
	ps := scanner.PortState{Port: 8080, Protocol: "TCP", Address: "127.0.0.1"}
	expected := "127.0.0.1:8080 (TCP)"
	if ps.String() != expected {
		t.Errorf("expected %q, got %q", expected, ps.String())
	}
}
