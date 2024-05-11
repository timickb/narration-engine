package core

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/domain"
)

// Осуществить переход экземпляра сценария в новое состояние.
func (w *AsyncWorker) performTransition(
	ctx context.Context, instance *domain.Instance, event *domain.PendingEvent,
) (result domain.TransitionResult, err error) {

	logger := log.WithContext(ctx).
		WithField("instanceId", instance.Id).
		WithField("pendingEvent", event.EventName).
		WithField("currentState", instance.CurrentState.Name)
	defer func() {
		if err = w.instanceRepo.Update(ctx, instance); err != nil {
			logger.Errorf("Instance update failed: %s", err.Error())
		}
	}()

	// 1. Найти нужный переход в сценарии.
	var transition *domain.Transition
	if event.EventName == domain.EventStart.Name {
		transition = domain.TransitionToStart
	} else {
		for _, t := range instance.Scenario.Transitions {
			if t.From.Name == instance.CurrentState.Name && t.Event.Name == event.EventName {
				transition = t
				break
			}
		}
	}
	if transition == nil {
		logger.Warn("No transition found for event, break execution")
		// Убрать событие из очереди - переход по нему уже никогда не произойдет.
		instance.PendingEvents.Dequeue()
		return domain.TransitionResultBreak, err
	}
	nextState := transition.To

	// 2. Добавить задержку в экземпляр, если она есть у состояния.
	if nextState.Delay > 0 {
		logger.Infof("Next state %s has execution delay", nextState.Name)
		instance.SetDelay(nextState.Delay)
	}

	// 3. Осуществить переход и сохранить его в историю.
	transitionId, err := w.transitionRepo.Save(ctx, &domain.SaveTransitionDto{
		InstanceId:  instance.Id,
		StateFrom:   instance.CurrentState.Name,
		StateTo:     nextState.Name,
		EventName:   event.EventName,
		EventParams: event.EventParams,
	})
	instance.DropRetires()
	instance.PerformTransition(nextState, transitionId)
	instance.PendingEvents.Dequeue()

	// 4. Достигнуто ли терминальное состояние сценария?
	if nextState.Name == domain.StateEnd.Name {
		logger.Infof("Instance has reached terminal state. Break execution.")
		return domain.TransitionResultBreak, err
	}

	// 5. Есть ли обработчик у нового состояния?
	if nextState.Handler != "" {
		logger.Infof("State %s has handler %s", nextState.Name, nextState.Handler)
		return domain.TransitionResultPendingHandler, err
	}
	logger.Infof("No handler for state %s", nextState.Name)

	return domain.TransitionResultCompleted, err
}
