package handler

import (
	"context"
	"fmt"
	"github.com/timickb/narration-engine/internal/domain"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

// CallHandler Вызов внешнего обработчика состояния.
func (h *Handler) CallHandler(ctx context.Context, dto *domain.CallHandlerDto) (*domain.CallHandlerResult, error) {
	resp, err := h.client.Handle(ctx, &schema.HandleRequest{
		InstanceId: dto.InstanceId.String(),
		Context:    dto.Context,
		State:      dto.StateName,
	})
	if err != nil {
		return nil, fmt.Errorf("client.Handle: %w", err)
	}

	if resp.Status.Error != nil {
		return nil, fmt.Errorf(
			"client.Handle: Error %d: %s",
			resp.Status.Error.Code,
			resp.Status.Error.Message,
		)
	}

	return &domain.CallHandlerResult{
		NextEvent:        domain.Event{Name: resp.NextEvent},
		DataToContext:    resp.DataToContext,
		NextEventPayload: resp.NextEventPayload,
	}, nil
}
