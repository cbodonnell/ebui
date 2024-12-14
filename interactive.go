package ebui

// Interactive is an interface that can receive input events
type Interactive interface {
	HandleEvent(event Event)
	GetEventDispatcher() *EventDispatcher
}

type InteractiveComponent interface {
	Component
	Interactive
}

// BaseInteractive provides common interactive functionality
type BaseInteractive struct {
	eventDispatcher *EventDispatcher
}

func NewBaseInteractive() *BaseInteractive {
	return &BaseInteractive{
		eventDispatcher: NewEventDispatcher(),
	}
}

func (bi *BaseInteractive) HandleEvent(event Event) {
	bi.eventDispatcher.DispatchEvent(event)
}

func (bi *BaseInteractive) GetEventDispatcher() *EventDispatcher {
	return bi.eventDispatcher
}
