package parser

import "github.com/timickb/narration-engine/internal/domain"

// Transitions Переходы между состояниями.
type Transitions struct {
	Transitions []*domain.Transition
}

// States Состояния.
type States struct {
	States []*domain.State
}

// AddTransitionDto Данные для добавления нового перехода.
type AddTransitionDto struct {
	StateFrom string
	StateTo   string
	Event     string
}

// AddStateDto Данные для добавления нового состояния.
type AddStateDto struct {
	StateName string
	Handler   string
	Delay     string
	Retry     string
	Params    map[string]domain.StateParamValue
}

// Init Инициалзировать список переходов.
func (t *Transitions) Init() {
	t.Transitions = make([]*domain.Transition, 0)
}

// Init Инициализировать список состояний.
func (s *States) Init() {
	s.States = make([]*domain.State, 0)
}
