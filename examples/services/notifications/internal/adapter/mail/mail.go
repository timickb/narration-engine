package mail

import (
	"context"
	log "github.com/sirupsen/logrus"
	"notifications/internal/domain"
)

// Adapter - заглушка.
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) SendMail(ctx context.Context, msg *domain.Message) error {
	log.Info("Sending mail to %s from %s...", msg.MailTo, msg.MailFrom)
	// Something happens...
	log.Info("Mail was sent.")
	return nil
}
