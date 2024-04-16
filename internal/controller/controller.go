package controller

import (
	"github.com/timickb/go-stateflow/internal/domain"
	schema "github.com/timickb/go-stateflow/schema/v1/gen"
)

type grpcController struct {
	schema.StateflowServiceServer
	usecase domain.Usecase
}

func New(usecase domain.Usecase) *grpcController {
	return &grpcController{usecase: usecase}
}
