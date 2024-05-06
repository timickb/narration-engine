package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/domain"
)

// Start Создать экземпляр сценария.
func (u *Usecase) Start(ctx context.Context, dto *domain.ScenarioStartDto) (uuid.UUID, error) {
	logger := log.WithContext(ctx).WithField("start", dto.ScenarioName+":"+dto.ScenarioVersion)

	scenario, err := u.config.GetLoadedScenario(dto.ScenarioName, dto.ScenarioVersion)
	if err != nil {
		return uuid.Nil, fmt.Errorf("config.GetLoadedScenario: %w", err)
	}

	if dto.BlockingKey != nil {
		blocked, err := u.instanceRepo.IsKeyBlocked(ctx, *dto.BlockingKey)
		if err != nil {
			return uuid.Nil, fmt.Errorf("instanceRepo.IsKeyBlocked: %w", err)
		}
		if blocked {
			return uuid.Nil, fmt.Errorf("key %s is blocked by another instance", *dto.BlockingKey)
		}
	}

	var instanceId, startEventId uuid.UUID

	err = u.transactor.Transaction(ctx, func(ctx context.Context) error {
		// Создать запись экземпляра.
		instanceId, err = u.instanceRepo.Create(ctx, &domain.CreateInstanceDto{
			ScenarioName:    scenario.Name,
			ScenarioVersion: scenario.Version,
			BlockingKey:     dto.BlockingKey,
			Context:         dto.Context,
		})
		if err != nil {
			return fmt.Errorf("instanceRepo.Create: %w", err)
		}

		// Создать запись события Start.
		startEventId, err = u.eventRepo.Create(ctx, &domain.CreatePendingEventDto{
			InstanceId: instanceId,
			Name:       domain.EventStart.Name,
			Params:     []byte("{}"),
		})
		if err != nil {
			return fmt.Errorf("eventRepo.Create: %w", err)
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("transactor.Transaction: %w", err)
	}

	logger.Infof(
		"Created new instance (id=%s, startEventId=%s)",
		instanceId.String(),
		startEventId.String(),
	)

	return instanceId, nil
}
