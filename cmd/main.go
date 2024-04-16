package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/timickb/go-stateflow/internal/builder"
	"github.com/timickb/go-stateflow/internal/config"
	"github.com/timickb/go-stateflow/pkg"
	"log"
)

func main() {
	configPath := flag.String("cfg", "config.yaml", "config file")
	if pkg.IsStrNilOrEmpty(configPath) {
		log.Fatal("empty config path")
	}
	if pkg.IsStrNilOrEmpty(configPath) {
		log.Fatal("empty config path")
	}

	if err := mainNoExit(*configPath); err != nil {
		log.Fatalf("Application start failed: %s", err.Error())
	}
}

func mainNoExit(cfgPath string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewFromFile(cfgPath)
	if err != nil {
		return fmt.Errorf("parse config from file: %w", err)
	}

	b, err := builder.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("create app builder: %w", err)
	}

	if err = b.ServeGrpc(); err != nil {
		return fmt.Errorf("server grpc: %w", err)
	}

	return nil
}
