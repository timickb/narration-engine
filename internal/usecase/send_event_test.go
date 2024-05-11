package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/internal/domain"
	"github.com/timickb/narration-engine/internal/mocks"
	"go.uber.org/mock/gomock"
	"testing"
)

var (
	testEventId      = uuid.MustParse("3dd6227a-0fea-11ef-84b0-3372c8e052fe")
	testEventSendDto = &domain.EventSendDto{
		InstanceId:     testInstanceId1,
		Event:          domain.Event{Name: "approve_delivery"},
		PayloadToMerge: []byte("{\"comment\": \"good\"}"),
	}
	testCreatePendingEventDto = &domain.CreatePendingEventDto{
		InstanceId: testEventSendDto.InstanceId,
		Name:       testEventSendDto.Event.Name,
		Params:     testEventSendDto.PayloadToMerge,
		External:   true,
	}
)

func TestSendEvent(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(m *mocks.Mocker)
		dto     *domain.EventSendDto
		wantErr bool
	}{
		{
			name: "success",
			setup: func(m *mocks.Mocker) {
				m.UsecasePendingEventRepositoryMock.EXPECT().
					Create(gomock.Any(), testCreatePendingEventDto).
					Return(testEventId, nil)
			},
			dto:     testEventSendDto,
			wantErr: false,
		},
		{
			name: "fail create",
			setup: func(m *mocks.Mocker) {
				m.UsecasePendingEventRepositoryMock.EXPECT().
					Create(gomock.Any(), testCreatePendingEventDto).
					Return(testEventId, errors.New(""))
			},
			dto:     testEventSendDto,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := mocks.NewMocker(t)
			defer m.Finish()
			tc.setup(m)
			uc := createTestUsecase(m)

			_, err := uc.SendEvent(context.Background(), tc.dto)
			if tc.wantErr != (err != nil) {
				t.Fatalf("SendEvent: wantErr = %v, err = %v", tc.wantErr, err)
			}
		})
	}
}
