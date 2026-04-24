package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// PortState represents the state of a single open port.
type PortState struct {
	Port     int
	Protocol string
	Address  string
}

// String returns a human-readable representation of a PortState.
func (p PortState) String() string {
	return fmt.Sprintf("%s:%d (%s)", p.Address, p.Port, p.Protocol)
}

// Scanner scans for open ports on the local machine.
type Scanner struct {
	Timeout time.Duration
}

// New creates a new Scanner with the given timeout.
func New(timeout time.Duration) *Scanner {
	return &Scanner{Timeout: timeout}
}

// Scan checks all ports in the given range and returns those that are open.
func (s *Scanner) Scan(protocol string, startPort, endPort int) ([]PortState, error) {
	if !isValidProtocol(protocol) {
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}
	if startPort < 1 || endPort > 65535 || startPort > endPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", startPort, endPort)
	}

	var open []PortState
	for port := startPort; port <= endPort; port++ {
		address := "localhost:" + strconv.Itoa(port)
		conn, err := net.DialTimeout(protocol, address, s.Timeout)
		if err == nil {
			conn.Close()
			host, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
			open = append(open, PortState{
				Port:     port,
				Protocol: strings.ToUpper(protocol),
				Address:  host,
			})
		}
	}
	return open, nil
}

func isValidProtocol(p string) bool {
	switch strings.ToLower(p) {
	case "tcp", "tcp4", "tcp6":
		return true
	}
	return false
}
