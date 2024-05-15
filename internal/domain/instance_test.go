package domain

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/timickb/narration-engine/pkg/utils"
	"testing"
	"time"
)

var (
	testInstance1 = &Instance{
		Id:                 uuid.MustParse("bc3269fe-0afd-11ef-a631-33cb6c3f7c28"),
		CurrentState:       testCurrentState,
		PreviousState:      nil,
		Context:            &InstanceContext{testInstanceContext},
		Retries:            0,
		Failed:             false,
		CurrentStateStatus: StateStatusHandlerExecuted,
	}
	testCurrentState = &State{
		Name:    "some_state",
		Handler: "service.some_handler",
	}
	testNextState = &State{
		Name:    "some_other_state",
		Handler: "serivce.some_other_handler",
		Params: map[string]StateParamValue{
			"masked_pan": {
				Value:       "ctx.order_data.masked_pan",
				FromContext: true,
			},
		},
	}
	testNextStateWithoutHandler = &State{
		Name: "some_other_state_no_handler",
	}
)

func TestInstance_PerformTransition(t *testing.T) {
	testInstance1.PerformTransition(testNextState, uuid.New())

	assert.Equal(t, testCurrentState, testInstance1.PreviousState)
	assert.Equal(t, testNextState, testInstance1.CurrentState)
	assert.Equal(t, StateStatusWaitingForHandler, testInstance1.CurrentStateStatus)

	testInstance1.PerformTransition(testNextStateWithoutHandler, uuid.New())
	assert.Equal(t, testNextState, testInstance1.PreviousState)
	assert.Equal(t, testNextStateWithoutHandler, testInstance1.CurrentState)
	assert.Equal(t, StateStatusHandlerExecuted, testInstance1.CurrentStateStatus)
}

func TestInstance_SetDelay(t *testing.T) {
	testInstance1.SetDelay(time.Hour)
	assert.NotNil(t, testInstance1.startAfter)
	assert.True(t, (*testInstance1.startAfter).After(time.Now()))
}

func TestInstance_IsDelayAccomplished(t *testing.T) {
	testInstance1.startAfter = utils.Ptr(time.Now().Add(time.Hour))
	assert.False(t, testInstance1.IsDelayAccomplished())
	testInstance1.startAfter = utils.Ptr(time.Now().Add(-time.Hour))
	assert.True(t, testInstance1.IsDelayAccomplished())
	testInstance1.startAfter = nil
	assert.False(t, testInstance1.IsDelayAccomplished())
}
