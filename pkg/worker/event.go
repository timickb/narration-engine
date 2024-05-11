package worker

import (
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

const (
	EventContinue    = "continue"
	EventBreak       = "break"
	EventHandlerFail = "handler_fail"
)

func DefaultContinueResponse() (*schema.HandleResponse, error) {
	return &schema.HandleResponse{
		Status:           &schema.Status{},
		NextEvent:        EventContinue,
		NextEventPayload: "{}",
		DataToContext:    "{}",
	}, nil
}
