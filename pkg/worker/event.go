package worker

import (
	"encoding/json"
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

func SendEventAndMergeData(event string, data interface{}) (*schema.HandleResponse, error) {
	marshalled, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &schema.HandleResponse{
		Status:           &schema.Status{},
		NextEvent:        event,
		NextEventPayload: "{}",
		DataToContext:    string(marshalled),
	}, nil
}
