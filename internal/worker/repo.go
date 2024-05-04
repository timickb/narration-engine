package worker

import (
	"context"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
)

// InstanceRepository Репозиторий над экземплярами сценариев.
type InstanceRepository interface {
	FetchWithLock(ctx context.Context, dto *domain.FetchInstanceDto) (*domain.Instance, error)
	GetWaitingIds(ctx context.Context, limit int) ([]uuid.UUID, error)
	Unlock(ctx context.Context, id uuid.UUID) error
}
