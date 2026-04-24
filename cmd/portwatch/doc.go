// Package main is the entry point for the portwatch CLI daemon.
//
// Usage:
//
//	portwatch [flags]
//
// Flags:
//
//	-config string
//	      path to a TOML configuration file (default: built-in defaults)
//	-json
//	      emit alerts as JSON objects instead of human-readable text
//
// Portwatch scans the configured port range at a regular interval and
// prints an alert whenever a port is opened or closed since the previous
// scan.  Send SIGINT or SIGTERM to stop the daemon gracefully.
//
// Example:
//
//	# Watch TCP ports 1–1024 every 5 s, output JSON
//	portwatch -config /etc/portwatch.toml -json
package main
