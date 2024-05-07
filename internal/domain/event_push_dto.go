package domain

import "github.com/google/uuid"

// EventPushDto Запрос постановки события в очередь.
type EventPushDto struct {
	Id        uuid.UUID
	EventName string
	Params    string
	External  bool
	FromDb    bool
}
