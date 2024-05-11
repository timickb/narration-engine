package mocks

import (
	"go.uber.org/mock/gomock"
	"testing"
)

type Mocker struct {
	ctrl *gomock.Controller

	TransactorMock                    *TransactorMock
	AsyncWorkerConfigMock             *AsyncWorkerConfigMock
	HandlerAdapterMock                *HandlerAdapterMock
	CoreInstanceRepositoryMock        *CoreInstanceRepositoryMock
	CoreTransitionRepositoryMock      *CoreTransitionRepositoryMock
	UsecaseInstanceRepositoryMock     *UsecaseInstanceRepositoryMock
	UsecasePendingEventRepositoryMock *UsecasePendingEventRepositoryMock
	UsecaseConfigMock                 *UsecaseConfigMock
}

func NewMocker(t *testing.T) *Mocker {
	ctrl := gomock.NewController(t)
	return &Mocker{
		ctrl:                              ctrl,
		TransactorMock:                    NewTransactorMock(ctrl),
		AsyncWorkerConfigMock:             NewAsyncWorkerConfigMock(ctrl),
		HandlerAdapterMock:                NewHandlerAdapterMock(ctrl),
		CoreInstanceRepositoryMock:        NewCoreInstanceRepositoryMock(ctrl),
		CoreTransitionRepositoryMock:      NewCoreTransitionRepositoryMock(ctrl),
		UsecaseInstanceRepositoryMock:     NewUsecaseInstanceRepositoryMock(ctrl),
		UsecasePendingEventRepositoryMock: NewUsecasePendingEventRepositoryMock(ctrl),
		UsecaseConfigMock:                 NewUsecaseConfigMock(ctrl),
	}
}

func (m *Mocker) Finish() {
	m.ctrl.Finish()
}
