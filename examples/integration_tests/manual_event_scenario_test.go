package integration_tests

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	stateflow "github.com/timickb/narration-engine/schema/v1/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

type message struct {
	MailFrom string `json:"mail_from"`
	MailTo   string `json:"mail_to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
}

func TestManualEventScenario(t *testing.T) {
	conn, err := grpc.Dial("localhost:2140", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	client := stateflow.NewStateflowServiceClient(conn)
	ctx := context.Background()

	instanceCtx, err := json.Marshal(map[string]interface{}{
		"blog_id":        "66071ee8-0eeb-11ef-8aa6-c7962b59b80e", // Finance blog
		"publication_id": "f27b7b1a-12da-11ef-89a6-cb3a18c1abb8",
		"moderator_message_data": &message{
			MailFrom: "no-reply@service.com",
			MailTo:   "moderation@service.com",
			Subject:  "Гляньте публикацию тут",
			Body:     "...",
		},
		"user_success_data": &message{
			MailFrom: "no-reply@service.com",
			MailTo:   "somebody1@somewhere.com",
			Subject:  "Вам одобрили публикацию",
			Body:     "...",
		},
		"user_decline_data": &message{
			MailFrom: "no-reply@service.com",
			MailTo:   "moderation@service.com",
			Subject:  "Ваша публикация отклонена модератором",
			Body:     "...",
		},
		"user_rework_data": &message{
			MailFrom: "no-reply@service.com",
			MailTo:   "moderation@service.com",
			Subject:  "Плохо, переделывай",
			Body:     "...",
		},
	})
	assert.NoError(t, err)

	// Стартовать сценарий модерации
	resp, err := client.Start(ctx, &stateflow.StartRequest{
		ScenarioName:    "moderation_scenario",
		ScenarioVersion: "1.0",
		Context:         string(instanceCtx),
	})
	assert.NoError(t, err)

	time.Sleep(time.Second * 5)

	// Подождать и проверить, что он остановился на ожидании события от модератора
	state, err := client.GetState(ctx, &stateflow.GetStateRequest{InstanceId: resp.InstanceId})
	assert.NoError(t, err)
	assert.Equal(t, "ОжиданиеОтветаМодератора", state.State.CurrentName)

	// Послать событие от лица модератора (одобрить публикацию)
	_, err = client.SendEvent(ctx, &stateflow.SendEventRequest{
		InstanceId:  resp.InstanceId,
		Event:       "approve",
		EventParams: "{}",
	})
	assert.NoError(t, err)

	// Подождать еще и проверить, что сценарий завершился
	time.Sleep(time.Second * 5)
	state, err = client.GetState(ctx, &stateflow.GetStateRequest{InstanceId: resp.InstanceId})
	assert.NoError(t, err)
	assert.Equal(t, "END", state.State.CurrentName)
}
