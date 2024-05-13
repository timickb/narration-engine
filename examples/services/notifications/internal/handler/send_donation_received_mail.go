package handler

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"github.com/timickb/notifications-example/internal/domain"
)

type sendDonationReceivedMail struct {
	adapter MailAdapter
}

func NewSendDonationReceivedMail(adapter MailAdapter) *sendDonationReceivedMail {
	return &sendDonationReceivedMail{adapter: adapter}
}

func (h *sendDonationReceivedMail) Name() string {
	return "notifications.send_donation_received_mail"
}

type SendDonationReceivedMailRequest struct {
	AuthorEmail string
	Amount      string
}

func (h *sendDonationReceivedMail) HandleState(
	ctx context.Context, req *schema.HandleRequest,
) (*schema.HandleResponse, error) {

	log.Info("Called notifications.send_donation_received_mail")
	log.Infof("Context received: %s", req.Context)
	log.Infof("Event params received: %s", req.EventParams)

	parsedReq, err := worker.UnmarshallRequestBody[SendDonationReceivedMailRequest](req)
	if err != nil {
		return nil, fmt.Errorf("unmarshall SendDonationReceivedMailRequest")
	}

	err = h.adapter.SendMail(ctx, &domain.Message{
		MailFrom: "no-reply@service.com",
		MailTo:   parsedReq.AuthorEmail,
		Subject:  "Вам пришел новый донат",
		Body:     fmt.Sprintf("На ваш счет зачислено %s р. от нового подписчика", parsedReq.Amount),
	})
	if err != nil {
		return nil, fmt.Errorf("adapter.SendMail: %w", err)
	}

	return worker.DefaultContinueResponse()
}
