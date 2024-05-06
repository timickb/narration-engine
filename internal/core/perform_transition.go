package core

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/domain"
)

// Осуществить переход экземпляра сценария между состояниями.
func (w *AsyncWorker) performTransition(
	ctx context.Context, instance *domain.Instance, event *domain.PendingEvent,
) (result domain.TransitionResult, err error) {

	logger := log.WithContext(ctx).
		WithField("instanceId", instance.Id).
		WithField("pendingEvent", event.EventName).
		WithField("currentState", instance.CurrentState.Name)

	// В конце нужно обновить экземпляр в БД.
	defer func() {
		err = w.instanceRepo.Update(ctx, instance)
	}()

	// Найти нужный переход в сценарии.
	var transition *domain.Transition
	for _, t := range instance.Scenario.Transitions {
		if t.From.Name == instance.CurrentState.Name && t.Event.Name == event.EventName {
			transition = t
			break
		}
	}
	if transition == nil {
		instance.Failed = true
		logger.Errorf("No transition found in scenario, break execution")
		return domain.TransitionResultBreak, err
	}
	nextState := transition.To

	// Добавить задержку в экземпляр, если она есть у состояния.
	if nextState.Delay > 0 {
		logger.Infof("Next state %s has execution delay", nextState.Name)
		if instance.IsDelayAccomplished() {
			logger.Infof("Delay is already accomplished")
			instance.RemoveDelay()
		} else {
			logger.Info("Break execution due to delay")
			instance.SetDelay(nextState.Delay)
			return domain.TransitionResultBreak, err
		}
	}

	instance.PerformTransition(nextState)
	instance.PendingEvents.Dequeue()

	// Вернуть результат в зависимости от наличия обработчика у состояния.
	if nextState.Handler != "" {
		logger.Infof("State %s has handler %s", nextState.Name, nextState.Handler)
		return domain.TransitionResultHandlerStarted, err
	}
	logger.Infof("No handler for state %s", nextState.Name)

	return domain.TransitionResultCompleted, err
}
