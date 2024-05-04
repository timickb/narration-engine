package handler

import (
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

type Handler struct {
	client schema.WorkerServiceClient
	name   string
}

func NewHandlerClient(client schema.WorkerServiceClient, name string) *Handler {
	return &Handler{
		client: client,
		name:   name,
	}
}
