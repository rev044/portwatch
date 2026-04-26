package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/history"
	"portwatch/internal/monitor"
)

// TestSnapshotRoundTrip verifies that a full save/load cycle preserves all
// entry fields and ordering across a realistic sequence of port changes.
func TestSnapshotRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "roundtrip.json")

	h := history.New(16)

	changes := []monitor.PortChange{
		{Port: 8080, Proto: "tcp", State: monitor.Opened, DetectedAt: time.Unix(1700000001, 0).UTC()},
		{Port: 8080, Proto: "tcp", State: monitor.Closed, DetectedAt: time.Unix(1700000060, 0).UTC()},
		{Port: 443, Proto: "tcp", State: monitor.Opened, DetectedAt: time.Unix(1700000120, 0).UTC()},
		{Port: 53, Proto: "udp", State: monitor.Opened, DetectedAt: time.Unix(1700000180, 0).UTC()},
	}

	for _, c := range changes {
		h.Record(c)
	}

	if err := history.SaveSnapshot(path, h); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	loaded := history.New(16)
	if err := history.LoadSnapshot(path, loaded); err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	orig := h.Entries()
	restored := loaded.Entries()

	if len(orig) != len(restored) {
		t.Fatalf("entry count mismatch: got %d, want %d", len(restored), len(orig))
	}

	for i, e := range orig {
		r := restored[i]
		if e.Change.Port != r.Change.Port {
			t.Errorf("[%d] port: got %d, want %d", i, r.Change.Port, e.Change.Port)
		}
		if e.Change.Proto != r.Change.Proto {
			t.Errorf("[%d] proto: got %q, want %q", i, r.Change.Proto, e.Change.Proto)
		}
		if e.Change.State != r.Change.State {
			t.Errorf("[%d] state: got %v, want %v", i, r.Change.State, e.Change.State)
		}
		if !e.Change.DetectedAt.Equal(r.Change.DetectedAt) {
			t.Errorf("[%d] detectedAt: got %v, want %v", i, r.Change.DetectedAt, e.Change.DetectedAt)
		}
	}
}

// TestSnapshotRoundTrip_OverCapacity ensures that when a snapshot file holds
// more entries than the target History capacity, only the most recent N are
// loaded — matching the ring-buffer semantics of New.
func TestSnapshotRoundTrip_OverCapacity(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "overflow.json")

	// Record 8 entries into a large history.
	big := history.New(32)
	for i := 0; i < 8; i++ {
		big.Record(monitor.PortChange{
			Port:       uint16(9000 + i),
			Proto:      "tcp",
			State:      monitor.Opened,
			DetectedAt: time.Unix(int64(1700000000+i*60), 0).UTC(),
		})
	}

	if err := history.SaveSnapshot(path, big); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	// Load into a history with capacity 4 — should keep the 4 most recent.
	small := history.New(4)
	if err := history.LoadSnapshot(path, small); err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	entries := small.Entries()
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}

	// The last 4 ports recorded were 9004–9007.
	for i, e := range entries {
		want := uint16(9004 + i)
		if e.Change.Port != want {
			t.Errorf("entry[%d]: port %d, want %d", i, e.Change.Port, want)
		}
	}
}

// TestSaveSnapshot_FilePermissions checks that the snapshot file is created
// with restrictive permissions (0600) so credentials/port data are not
// world-readable.
func TestSaveSnapshot_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "perms.json")

	h := history.New(8)
	h.Record(monitor.PortChange{Port: 22, Proto: "tcp", State: monitor.Opened, DetectedAt: time.Now().UTC()})

	if err := history.SaveSnapshot(path, h); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}

	const want = os.FileMode(0600)
	if got := info.Mode().Perm(); got != want {
		t.Errorf("file permissions: got %04o, want %04o", got, want)
	}
}
