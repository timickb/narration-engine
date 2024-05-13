package handler

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"github.com/timickb/notifications-example/internal/domain"
)

// Отправка письма автору блога, которому пришел донат
type sendBlogCreatedMailHandler struct {
	adapter MailAdapter
}

func NewSendBlogCreateMailHandler(adapter MailAdapter) *sendBlogCreatedMailHandler {
	return &sendBlogCreatedMailHandler{adapter: adapter}
}

func (h *sendBlogCreatedMailHandler) Name() string {
	return "notifications.send_blog_created_mail"
}

type SendBlogCreatedMailRequest struct {
	Email    string `json:"email,omitempty"`
	BlogName string `json:"blog_name,omitempty"`
}

func (h *sendBlogCreatedMailHandler) HandleState(
	ctx context.Context, req *schema.HandleRequest,
) (*schema.HandleResponse, error) {

	log.Info("Called notifications.send_blog_created_mail")
	log.Infof("Context received: %s", req.Context)
	log.Infof("Event params received: %s", req.EventParams)

	parsedReq, err := worker.UnmarshallRequestBody[SendBlogCreatedMailRequest](req)
	if err != nil {
		return nil, fmt.Errorf("unmarshall SendBlogCreatedMailRequest")
	}

	err = h.adapter.SendMail(ctx, &domain.Message{
		MailFrom: "no-reply@service.com",
		MailTo:   parsedReq.Email,
		Subject:  "Вы успешно создали блог",
		Body:     fmt.Sprintf("Его название: %s", parsedReq.BlogName),
	})
	if err != nil {
		return nil, fmt.Errorf("adapter.SendMail: %w", err)
	}

	return worker.DefaultContinueResponse()
}
