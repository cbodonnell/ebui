package ebui

// EventType represents different types of UI events
type EventType int

const (
	EventMouseMove EventType = iota
	EventMouseDown
	EventMouseUp
	EventMouseEnter
	EventMouseLeave
	EventClick
	EventMouseWheel
)

// Event represents a UI event
type Event struct {
	Type      EventType
	X, Y      float64
	Component Component
}

// EventHandler is a function that handles events
type EventHandler func(event Event)

// EventDispatcher manages event subscriptions and dispatching
type EventDispatcher struct {
	handlers map[EventType][]EventHandler
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[EventType][]EventHandler),
	}
}

func (ed *EventDispatcher) AddEventListener(eventType EventType, handler EventHandler) {
	ed.handlers[eventType] = append(ed.handlers[eventType], handler)
}

func (ed *EventDispatcher) RemoveEventListener(eventType EventType, handler EventHandler) {
	handlers := ed.handlers[eventType]
	for i, h := range handlers {
		if &h == &handler {
			ed.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

func (ed *EventDispatcher) DispatchEvent(event Event) {
	for _, handler := range ed.handlers[event.Type] {
		handler(event)
	}
}
