package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/timickb/go-stateflow/internal/builder"
	"github.com/timickb/go-stateflow/internal/config"
	"github.com/timickb/go-stateflow/pkg/utils"
	"log"
)

func main() {
	configPath := flag.String("cfg", "config.yaml", "config file")
	if utils.IsStrNilOrEmpty(configPath) {
		log.Fatal("empty config path")
	}
	if utils.IsStrNilOrEmpty(configPath) {
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

	log.Printf("Starting gRPC server on port %d", b.ServerPort())
	if err = b.ServeGrpc(); err != nil {
		log.Fatalf("Fail to serve grpc: %s", err.Error())
	}

	return nil
}
