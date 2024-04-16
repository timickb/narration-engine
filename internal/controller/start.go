package controller

import (
	"context"
	"fmt"
	schema "github.com/timickb/go-stateflow/schema/v1/gen"
)

// Start - создать и запустить экземпляр сценария.
func (c *grpcController) Start(ctx context.Context, req *schema.StartRequest) (*schema.StartResponse, error) {
	resp := &schema.StartResponse{Status: &schema.Status{}}

	instanceId, err := c.usecase.Start(ctx, MapStartRequestToDomain(req))
	if err != nil {
		return resp, fmt.Errorf("usecase.Start: %w", err)
	}

	resp.InstanceId = instanceId.String()
	return resp, nil
}
