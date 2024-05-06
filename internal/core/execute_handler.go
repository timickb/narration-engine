package core

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/domain"
)

func (w *AsyncWorker) executeHandler(ctx context.Context, instance *domain.Instance) error {
	logger := log.WithContext(ctx).
		WithField("instance", instance.Id).
		WithField("stateFrom", instance.PreviousState.Name).
		WithField("stateTo", instance.CurrentState.Name).
		WithField("handler", instance.CurrentState.Name)

	nextState := instance.CurrentState

	// 1. Выполнить обработчик.
	logger.Info("Calling handler...")
	result, err := w.handlerAdapter.CallHandler(ctx, &domain.CallHandlerDto{
		HandlerName: nextState.Handler,
		StateName:   nextState.Name,
		InstanceId:  instance.Id,
		Context:     instance.Context.String(),
	})
	if err != nil {
		logger.Error("Handler invocation failed")
		return fmt.Errorf("handlerAdapter.CallHandler: %w", err)
	}

	// 2. Добавить сгенерированные обработчиком данные в контекст экземпляра.
	if err = instance.Context.MergeData([]byte(result.DataToContext)); err != nil {
		return fmt.Errorf("instance.Context.MergeData: %w", err)
	}

	// 3. Добавить сгенерированное обработчиком событие в очередь событий.
	if err = instance.PendingEvents.Enqueue(&domain.EventPushDto{
		EventName: result.NextEvent.Name,
		Params:    result.NextEventPayload,
		External:  false,
	}); err != nil {
		return fmt.Errorf("instance.PendingEvents.Enqueue: %w", err)
	}

	return nil
}
