package parser

import (
	"fmt"
	"github.com/timickb/narration-engine/internal/domain"
	"os"
)

// StateDiagramParser - реализация парсера Plant UML диаграмм.
type StateDiagramParser struct{}

func New() *StateDiagramParser {
	return &StateDiagramParser{}
}

// Parse - прочитать и разобрать сценарий из файла.
func (p *StateDiagramParser) Parse(filePath string) (*domain.Scenario, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	diagram := &StateDiagram{Buffer: string(content)}
	if err = diagram.Init(); err != nil {
		return nil, fmt.Errorf("StateDiagram.Init: %w", err)
	}
	if err = diagram.Parse(); err != nil {
		return nil, fmt.Errorf("StateDiagram.Parse: %w", err)
	}

	diagram.Execute()
	return diagram.ToDomain(), nil
}

// AddState Добавить состояние.
func (d *StateDiagram) AddState(dto *AddStateDto) {
	if d.States.States == nil {
		d.States.Init()
	}

	for _, state := range d.States.States {
		if state.Name == dto.StateName {
			panic(fmt.Sprintf("duplicate state name %s", dto.StateName))
		}
	}

	d.States.States = append(d.States.States, &domain.State{
		Name:    dto.StateName,
		Handler: dto.Handler,
		Params:  dto.Params,
	})
}

// AddTransition Добавить переход между состояниями.
func (d *StateDiagram) AddTransition(dto *AddTransitionDto) {
	var event domain.Event

	if dto.Event == "" {
		event = domain.EventContinue
	} else {
		event = domain.Event{Name: dto.Event}
	}

	var stateFrom, stateTo *domain.State
	for _, state := range d.States.States {
		if state.Name == dto.StateFrom {
			stateFrom = state
		}
		if state.Name == dto.StateTo {
			stateTo = state
		}
	}

	if stateFrom == nil {
		if dto.StateFrom == domain.StateStart.Name {
			stateFrom = domain.StateStart
		} else {
			panic(fmt.Sprintf("stateFrom %s doesn't declared in the scenario", dto.StateFrom))
		}
	}
	if stateTo == nil {
		if dto.StateTo == domain.StateEnd.Name {
			stateTo = domain.StateEnd
		} else {
			panic(fmt.Sprintf("stateTo %s doesn't declared in the scenario", dto.StateTo))
		}
	}

	d.Transitions.Transitions = append(d.Transitions.Transitions, &domain.Transition{
		From:  stateFrom,
		To:    stateTo,
		Event: event,
	})
}

// AddTag Добавить тег.
func (d *StateDiagram) AddTag(tag string) {
	if d.tags == nil {
		d.tags = []string{tag}
		return
	}
	d.tags = append(d.tags, tag)
}

func (d *StateDiagram) AddRetryLabel(name, value string) {

}

func (d *StateDiagram) appendParam() {
	if d.params == nil {
		d.params = make(map[string]domain.StateParamValue)
	}
	if d.contextVarPath == "" {
		// Параметр - конкретное значение
		d.params[d.paramName] = domain.StateParamValue{
			Value:       d.paramValue,
			FromContext: false,
		}
	} else {
		// Параметр - путь к переменной в контексте сценария
		d.params[d.paramName] = domain.StateParamValue{
			Value:       d.contextVarPath,
			FromContext: true,
		}
	}
}

func (d *StateDiagram) setName(name string) {
	d.Name = name
}

func (d *StateDiagram) setVersion(version string) {
	d.Version = version
}

func (d *StateDiagram) clearState() {
	d.stateFrom = ""
	d.stateTo = ""
	d.stateName = ""
	d.eventName = ""
	d.tags = nil
	d.params = nil
	d.contextVarPath = ""
	d.word = ""
	d.delay = ""
	d.handlerName = ""
	d.paramName = ""
	d.paramValue = ""
}
