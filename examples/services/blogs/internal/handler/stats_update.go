package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
	BlogId          uuid.UUID `json:"blog_id"`
	IncPublications string    `json:"inc_publications"`
	IncDonations    string    `json:"inc_donations"`
	IncSubscribers  string    `json:"inc_subscribers"`
}

func (h *statsUpdateHandler) HandleState(
	ctx context.Context, req *schema.HandleRequest,
) (*schema.HandleResponse, error) {

	log.Info("Called blogs.stats_update")
	log.Infof("Context received: %s", req.Context)
	log.Infof("Event params received: %s", req.EventParams)

	payload, err := worker.UnmarshallRequestBody[StatsUpdateRequest](req)
	if err != nil {
		return nil, fmt.Errorf("unmarshall stats update request: %w", err)
	}

	blog, err := h.blogUsecase.BlogGetById(ctx, payload.BlogId)
	if err != nil {
		return nil, fmt.Errorf("blogUsecase.BlogGetById: %w", err)
	}

	var incSubscribers, incPublications, incDonations bool
	if payload.IncSubscribers == "true" {
		incSubscribers = true
	}
	if payload.IncPublications == "true" {
		incPublications = true
	}
	if payload.IncDonations == "true" {
		incDonations = true
	}

	err = h.blogUsecase.BlogUpdateStats(ctx, &domain.BlogUpdateStatsDto{
		Id:              blog.Id,
		IncSubscribers:  incSubscribers,
		IntPublications: incPublications,
		IncDonations:    incDonations,
	})
	if err != nil {
		return nil, fmt.Errorf("blogUsecase.BlogUpdateStats: %w", err)
	}

	return worker.DefaultContinueResponse()
}
