package controller

import (
	"context"
	"fmt"
	schema "github.com/timickb/go-stateflow/schema/v1/gen"
)

// SendEvent - отправить событие в экземпляр.
func (c *grpcController) SendEvent(
	ctx context.Context, req *schema.SendEventRequest,
) (*schema.SendEventResponse, error) {
	resp := &schema.SendEventResponse{Status: &schema.Status{}}

	dto, err := MapSendEventRequestToDomain(req)
	if err != nil {
		return resp, fmt.Errorf("MapSendEventRequestToDomain: %w", err)
	}

	if err = c.usecase.SendEvent(ctx, dto); err != nil {
		return resp, fmt.Errorf("usecase.SendEvent: %w", err)
	}

	return resp, nil
}
