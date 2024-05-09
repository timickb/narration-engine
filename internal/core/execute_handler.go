package core

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/domain"
	"strings"
	"time"
)

func (w *AsyncWorker) executeHandler(
	ctx context.Context, instance *domain.Instance, eventName, eventParams string,
) error {

	logger := log.WithContext(ctx).
		WithField("instance", instance.Id).
		WithField("stateFrom", instance.PreviousState.Name).
		WithField("stateTo", instance.CurrentState.Name).
		WithField("handler", instance.CurrentState.Name)
	nextState := instance.CurrentState
	handlerParts := strings.Split(nextState.Handler, ".")
	if len(handlerParts) != 2 {
		return fmt.Errorf("invalid handler notation")
	}
	serviceName := handlerParts[0]
	handlerName := handlerParts[1]

	// 1. Найти сервис
	service, ok := w.handlerAdapters[serviceName]
	if !ok {
		return fmt.Errorf("service %s is not registered in handlers", serviceName)
	}

	// 2. Выполнить обработчик.
	logger.Info("Calling handler...")
	result, err := service.CallHandler(ctx, &domain.CallHandlerDto{
		HandlerName: handlerName,
		StateName:   nextState.Name,
		InstanceId:  instance.Id,
		Context:     instance.Context.String(),
		EventParams: eventParams,
		EventName:   eventName,
	})
	if err != nil {
		logger.Error("Handler invocation failed")
		return fmt.Errorf("handlerAdapter.CallHandler: %w", err)
	}
	logger.Info("Handler had invoked with success.")

	// 3. Добавить сгенерированные обработчиком данные в контекст экземпляра.
	if err = instance.Context.MergeData([]byte(result.DataToContext)); err != nil {
		return fmt.Errorf("instance.Context.MergeData: %w", err)
	}

	// 4. Добавить сгенерированное обработчиком событие в очередь событий.
	instance.PendingEvents.PushToFront(&domain.PendingEvent{
		Id:          uuid.New(),
		EventName:   result.NextEvent.Name,
		EventParams: result.NextEventPayload,
		CreatedAt:   time.Now(),
		ExecutedAt:  time.Now(),
		External:    false,
	})

	return nil
}
