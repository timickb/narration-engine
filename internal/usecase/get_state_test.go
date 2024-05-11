package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/mocks"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGetState(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m *mocks.Mocker)
		instanceId uuid.UUID
		wantErr    bool
	}{
		{
			name: "success",
			setup: func(m *mocks.Mocker) {
				m.UsecaseInstanceRepositoryMock.EXPECT().GetById(gomock.Any(), testInstanceId1).
					Return(testInstance1, nil)
			},
			instanceId: testInstanceId1,
			wantErr:    false,
		},
		{
			name: "fail get by id",
			setup: func(m *mocks.Mocker) {
				m.UsecaseInstanceRepositoryMock.EXPECT().GetById(gomock.Any(), testInstanceId1).
					Return(nil, errors.New(""))
			},
			instanceId: testInstanceId1,
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := mocks.NewMocker(t)
			defer m.Finish()
			tc.setup(m)
			uc := createTestUsecase(m)

			actual, err := uc.GetState(context.Background(), tc.instanceId)
			if tc.wantErr != (err != nil) {
				t.Fatalf("GetState: wantErr = %v, err = %v", tc.wantErr, err)
			}
			if tc.wantErr != (actual == nil) {
				t.Fatalf("GetState: wantErr = %v, actual = %v", tc.wantErr, actual)
			}
		})
	}
}
