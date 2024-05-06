package core

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/domain"
)

func (w *AsyncWorker) startEventLoop(ctx context.Context, id uuid.UUID) error {
	logger := log.WithContext(ctx).WithField("instance", id.String())

	// Взять экземпляр из БД с блокировкой.
	instance, err := w.instanceRepo.FetchWithLock(ctx, &domain.FetchInstanceDto{
		Id:          id,
		LockerId:    w.config.GetLockerId(),
		LockTimeout: w.config.GetInstanceLockTimeout(),
	})
	if err != nil {
		return fmt.Errorf("instanceRepo.FetchWithLock: %w", err)
	}

	defer func() {
		if err := w.instanceRepo.Unlock(ctx, id); err != nil {
			logger.Errorf("Failed to unlock instance: %s", err.Error())
		}
	}()

	// Вытянуть модель сценария.
	scenario, err := w.config.GetLoadedScenario(instance.Scenario.Name, instance.Scenario.Version)
	if err != nil {
		return fmt.Errorf("config.GetLoadedScenario: %w", err)
	}
	instance.Scenario = scenario

	// Если в прошлый раз не удалось выполнить обработчик - вызываем его.
	if instance.CurrentStateStatus == domain.StateStatusWaitingForHandler {
		logger.Info("Last event is not handled yet - finish it...")
		if err := w.executeHandler(ctx, instance); err != nil {
			return fmt.Errorf("execute handler: %w", err)
		}
	}

	logger.Info("Starting event loop")

	pendingEvent := instance.PendingEvents.Front()
	for pendingEvent != nil {
		log.Infof("Event name: %s; event params: %s", pendingEvent.EventName, pendingEvent.EventParams)

		transitionResult, err := w.performTransition(ctx, instance, pendingEvent)
		if err != nil {
			return fmt.Errorf("perfrom transition: %w", err)
		}
		switch transitionResult {
		case domain.TransitionResultBreak:
			break
		case domain.TransitionResultCompleted:
			continue
		case domain.TransitionResultHandlerStarted:
			if err := w.executeHandler(ctx, instance); err != nil {
				return fmt.Errorf("execute handler: %w", err)
			}
		}
	}

	if instance.CurrentState == domain.StateEnd {
		logger.Infof("Instance execution is finished.")
	}

	return nil
}
