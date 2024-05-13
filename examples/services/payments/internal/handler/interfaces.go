package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/timickb/payments-example/internal/domain"
)

type Usecase interface {
	AccountCreate(ctx context.Context, account *domain.Account) error
	AccountGetById(ctx context.Context, accountId uuid.UUID) (*domain.Account, error)
	AccountAddFunds(ctx context.Context, accountId uuid.UUID, amount decimal.Decimal) error
	AccountRemoveFunds(ctx context.Context, accountId uuid.UUID, amount decimal.Decimal) error
	InvoiceCreate(ctx context.Context, invoice *domain.Invoice) error
	InvoiceGetById(ctx context.Context, invoiceId uuid.UUID) (*domain.Invoice, error)
}
