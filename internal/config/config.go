package config

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/timickb/go-stateflow/pkg/db"
	"os"
)

// Config Конфигурация сервиса
type Config struct {
	GrpcPort int                `json:"grpc_port" yaml:"grpc_port"`
	Database *db.PostgresConfig `json:"database" yaml:"database"`
}

func NewFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	cfg := Config{}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshall yaml config: %w", err)
	}
	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.GrpcPort <= 0 {
		return fmt.Errorf("invalid grpc port value")
	}
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("invalid database section: %w", err)
	}
	return nil
}
