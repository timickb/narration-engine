package domain

// StateStatus Статус текущего состояния сценария.
type StateStatus string

const (
	// StateStatusWaitingForHandler Ожидание выполнения обработчика.
	StateStatusWaitingForHandler StateStatus = "WAITING_FOR_HANDLER"
	// StateStatusHandlerExecuted Обработчик выполнен.
	StateStatusHandlerExecuted StateStatus = "HANDLER_EXECUTED"
)
