package core

import (
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/internal/mocks"
	"sync"
)

var (
	testScenario1States = []*domain.State{
		domain.StateStart,
		{
			Name: "state1",
		},
		{
			Name: "state2",
		},
		{
			Name:    "state3",
			Handler: "blogs.stub",
		},
		domain.StateEnd,
	}

	testScenario1 = &domain.Scenario{
		Name:    "scenario1",
		Version: "1.0",
		States:  testScenario1States,
		Transitions: []*domain.Transition{
			{
				From:  testScenario1States[0],
				To:    testScenario1States[1],
				Event: domain.EventContinue,
			},
			{
				From:  testScenario1States[1],
				To:    testScenario1States[2],
				Event: domain.EventContinue,
			},
			{
				From:  testScenario1States[2],
				To:    testScenario1States[3],
				Event: domain.EventContinue,
			},
			{
				From:  testScenario1States[3],
				To:    testScenario1States[4],
				Event: domain.EventContinue,
			},
		},
	}
)

func createdMockedWorker(m *mocks.Mocker) *AsyncWorker {
	return &AsyncWorker{
		transactor:     m.TransactorMock,
		instanceRepo:   m.CoreInstanceRepositoryMock,
		transitionRepo: m.CoreTransitionRepositoryMock,
		config:         m.AsyncWorkerConfigMock,
		handlerAdapters: map[string]HandlerAdapter{
			"notifications": m.HandlerAdapterMock,
			"blogs":         m.HandlerAdapterMock,
			"orders":        m.HandlerAdapterMock,
		},
		instanceChan: make(chan uuid.UUID),
		waitGroup:    &sync.WaitGroup{},
		orderNumber:  1,
	}
}
