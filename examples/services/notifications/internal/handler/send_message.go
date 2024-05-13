package handler

import (
	"context"
	"fmt"
	"github.com/timickb/narration-engine/pkg/worker"
	schema "github.com/timickb/narration-engine/schema/v1/gen"
	"github.com/timickb/notifications-example/internal/domain"
)

type sendMessageHandler struct {
	mailer MailAdapter
}

func NewSendMessageHandler(mailer MailAdapter) *sendMessageHandler {
	return &sendMessageHandler{
		mailer: mailer,
	}
}

// Name Наименование обработчика в сценарии.
func (h *sendMessageHandler) Name() string {
	return "notifications.send_message"
}

// SendMessageRequest Структура с данными из контекста сценария для отправки сообщения.
type SendMessageRequest struct {
	UserMessage *domain.Message `json:"user_message"`
}

func (h *sendMessageHandler) HandleState(ctx context.Context, req *schema.HandleRequest) (*schema.HandleResponse, error) {
	parsed, err := worker.UnmarshallRequestBody[SendMessageRequest](req)
	if err != nil {
		return nil, fmt.Errorf("unmarshall request: %w", err)
	}

	if err := h.mailer.SendMail(ctx, parsed.UserMessage); err != nil {
		return nil, fmt.Errorf("mailer.SendMail: %w", err)
	}

	return &schema.HandleResponse{
		Status:    &schema.Status{},
		NextEvent: worker.EventContinue,
	}, nil
}
