package monitor

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// ChangeType describes the kind of port change detected.
type ChangeType string

const (
	ChangeOpened ChangeType = "OPENED"
	ChangeClosed ChangeType = "CLOSED"
)

// PortChange represents a detected change in port state.
type PortChange struct {
	Port      int
	Protocol  string
	Change    ChangeType
	DetectedAt time.Time
}

func (c PortChange) String() string {
	return fmt.Sprintf("[%s] %s/%d %s", c.DetectedAt.Format(time.RFC3339), c.Protocol, c.Port, c.Change)
}

// Monitor watches a set of ports and reports changes.
type Monitor struct {
	scanner  *scanner.Scanner
	interval time.Duration
	previous map[string]bool
	Alerts   chan PortChange
}

// New creates a Monitor that polls the given ports every interval.
func New(s *scanner.Scanner, interval time.Duration) *Monitor {
	return &Monitor{
		scanner:  s,
		interval: interval,
		previous: make(map[string]bool),
		Alerts:   make(chan PortChange, 64),
	}
}

// Start begins the polling loop; it blocks until stop is closed.
func (m *Monitor) Start(stop <-chan struct{}) error {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			close(m.Alerts)
			return nil
		case <-ticker.C:
			if err := m.poll(); err != nil {
				return err
			}
		}
	}
}

func (m *Monitor) poll() error {
	states, err := m.scanner.Scan()
	if err != nil {
		return fmt.Errorf("monitor poll: %w", err)
	}

	current := make(map[string]bool, len(states))
	for _, ps := range states {
		key := fmt.Sprintf("%s/%d", ps.Protocol, ps.Port)
		current[key] = ps.Open
		if ps.Open && !m.previous[key] {
			m.Alerts <- PortChange{Port: ps.Port, Protocol: ps.Protocol, Change: ChangeOpened, DetectedAt: time.Now()}
		}
	}
	for key, wasOpen := range m.previous {
		if wasOpen && !current[key] {
			var proto string
			var port int
			fmt.Sscanf(key, "%3s/%d", &proto, &port)
			m.Alerts <- PortChange{Port: port, Protocol: proto, Change: ChangeClosed, DetectedAt: time.Now()}
		}
	}
	m.previous = current
	return nil
}
