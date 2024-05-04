package worker

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

// InstanceRunner Главный воркер сервиса, ответственный за вытягивание экземпляров к выполнению из БД.
type InstanceRunner struct {
	instanceRunnerCfg InstanceRunnerConfig
	asyncWorkerCfg    AsyncWorkerConfig
	instanceRepo      InstanceRepository
	instanceChan      chan uuid.UUID
}

func NewInstanceRunner(
	instanceRunnerCfg InstanceRunnerConfig,
	asyncWorkerCfg AsyncWorkerConfig,
	instanceRepo InstanceRepository,
	instanceChan chan uuid.UUID,
) *InstanceRunner {
	return &InstanceRunner{
		instanceRunnerCfg: instanceRunnerCfg,
		asyncWorkerCfg:    asyncWorkerCfg,
		instanceRepo:      instanceRepo,
		instanceChan:      instanceChan,
	}
}

func (r *InstanceRunner) Start(ctx context.Context) {
	logger := log.WithContext(ctx).WithField("process", "InstanceRunner")
	logger.Info("Instance runner started")

	for i := 0; i < r.instanceRunnerCfg.GetAsyncWorkersCount(); i++ {
		go func(orderNumber int) {
			createAsyncWorker(r.instanceRepo, r.asyncWorkerCfg, r.instanceChan, orderNumber).Start(ctx)
		}(i)
	}

	ticker := time.NewTicker(r.instanceRunnerCfg.GetInstanceFetchingInterval())
	defer ticker.Stop()

	for {
		instancesIds, err := r.fetchWaitingInstancesIds(ctx)
		if err != nil {
			logger.Printf("Failed to fetch instances, err: %s", err.Error())
		}
		logger.Printf("Fetched %d instances to handle", len(instancesIds))
		for _, id := range instancesIds {
			r.instanceChan <- id
		}

		select {
		case <-ctx.Done():
			logger.Info("Received context done")
			return
		case <-ticker.C:
			// Stub
		}
	}
}

func (r *InstanceRunner) fetchWaitingInstancesIds(ctx context.Context) ([]uuid.UUID, error) {
	ids, err := r.instanceRepo.GetWaitingIds(ctx, r.instanceRunnerCfg.GetInstancesBatchSize())
	if err != nil {
		return nil, fmt.Errorf("instanceRepo.GetWaitingIds: %w", err)
	}
	return ids, nil
}
