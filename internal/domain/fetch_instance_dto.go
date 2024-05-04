package domain

import (
	"github.com/google/uuid"
	"time"
)

// FetchInstanceDto Данные для взятия экземпляра с блокировкой.
type FetchInstanceDto struct {
	// Id Идентификатор экземпляра.
	Id uuid.UUID
	// LockerId Идентификатор ноды, которая берет блокировку.
	LockerId string
	// LockTimeout Таймаут блокировки.
	LockTimeout time.Duration
}
