package domain

import "context"

// Transactor Контракт исполнителя транзакций БД.
type Transactor interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
