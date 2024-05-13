package core

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/internal/domain"
	"sync"
	"time"
)

// InstanceRunner Главный воркер сервиса, ответственный за вытягивание экземпляров к выполнению из БД.
type InstanceRunner struct {
	transactor        domain.Transactor
	instanceRunnerCfg InstanceRunnerConfig
	asyncWorkerCfg    AsyncWorkerConfig
	instanceRepo      InstanceRepository
	transitionRepo    TransitionRepository
	handlerAdapters   map[string]HandlerAdapter
	instanceChan      chan uuid.UUID
	waitGroup         *sync.WaitGroup
}

func NewInstanceRunner(
	instanceRunnerCfg InstanceRunnerConfig,
	asyncWorkerCfg AsyncWorkerConfig,
	transactor domain.Transactor,
	instanceRepo InstanceRepository,
	transitionRepo TransitionRepository,
	handlerAdapters map[string]HandlerAdapter,
	instanceChan chan uuid.UUID,
	waitGroup *sync.WaitGroup,
) *InstanceRunner {
	return &InstanceRunner{
		instanceRunnerCfg: instanceRunnerCfg,
		asyncWorkerCfg:    asyncWorkerCfg,
		transactor:        transactor,
		instanceRepo:      instanceRepo,
		transitionRepo:    transitionRepo,
		handlerAdapters:   handlerAdapters,
		instanceChan:      instanceChan,
		waitGroup:         waitGroup,
	}
}

func (r *InstanceRunner) Start(ctx context.Context) {
	logger := log.WithContext(ctx).WithField("process", "InstanceRunner")
	logger.Info("Instance runner started")

	for i := 0; i < r.instanceRunnerCfg.GetAsyncWorkersCount(); i++ {
		go func(orderNumber int) {
			createAsyncWorker(
				r.transactor,
				r.instanceRepo,
				r.transitionRepo,
				r.asyncWorkerCfg,
				r.handlerAdapters,
				r.instanceChan,
				r.waitGroup,
				orderNumber,
			).Start(ctx)
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
			r.waitGroup.Done()
			return
		case <-ticker.C:
			// Stub
		}
	}
}

func (r *InstanceRunner) Stop(shutdownCtx context.Context) {

}

func (r *InstanceRunner) fetchWaitingInstancesIds(ctx context.Context) ([]uuid.UUID, error) {
	ids, err := r.instanceRepo.GetWaitingIds(ctx, r.instanceRunnerCfg.GetInstancesBatchSize())
	if err != nil {
		return nil, fmt.Errorf("instanceRepo.GetWaitingIds: %w", err)
	}
	return ids, nil
}
