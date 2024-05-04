package controller

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
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

func MapInstanceToPbState(src *domain.Instance) *schema.State {
	return &schema.State{
		CurrentName:  src.CurrentState.Name,
		PreviousName: src.PreviousState.Name,
		LastEvent:    "",
		Context:      src.Context,
	}
}
