package domain

import "github.com/google/uuid"

type CallHandlerDto struct {
	HandlerName string
	StateName   string
	InstanceId  uuid.UUID
	Context     string
	EventName   string
	EventParams string
}
