package domain

// State Состояние сценария.
type State struct {
	Name    string
	Handler string
	Params  map[string]StateParamValue
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
