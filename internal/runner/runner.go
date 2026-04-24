// Package runner wires together the scanner, monitor, and alert components
// and exposes a single Run function suitable for use in tests and the CLI.
package runner

import (
	"context"
	"fmt"
	"io"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

// Options controls runtime behaviour of the runner.
type Options struct {
	Config  *config.Config
	Output  io.Writer
	UseJSON bool
}

// Runner holds the wired-up components.
type Runner struct {
	mon *monitor.Monitor
}

// New constructs a Runner from the supplied options.
func New(opts Options) (*Runner, error) {
	sc, err := scanner.New(opts.Config.Protocol, opts.Config.StartPort, opts.Config.EndPort)
	if err != nil {
		return nil, fmt.Errorf("runner: scanner: %w", err)
	}

	var n alert.Notifier
	if opts.UseJSON {
		n = alert.NewJSON(opts.Output)
	} else {
		n = alert.New(opts.Output)
	}

	return &Runner{mon: monitor.New(sc, n, opts.Config.Interval)}, nil
}

// Start begins monitoring in a background goroutine and blocks until ctx
// is cancelled.
func (r *Runner) Start(ctx context.Context) {
	go r.mon.Start()
	<-ctx.Done()
	r.mon.Stop()
}
