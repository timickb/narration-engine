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

type invoiceCreateHandler struct {
	usecase Usecase
}

func NewInvoiceCreateHandler(usecase Usecase) *invoiceCreateHandler {
	return &invoiceCreateHandler{usecase: usecase}
}

func (h *invoiceCreateHandler) Name() string {
	return "payments.invoice_create"
}

// InvoiceCreateRequest Данные на вход обработчику
type InvoiceCreateRequest struct {
	Amount      string `json:"amount"`
	Description string `json:"description"`
}

// InvoiceCreateResponse Данные, которые обработчик отдает в контекст экземпляра в качестве результата работы.
type InvoiceCreateResponse struct {
	InvoiceId uuid.UUID `json:"invoice_id"`
}

func (h *invoiceCreateHandler) HandleState(
	ctx context.Context, req *schema.HandleRequest,
) (*schema.HandleResponse, error) {

	log.Infof("Called %s", h.Name())
	log.Infof("Context received: %s", req.Context)
	log.Infof("Event params received: %s", req.EventParams)

	parsedReq, err := worker.UnmarshallRequestBody[InvoiceCreateRequest](req)
	if err != nil {
		return nil, fmt.Errorf("unmarshall InvoiceCreateRequest: %w", err)
	}

	amount, err := decimal.NewFromString(parsedReq.Amount)
	if err != nil {
		return nil, fmt.Errorf("parse invoice amount: %w", err)
	}

	invoiceId := uuid.New()
	err = h.usecase.InvoiceCreate(ctx, &domain.Invoice{
		Id:          invoiceId,
		Description: parsedReq.Description,
		Amount:      amount,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("usecase.InvoiceCreate: %w", err)
	}

	return worker.SendEventAndMergeData(
		worker.EventContinue,
		&InvoiceCreateResponse{InvoiceId: invoiceId},
	)
}
