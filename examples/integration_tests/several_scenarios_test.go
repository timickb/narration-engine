package integration_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	stateflow "github.com/timickb/narration-engine/schema/v1/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

func TestSeveralScenarios(t *testing.T) {
	conn, err := grpc.Dial("localhost:2140", grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	client := stateflow.NewStateflowServiceClient(conn)
	ctx := context.Background()

	count := 3

	startedInstances := make([]string, count)

	for i := 0; i < count; i++ {
		userId := uuid.New()

		instanceCtx, _ := json.Marshal(map[string]interface{}{
			"user_id":   userId,
			"blog_name": fmt.Sprintf("Мемы из 201%d года", i),
			"email":     fmt.Sprintf("somebody%d@somewhere.com", i),
		})
		resp, err := client.Start(ctx, &stateflow.StartRequest{
			ScenarioName:    "blog_create",
			ScenarioVersion: "1.0",
			Context:         string(instanceCtx),
		})
		assert.NoError(t, err)
		startedInstances[i] = resp.InstanceId
	}

	time.Sleep(time.Second * 5)

	for _, instanceId := range startedInstances {
		state, err := client.GetState(ctx, &stateflow.GetStateRequest{InstanceId: instanceId})
		assert.NoError(t, err)
		assert.Equal(t, "END", state.State.CurrentName)
	}
}
