package history_test

import (
	"testing"

	"portwatch/internal/history"
	"portwatch/internal/monitor"
	"portwatch/internal/scanner"
)

func makeChange(port int, state scanner.PortState) monitor.PortChange {
	return monitor.PortChange{
		Port:  port,
		Proto: "tcp",
		State: state,
	}
}

func TestNew_DefaultCapacity(t *testing.T) {
	h := history.New(0)
	if h == nil {
		t.Fatal("expected non-nil History")
	}
}

func TestRecord_And_Entries(t *testing.T) {
	h := history.New(10)
	h.Record(makeChange(8080, scanner.StateOpen))
	h.Record(makeChange(9090, scanner.StateOpen))

	entries := h.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Change.Port != 8080 {
		t.Errorf("expected first entry port 8080, got %d", entries[0].Change.Port)
	}
	if entries[1].Change.Port != 9090 {
		t.Errorf("expected second entry port 9090, got %d", entries[1].Change.Port)
	}
}

func TestHistory_RingOverwrite(t *testing.T) {
	h := history.New(3)
	for i := 1; i <= 5; i++ {
		h.Record(makeChange(i, scanner.StateOpen))
	}

	if h.Len() != 3 {
		t.Fatalf("expected 3 entries after overflow, got %d", h.Len())
	}

	entries := h.Entries()
	// Oldest surviving entry should be port 3
	if entries[0].Change.Port != 3 {
		t.Errorf("expected oldest port 3, got %d", entries[0].Change.Port)
	}
	if entries[2].Change.Port != 5 {
		t.Errorf("expected newest port 5, got %d", entries[2].Change.Port)
	}
}

func TestEntries_EmptyHistory(t *testing.T) {
	h := history.New(10)
	if h.Entries() != nil {
		t.Error("expected nil slice for empty history")
	}
}

func TestHistory_TimestampsSet(t *testing.T) {
	h := history.New(5)
	h.Record(makeChange(443, scanner.StateClosed))
	entries := h.Entries()
	if entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
