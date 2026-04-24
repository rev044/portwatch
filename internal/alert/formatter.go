package alert

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// JSONNotifier writes alerts as newline-delimited JSON.
type JSONNotifier struct {
	w io.Writer
}

// NewJSON creates a JSONNotifier that writes to w.
func NewJSON(w io.Writer) *JSONNotifier {
	return &JSONNotifier{w: w}
}

type jsonAlert struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	ChangeType string `json:"change_type"`
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"`
	Message   string `json:"message"`
}

// Notify serialises the PortChange as a JSON line.
func (j *JSONNotifier) Notify(change monitor.PortChange) error {
	level := LevelInfo
	var msg string
	switch change.Type {
	case monitor.ChangeOpened:
		level = LevelAlert
		msg = fmt.Sprintf("port opened: %s", change.Port)
	case monitor.ChangeClosed:
		level = LevelWarn
		msg = fmt.Sprintf("port closed: %s", change.Port)
	}

	a := jsonAlert{
		Timestamp:  time.Now().Format(time.RFC3339),
		Level:      string(level),
		ChangeType: string(change.Type),
		Port:       change.Port.Port,
		Protocol:   change.Port.Protocol,
		Message:    msg,
	}

	b, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}
	_, err = fmt.Fprintf(j.w, "%s\n", b)
	return err
}
