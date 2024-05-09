package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

type StateHandler struct {
	schema.WorkerServiceServer
	handlers map[string]worker.Worker
}

func New(handlers map[string]worker.Worker) *StateHandler {
	return &StateHandler{
		handlers: handlers,
	}
}

// Handle Вызвать нужный обработчик для состояния.
func (h *StateHandler) Handle(ctx context.Context, req *schema.HandleRequest) (*schema.HandleResponse, error) {
	resp := &schema.HandleResponse{Status: &schema.Status{}}

	handler, ok := h.handlers[req.State]
	if !ok {
		return resp, errors.New(fmt.Sprintf("handler for state %s is not registered", req.State))
	}

	return handler.HandleState(ctx, req)
}
