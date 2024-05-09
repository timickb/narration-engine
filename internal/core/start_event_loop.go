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

	// Если обработчик состояния еще не выполнен - вызываем его.
	if instance.CurrentStateStatus == domain.StateStatusWaitingForHandler {
		logger.Info("Last event is not handled yet - finish it...")
		lastTransition, err := w.transitionRepo.GetLastForInstance(ctx, instance.Id)
		if err != nil {
			return fmt.Errorf("transitionRepo.GetLastForInstance: %w", err)
		}

		if err := w.executeHandler(ctx, instance, lastTransition.EventName, lastTransition.EventParams); err != nil {
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

		// switch..case не используется, т.к. внутри нужен break для цикла
		if transitionResult == domain.TransitionResultBreak {
			// Прервать цикл событий по одной из причин:
			// 1. Сценарий достиг терминального состояния;
			// 2. Новым состоянием установлена задержка выполнения;
			// 3. Не найден переход из текущего состояния для события.
			break
		} else if transitionResult == domain.TransitionResultCompleted {
			// У нового состояния не было обработчика -> закончить итерацию.
			if err = w.pushEventAndUpdate(ctx, instance, domain.EventContinue); err != nil {
				return fmt.Errorf("push event continue: %w", err)
			}
			continue
		} else if transitionResult == domain.TransitionResultPendingHandler {
			// У нового состояния есть обработчик -> вызвать его.
			handlerErr := w.executeHandler(ctx, instance, pendingEvent.EventName, pendingEvent.EventParams)
			if handlerErr != nil {
				logger.Error("Push event handler_fail due to handler execution error")
				instance.Failed = true
				if err = w.pushEventAndUpdate(ctx, instance, domain.EventHandlerFail); err != nil {
					return fmt.Errorf("push event handler fail: %w", err)
				}
				break
			}
		}
	}

	return nil
}

func (w *AsyncWorker) pushEventAndUpdate(ctx context.Context, instance *domain.Instance, event domain.Event) error {
	return w.transactor.Transaction(ctx, func(ctx context.Context) error {
		instance.PendingEvents.Enqueue(&domain.PendingEvent{
			Id:          uuid.New(),
			EventName:   event.Name,
			EventParams: "{}",
			External:    false,
			FromDb:      false,
			CreatedAt:   time.Now(),
			ExecutedAt:  time.Now(),
		})
		if event == domain.EventHandlerFail && instance.LastTransitionId != nil {
			err := w.transitionRepo.SetError(ctx, *instance.LastTransitionId, "handler failed")
			if err != nil {
				return err
			}
		}
		return w.instanceRepo.Update(ctx, instance)
	})
}
