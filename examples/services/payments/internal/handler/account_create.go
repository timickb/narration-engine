package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"github.com/timickb/payments-example/internal/domain"
	"time"
)

type accountCreateHandler struct {
	usecase Usecase
}

func NewAccountCreateHandler(usecase Usecase) *accountCreateHandler {
	return &accountCreateHandler{usecase: usecase}
}

func (h *accountCreateHandler) Name() string {
	return "payments.account_create"
}

type AccountCreateRequest struct {
	UserId uuid.UUID `json:"user_id,omitempty"`
}
type AccountCreateResponse struct {
	AccountId uuid.UUID `json:"account_id,omitempty"`
}

func (h *accountCreateHandler) HandleState(
	ctx context.Context, req *schema.HandleRequest,
) (*schema.HandleResponse, error) {

	log.Infof("Called %s", h.Name())
	log.Infof("Context received: %s", req.Context)
	log.Infof("Event params received: %s", req.EventParams)

	parsedReq, err := worker.UnmarshallRequestBody[AccountCreateRequest](req)
	if err != nil {
		return nil, fmt.Errorf("unmarshall AccountCreateRequest")
	}

	accountId := uuid.New()
	err = h.usecase.AccountCreate(ctx, &domain.Account{
		Id:        accountId,
		UserId:    parsedReq.UserId,
		Amount:    decimal.Zero,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("usecase.AccountCreate: %w", err)
	}

	return worker.SendEventAndMergeData(
		worker.EventContinue,
		&AccountCreateResponse{AccountId: accountId},
	)
}
