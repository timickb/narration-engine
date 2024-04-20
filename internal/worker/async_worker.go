package worker

import (
	"context"
	"github.com/google/uuid"
	"log"
)

type AsyncWorker struct {
	instanceChan chan uuid.UUID
	// Порядковый номер обработчика у InstanceRunner'а.
	orderNumber int
}

func createAsyncWorker(instanceChan chan uuid.UUID, orderNumber int) *AsyncWorker {
	return &AsyncWorker{
		instanceChan: instanceChan,
		orderNumber:  orderNumber,
	}
}

func (w *AsyncWorker) Start(ctx context.Context) {
	log.Printf("[AsyncWorker #%d] Started", w.orderNumber)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[AsyncWorker #%d] Received context done", w.orderNumber)
			return
		case instanceId, ok := <-w.instanceChan:
			if !ok {
				log.Printf("[AsyncWorker #%d] Instance chan closed, stop work", w.orderNumber)
				return
			}
			if err := w.handleInstance(instanceId); err != nil {
				log.Printf("[AsyncWorker #%d], Failed to handle instance: %s", w.orderNumber, err.Error())
			}
		}
	}
}

func (w *AsyncWorker) handleInstance(id uuid.UUID) error {
	// TODO: implement
	panic("implement me")
}
