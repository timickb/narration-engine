package core

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/pkg/utils"
	"strings"
	"time"
)

const (
	errorMsgCtxKey = "error_message"
)

func (w *AsyncWorker) executeHandler(
	ctx context.Context, instance *domain.Instance, eventName, eventParams string,
) error {

	logger := log.WithContext(ctx).
		WithField("instance", instance.Id).
		WithField("stateFrom", instance.PreviousState.Name).
		WithField("stateTo", instance.CurrentState.Name).
		WithField("event", eventName).
		WithField("handler", instance.CurrentState.Handler)
	nextState := instance.CurrentState
	serviceName := strings.Split(nextState.Handler, ".")[0]

	// 1. Найти сервис
	service, ok := w.handlerAdapters[serviceName]
	if !ok {
		return fmt.Errorf("service %s is not registered in handlers", serviceName)
	}

	// 2. Выполнить обработчик.
	logger.Info("Calling handler...")
	result, err := service.CallHandler(ctx, &domain.CallHandlerDto{
		HandlerName: nextState.Handler,
		StateName:   nextState.Name,
		InstanceId:  instance.Id,
		Context:     instance.Context.String(),
		EventParams: eventParams,
		EventName:   eventName,
	})
	if err != nil {
		instance.Context.SetRootValue(errorMsgCtxKey, err.Error())
		logger.Errorf("Handler invocation failed: %s", err.Error())
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

	// 5. Установить флаг выполненности обработчика состояния.
	instance.CurrentStateStatus = domain.StateStatusHandlerExecuted

	logger.Debugf("Context after handler execution: %s", instance.Context.String())
	return nil
}

func (w *AsyncWorker) handleHandlerErr(
	ctx context.Context, logger *log.Entry, instance *domain.Instance, handlerErr error,
) (err error, breakLoop bool) {
	retry, retryPresents := instance.CurrentState.GetNextRetryIfPresents(instance)

	if retryPresents {
		logger.Info("Handler failed, but one more retry presents. Use it and delay execution.")
		instance.IncRetries()
		instance.SetDelay(retry)
		if err = w.instanceRepo.Update(ctx, instance); err != nil {
			return fmt.Errorf("instanceRepo.Update: %w", err), false
		}
		return nil, true
	}

	logger.Error("No retries left for state. Push event handler_fail due to handler execution error")
	err = w.pushEventAndUpdate(ctx, instance, domain.EventHandlerFail, utils.Ptr(handlerErr.Error()))
	if err != nil {
		return fmt.Errorf("push event handler fail: %w", err), false
	}
	return nil, false
}
