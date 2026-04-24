// Package alert provides notification primitives for portwatch.
//
// It consumes monitor.PortChange events and dispatches human-readable
// or machine-readable (JSON) alerts to one or more io.Writer sinks.
//
// Basic usage:
//
//	notifier := alert.New(os.Stdout)
//	notifier.Notify(change)
//
// For structured logging pipelines use JSONNotifier:
//
//	jn := alert.NewJSON(logFile)
//	jn.Notify(change)
package alert
