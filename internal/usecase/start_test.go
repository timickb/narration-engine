package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/internal/mocks"
	"github.com/timickb/narration-engine/pkg/utils"
	"go.uber.org/mock/gomock"
	"testing"
)

var (
	testStartDto1 = &domain.ScenarioStartDto{
		ScenarioName:    testScenario.Name,
		ScenarioVersion: testScenario.Version,
		BlockingKey:     nil,
		Context:         nil,
	}
	testStartDto2 = &domain.ScenarioStartDto{
		ScenarioName:    testScenario.Name,
		ScenarioVersion: testScenario.Version,
		BlockingKey:     utils.Ptr("1e0ec0e4-0fe7-11ef-86ad-d7a9093b6773"),
		Context:         nil,
	}

	testCreateInstanceDto1 = &domain.CreateInstanceDto{
		ScenarioName:    testScenario.Name,
		ScenarioVersion: testScenario.Version,
	}
	testCreateInstanceDto2 = &domain.CreateInstanceDto{
		ScenarioName:    testScenario.Name,
		ScenarioVersion: testScenario.Version,
		BlockingKey:     testStartDto2.BlockingKey,
	}

	testCreatePendingEventDto1 = &domain.CreatePendingEventDto{
		InstanceId: testInstanceId1,
		Name:       domain.EventStart.Name,
		Params:     []byte("{}"),
	}

	testStartEventId = uuid.MustParse("38aa6d6c-0fe8-11ef-90c7-d3e74b0a5b1b")
)

func TestStart(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(m *mocks.Mocker)
		dto     *domain.ScenarioStartDto
		wantErr bool
	}{
		{
			name: "success no blocking key",
			setup: func(m *mocks.Mocker) {
				gomock.InOrder(
					m.UsecaseConfigMock.EXPECT().
						GetLoadedScenario(testScenario.Name, testScenario.Version).
						Return(testScenario, nil),
					m.TransactorMock.EXPECT().
						Transaction(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
							return fn(ctx)
						}),
					m.UsecaseInstanceRepositoryMock.EXPECT().
						Create(gomock.Any(), testCreateInstanceDto1).
						Return(testInstanceId1, nil),
					m.UsecasePendingEventRepositoryMock.EXPECT().
						Create(gomock.Any(), testCreatePendingEventDto1).
						Return(testStartEventId, nil),
				)
			},
			dto:     testStartDto1,
			wantErr: false,
		},
		{
			name: "success blocking key not blocked",
			setup: func(m *mocks.Mocker) {
				gomock.InOrder(
					m.UsecaseConfigMock.EXPECT().
						GetLoadedScenario(testScenario.Name, testScenario.Version).
						Return(testScenario, nil),
					m.UsecaseInstanceRepositoryMock.EXPECT().
						IsKeyBlocked(gomock.Any(), *testStartDto2.BlockingKey).
						Return(false, nil),
					m.TransactorMock.EXPECT().
						Transaction(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
							return fn(ctx)
						}),
					m.UsecaseInstanceRepositoryMock.EXPECT().
						Create(gomock.Any(), testCreateInstanceDto2).
						Return(testInstanceId1, nil),
					m.UsecasePendingEventRepositoryMock.EXPECT().
						Create(gomock.Any(), testCreatePendingEventDto1).
						Return(testStartEventId, nil),
				)
			},
			dto:     testStartDto2,
			wantErr: false,
		},
		{
			name: "fail blocking key blocked",
			setup: func(m *mocks.Mocker) {
				gomock.InOrder(
					m.UsecaseConfigMock.EXPECT().
						GetLoadedScenario(testScenario.Name, testScenario.Version).
						Return(testScenario, nil),
					m.UsecaseInstanceRepositoryMock.EXPECT().
						IsKeyBlocked(gomock.Any(), *testStartDto2.BlockingKey).
						Return(true, nil),
				)
			},
			dto:     testStartDto2,
			wantErr: true,
		},
		{
			name: "fail get loaded scenario",
			setup: func(m *mocks.Mocker) {
				gomock.InOrder(
					m.UsecaseConfigMock.EXPECT().
						GetLoadedScenario(testScenario.Name, testScenario.Version).
						Return(nil, errors.New("")),
				)
			},
			dto:     testStartDto2,
			wantErr: true,
		},
		{
			name: "fail no blocking key create instance",
			setup: func(m *mocks.Mocker) {
				gomock.InOrder(
					m.UsecaseConfigMock.EXPECT().
						GetLoadedScenario(testScenario.Name, testScenario.Version).
						Return(testScenario, nil),
					m.TransactorMock.EXPECT().
						Transaction(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
							return fn(ctx)
						}),
					m.UsecaseInstanceRepositoryMock.EXPECT().
						Create(gomock.Any(), testCreateInstanceDto1).
						Return(testInstanceId1, errors.New("")),
				)
			},
			dto:     testStartDto1,
			wantErr: true,
		},
		{
			name: "fail no blocking key create event",
			setup: func(m *mocks.Mocker) {
				gomock.InOrder(
					m.UsecaseConfigMock.EXPECT().
						GetLoadedScenario(testScenario.Name, testScenario.Version).
						Return(testScenario, nil),
					m.TransactorMock.EXPECT().
						Transaction(gomock.Any(), gomock.Any()).
						DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
							return fn(ctx)
						}),
					m.UsecaseInstanceRepositoryMock.EXPECT().
						Create(gomock.Any(), testCreateInstanceDto1).
						Return(testInstanceId1, nil),
					m.UsecasePendingEventRepositoryMock.EXPECT().
						Create(gomock.Any(), testCreatePendingEventDto1).
						Return(testStartEventId, errors.New("")),
				)
			},
			dto:     testStartDto1,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := mocks.NewMocker(t)
			defer m.Finish()
			tc.setup(m)
			uc := createTestUsecase(m)

			_, err := uc.Start(context.Background(), tc.dto)
			if tc.wantErr != (err != nil) {
				t.Fatalf("Start: wantErr = %v, err = %v", tc.wantErr, err)
			}
		})
	}
}
