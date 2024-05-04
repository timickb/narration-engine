package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/pkg/db"
	"github.com/timickb/narration-engine/pkg/utils"
	"os"
	"time"
)

// HandlerConfig Конфигурация для внешнего обработчика состояний.
type HandlerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// AsyncWorkerConfig Конфигурация для асинхронного обработчика экземпляров.
type AsyncWorkerConfig struct {
	Count                    int    `json:"count"`
	InstancesBatchSize       int    `json:"instances_batch_size"`
	InstanceFetchingInterval string `json:"instance_fetching_interval"`
	InstanceLockTimeout      string `json:"instance_lock_timeout"`
	LockerId                 string `json:"locker_id"`

	parsedFetchingInterval time.Duration
	parsedLockTimeout      time.Duration
}

// Config Конфигурация сервиса.
type Config struct {
	GrpcPort      int                      `json:"grpc_port" yaml:"grpc_port"`
	ScenariosPath string                   `json:"scenarios_path" yaml:"scenarios_path"`
	Database      *db.PostgresConfig       `json:"database" yaml:"database"`
	Handlers      map[string]HandlerConfig `json:"handlers" yaml:"handlers"`
	AsyncWorker   *AsyncWorkerConfig       `json:"async_worker" yaml:"async_worker"`

	ctx             context.Context
	loadedScenarios map[string]*domain.Scenario
}

// NewFromFile Создать конфиг из YAML файла.
func NewFromFile(ctx context.Context, path string) (*Config, error) {
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
	logger := log.WithContext(c.ctx)

	if c.GrpcPort <= 0 {
		return fmt.Errorf("invalid grpc port value")
	}

	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("invalid database section: %w", err)
	}

	for k, v := range c.Handlers {
		if v.Host == "" {
			return fmt.Errorf("invalid host for worker %s", k)
		}
		if v.Port <= 0 {
			return fmt.Errorf("invalid port for worker %s", v)
		}
	}

	if c.AsyncWorker.Count < 1 {
		return fmt.Errorf("invalid async workers count")
	}
	if c.AsyncWorker.InstancesBatchSize < 1 {
		return fmt.Errorf("invalid instances batch size")
	}
	if c.AsyncWorker.LockerId == "" {
		return fmt.Errorf("empty locker id")
	}

	fetchInterval, err := time.ParseDuration(c.AsyncWorker.InstanceFetchingInterval)
	if err != nil {
		return fmt.Errorf("invalid instance fetching interval: %w", err)
	} else {
		c.AsyncWorker.parsedFetchingInterval = fetchInterval
	}
	lockTimeout, err := time.ParseDuration(c.AsyncWorker.InstanceLockTimeout)
	if err != nil {
		return fmt.Errorf("invalid instance lock timeout: %w", err)
	} else {
		c.AsyncWorker.parsedLockTimeout = lockTimeout
	}

	scenarios, err := readScenarios(c.ScenariosPath)
	if err != nil {
		return fmt.Errorf("read scenarios: %w", err)
	}
	logger.Infof("Loaded %d scenarios from dir %s", len(scenarios), c.ScenariosPath)
	c.loadedScenarios = utils.SliceToMap(scenarios, func(s *domain.Scenario) (string, *domain.Scenario) {
		logger.Infof("Scneario %s, version %s", s.Name, s.Version)
		name := fmt.Sprintf("%s:%s", s.Name, s.Version)
		return name, s
	})

	return nil
}

func (c *Config) GetInstancesBatchSize() int {
	return c.AsyncWorker.InstancesBatchSize
}

func (c *Config) GetAsyncWorkersCount() int {
	return c.AsyncWorker.Count
}

func (c *Config) GetLockerId() string {
	return c.AsyncWorker.LockerId
}

func (c *Config) GetInstanceFetchingInterval() time.Duration {
	return c.AsyncWorker.parsedFetchingInterval
}

func (c *Config) GetInstanceLockTimeout() time.Duration {
	return c.AsyncWorker.parsedLockTimeout
}

func (c *Config) GetLoadedScenario(name string, version string) (*domain.Scenario, error) {
	key := fmt.Sprintf("%s:%s", name, version)
	scenario, ok := c.loadedScenarios[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("scenario %s does not exist", key))
	}
	return scenario, nil
}

func (c *Config) GetHandlerAddr(service string) (string, error) {
	conf, ok := c.Handlers[service]
	if !ok {
		return "", errors.New(fmt.Sprintf("service %s is not registered", service))
	}
	return fmt.Sprintf("%s:%d", conf.Host, conf.Port), nil
}
