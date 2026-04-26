package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SnapshotMeta holds metadata written alongside a snapshot file.
type SnapshotMeta struct {
	CreatedAt time.Time `json:"created_at"`
	Entries   int       `json:"entries"`
	Version   int       `json:"version"`
}

const snapshotVersion = 1

// SaveSnapshot writes the current history entries to a JSON file at the given
// path. The file is written atomically via a temp file + rename.
func SaveSnapshot(h *History, path string) error {
	entries := h.Entries()

	payload := struct {
		Meta    SnapshotMeta `json:"meta"`
		Entries []Entry      `json:"entries"`
	}{
		Meta: SnapshotMeta{
			CreatedAt: time.Now().UTC(),
			Entries:   len(entries),
			Version:   snapshotVersion,
		},
		Entries: entries,
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("snapshot: mkdir: %w", err)
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write temp: %w", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("snapshot: rename: %w", err)
	}

	return nil
}

// LoadSnapshot reads a previously saved snapshot and replays its entries into
// a new History with the given capacity. Returns an error if the file does not
// exist or is malformed.
func LoadSnapshot(path string, capacity int) (*History, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}

	var payload struct {
		Meta    SnapshotMeta `json:"meta"`
		Entries []Entry      `json:"entries"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}

	h := New(capacity)
	for _, e := range payload.Entries {
		h.Record(e.Change)
	}
	return h, nil
}
