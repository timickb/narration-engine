package domain

import (
	"github.com/google/uuid"
	"time"
)

// PendingEvent Событие из очереди на обработку.
type PendingEvent struct {
	// Id Идентификатор события.
	Id uuid.UUID
	// EventName Имя события.
	EventName string
	// EventParams Параметры события.
	EventParams string
	// External Пришло ли событие из API (метод SendEvent)
	External bool

	// CreatedAt Дата постановки в очередь.
	CreatedAt time.Time
	// Executed Дата выполнения.
	ExecutedAt *time.Time

	// Next Следующее событие.
	Next *PendingEvent
}

func (e *PendingEvent) String() string {
	return e.EventName
}
