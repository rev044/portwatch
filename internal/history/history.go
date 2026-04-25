// Package history provides a simple in-memory ring buffer for recording
// port change events over time, allowing portwatch to report recent activity.
package history

import (
	"sync"
	"time"

	"portwatch/internal/monitor"
)

// Entry records a single port change event with a timestamp.
type Entry struct {
	Timestamp time.Time
	Change    monitor.PortChange
}

// History is a thread-safe ring buffer of port change entries.
type History struct {
	mu      sync.RWMutex
	entries []Entry
	cap     int
	head    int
	size    int
}

// New creates a History that retains at most capacity entries.
// If capacity is less than 1 it defaults to 100.
func New(capacity int) *History {
	if capacity < 1 {
		capacity = 100
	}
	return &History{
		entries: make([]Entry, capacity),
		cap:     capacity,
	}
}

// Record appends a new change event, overwriting the oldest entry when full.
func (h *History) Record(c monitor.PortChange) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.entries[h.head] = Entry{Timestamp: time.Now(), Change: c}
	h.head = (h.head + 1) % h.cap
	if h.size < h.cap {
		h.size++
	}
}

// Entries returns a snapshot of recorded entries ordered oldest-first.
func (h *History) Entries() []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.size == 0 {
		return nil
	}

	out := make([]Entry, h.size)
	start := (h.head - h.size + h.cap) % h.cap
	for i := 0; i < h.size; i++ {
		out[i] = h.entries[(start+i)%h.cap]
	}
	return out
}

// Len returns the number of entries currently stored.
func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.size
}
