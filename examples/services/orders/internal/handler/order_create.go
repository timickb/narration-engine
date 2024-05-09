package handler

import (
	"context"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

type orderCreateHandler struct{}

func NewOrderCreateHandler() *orderCreateHandler {
	return &orderCreateHandler{}
}

func (h *orderCreateHandler) Name() string {
	return "orders.create"
}

func (h *orderCreateHandler) HandleState(
	ctx context.Context, req *schema.HandleRequest,
) (*schema.HandleResponse, error) {

	// TODO: implement
	return &schema.HandleResponse{
		Status:           &schema.Status{},
		NextEvent:        "continue",
		NextEventPayload: "{}",
		DataToContext:    "{}",
	}, nil
}
