package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
)

// Start Создать экземпляр сценария.
func (u *Usecase) Start(ctx context.Context, dto *domain.ScenarioStartDto) (uuid.UUID, error) {
	scenario, err := u.config.GetLoadedScenario(dto.ScenarioName, dto.ScenarioVersion)
	if err != nil {
		return uuid.Nil, fmt.Errorf("config.GetLoadedScenario: %w", err)
	}

	instanceId, err := u.instanceRepo.Create(ctx, &domain.CreateInstanceDto{
		ScenarioName:    scenario.Name,
		ScenarioVersion: scenario.Version,
		BlockingKey:     dto.BlockingKey,
		Context:         dto.Context,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("instanceRepo.Create: %w", err)
	}
	return instanceId, nil
}
