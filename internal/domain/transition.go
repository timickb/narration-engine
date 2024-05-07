package domain

// Transition Переход между двумя состояниями сценария.
type Transition struct {
	From  *State
	To    *State
	Event Event
}

var (
	// TransitionToStart Переход-заглушка для попадания в терминальное состояние start.
	TransitionToStart = &Transition{
		From:  nil,
		To:    StateStart,
		Event: EventStart,
	}
)
