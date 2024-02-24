package usecase

import (
	"context"
	"github.com/timickb/go-stateflow/internal/domain"
)

// SendEvent Отправить событие в экземпляр сценария.
func (u *Usecase) SendEvent(ctx context.Context, dto *domain.EventSendDto) error {
	panic("implement me")
}
