package worker

import (
	"context"
	"github.com/google/uuid"
)

// InstanceRepository Репозиторий над экземплярами сценариев.
type InstanceRepository interface {
	GetWaitingIds(ctx context.Context, limit int) ([]uuid.UUID, error)
}
