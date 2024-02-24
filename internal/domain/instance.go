package domain

import (
	"github.com/google/uuid"
	"time"
)

// Instance Модель экземпляра сценария.
type Instance struct {
	Id           uuid.UUID
	ScenarioName string

	CurrentState *State

	// BlockingKey Ключ сущности, по которой блокируется экземпляр. Пока он присутствует,
	// не могут создаваться другие экземпляры с таким же ключом.
	BlockingKey *string

	CreatedAt        time.Time
	LastTransitionAt time.Time
}
