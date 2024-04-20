package worker

import "time"

// InstanceRunnerConfig Конфигурация для InstanceRunner.
type InstanceRunnerConfig interface {
	InstancesBatchSize() int
	InstanceFetchingInterval() time.Duration
	AsyncWorkersCount() int
}
