package parser

import "github.com/timickb/narration-engine/internal/domain"

// ToDomain - отобразить сценарий в доменную сущность.
func (d *StateDiagram) ToDomain() *domain.Scenario {
	return &domain.Scenario{
		Name:        d.Name,
		Version:     d.Version,
		States:      d.States.States,
		Transitions: d.Transitions.Transitions,
	}
}
