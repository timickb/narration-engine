package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/examples/services/blogs/internal/domain"
	"github.com/timickb/narration-engine/pkg/utils"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

type publicationUpdateHandler struct {
	blogUsecase domain.BlogUsecase
}

func NewPublicationUpdateHandler(blogUsecase domain.BlogUsecase) *publicationUpdateHandler {
	return &publicationUpdateHandler{
		blogUsecase: blogUsecase,
	}
}

// Name Имя обработчика в сценарии.
func (h *publicationUpdateHandler) Name() string {
	return "blogs.publication_update"
}

type PublicationUpdateRequest struct {
	PublicationId uuid.UUID                `json:"publication_id"`
	Status        domain.PublicationStatus `json:"publication_status"`
}

func (h *publicationUpdateHandler) HandleState(
	ctx context.Context, req *schema.HandleRequest,
) (*schema.HandleResponse, error) {

	data, err := worker.UnmarshallRequestBody[PublicationUpdateRequest](req)
	if err != nil {
		return nil, fmt.Errorf("unmarshall request body: %w", err)
	}

	err = h.blogUsecase.PublicationUpdate(ctx, &domain.PublicationUpdateDto{
		Id:     data.PublicationId,
		Status: utils.Ptr(data.Status),
	})
	if err != nil {
		return nil, fmt.Errorf("blogUsecase.PublicationUpdate: %w", err)
	}

	return &schema.HandleResponse{
		Status:           &schema.Status{},
		NextEvent:        "continue",
		NextEventPayload: "{}",
		DataToContext:    "{}",
	}, nil
}
