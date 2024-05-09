package core

import (
	"context"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
	"time"
)

// AsyncWorkerConfig Конфигурация для AsyncWorker.
type AsyncWorkerConfig interface {
	GetInstanceLockTimeout() time.Duration
	GetLockerId() string
	GetLoadedScenario(name string, version string) (*domain.Scenario, error)
}

// HandlerAdapter Интерфейс адаптера для взаимодействия со внешними обработчиками.
type HandlerAdapter interface {
	CallHandler(ctx context.Context, dto *domain.CallHandlerDto) (*domain.CallHandlerResult, error)
}

// InstanceRepository Репозиторий над экземплярами сценариев.
type InstanceRepository interface {
	FetchWithLock(ctx context.Context, dto *domain.FetchInstanceDto) (*domain.Instance, error)
	GetWaitingIds(ctx context.Context, limit int) ([]uuid.UUID, error)
	Unlock(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, instance *domain.Instance) error
}

// TransitionRepository Репозиторий над историей переходов.
type TransitionRepository interface {
	Save(ctx context.Context, dto *domain.SaveTransitionDto) (uuid.UUID, error)
	SetError(ctx context.Context, transitionId uuid.UUID, errText string) error
	GetLastForInstance(ctx context.Context, instanceId uuid.UUID) (*domain.SavedTransition, error)
}
