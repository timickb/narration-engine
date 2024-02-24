package usecase

type InstanceRunner interface {
}

type PlantUMLParser interface {
}

type Usecase struct {
	instanceRunner InstanceRunner
	plantUMLParser PlantUMLParser
}

func New(runner InstanceRunner, parser PlantUMLParser) *Usecase {
	return &Usecase{
		instanceRunner: runner,
		plantUMLParser: parser,
	}
}
