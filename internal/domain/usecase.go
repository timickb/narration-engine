package domain

import (
	"context"
	"github.com/google/uuid"
)

// Usecase Интерфейс основного юзкейса сервиса.
type Usecase interface {
	Start(ctx context.Context, dto *ScenarioStartDto) (uuid.UUID, error)
	SendEvent(ctx context.Context, dto *EventSendDto) (uuid.UUID, error)
	GetState(ctx context.Context, instanceId uuid.UUID) (*Instance, error)
}
