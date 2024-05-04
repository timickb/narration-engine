package domain

// EventPushDto Запрос постановки события в очередь.
type EventPushDto struct {
	EventName string
	Params    string
	External  bool
}
