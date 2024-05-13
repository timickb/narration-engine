package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/timickb/payments-example/internal/domain"
	"sync"
)

// Usecase Бизнес-логика сервиса платежей.
type Usecase struct {
	invoicesMu sync.Mutex
	accountsMu sync.Mutex

	invoicesRepo map[uuid.UUID]*domain.Invoice
	accountsRepo map[uuid.UUID]*domain.Account
}

func New() *Usecase {
	return &Usecase{
		invoicesRepo: make(map[uuid.UUID]*domain.Invoice),
		accountsRepo: make(map[uuid.UUID]*domain.Account),
	}
}

// AccountGetById Получить счет по идентификатору.
func (u *Usecase) AccountGetById(ctx context.Context, accountId uuid.UUID) (*domain.Account, error) {
	account, ok := u.accountsRepo[accountId]
	if !ok {
		return nil, fmt.Errorf("account %s not found", accountId.String())
	}
	return account, nil
}

func (u *Usecase) AccountCreate(ctx context.Context, account *domain.Account) error {
	u.accountsMu.Lock()
	defer u.accountsMu.Unlock()

	if _, ok := u.accountsRepo[account.Id]; ok {
		return fmt.Errorf("account %s already exists", account.Id.String())
	}
	u.accountsRepo[account.Id] = account
	return nil
}

// AccountAddFunds Увеличить баланс счета.
func (u *Usecase) AccountAddFunds(ctx context.Context, accountId uuid.UUID, amount decimal.Decimal) error {
	u.accountsMu.Lock()
	defer u.accountsMu.Unlock()

	account, ok := u.accountsRepo[accountId]
	if !ok {
		return fmt.Errorf("account %s not found", accountId.String())
	}
	u.accountsRepo[accountId].Amount = account.Amount.Add(amount)

	return nil
}

// AccountRemoveFunds Уменьшить баланс счета.
func (u *Usecase) AccountRemoveFunds(ctx context.Context, accountId uuid.UUID, amount decimal.Decimal) error {
	u.accountsMu.Lock()
	defer u.accountsMu.Unlock()

	account, ok := u.accountsRepo[accountId]
	if !ok {
		return fmt.Errorf("account %s not found", accountId.String())
	}
	u.accountsRepo[accountId].Amount = account.Amount.Sub(amount)

	return nil
}

// InvoiceCreate Создать платеж.
func (u *Usecase) InvoiceCreate(ctx context.Context, invoice *domain.Invoice) error {
	u.invoicesMu.Lock()
	defer u.invoicesMu.Unlock()

	if _, ok := u.invoicesRepo[invoice.Id]; ok {
		return fmt.Errorf("invoice %s already exists", invoice.Id.String())
	}
	return nil
}

// InvoiceGetById Получить платеж по идентификатору.
func (u *Usecase) InvoiceGetById(ctx context.Context, invoiceId uuid.UUID) (*domain.Invoice, error) {
	invoice, ok := u.invoicesRepo[invoiceId]
	if !ok {
		return nil, fmt.Errorf("invoice %s not found", invoiceId.String())
	}
	return invoice, nil
}
