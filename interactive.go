package ebui

// Interactive is an interface that can receive input events
type Interactive interface {
	HandleEvent(event *Event)
}

// InteractiveComponent is an interface that combines the Component and Interactive interfaces
type InteractiveComponent interface {
	Component
	Interactive
}

// BaseInteractive is a base struct that implements the Interactive interface
type BaseInteractive struct {
	eventDispatcher *EventDispatcher
}

func NewBaseInteractive() *BaseInteractive {
	return &BaseInteractive{
		eventDispatcher: NewEventDispatcher(),
	}
}

func (bi *BaseInteractive) HandleEvent(event *Event) {
	bi.eventDispatcher.DispatchEvent(event)
}
