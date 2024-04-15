package domain

// Transition Переход между двумя состояниями сценария.
type Transition struct {
	From  *State
	To    *State
	Event Event
}
