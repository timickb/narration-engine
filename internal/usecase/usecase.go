package usecase

import (
	"github.com/timickb/narration-engine/internal/domain"
)

// Usecase Бизнес-логика, реализующая API сервиса.
type Usecase struct {
	instanceRepo InstanceRepository
	eventRepo    PendingEventRepository
	transactor   domain.Transactor
	config       Config
}

func New(
	instanceRepo InstanceRepository,
	eventRepo PendingEventRepository,
	transactor domain.Transactor,
	config Config,
) *Usecase {
	return &Usecase{
		instanceRepo: instanceRepo,
		eventRepo:    eventRepo,
		transactor:   transactor,
		config:       config,
	}
}
