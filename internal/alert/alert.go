package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single alert event.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Change    monitor.PortChange
	Message   string
}

// Notifier sends alerts to one or more outputs.
type Notifier struct {
	writers []io.Writer
}

// New creates a Notifier that writes to the provided writers.
// If no writers are given, os.Stdout is used.
func New(writers ...io.Writer) *Notifier {
	if len(writers) == 0 {
		writers = []io.Writer{os.Stdout}
	}
	return &Notifier{writers: writers}
}

// Notify formats and dispatches an alert for the given PortChange.
func (n *Notifier) Notify(change monitor.PortChange) error {
	a := Alert{
		Timestamp: time.Now(),
		Change:    change,
	}

	switch change.Type {
	case monitor.ChangeOpened:
		a.Level = LevelAlert
		a.Message = fmt.Sprintf("port opened: %s", change.Port)
	case monitor.ChangeClosed:
		a.Level = LevelWarn
		a.Message = fmt.Sprintf("port closed: %s", change.Port)
	default:
		a.Level = LevelInfo
		a.Message = fmt.Sprintf("port change: %s", change.Port)
	}

	line := fmt.Sprintf("%s [%s] %s\n",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Message,
	)

	for _, w := range n.writers {
		if _, err := fmt.Fprint(w, line); err != nil {
			return fmt.Errorf("alert write failed: %w", err)
		}
	}
	return nil
}
