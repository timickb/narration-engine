package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
)

// GetState Получить текущее состояние экземпляра сценария.
func (u *Usecase) GetState(ctx context.Context, instanceId uuid.UUID) (*domain.Instance, error) {
	instance, err := u.instanceRepo.GetById(ctx, instanceId)
	if err != nil {
		return nil, fmt.Errorf("instanceRepo.GetById: %w", err)
	}
	return instance, nil
}
