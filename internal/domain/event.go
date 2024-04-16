package domain

// Event Событие, по которому осуществляется переход между состояниями.
type Event struct {
	Name string
}

var (
	EventStart       Event = Event{Name: "start"}
	EventContinue    Event = Event{Name: "continue"}
	EventJoin        Event = Event{Name: "join"}
	EventParent      Event = Event{Name: "parent"}
	EventChild       Event = Event{Name: "child"}
	EventBreak       Event = Event{Name: "break"}
	EventHandlerFail Event = Event{Name: "handler_fail"}
)
