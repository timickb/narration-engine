package handler

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/timickb/narration-engine/pkg/utils"
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

	handler, ok := h.handlers[req.Handler]
	if !ok {
		fmt.Printf("Handler %s is not registered\n", req.Handler)
		fmt.Printf("Registered handlers: %v\n", utils.MapToKeysSlice(h.handlers))
		return resp, errors.New(fmt.Sprintf("handler %s is not registered", req.Handler))
	}

	return handler.HandleState(ctx, req)
}
