package history

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// ExportJSON writes all history entries as a JSON array to w.
func (h *History) ExportJSON(w io.Writer) error {
	h.mu.Lock()
	entries := h.Entries()
	h.mu.Unlock()

	type jsonEntry struct {
		Timestamp time.Time `json:"timestamp"`
		Protocol  string    `json:"protocol"`
		Port      int       `json:"port"`
		Kind      string    `json:"kind"`
	}

	out := make([]jsonEntry, 0, len(entries))
	for _, e := range entries {
		out = append(out, jsonEntry{
			Timestamp: e.Timestamp,
			Protocol:  e.Change.Protocol,
			Port:      e.Change.Port,
			Kind:      e.Change.Kind,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// ExportTable writes all history entries as a human-readable table to w.
func (h *History) ExportTable(w io.Writer) error {
	h.mu.Lock()
	entries := h.Entries()
	h.mu.Unlock()

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tPROTOCOL\tPORT\tEVENT")

	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\n",
			e.Timestamp.Format(time.RFC3339),
			e.Change.Protocol,
			e.Change.Port,
			e.Change.Kind,
		)
	}

	return tw.Flush()
}
