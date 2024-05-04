package domain

import "github.com/google/uuid"

// CreateInstanceDto Структура для создания нового экземпляра сценария.
type CreateInstanceDto struct {
	ParentId        *uuid.UUID
	ScenarioName    string
	ScenarioVersion string
	BlockingKey     *string
	Context         []byte
}
