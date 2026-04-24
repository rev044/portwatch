package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfgPath := flag.String("config", "", "path to config file (optional)")
	jsonOutput := flag.Bool("json", false, "output alerts as JSON")
	flag.Parse()

	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	sc, err := scanner.New(cfg.Protocol, cfg.StartPort, cfg.EndPort)
	if err != nil {
		return fmt.Errorf("creating scanner: %w", err)
	}

	var notifier alert.Notifier
	if *jsonOutput {
		notifier = alert.NewJSON(os.Stdout)
	} else {
		notifier = alert.New(os.Stdout)
	}

	mon := monitor.New(sc, notifier, cfg.Interval)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	fmt.Fprintf(os.Stderr, "portwatch started (proto=%s ports=%d-%d interval=%s)\n",
		cfg.Protocol, cfg.StartPort, cfg.EndPort, cfg.Interval)

	go mon.Start()

	<-sigs
	fmt.Fprintln(os.Stderr, "shutting down")
	mon.Stop()
	return nil
}

func loadConfig(path string) (*config.Config, error) {
	if path == "" {
		cfg := config.DefaultConfig()
		return cfg, nil
	}
	return config.Load(path)
}
