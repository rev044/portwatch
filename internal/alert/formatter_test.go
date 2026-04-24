package alert_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func TestJSONNotifier_OutputIsValidJSON(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewJSON(&buf)

	change := monitor.PortChange{
		Type: monitor.ChangeOpened,
		Port: scanner.PortState{Port: 443, Protocol: "tcp", Open: true},
	}

	if err := n.Notify(change); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v\nraw: %s", err, buf.String())
	}
}

func TestJSONNotifier_FieldsPresent(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewJSON(&buf)

	change := monitor.PortChange{
		Type: monitor.ChangeClosed,
		Port: scanner.PortState{Port: 22, Protocol: "tcp", Open: false},
	}

	if err := n.Notify(change); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &out)

	for _, field := range []string{"timestamp", "level", "port", "protocol", "message", "change_type"} {
		if _, ok := out[field]; !ok {
			t.Errorf("missing field %q in JSON output", field)
		}
	}

	if out["level"] != "WARN" {
		t.Errorf("expected level WARN for closed port, got %v", out["level"])
	}
	if int(out["port"].(float64)) != 22 {
		t.Errorf("expected port 22, got %v", out["port"])
	}
}
