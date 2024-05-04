package controller

import (
	"github.com/timickb/narration-engine/internal/domain"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

type grpcController struct {
	schema.StateflowServiceServer
	usecase domain.Usecase
}

func New(usecase domain.Usecase) *grpcController {
	return &grpcController{usecase: usecase}
}
