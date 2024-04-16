package domain

// ScenarioStartDto Структура с данными для старта экземпляра сценария.
type ScenarioStartDto struct {
	ScenarioName    string
	ScenarioVersion string
	BlockingKey     *string
	Context         []byte
}
