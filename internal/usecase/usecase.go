package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
)

type InstanceRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*domain.Instance, error)
	Create(ctx context.Context, dto *domain.CreateInstanceDto) (uuid.UUID, error)
}

type PendingEventRepository interface {
	Create(ctx context.Context, dto *domain.CreatePendingEventDto) (uuid.UUID, error)
}

type Config interface {
	GetLoadedScenario(name string, version string) (*domain.Scenario, error)
}

type Usecase struct {
	instanceRepo InstanceRepository
	eventRepo    PendingEventRepository
	config       Config
}

func New(instanceRepo InstanceRepository, eventRepo PendingEventRepository, config Config) *Usecase {
	return &Usecase{
		instanceRepo: instanceRepo,
		eventRepo:    eventRepo,
		config:       config,
	}
}
