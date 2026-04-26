package history

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeExportChange(kind, proto string, port int) Entry {
	return Entry{
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Change: ChangeRecord{
			Kind:     kind,
			Protocol: proto,
			Port:     port,
		},
	}
}

func TestExportJSON_ValidOutput(t *testing.T) {
	h := New(10)
	h.Record(makeExportChange("opened", "tcp", 8080).Change)
	h.Record(makeExportChange("closed", "udp", 53).Change)

	var buf bytes.Buffer
	if err := h.ExportJSON(&buf); err != nil {
		t.Fatalf("ExportJSON returned error: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}

	if result[0]["protocol"] != "tcp" || result[0]["kind"] != "opened" {
		t.Errorf("unexpected first entry: %v", result[0])
	}
}

func TestExportTable_ContainsHeaders(t *testing.T) {
	h := New(10)
	h.Record(makeExportChange("opened", "tcp", 9090).Change)

	var buf bytes.Buffer
	if err := h.ExportTable(&buf); err != nil {
		t.Fatalf("ExportTable returned error: %v", err)
	}

	out := buf.String()
	for _, col := range []string{"TIMESTAMP", "PROTOCOL", "PORT", "EVENT"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected column %q in table output", col)
		}
	}

	if !strings.Contains(out, "tcp") || !strings.Contains(out, "9090") {
		t.Errorf("expected entry data in table output, got:\n%s", out)
	}
}

func TestExportJSON_EmptyHistory(t *testing.T) {
	h := New(10)

	var buf bytes.Buffer
	if err := h.ExportJSON(&buf); err != nil {
		t.Fatalf("ExportJSON on empty history returned error: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected empty array, got %d entries", len(result))
	}
}
