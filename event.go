package ebui

import "github.com/hajimehoshi/ebiten/v2"

type EventType string

const (
	MouseDown  EventType = "mousedown"
	MouseUp    EventType = "mouseup"
	MouseMove  EventType = "mousemove"
	MouseEnter EventType = "mouseenter"
	MouseLeave EventType = "mouseleave"
	Wheel      EventType = "wheel"
	DragStart  EventType = "dragstart"
	Drag       EventType = "drag"
	DragOver   EventType = "dragover"
	DragEnd    EventType = "dragend"
	Drop       EventType = "drop"
	Focus      EventType = "focus"
	Blur       EventType = "blur"
)

type EventPhase int

const (
	PhaseNone    EventPhase = 0
	PhaseCapture EventPhase = 1
	PhaseTarget  EventPhase = 2
	PhaseBubble  EventPhase = 3
)

type Event struct {
	Type                     EventType
	Target                   InteractiveComponent
	CurrentTarget            InteractiveComponent
	RelatedTarget            InteractiveComponent
	MouseX, MouseY           float64
	MouseDeltaX, MouseDeltaY float64
	WheelDeltaX, WheelDeltaY float64
	MouseButton              ebiten.MouseButton
	Timestamp                int64
	Bubbles                  bool
	Phase                    EventPhase
	Path                     []InteractiveComponent
}

// EventBoundary represents a component that controls event propagation
type EventBoundary interface {
	IsWithinBounds(x, y float64) bool
}

// Interactive is an interface that can receive input events
type Interactive interface {
	HandleEvent(event *Event)
	AddEventListener(eventType EventType, handler EventHandler)
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

func (bi *BaseInteractive) AddEventListener(eventType EventType, handler EventHandler) {
	bi.eventDispatcher.AddEventListener(eventType, handler)
}

// EventHandler is a function that handles events
type EventHandler func(event *Event)

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

func (ed *EventDispatcher) DispatchEvent(event *Event) {
	for _, handler := range ed.handlers[event.Type] {
		handler(event)
	}
}
