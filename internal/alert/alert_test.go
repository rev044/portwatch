package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func makeChange(t monitor.ChangeType) monitor.PortChange {
	return monitor.PortChange{
		Type: t,
		Port: scanner.PortState{Port: 8080, Protocol: "tcp", Open: true},
	}
}

func TestNotifier_OpenedPort(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	if err := n.Notify(makeChange(monitor.ChangeOpened)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[ALERT]") {
		t.Errorf("expected ALERT level, got: %s", out)
	}
	if !strings.Contains(out, "port opened") {
		t.Errorf("expected 'port opened' in output, got: %s", out)
	}
}

func TestNotifier_ClosedPort(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	if err := n.Notify(makeChange(monitor.ChangeClosed)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("expected WARN level, got: %s", out)
	}
	if !strings.Contains(out, "port closed") {
		t.Errorf("expected 'port closed' in output, got: %s", out)
	}
}

func TestNotifier_MultipleWriters(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	n := alert.New(&buf1, &buf2)

	if err := n.Notify(makeChange(monitor.ChangeOpened)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf1.String() != buf2.String() {
		t.Errorf("writers received different output")
	}
}

func TestNotifier_DefaultsToStdout(t *testing.T) {
	// Should not panic when no writers provided
	n := alert.New()
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
