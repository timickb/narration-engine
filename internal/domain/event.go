package domain

// Event Событие, по которому осуществляется переход между состояниями.
type Event struct {
	Name string
}

var (
	EventStart       = Event{Name: "start"}
	EventContinue    = Event{Name: "continue"}
	EventJoin        = Event{Name: "join"}
	EventParent      = Event{Name: "parent"}
	EventChild       = Event{Name: "child"}
	EventBreak       = Event{Name: "break"}
	EventHandlerFail = Event{Name: "handler_fail"}
)
