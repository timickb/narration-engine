package integration_tests

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	stateflow "github.com/timickb/narration-engine/schema/v1/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

func TestSingleScenario(t *testing.T) {
	conn, err := grpc.Dial("localhost:2140", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	client := stateflow.NewStateflowServiceClient(conn)
	ctx := context.Background()
	userId := uuid.New()

	instanceCtx, err := json.Marshal(map[string]interface{}{
		"user_id":   userId,
		"blog_name": "Мемы из 2013 года",
		"email":     "somebody@somewhere.com",
	})
	assert.NoError(t, err)

	resp, err := client.Start(ctx, &stateflow.StartRequest{
		ScenarioName:    "blog_create",
		ScenarioVersion: "1.0",
		Context:         string(instanceCtx),
	})
	assert.NoError(t, err)

	instanceId := resp.InstanceId
	time.Sleep(time.Second * 5)

	state, err := client.GetState(ctx, &stateflow.GetStateRequest{InstanceId: instanceId})
	assert.NoError(t, err)
	assert.Equal(t, "END", state.State.CurrentName)
}
