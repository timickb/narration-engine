package domain

// TransitionResult Результат выполнения перехода в новое состояние.
type TransitionResult string

const (
	// TransitionResultBreak Выполнение экземпляра закончено или прервано.
	TransitionResultBreak = "BREAK"
	// TransitionResultHandlerStarted Переход осуществлен, обработчик еще не выполнен.
	TransitionResultHandlerStarted = "HANDLER_STARTED"
	// TransitionResultCompleted Переход осуществлен, можно перейти к следующему событию.
	TransitionResultCompleted = "COMPLETED"
)
