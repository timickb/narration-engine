package domain

import "github.com/google/uuid"

// EventSendDto Структура с данными для отправки события в экземпляр сценария.
type EventSendDto struct {
	InstanceId     uuid.UUID
	Event          Event
	PayloadToMerge []byte
}
