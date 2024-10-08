package handler

import (
	"context"
	log "github.com/sirupsen/logrus"
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

	log.Info("Called blogs.stub")
	log.Infof("Context received: %s", req.Context)
	log.Infof("Event params received: %s", req.EventParams)

	return &schema.HandleResponse{
		Status:           &schema.Status{},
		NextEvent:        "continue",
		NextEventPayload: "{}",
		DataToContext:    "{}",
	}, nil
}
