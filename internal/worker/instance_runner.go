package worker

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

// InstanceRunner Главный воркер сервиса, ответственный за вытягивание экземпляров к выполнению из БД.
type InstanceRunner struct {
	config       InstanceRunnerConfig
	instanceRepo InstanceRepository
	instanceChan chan uuid.UUID
}

func NewInstanceRunner(
	config InstanceRunnerConfig,
	instanceRepo InstanceRepository,
	instanceChan chan uuid.UUID,
) *InstanceRunner {
	return &InstanceRunner{
		config:       config,
		instanceRepo: instanceRepo,
		instanceChan: instanceChan,
	}
}

func (r *InstanceRunner) Start(ctx context.Context) {
	for i := 0; i < r.config.AsyncWorkersCount(); i++ {
		go func(orderNumber int) {
			createAsyncWorker(r.instanceChan, orderNumber).Start(ctx)
		}(i)
	}

	ticker := time.NewTicker(r.config.InstanceFetchingInterval())
	defer ticker.Stop()

	for {
		instancesIds, err := r.fetchWaitingInstancesIds(ctx)
		if err != nil {
			log.Printf("[InstanceRunner] Failed to fetch instances, err: %s", err.Error())
		}
		for _, id := range instancesIds {
			r.instanceChan <- id
		}

		select {
		case <-ctx.Done():
			log.Println("[InstanceRunner] received context done")
			return
		case <-ticker.C:
			// Stub
		}
	}
}

func (r *InstanceRunner) fetchWaitingInstancesIds(ctx context.Context) ([]uuid.UUID, error) {
	ids, err := r.instanceRepo.GetWaitingIds(ctx, r.config.InstancesBatchSize())
	if err != nil {
		return nil, fmt.Errorf("instanceRepo.GetWaitingIds: %w", err)
	}
	return ids, nil
}
