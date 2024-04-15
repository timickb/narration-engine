package usecase

import "github.com/timickb/go-stateflow/internal/domain"

type InstanceRunner interface {
}

type StateDiagramParser interface {
	Parse(filePath string) (*domain.Scenario, error)
}

type Usecase struct {
	instanceRunner InstanceRunner
	scenarioParser StateDiagramParser
}

func New(runner InstanceRunner, parser StateDiagramParser) *Usecase {
	return &Usecase{
		instanceRunner: runner,
		scenarioParser: parser,
	}
}
