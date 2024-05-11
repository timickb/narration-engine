package domain

import "time"

// State Состояние сценария.
type State struct {
	Name    string
	Handler string
	Params  map[string]StateParamValue
	Delay   time.Duration
	Retries []time.Duration
}

// StateParamValue Значение параметра состояния.
type StateParamValue struct {
	Value       string
	FromContext bool
}

var (
	// StateStart Терминальное состояние в начале сценария.
	StateStart = &State{
		Name: "START",
	}

	// StateEnd Терминальное состояние в конце сценария.
	StateEnd = &State{
		Name: "END",
	}
)

func (s *State) HasRetries() bool {
	return len(s.Retries) > 0
}

func (s *State) GetNextRetryIfPresents(instance *Instance) (time.Duration, bool) {
	if instance.Retries >= len(s.Retries) {
		return 0, false
	}
	return s.Retries[instance.Retries], true
}
