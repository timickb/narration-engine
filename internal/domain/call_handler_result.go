package domain

// CallHandlerResult Ответ внешнего обработчика после обработки состояния.
type CallHandlerResult struct {
	NextEvent        Event
	DataToContext    string
	NextEventPayload string
}
