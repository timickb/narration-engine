package handler

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"github.com/timickb/notifications-example/internal/domain"
)

// Отправка письма пользователю, оплатившему донат для блога
type sendDonationSentMail struct {
	adapter MailAdapter
}

func NewSendDonationSentMail(adapter MailAdapter) *sendDonationSentMail {
	return &sendDonationSentMail{adapter: adapter}
}

func (h *sendDonationSentMail) Name() string {
	return "notifications.send_donation_sent_mail"
}

type SendDonationSentRequest struct {
	BuyerEmail string `json:"buyer_email"`
	BlogName   string `json:"blog_name"`
	Amount     string `json:"amount"`
}

func (h *sendDonationSentMail) HandleState(
	ctx context.Context, req *schema.HandleRequest,
) (*schema.HandleResponse, error) {

	log.Info("Called notifications.send_donation_sent_mail")
	log.Infof("Context received: %s", req.Context)
	log.Infof("Event params received: %s", req.EventParams)

	parsedReq, err := worker.UnmarshallRequestBody[SendDonationSentRequest](req)
	if err != nil {
		return nil, fmt.Errorf("unmarshall SendDonationSentRequest")
	}

	err = h.adapter.SendMail(ctx, &domain.Message{
		MailFrom: "no-reply@service.com",
		MailTo:   parsedReq.BuyerEmail,
		Subject:  "Вы оплатили пожертвование",
		Body:     fmt.Sprintf("Название блога: %s, сумма: %s", parsedReq.BlogName, parsedReq.Amount),
	})
	if err != nil {
		return nil, fmt.Errorf("adapter.SendMail: %w", err)
	}

	return worker.DefaultContinueResponse()
}
