package handler

import (
	"context"
	"notifications/internal/domain"
)

// MailAdapter Адаптер почтового сервиса.
type MailAdapter interface {
	SendMail(ctx context.Context, msg *domain.Message) error
}
