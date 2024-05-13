package domain

// TransitionResult Результат выполнения перехода в новое состояние.
type TransitionResult string

const (
	// TransitionResultBreak Выполнение экземпляра закончено или прервано.
	TransitionResultBreak = "BREAK"
	// TransitionResultPendingHandler Переход осуществлен, обработчик еще не выполнен.
	TransitionResultPendingHandler = "PENDING_HANDLER"
	// TransitionResultCompleted Переход осуществлен, можно перейти к следующему событию.
	TransitionResultCompleted = "COMPLETED"
)
