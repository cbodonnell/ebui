package ebui

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
)

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

// HandlerID is a unique identifier for event handlers
type HandlerID string

// Interactive is an interface that can receive input events
type Interactive interface {
	HandleEvent(event *Event)
	AddEventListener(eventType EventType, handler EventHandler) HandlerID
	RemoveEventListener(eventType EventType, handlerID HandlerID)
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

func (bi *BaseInteractive) AddEventListener(eventType EventType, handler EventHandler) HandlerID {
	return bi.eventDispatcher.AddEventListener(eventType, handler)
}

func (bi *BaseInteractive) RemoveEventListener(eventType EventType, handlerID HandlerID) {
	bi.eventDispatcher.RemoveEventListener(eventType, handlerID)
}

// EventHandler is a function that handles events
type EventHandler func(event *Event)

// HandlerEntry represents an event handler with its ID
type HandlerEntry struct {
	ID      HandlerID
	Handler EventHandler
}

// EventDispatcher manages event subscriptions and dispatching
type EventDispatcher struct {
	handlers  map[EventType][]HandlerEntry
	nextID    int
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[EventType][]HandlerEntry),
		nextID:   1,
	}
}

func (ed *EventDispatcher) AddEventListener(eventType EventType, handler EventHandler) HandlerID {
	id := HandlerID(fmt.Sprintf("handler_%d", ed.nextID))
	ed.nextID++
	
	entry := HandlerEntry{
		ID:      id,
		Handler: handler,
	}
	
	ed.handlers[eventType] = append(ed.handlers[eventType], entry)
	return id
}

func (ed *EventDispatcher) RemoveEventListener(eventType EventType, handlerID HandlerID) {
	handlers := ed.handlers[eventType]
	
	for i, entry := range handlers {
		if entry.ID == handlerID {
			ed.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

func (ed *EventDispatcher) DispatchEvent(event *Event) {
	for _, entry := range ed.handlers[event.Type] {
		entry.Handler(event)
	}
}
