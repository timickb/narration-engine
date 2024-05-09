package handler

import (
	"context"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

type stubHandler struct{}

func NewStubHandler() *stubHandler {
	return &stubHandler{}
}

func (h *stubHandler) Name() string {
	return "blogs.stub"
}

func (h *stubHandler) HandleState(
	ctx context.Context, req *schema.HandleRequest,
) (*schema.HandleResponse, error) {
	return &schema.HandleResponse{
		Status:           &schema.Status{},
		NextEvent:        "continue",
		NextEventPayload: "{}",
		DataToContext:    "{}",
	}, nil
}
