package domain

// Scenario Модель сценария.
type Scenario struct {
	// Name Название сценария.
	Name string
	// Version Версия сценария.
	Version string
	// States Набор состояний сценария.
	States []*State
	// Transitions Набор переходов между состояниями.
	Transitions []*Transition
}
