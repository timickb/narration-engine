package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/timickb/narration-engine/internal/domain"
	"reflect"
	"testing"
)

var (
	want1 = &domain.Scenario{
		Name:    "test_scenario_1",
		Version: "1.0",
		States:  []*domain.State{testState11, testState12, testState13},
		Transitions: []*domain.Transition{
			{
				From:  domain.StateStart,
				To:    testState11,
				Event: domain.EventContinue,
			},
			{
				From:  testState11,
				To:    testState12,
				Event: domain.EventContinue,
			},
			{
				From:  testState12,
				To:    testState13,
				Event: domain.EventContinue,
			},
			{
				From:  testState13,
				To:    domain.StateEnd,
				Event: domain.EventContinue,
			},
		},
	}

	want2 = &domain.Scenario{
		Name:    "test_scenario_2",
		Version: "1.0",
		States:  []*domain.State{testState2Fetch, testState2Good, testState2Send, testState2Bad},
		Transitions: []*domain.Transition{
			{
				From:  domain.StateStart,
				To:    testState2Fetch,
				Event: domain.EventContinue,
			},
			{
				From:  testState2Fetch,
				To:    testState2Good,
				Event: domain.Event{Name: "yes"},
			},
			{
				From:  testState2Fetch,
				To:    testState2Bad,
				Event: domain.Event{Name: "no"},
			},
			{
				From:  testState2Good,
				To:    testState2Send,
				Event: domain.EventContinue,
			},
			{
				From:  testState2Send,
				To:    domain.StateEnd,
				Event: domain.EventContinue,
			},
			{
				From:  testState2Bad,
				To:    domain.StateEnd,
				Event: domain.Event{Name: "something_performed"},
			},
		},
	}

	want3 = &domain.Scenario{
		Name:    "test_scenario_3",
		Version: "1.0",
		States: []*domain.State{
			testState3FetchOrderInfo,
			testState3CreateInvoice,
			testState3SendSuccessEmail,
			testState3SendFailEmail,
		},
		Transitions: []*domain.Transition{
			{
				From:  domain.StateStart,
				To:    testState3FetchOrderInfo,
				Event: domain.EventContinue,
			},
			{
				From:  testState3FetchOrderInfo,
				To:    testState3CreateInvoice,
				Event: domain.EventContinue,
			},
			{
				From:  testState3FetchOrderInfo,
				To:    testState3SendFailEmail,
				Event: domain.EventHandlerFail,
			},
			{
				From:  testState3CreateInvoice,
				To:    testState3SendSuccessEmail,
				Event: domain.EventContinue,
			},
			{
				From:  testState3CreateInvoice,
				To:    testState3SendFailEmail,
				Event: domain.EventHandlerFail,
			},
			{
				From:  testState3CreateInvoice,
				To:    testState3SendFailEmail,
				Event: domain.Event{Name: "not_enough_money"},
			},
			{
				From:  testState3SendSuccessEmail,
				To:    domain.StateEnd,
				Event: domain.EventContinue,
			},
			{
				From:  testState3SendFailEmail,
				To:    domain.StateEnd,
				Event: domain.EventContinue,
			},
		},
	}
)

func TestParse(t *testing.T) {
	cases := []struct {
		name     string
		filePath string
		want     *domain.Scenario
		wantErr  bool
	}{
		{
			name:     "success simple scenario",
			filePath: "../../examples/test_scenario_1.puml",
			want:     want1,
			wantErr:  false,
		},
		{
			name:     "success scenario with branches",
			filePath: "../../examples/test_scenario_2.puml",
			want:     want2,
			wantErr:  false,
		},
		{
			name:     "success scenario with branches and params",
			filePath: "../../examples/test_scenario_3.puml",
			want:     want3,
			wantErr:  false,
		},
		{
			name:     "publication moderation scenario",
			filePath: "../../examples/moderation_scenario.puml",
			wantErr:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parser := New()
			actual, err := parser.Parse(tc.filePath)

			if tc.wantErr != (err != nil) {
				t.Errorf("Parse() want err = %v, got err = %v", tc.wantErr, err)
			}

			assert.True(t, reflect.DeepEqual(tc.want, actual))
		})
	}
}
