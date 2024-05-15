package main

import (
	"context"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/builder"
	"github.com/timickb/narration-engine/internal/config"
	"github.com/timickb/narration-engine/pkg/utils"
	"os/signal"
	"syscall"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

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
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg, err := config.NewFromFile(ctx, cfgPath)
	if err != nil {
		return fmt.Errorf("parse config from file: %w", err)
	}

	b, err := builder.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("create app builder: %w", err)
	}

	// Запуск асинхронного обработчика
	go func() {
		b.StartInstanceRunner(ctx)
	}()

	// Запуск gRPC сервера
	go func() {
		log.Printf("Starting gRPC server on port %d", b.ServerPort())
		if err = b.ServeGrpc(); err != nil {
			log.Fatalf("Fail to serve grpc: %s", err.Error())
		}
	}()

	// Завершение работы сервиса
	<-ctx.Done()
	log.Info("Stopping engine gracefully..")
	b.WaitGroup().Wait()
	log.Info("Engine stopped.")
	return nil
}
