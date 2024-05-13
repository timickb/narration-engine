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

type accountAddFundsHandler struct {
	usecase Usecase
}

func NewAccountAddFundsHandler(usecase Usecase) *accountAddFundsHandler {
	return &accountAddFundsHandler{usecase: usecase}
}

func (h *accountAddFundsHandler) Name() string {
	return "payments.account_add_funds"
}

type AccountAddFundsRequest struct {
	AccountId uuid.UUID `json:"account_id"`
	Amount    string    `json:"amount"`
}

func (h *accountAddFundsHandler) HandleState(
	ctx context.Context, req *schema.HandleRequest,
) (*schema.HandleResponse, error) {

	log.Infof("Called %s", h.Name())
	log.Infof("Context received: %s", req.Context)
	log.Infof("Event params received: %s", req.EventParams)

	parsedReq, err := worker.UnmarshallRequestBody[AccountAddFundsRequest](req)
	if err != nil {
		return nil, fmt.Errorf("unmarshall AccountAddFundsRequest: %w", err)
	}

	amount, err := decimal.NewFromString(parsedReq.Amount)
	if err != nil {
		return nil, fmt.Errorf("parse decimal: %w", err)
	}

	if err = h.usecase.AccountAddFunds(ctx, parsedReq.AccountId, amount); err != nil {
		return nil, fmt.Errorf("usecase.AccountAddFunds: %w", err)
	}

	return worker.DefaultContinueResponse()
}
