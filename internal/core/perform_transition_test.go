package core

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/internal/mocks"
	"github.com/timickb/narration-engine/pkg/utils"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

var (
	emptyInstanceContext, _ = domain.NewInstanceContext([]byte("{}"))

	testInstance1 = &domain.Instance{
		Id:                 uuid.New(),
		Scenario:           testScenario1,
		CurrentState:       testScenario1States[1],
		PreviousState:      domain.StateStart,
		Context:            emptyInstanceContext,
		Retries:            0,
		Failed:             false,
		CurrentStateStatus: domain.StateStatusHandlerExecuted,
		LockedBy:           utils.Ptr("node"),
		LockedTill:         utils.Ptr(time.Now().Add(time.Second * 5)),
		BlockingKey:        nil,
		PendingEvents:      &domain.EventsQueue{},
		CreatedAt:          time.Now().Add(-time.Minute),
	}
	testInstance2 = &domain.Instance{
		Id:                 uuid.New(),
		Scenario:           testScenario1,
		CurrentState:       testScenario1States[2],
		PreviousState:      testScenario1States[1],
		Context:            emptyInstanceContext,
		Retries:            0,
		Failed:             false,
		CurrentStateStatus: domain.StateStatusHandlerExecuted,
		LockedBy:           utils.Ptr("node"),
		LockedTill:         utils.Ptr(time.Now().Add(time.Second * 5)),
		BlockingKey:        nil,
		PendingEvents:      &domain.EventsQueue{},
		CreatedAt:          time.Now().Add(-time.Minute),
	}
	testInstance3 = &domain.Instance{
		Id:                 uuid.New(),
		Scenario:           testScenario1,
		CurrentState:       testScenario1States[3],
		PreviousState:      testScenario1States[2],
		Context:            emptyInstanceContext,
		Retries:            0,
		Failed:             false,
		CurrentStateStatus: domain.StateStatusHandlerExecuted,
		LockedBy:           utils.Ptr("node"),
		LockedTill:         utils.Ptr(time.Now().Add(time.Second * 5)),
		BlockingKey:        nil,
		PendingEvents:      &domain.EventsQueue{},
		CreatedAt:          time.Now().Add(-time.Minute),
	}
	testInstance4 = &domain.Instance{
		Id:                 uuid.New(),
		Scenario:           testScenario1,
		CurrentState:       testScenario1States[1],
		PreviousState:      domain.StateStart,
		Context:            emptyInstanceContext,
		Retries:            0,
		Failed:             false,
		CurrentStateStatus: domain.StateStatusHandlerExecuted,
		LockedBy:           utils.Ptr("node"),
		LockedTill:         utils.Ptr(time.Now().Add(time.Second * 5)),
		BlockingKey:        nil,
		PendingEvents:      &domain.EventsQueue{},
		CreatedAt:          time.Now().Add(-time.Minute),
	}
	testInstance5 = &domain.Instance{
		Id:                 uuid.New(),
		Scenario:           testScenario1,
		CurrentState:       testScenario1States[1],
		PreviousState:      domain.StateStart,
		Context:            emptyInstanceContext,
		Retries:            0,
		Failed:             false,
		CurrentStateStatus: domain.StateStatusHandlerExecuted,
		LockedBy:           utils.Ptr("node"),
		LockedTill:         utils.Ptr(time.Now().Add(time.Second * 5)),
		BlockingKey:        nil,
		PendingEvents:      &domain.EventsQueue{},
		CreatedAt:          time.Now().Add(-time.Minute),
	}

	testPendingEventContinue = &domain.PendingEvent{
		Id:          uuid.New(),
		EventName:   domain.EventContinue.Name,
		EventParams: "{}",
		External:    false,
		FromDb:      false,
		CreatedAt:   time.Now().Add(-time.Second),
		ExecutedAt:  time.Now().Add(-time.Second),
	}
	testPendingEventStub = &domain.PendingEvent{
		Id:          uuid.New(),
		EventName:   "stub",
		EventParams: "{}",
		External:    false,
		FromDb:      false,
		CreatedAt:   time.Now().Add(-time.Second),
		ExecutedAt:  time.Now().Add(-time.Second),
	}
	testPendingEventHandlerFail = &domain.PendingEvent{
		Id:          uuid.New(),
		EventName:   domain.EventHandlerFail.Name,
		EventParams: "{}",
		External:    false,
		FromDb:      false,
		CreatedAt:   time.Now().Add(-time.Second),
		ExecutedAt:  time.Now().Add(-time.Second),
	}
)

func TestPerformTransition(t *testing.T) {
	tests := []struct {
		name                   string
		setup                  func(m *mocks.Mocker)
		instance               *domain.Instance
		event                  *domain.PendingEvent
		wantCurrentStateStatus domain.StateStatus
		wantCurrentState       string
		wantResult             domain.TransitionResult
		wantInstanceFailed     bool
		wantErr                bool
	}{
		{
			name: "success transition to state without handler",
			setup: func(m *mocks.Mocker) {
				m.CoreTransitionRepositoryMock.EXPECT().Save(gomock.Any(), &domain.SaveTransitionDto{
					InstanceId:  testInstance1.Id,
					StateFrom:   testScenario1States[1].Name,
					StateTo:     testScenario1States[2].Name,
					EventName:   testPendingEventContinue.EventName,
					EventParams: testPendingEventContinue.EventParams,
				})
				m.CoreInstanceRepositoryMock.EXPECT().Update(gomock.Any(), testInstance1)
			},
			instance:               testInstance1,
			event:                  testPendingEventContinue,
			wantCurrentStateStatus: domain.StateStatusHandlerExecuted,
			wantCurrentState:       testScenario1States[2].Name,
			wantResult:             domain.TransitionResultCompleted,
			wantInstanceFailed:     false,
			wantErr:                false,
		},
		{
			name: "success transition to state with handler",
			setup: func(m *mocks.Mocker) {
				m.CoreTransitionRepositoryMock.EXPECT().Save(gomock.Any(), &domain.SaveTransitionDto{
					InstanceId:  testInstance2.Id,
					StateFrom:   testScenario1States[2].Name,
					StateTo:     testScenario1States[3].Name,
					EventName:   testPendingEventContinue.EventName,
					EventParams: testPendingEventContinue.EventParams,
				})
				m.CoreInstanceRepositoryMock.EXPECT().Update(gomock.Any(), testInstance2)
			},
			instance:               testInstance2,
			event:                  testPendingEventContinue,
			wantCurrentStateStatus: domain.StateStatusWaitingForHandler,
			wantCurrentState:       testScenario1States[3].Name,
			wantResult:             domain.TransitionResultPendingHandler,
			wantInstanceFailed:     false,
			wantErr:                false,
		},
		{
			name: "success transition to terminal state",
			setup: func(m *mocks.Mocker) {
				m.CoreTransitionRepositoryMock.EXPECT().Save(gomock.Any(), &domain.SaveTransitionDto{
					InstanceId:  testInstance3.Id,
					StateFrom:   testScenario1States[3].Name,
					StateTo:     testScenario1States[4].Name,
					EventName:   testPendingEventContinue.EventName,
					EventParams: testPendingEventContinue.EventParams,
				})
				m.CoreInstanceRepositoryMock.EXPECT().Update(gomock.Any(), testInstance3)
			},
			instance:               testInstance3,
			event:                  testPendingEventContinue,
			wantCurrentStateStatus: domain.StateStatusHandlerExecuted,
			wantCurrentState:       testScenario1States[4].Name,
			wantResult:             domain.TransitionResultBreak,
			wantInstanceFailed:     false,
			wantErr:                false,
		},
		{
			name: "transition not found",
			setup: func(m *mocks.Mocker) {
				m.CoreInstanceRepositoryMock.EXPECT().Update(gomock.Any(), testInstance4)
			},
			instance:               testInstance4,
			event:                  testPendingEventStub,
			wantCurrentStateStatus: domain.StateStatusHandlerExecuted,
			wantCurrentState:       testScenario1States[1].Name,
			wantResult:             domain.TransitionResultBreak,
			wantInstanceFailed:     false,
			wantErr:                false,
		},
		{
			name: "transition not found and event is handler_fail",
			setup: func(m *mocks.Mocker) {
				m.CoreInstanceRepositoryMock.EXPECT().Update(gomock.Any(), testInstance5)
			},
			instance:               testInstance5,
			event:                  testPendingEventHandlerFail,
			wantCurrentStateStatus: domain.StateStatusHandlerExecuted,
			wantCurrentState:       testScenario1States[1].Name,
			wantResult:             domain.TransitionResultBreak,
			wantInstanceFailed:     true,
			wantErr:                false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := mocks.NewMocker(t)
			tc.setup(m)
			defer m.Finish()
			worker := createdMockedWorker(m)

			actual, err := worker.performTransition(context.Background(), tc.instance, tc.event)
			if tc.wantErr != (err != nil) {
				t.Errorf("performTransition() wantErr = %v, err = %v", tc.wantErr, err)
			}

			assert.Equal(t, tc.wantResult, actual)
			assert.Equal(t, tc.wantCurrentState, tc.instance.CurrentState.Name)
			assert.Equal(t, tc.wantCurrentStateStatus, tc.instance.CurrentStateStatus)
			assert.Equal(t, tc.wantInstanceFailed, tc.instance.Failed)
		})
	}
}
