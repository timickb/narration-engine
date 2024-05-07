package core

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/domain"
	"time"
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
		logger.Error("Instance fetch failed")
		return fmt.Errorf("instanceRepo.FetchWithLock: %w", err)
	}
	defer func() {
		if err = w.instanceRepo.Unlock(ctx, id); err != nil {
			logger.Errorf("Instance unlock failed: %s", err.Error())
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

	for pendingEvent := instance.PendingEvents.Front(); pendingEvent != nil; pendingEvent = instance.PendingEvents.Front() {
		log.Infof("Event name: %s; event params: %s", pendingEvent.EventName, pendingEvent.EventParams)

		transitionResult, err := w.performTransition(ctx, instance, pendingEvent)
		if err != nil {
			return fmt.Errorf("perfrom transition: %w", err)
		}

		switch transitionResult {
		case domain.TransitionResultBreak:
			// Прервать цикл событий по одной из причин:
			// 1. Сценарий достиг терминального состояния;
			// 2. Новым состоянием установлена задержка выполнения;
			// 3. Не найден переход из текущего состояния для события.
			break
		case domain.TransitionResultCompleted:
			// У нового состояния нет обработчика -> сгенерировать событие continue.
			instance.PendingEvents.PushToFront(&domain.PendingEvent{
				Id:          uuid.New(),
				EventName:   domain.EventContinue.Name,
				EventParams: "{}",
				External:    false,
				FromDb:      false,
				CreatedAt:   time.Now(),
				ExecutedAt:  time.Now(),
			})
		case domain.TransitionResultHandlerStarted:
			// У нового состояния есть обработчик -> вызвать его.
			if err := w.executeHandler(ctx, instance); err != nil {
				return fmt.Errorf("execute handler: %w", err)
			}
		}
	}

	return nil
}
