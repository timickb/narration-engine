package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
)

// InstanceRepository Контракт репозитория над экземплярами сценариев.
type InstanceRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*domain.Instance, error)
	Create(ctx context.Context, dto *domain.CreateInstanceDto) (uuid.UUID, error)
	IsKeyBlocked(ctx context.Context, key string) (bool, error)
}

// PendingEventRepository Контракт репозитория над очередью событий.
type PendingEventRepository interface {
	Create(ctx context.Context, dto *domain.CreatePendingEventDto) (uuid.UUID, error)
}

// Config Контракт конфигурации сервиса.
type Config interface {
	GetLoadedScenario(name string, version string) (*domain.Scenario, error)
}
