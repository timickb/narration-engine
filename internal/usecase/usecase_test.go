package usecase

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/internal/mocks"
	"github.com/timickb/narration-engine/pkg/utils"
	"testing"
	"time"
)

var (
	testScenario = &domain.Scenario{
		Name:    "test_scenario",
		Version: "1.0",
	}

	testInstanceId1 = uuid.MustParse("65127a46-0fe5-11ef-92e7-6beef4928c27")

	testInstance1 = &domain.Instance{
		Id:                 testInstanceId1,
		Scenario:           testScenario,
		CurrentState:       &domain.State{Name: "State2"},
		PreviousState:      &domain.State{Name: "State1"},
		Context:            &domain.InstanceContext{},
		Retries:            0,
		Failed:             false,
		CurrentStateStatus: domain.StateStatusWaitingForHandler,
		LockedBy:           utils.Ptr("node"),
		LockedTill:         utils.Ptr(time.Now().Add(time.Minute)),
		BlockingKey:        utils.Ptr("b92f9fb4-0fe5-11ef-ba62-bf13775b208a"),
		CreatedAt:          time.Now().Add(-time.Second * 2),
	}
)

func TestNew(t *testing.T) {
	m := mocks.NewMocker(t)

	uc := createTestUsecase(m)

	assert.NotNil(t, uc)
	assert.NotNil(t, uc.instanceRepo)
	assert.NotNil(t, uc.eventRepo)
	assert.NotNil(t, uc.transactor)
	assert.NotNil(t, uc.config)
}

func createTestUsecase(m *mocks.Mocker) *Usecase {
	return New(
		m.UsecaseInstanceRepositoryMock,
		m.UsecasePendingEventRepositoryMock,
		m.TransactorMock,
		m.UsecaseConfigMock,
	)
}
