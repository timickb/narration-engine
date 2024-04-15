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
