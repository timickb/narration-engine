package domain

import (
	"github.com/google/uuid"
	"time"
)

// Instance Модель экземпляра сценария.
type Instance struct {
	Id            uuid.UUID `json:"id"`
	Scenario      *Scenario `json:"scenario"`
	CurrentState  *State    `json:"current_state"`
	PreviousState *State    `json:"previous_state"`
	Context       string    `json:"context"`
	Retries       int       `json:"retries"`
	Failed        bool      `json:"failed"`

	LockedBy   *string    `json:"locked_by,omitempty"`
	LockedTill *time.Time `json:"locked_till,omitempty"`

	// BlockingKey Ключ сущности, по которой блокируется экземпляр. Пока он присутствует,
	// не могут создаваться другие экземпляры с таким же ключом.
	BlockingKey *string `json:"blocking_key,omitempty"`

	// PendingEvents Очередь ожидающих обработки событий
	PendingEvents *EventsQueue `json:"pending_events"`

	CreatedAt        time.Time  `json:"created_at"`
	LastTransitionAt *time.Time `json:"last_transition_at"`
}
