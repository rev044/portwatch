package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/monitor"
)

func makeSnapChange(port int) monitor.PortChange {
	return monitor.PortChange{
		Port:     port,
		Protocol: "tcp",
		State:    monitor.Opened,
		At:       time.Now().UTC(),
	}
}

func TestSaveSnapshot_CreatesFile(t *testing.T) {
	h := New(10)
	h.Record(makeSnapChange(8080))
	h.Record(makeSnapChange(9090))

	path := filepath.Join(t.TempDir(), "snap", "history.json")
	if err := SaveSnapshot(h, path); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("snapshot file not found: %v", err)
	}
}

func TestLoadSnapshot_RestoresEntries(t *testing.T) {
	h := New(10)
	h.Record(makeSnapChange(1234))
	h.Record(makeSnapChange(5678))

	path := filepath.Join(t.TempDir(), "history.json")
	if err := SaveSnapshot(h, path); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	h2, err := LoadSnapshot(path, 10)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	entries := h2.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Change.Port != 1234 {
		t.Errorf("expected port 1234, got %d", entries[0].Change.Port)
	}
	if entries[1].Change.Port != 5678 {
		t.Errorf("expected port 5678, got %d", entries[1].Change.Port)
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/history.json", 10)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadSnapshot_RespectCapacity(t *testing.T) {
	h := New(20)
	for i := 0; i < 15; i++ {
		h.Record(makeSnapChange(3000 + i))
	}

	path := filepath.Join(t.TempDir(), "history.json")
	if err := SaveSnapshot(h, path); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	h2, err := LoadSnapshot(path, 5)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	// capacity 5 means only last 5 entries survive the ring buffer
	if got := len(h2.Entries()); got > 5 {
		t.Errorf("expected at most 5 entries with capacity 5, got %d", got)
	}
}
