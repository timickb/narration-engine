package core

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/internal/mocks"
	"github.com/timickb/narration-engine/pkg/utils"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
	"time"
)

var (
	testInstance6 = &domain.Instance{
		Id:                 uuid.New(),
		Scenario:           testScenario1,
		CurrentState:       testScenario1States[3],
		PreviousState:      testScenario1States[2],
		Context:            emptyInstanceContext,
		Retries:            0,
		Failed:             false,
		CurrentStateStatus: domain.StateStatusWaitingForHandler,
		LockedBy:           utils.Ptr("node"),
		LockedTill:         utils.Ptr(time.Now().Add(time.Second * 5)),
		BlockingKey:        nil,
		PendingEvents:      &domain.EventsQueue{},
		CreatedAt:          time.Now().Add(-time.Minute),
	}
	testInstance7 = &domain.Instance{
		Id:                 uuid.New(),
		Scenario:           testScenario1,
		CurrentState:       testScenario1States[3],
		PreviousState:      testScenario1States[2],
		Context:            emptyInstanceContext,
		Retries:            0,
		Failed:             false,
		CurrentStateStatus: domain.StateStatusWaitingForHandler,
		LockedBy:           utils.Ptr("node"),
		LockedTill:         utils.Ptr(time.Now().Add(time.Second * 5)),
		BlockingKey:        nil,
		PendingEvents:      &domain.EventsQueue{},
		CreatedAt:          time.Now().Add(-time.Minute),
	}
	testCallHandlerResult1 = &domain.CallHandlerResult{
		NextEvent:        domain.Event{Name: "some_event"},
		DataToContext:    "{\"key\": \"value\"}",
		NextEventPayload: "{}",
	}
	testInstanceContext1, _ = domain.NewInstanceContext([]byte("{\"key\": \"value\"}"))
)

func TestExecuteHandler(t *testing.T) {
	tests := []struct {
		name                   string
		setup                  func(m *mocks.Mocker)
		instance               *domain.Instance
		eventName              string
		eventParams            string
		wantErr                bool
		wantCtxErr             bool
		wantContext            *domain.InstanceContext
		wantCurrentStateStatus domain.StateStatus
		wantFrontEvent         string
	}{
		{
			name: "success",
			setup: func(m *mocks.Mocker) {
				m.HandlerAdapterMock.EXPECT().CallHandler(gomock.Any(), &domain.CallHandlerDto{
					HandlerName: testInstance6.CurrentState.Handler,
					StateName:   testInstance6.CurrentState.Name,
					InstanceId:  testInstance6.Id,
					Context:     testInstance6.Context.String(),
					EventName:   "continue",
					EventParams: "{}",
				}).Return(testCallHandlerResult1, nil)

			},
			instance:               testInstance6,
			eventName:              "continue",
			eventParams:            "{}",
			wantErr:                false,
			wantContext:            testInstanceContext1,
			wantCurrentStateStatus: domain.StateStatusHandlerExecuted,
			wantFrontEvent:         testCallHandlerResult1.NextEvent.Name,
		},
		{
			name: "fail handler",
			setup: func(m *mocks.Mocker) {
				m.HandlerAdapterMock.EXPECT().CallHandler(gomock.Any(), &domain.CallHandlerDto{
					HandlerName: testInstance7.CurrentState.Handler,
					StateName:   testInstance7.CurrentState.Name,
					InstanceId:  testInstance7.Id,
					Context:     testInstance7.Context.String(),
					EventName:   "continue",
					EventParams: "{}",
				}).Return(nil, errors.New("failed"))

			},
			instance:               testInstance7,
			eventName:              "continue",
			eventParams:            "{}",
			wantErr:                true,
			wantCtxErr:             true,
			wantContext:            emptyInstanceContext,
			wantCurrentStateStatus: domain.StateStatusWaitingForHandler,
			wantFrontEvent:         "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := mocks.NewMocker(t)
			tc.setup(m)
			defer m.Finish()
			worker := createdMockedWorker(m)

			err := worker.executeHandler(context.Background(), tc.instance, tc.eventName, tc.eventParams)
			if tc.wantErr != (err != nil) {
				t.Errorf("executeHandler() wantErr = %v, err = %v", tc.wantErr, err)
			}
			if tc.wantCtxErr {
				errorMsg, err := tc.instance.Context.GetValue("error_message")
				assert.NoError(t, err)
				assert.NotEmpty(t, errorMsg)
			}
			if tc.wantFrontEvent != "" {
				assert.Equal(t, tc.wantFrontEvent, tc.instance.PendingEvents.Front().EventName)
			}
			assert.Equal(t, tc.wantCurrentStateStatus, tc.instance.CurrentStateStatus)
			reflect.DeepEqual(tc.wantContext, tc.instance.Context)
		})
	}
}
