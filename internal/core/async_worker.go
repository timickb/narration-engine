package core

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/domain"
)

type AsyncWorker struct {
	transactor      domain.Transactor
	instanceRepo    InstanceRepository
	transitionRepo  TransitionRepository
	config          AsyncWorkerConfig
	handlerAdapters map[string]HandlerAdapter
	instanceChan    chan uuid.UUID
	// Порядковый номер обработчика у InstanceRunner'а.
	orderNumber int
}

func createAsyncWorker(
	transactor domain.Transactor,
	instanceRepo InstanceRepository,
	transitionRepo TransitionRepository,
	config AsyncWorkerConfig,
	handlerAdapters map[string]HandlerAdapter,
	instanceChan chan uuid.UUID,
	orderNumber int,
) *AsyncWorker {
	return &AsyncWorker{
		transactor:      transactor,
		instanceRepo:    instanceRepo,
		transitionRepo:  transitionRepo,
		config:          config,
		handlerAdapters: handlerAdapters,
		instanceChan:    instanceChan,
		orderNumber:     orderNumber,
	}
}

func (w *AsyncWorker) Start(ctx context.Context) {
	logger := log.WithContext(ctx).WithField("process", fmt.Sprintf("AsyncWorker #%d", w.orderNumber))
	logger.Info("Async worker started")

	for {
		select {
		case <-ctx.Done():
			logger.Info("Received context done", w.orderNumber)
			return
		case instanceId, ok := <-w.instanceChan:
			if !ok {
				logger.Info("Instance chan closed, stop work", w.orderNumber)
				return
			}
			if err := w.startEventLoop(ctx, instanceId); err != nil {
				logger.Errorf("Failed to handle instance: %s", err.Error())
			}
		}
	}
}
