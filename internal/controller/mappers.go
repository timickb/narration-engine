package controller

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/go-stateflow/internal/domain"
	schema "github.com/timickb/go-stateflow/schema/v1/gen"
)

func MapStartRequestToDomain(src *schema.StartRequest) *domain.ScenarioStartDto {
	return &domain.ScenarioStartDto{
		ScenarioName:    src.ScenarioName,
		ScenarioVersion: src.ScenarioVersion,
		Context:         []byte(src.Context),
		BlockingKey:     src.BlockingKey,
	}
}

func MapSendEventRequestToDomain(src *schema.SendEventRequest) (*domain.EventSendDto, error) {
	instanceId, err := uuid.Parse(src.InstanceId)
	if err != nil {
		return nil, fmt.Errorf("parse instance_id: %w", err)
	}

	return &domain.EventSendDto{
		InstanceId:     instanceId,
		Event:          domain.Event{Name: src.Event},
		PayloadToMerge: []byte(src.EventParams),
	}, nil
}

func MapStateToPb(src *domain.State) *schema.State {
	// TODO
	return &schema.State{}
}
