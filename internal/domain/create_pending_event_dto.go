package domain

import "github.com/google/uuid"

// CreatePendingEventDto Структура для постановки события в очередь.
type CreatePendingEventDto struct {
	InstanceId uuid.UUID
	Name       string
	Params     []byte
	External   bool
}
