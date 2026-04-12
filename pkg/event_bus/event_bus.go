package event_bus

type Event struct {
	Name string
	Data any
}

const (
	EventCreateTask       = "Create_Task"
	EventDoneTask         = "Update_Task"
	EventDeleteActiveTask = "Delete_Active_Task"
	EventDeleteDoneTask   = "Delete_Done_Task"
)

type EventBus struct {
	Bus chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		Bus: make(chan Event),
	}
}
func (e *EventBus) Publish(event *Event) {
	e.Bus <- *event
}
func (e *EventBus) Subscribe() <-chan Event {
	return e.Bus
}
