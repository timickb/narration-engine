package domain

import "github.com/google/uuid"

// CallHandlerDto Структура для вызова внешнего обработчика.
type CallHandlerDto struct {
	HandlerName string
	StateName   string
	InstanceId  uuid.UUID
	Context     string
	EventName   string
	EventParams string
}
