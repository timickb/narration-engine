package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
)

type accountRemoveFundsHandler struct {
	usecase Usecase
}

func NewAccountRemoveFundsHandler(usecase Usecase) *accountRemoveFundsHandler {
	return &accountRemoveFundsHandler{usecase: usecase}
}

func (h *accountRemoveFundsHandler) Name() string {
	return "payments.account_remove_funds"
}

type AccountRemoveFundsRequest struct {
	AccountId uuid.UUID `json:"account_id"`
	Amount    string    `json:"amount"`
}

func (h *accountRemoveFundsHandler) HandleState(
	ctx context.Context, req *schema.HandleRequest,
) (*schema.HandleResponse, error) {

	log.Infof("Called %s", h.Name())
	log.Infof("Context received: %s", req.Context)
	log.Infof("Event params received: %s", req.EventParams)

	parsedReq, err := worker.UnmarshallRequestBody[AccountRemoveFundsRequest](req)
	if err != nil {
		return nil, fmt.Errorf("unmarshall AccountRemoveFundsRequest: %w", err)
	}

	amount, err := decimal.NewFromString(parsedReq.Amount)
	if err != nil {
		return nil, fmt.Errorf("parse decimal: %w", err)
	}

	if err = h.usecase.AccountRemoveFunds(ctx, parsedReq.AccountId, amount); err != nil {
		return nil, fmt.Errorf("usecase.AccountRemoveFunds")
	}

	return worker.DefaultContinueResponse()
}
