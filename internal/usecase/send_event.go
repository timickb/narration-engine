package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
)

// SendEvent Отправить событие в экземпляр сценария.
func (u *Usecase) SendEvent(ctx context.Context, dto *domain.EventSendDto) (uuid.UUID, error) {
	eventId, err := u.eventRepo.Create(ctx, &domain.CreatePendingEventDto{
		InstanceId: dto.InstanceId,
		Name:       dto.Event.Name,
		Params:     dto.PayloadToMerge,
		External:   true,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("eventRepo.Create: %w", err)
	}
	return eventId, nil
}
