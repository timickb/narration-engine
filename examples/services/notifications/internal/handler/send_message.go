package handler

import (
	"context"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

type sendMessageHandler struct{}

func NewSendMessageHandler() *sendMessageHandler {
	return &sendMessageHandler{}
}

func (h *sendMessageHandler) Name() string {
	return "notifications.send_message"
}

func (h *sendMessageHandler) HandleState(ctx context.Context, req *schema.HandleRequest) (*schema.HandleResponse, error) {
	// TODO: something happens...
	return &schema.HandleResponse{
		Status:    &schema.Status{},
		NextEvent: "continue",
	}, nil
}
