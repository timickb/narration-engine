package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/timickb/go-stateflow/internal/domain"
)

// Start Создать и запустить экземпляр сценария.
func (u *Usecase) Start(ctx context.Context, dto *domain.ScenarioStartDto) (uuid.UUID, error) {
	panic("implement me")
}
