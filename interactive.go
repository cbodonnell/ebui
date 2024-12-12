package ebui

// Interactive is an interface for components that can receive input events
type Interactive interface {
	Component
	HandleEvent(event Event)
	GetEventDispatcher() *EventDispatcher
}

// BaseInteractive provides common interactive functionality
type BaseInteractive struct {
	*BaseComponent
	eventDispatcher *EventDispatcher
}

func NewBaseInteractive() BaseInteractive {
	return BaseInteractive{
		BaseComponent:   NewBaseComponent(),
		eventDispatcher: NewEventDispatcher(),
	}
}

func (bi *BaseInteractive) HandleEvent(event Event) {
	bi.eventDispatcher.DispatchEvent(event)
}

func (bi *BaseInteractive) GetEventDispatcher() *EventDispatcher {
	return bi.eventDispatcher
}
