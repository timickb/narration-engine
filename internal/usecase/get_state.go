package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/timickb/go-stateflow/internal/domain"
)

// GetState Получить текущее состояние экземпляра сценария.
func (u *Usecase) GetState(ctx context.Context, instanceId uuid.UUID) (*domain.State, error) {
	panic("implement me")
}
