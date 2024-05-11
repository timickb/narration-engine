package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/blogs-example/internal/domain"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

type statsUpdateHandler struct {
	blogUsecase domain.BlogUsecase
}

func NewStatsUpdateHandler(blogUsecase domain.BlogUsecase) *statsUpdateHandler {
	return &statsUpdateHandler{blogUsecase: blogUsecase}
}

func (h *statsUpdateHandler) Name() string {
	return "blogs.stats_update"
}

type StatsUpdateRequest struct {
	PublicationId uuid.UUID `json:"publication_id"`
}

func (h *statsUpdateHandler) HandleState(
	ctx context.Context, req *schema.HandleRequest,
) (*schema.HandleResponse, error) {

	payload, err := worker.UnmarshallRequestBody[StatsUpdateRequest](req)
	if err != nil {
		return nil, fmt.Errorf("unmarshall stats update request: %w", err)
	}

	publication, err := h.blogUsecase.PublicationGetById(ctx, payload.PublicationId)
	if err != nil {
		return nil, fmt.Errorf("blogUsecase.PublicationGetById: %w", err)
	}

	err = h.blogUsecase.BlogUpdateStats(ctx, publication.BlogId, true, false)
	if err != nil {
		return nil, fmt.Errorf("blogUsecase.BlogUpdateStats: %w", err)
	}

	return worker.DefaultContinueResponse()
}
