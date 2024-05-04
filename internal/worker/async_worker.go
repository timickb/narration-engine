package worker

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/domain"
	"time"
)

type AsyncWorker struct {
	instanceRepo InstanceRepository
	config       AsyncWorkerConfig
	instanceChan chan uuid.UUID
	// Порядковый номер обработчика у InstanceRunner'а.
	orderNumber int
}

func createAsyncWorker(
	instanceRepo InstanceRepository,
	config AsyncWorkerConfig,
	instanceChan chan uuid.UUID,
	orderNumber int,
) *AsyncWorker {
	return &AsyncWorker{
		instanceRepo: instanceRepo,
		config:       config,
		instanceChan: instanceChan,
		orderNumber:  orderNumber,
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
			if err := w.handleInstance(ctx, instanceId); err != nil {
				logger.Info("Failed to handle instance: %s", w.orderNumber, err.Error())
			}
		}
	}
}

func (w *AsyncWorker) handleInstance(ctx context.Context, id uuid.UUID) error {
	logger := log.WithContext(ctx).WithField("instance", id.String())
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

	logger.Info("Starting event loop")

	time.Sleep(time.Second * 20)

	logger.Infof("Instance is terminated. State: %s", instance.CurrentState.Name)

	return nil
}
