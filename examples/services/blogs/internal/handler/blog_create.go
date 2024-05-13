package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/blogs-example/internal/domain"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"time"
)

type blogCreateHandler struct {
	usecase domain.BlogUsecase
}

func NewBlogCreateHandler(usecase domain.BlogUsecase) *blogCreateHandler {
	return &blogCreateHandler{usecase: usecase}
}

func (h *blogCreateHandler) Name() string {
	return "blogs.blog_create"
}

type BlogCreateRequest struct {
	UserId   uuid.UUID `json:"user_id,omitempty"`
	BlogName string    `json:"blog_name,omitempty"`
	Email    string    `json:"email,omitempty"`
}
type BlogCreateResponse struct {
	BlogId uuid.UUID `json:"blog_id,omitempty"`
}

func (h *blogCreateHandler) HandleState(
	ctx context.Context, req *schema.HandleRequest,
) (*schema.HandleResponse, error) {

	log.Info("Called blogs.blog_create")
	log.Infof("Context received: %s", req.Context)
	log.Infof("Event params received: %s", req.EventParams)

	parsedReq, err := worker.UnmarshallRequestBody[BlogCreateRequest](req)
	if err != nil {
		return nil, fmt.Errorf("unmarshall BlogCreateRequest")
	}

	// 1. Создать блог
	blogId := uuid.New()
	err = h.usecase.BlogCreate(ctx, &domain.Blog{
		Id:                blogId,
		AuthorId:          parsedReq.UserId,
		AuthorEmail:       parsedReq.Email,
		Name:              parsedReq.BlogName,
		SubscribersCount:  0,
		PublicationsCount: 0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("usecase.BlogCreate: %w", err)
	}

	return worker.SendEventAndMergeData(
		worker.EventContinue,
		&BlogCreateResponse{BlogId: blogId},
	)
}
