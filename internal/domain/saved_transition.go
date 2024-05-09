package domain

// SavedTransition Структура перехода, сохраненного в истории.
type SavedTransition struct {
	EventName   string
	EventParams string
	StateFrom   string
	StateTo     string
	Error       string
}
