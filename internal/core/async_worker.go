package core

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/domain"
	"time"
)

// HandlerAdapter Интерфейс адаптера для взаимодействия со внешними обработчиками.
type HandlerAdapter interface {
	CallHandler(ctx context.Context, dto *domain.CallHandlerDto) (*domain.CallHandlerResult, error)
}

// InstanceRepository Репозиторий над экземплярами сценариев.
type InstanceRepository interface {
	FetchWithLock(ctx context.Context, dto *domain.FetchInstanceDto) (*domain.Instance, error)
	GetWaitingIds(ctx context.Context, limit int) ([]uuid.UUID, error)
	Unlock(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, instance *domain.Instance) error
}

// AsyncWorkerConfig Конфигурация для AsyncWorker.
type AsyncWorkerConfig interface {
	GetInstanceLockTimeout() time.Duration
	GetLockerId() string
	GetLoadedScenario(name string, version string) (*domain.Scenario, error)
}

type AsyncWorker struct {
	instanceRepo   InstanceRepository
	config         AsyncWorkerConfig
	handlerAdapter HandlerAdapter
	instanceChan   chan uuid.UUID
	// Порядковый номер обработчика у InstanceRunner'а.
	orderNumber int
}

func createAsyncWorker(
	instanceRepo InstanceRepository,
	config AsyncWorkerConfig,
	handlerAdapter HandlerAdapter,
	instanceChan chan uuid.UUID,
	orderNumber int,
) *AsyncWorker {
	return &AsyncWorker{
		instanceRepo:   instanceRepo,
		config:         config,
		handlerAdapter: handlerAdapter,
		instanceChan:   instanceChan,
		orderNumber:    orderNumber,
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
				logger.Info("Failed to handle instance: %s", w.orderNumber, err.Error())
			}
		}
	}
}
