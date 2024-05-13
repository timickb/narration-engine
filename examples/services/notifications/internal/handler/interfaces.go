package handler

import (
	"context"
	"github.com/timickb/notifications-example/internal/domain"
)

// MailAdapter Адаптер почтового сервиса.
type MailAdapter interface {
	SendMail(ctx context.Context, msg *domain.Message) error
}
