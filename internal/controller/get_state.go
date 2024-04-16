package controller

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	schema "github.com/timickb/go-stateflow/schema/v1/gen"
)

// GetState Получить состояние экземпляра сценария.
func (c *grpcController) GetState(
	ctx context.Context, req *schema.GetStateRequest,
) (*schema.GetStateResponse, error) {

	resp := &schema.GetStateResponse{Status: &schema.Status{}}

	instanceId, err := uuid.Parse(req.InstanceId)
	if err != nil {
		return resp, fmt.Errorf("parse instance_id: %w", err)
	}

	state, err := c.usecase.GetState(ctx, instanceId)
	if err != nil {
		return resp, fmt.Errorf("usecase.GetState: %w", err)
	}

	resp.State = MapStateToPb(state)
	return resp, nil
}
