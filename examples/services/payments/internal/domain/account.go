package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// Account Счет пользователя
type Account struct {
	Id           uuid.UUID
	UserId       uuid.UUID
	Amount       decimal.Decimal
	CreatedAt    time.Time
	LastPayoutAt *time.Time
}
