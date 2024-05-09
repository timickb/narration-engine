package domain

import (
	"github.com/google/uuid"
	"github.com/timickb/narration-engine/pkg/utils"
	"time"
)

// Instance Модель экземпляра сценария.
type Instance struct {
	Id                 uuid.UUID        `json:"id"`
	Scenario           *Scenario        `json:"scenario"`
	CurrentState       *State           `json:"current_state"`
	PreviousState      *State           `json:"previous_state"`
	Context            *InstanceContext `json:"context"`
	Retries            int              `json:"retries"`
	Failed             bool             `json:"failed"`
	CurrentStateStatus StateStatus

	LockedBy   *string    `json:"locked_by,omitempty"`
	LockedTill *time.Time `json:"locked_till,omitempty"`

	// BlockingKey Ключ сущности, по которой блокируется экземпляр. Пока он присутствует,
	// не могут создаваться другие экземпляры с таким же ключом.
	BlockingKey *string `json:"blocking_key,omitempty"`

	// PendingEvents Очередь ожидающих обработки событий
	PendingEvents *EventsQueue `json:"pending_events"`

	CreatedAt        time.Time  `json:"created_at"`
	LastTransitionAt *time.Time `json:"last_transition_at"`
	LastTransitionId *uuid.UUID `json:"last_transition_id"`

	// Задержка выполнения экземпляра. Может быть установлена состоянием.
	startAfter *time.Time
}

// PerformTransition Выполнить переход в новое состояние.
func (i *Instance) PerformTransition(state *State, savedId uuid.UUID) {
	i.LastTransitionAt = utils.Ptr(time.Now())
	i.PreviousState = i.CurrentState
	i.CurrentState = state
	i.LastTransitionId = utils.Ptr(savedId)

	if state.Handler == "" {
		// Если обработчик за состоянием не закреплен - сразу поставить статус в "обработан".
		i.CurrentStateStatus = StateStatusHandlerExecuted
	} else {
		i.CurrentStateStatus = StateStatusWaitingForHandler
	}
}

// SetDelay Установить задержку выполнения.
func (i *Instance) SetDelay(delay time.Duration) {
	i.startAfter = utils.Ptr(time.Now().Add(delay))
}

// IsDelayAccomplished Прошло ли установленное задержкой время.
func (i *Instance) IsDelayAccomplished() bool {
	if i.startAfter != nil {
		return time.Now().After(*i.startAfter)
	}
	return false
}

// RemoveDelay Убрать задержку выполнения.
func (i *Instance) RemoveDelay() {
	i.startAfter = nil
}

// GetStartAfter Получить время окончания задержки.
func (i *Instance) GetStartAfter() *time.Time {
	return i.startAfter
}
