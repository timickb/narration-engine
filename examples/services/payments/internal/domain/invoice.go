package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// Invoice Сущность платежа
type Invoice struct {
	Id          uuid.UUID
	Description string
	Amount      decimal.Decimal
	CreatedAt   time.Time
	Completed   bool
	Failed      bool
}
