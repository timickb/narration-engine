package controller

import "github.com/timickb/go-stateflow/internal/domain"

type grpcController struct {
	usecase domain.Usecase
}

func New(usecase domain.Usecase) *grpcController {
	return &grpcController{usecase: usecase}
}
