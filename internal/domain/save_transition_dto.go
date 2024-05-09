package domain

import "github.com/google/uuid"

// SaveTransitionDto Структура для сохранения перехода в историю.
type SaveTransitionDto struct {
	InstanceId  uuid.UUID
	StateFrom   string
	StateTo     string
	EventName   string
	EventParams string
}
