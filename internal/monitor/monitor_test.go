package monitor_test

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func TestMonitor_DetectsOpenedPort(t *testing.T) {
	port := freePort(t)

	// Start listener after monitor's first poll so it appears as OPENED.
	s, err := scanner.New([]int{port}, "tcp")
	if err != nil {
		t.Fatalf("scanner.New: %v", err)
	}

	m := monitor.New(s, 50*time.Millisecond)
	stop := make(chan struct{})
	go m.Start(stop) //nolint:errcheck

	// Give the first poll (empty baseline) time to run.
	time.Sleep(80 * time.Millisecond)

	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()

	select {
	case change := <-m.Alerts:
		if change.Change != monitor.ChangeOpened {
			t.Errorf("expected OPENED, got %s", change.Change)
		}
		if change.Port != port {
			t.Errorf("expected port %d, got %d", port, change.Port)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for OPENED alert")
	}

	close(stop)
}

func TestPortChange_String(t *testing.T) {
	c := monitor.PortChange{
		Port:       8080,
		Protocol:   "tcp",
		Change:     monitor.ChangeOpened,
		DetectedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
	got := c.String()
	want := "[2024-01-15T10:00:00Z] tcp/8080 OPENED"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
