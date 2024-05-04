package worker

import "time"

// InstanceRunnerConfig Конфигурация для InstanceRunner.
type InstanceRunnerConfig interface {
	GetInstancesBatchSize() int
	GetInstanceFetchingInterval() time.Duration
	GetAsyncWorkersCount() int
}

// AsyncWorkerConfig Конфигурация для AsyncWorker.
type AsyncWorkerConfig interface {
	GetInstanceLockTimeout() time.Duration
	GetLockerId() string
}
