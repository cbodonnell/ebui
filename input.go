package ebui

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type InputManager struct {
	lastHoverTarget InteractiveComponent
	dragSource      InteractiveComponent
	isDragging      bool
	lastMouseX      float64
	lastMouseY      float64
	lastUpdateTime  int64
	buttonStates    map[ebiten.MouseButton]bool
}

func NewInputManager() *InputManager {
	return &InputManager{
		lastUpdateTime: time.Now().UnixNano(),
		buttonStates:   make(map[ebiten.MouseButton]bool),
	}
}

// buildEventPath constructs the event path for a given target component.
// It traverses the component hierarchy from the target to the root, collecting
// all interactive components along the way.
func (im *InputManager) buildEventPath(target InteractiveComponent) []InteractiveComponent {
	path := []InteractiveComponent{}
	var current Component = target
	for current != nil {
		if interactive, ok := current.(InteractiveComponent); ok {
			path = append([]InteractiveComponent{interactive}, path...)
		}
		current = current.GetParent()
	}
	return path
}

func findInteractiveComponentAt(root Component, x, y float64) InteractiveComponent {
	if component, ok := findComponentAt[InteractiveComponent](root, x, y); ok {
		return component
	}
	return nil
}

func findScrollableContainerAt(root Component, x, y float64) InteractiveComponent {
	if component, ok := findComponentAt[*ScrollableContainer](root, x, y); ok {
		return component
	}
	return nil
}

// findComponentAt recursively searches for an component of a given type at the given coordinates.
// It traverses the component hierarchy in reverse order to find the topmost component first.
// If an event boundary is encountered, it will stop the search and return nil.
func findComponentAt[T Component](root Component, x, y float64) (T, bool) {
	var zero T

	if root == nil {
		return zero, false
	}

	// Check if this component is an event boundary
	if eb, ok := root.(EventBoundary); ok && !eb.IsWithinBounds(x, y) {
		return zero, false
	}

	// Check children first (reverse order to get topmost first)
	if container, ok := root.(Container); ok {
		children := container.GetChildren()
		for i := len(children) - 1; i >= 0; i-- {
			if component, ok := findComponentAt[T](children[i], x, y); ok {
				return component, true
			}
		}
	}

	// Check if this component is of the desired type
	component, ok := root.(T)
	if !ok {
		return zero, false
	}

	// Check if the point is within the component's bounds
	if !root.Contains(x, y) {
		return zero, false
	}

	return component, true
}

// dispatchEvent dispatches the given event to the target component and its ancestors.
// It traverses the event path in capturing phase, at target phase, and bubbling phase.
// The event is dispatched to each component's event handlers.
func (im *InputManager) dispatchEvent(event *Event) bool {
	if event.Target == nil {
		return true
	}

	event.Path = im.buildEventPath(event.Target)

	// Capturing Phase
	event.Phase = PhaseCapture
	for i := 0; i < len(event.Path)-1; i++ {
		event.Path[i].HandleEvent(event)
	}

	// At Target Phase
	event.Phase = PhaseTarget
	event.Target.HandleEvent(event)

	// Bubbling Phase
	if event.Bubbles {
		event.Phase = PhaseBubble
		for i := len(event.Path) - 2; i >= 0; i-- {
			event.Path[i].HandleEvent(event)
		}
	}

	return true
}

// Update processes input events and dispatches them to the appropriate components.
// It handles mouse button events, mouse movement, wheel events, and drag events.
// The root component is used as the starting point for event propagation.
func (im *InputManager) Update(root Component) {
	currentTime := time.Now().UnixNano()
	x, y := ebiten.CursorPosition()
	fx, fy := float64(x), float64(y)

	deltaX := fx - im.lastMouseX
	deltaY := fy - im.lastMouseY

	target := findInteractiveComponentAt(root, fx, fy)

	// Base event properties
	baseEvent := Event{
		MouseX:      fx,
		MouseY:      fy,
		MouseDeltaX: deltaX,
		MouseDeltaY: deltaY,
		Timestamp:   currentTime,
		Bubbles:     true,
	}

	// Handle pointer input (mouse/touch)
	for _, btn := range []ebiten.MouseButton{
		ebiten.MouseButtonLeft,
		ebiten.MouseButtonRight,
		ebiten.MouseButtonMiddle,
	} {
		wasPressed := im.buttonStates[btn]
		isPressed := ebiten.IsMouseButtonPressed(btn)

		if isPressed != wasPressed {
			evt := baseEvent
			evt.MouseButton = btn
			evt.Target = target

			if isPressed {
				evt.Type = MouseDown
			} else {
				evt.Type = MouseUp
			}

			im.dispatchEvent(&evt)
			im.buttonStates[btn] = isPressed
		}
	}

	// Handle wheel
	wheelX, wheelY := ebiten.Wheel()
	if wheelX != 0 || wheelY != 0 {
		wheelEvent := baseEvent
		wheelEvent.Type = Wheel
		wheelEvent.Target = findScrollableContainerAt(root, fx, fy)
		wheelEvent.WheelDeltaX = float64(wheelX)
		wheelEvent.WheelDeltaY = float64(wheelY)
		im.dispatchEvent(&wheelEvent)
	}

	// Handle hover/pointer movement
	if target != im.lastHoverTarget {
		if im.lastHoverTarget != nil {
			leaveEvent := baseEvent
			leaveEvent.Type = MouseLeave
			leaveEvent.Target = im.lastHoverTarget
			leaveEvent.RelatedTarget = target
			im.dispatchEvent(&leaveEvent)
		}

		if target != nil {
			enterEvent := baseEvent
			enterEvent.Type = MouseEnter
			enterEvent.Target = target
			enterEvent.RelatedTarget = im.lastHoverTarget
			im.dispatchEvent(&enterEvent)
		}

		im.lastHoverTarget = target
	}

	// Handle pointer movement
	if deltaX != 0 || deltaY != 0 {
		moveEvent := baseEvent
		moveEvent.Type = MouseMove
		moveEvent.Target = target
		im.dispatchEvent(&moveEvent)
	}

	// Handle drag events
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !im.isDragging && target != nil {
			dragStartEvent := baseEvent
			dragStartEvent.Type = DragStart
			dragStartEvent.Target = target

			if im.dispatchEvent(&dragStartEvent) {
				im.isDragging = true
				im.dragSource = target
			}
		} else if im.isDragging {
			dragEvent := baseEvent
			dragEvent.Type = Drag
			dragEvent.Target = im.dragSource
			im.dispatchEvent(&dragEvent)

			if target != nil && target != im.dragSource {
				dragOverEvent := baseEvent
				dragOverEvent.Type = DragOver
				dragOverEvent.Target = target
				dragOverEvent.RelatedTarget = im.dragSource
				im.dispatchEvent(&dragOverEvent)
			}
		}
	} else if im.isDragging {
		dragEndEvent := baseEvent
		dragEndEvent.Type = DragEnd
		dragEndEvent.Target = im.dragSource
		im.dispatchEvent(&dragEndEvent)

		if target != nil && target != im.dragSource {
			dropEvent := baseEvent
			dropEvent.Type = Drop
			dropEvent.Target = target
			dropEvent.RelatedTarget = im.dragSource
			im.dispatchEvent(&dropEvent)
		}

		im.isDragging = false
		im.dragSource = nil
	}

	im.lastMouseX = fx
	im.lastMouseY = fy
	im.lastUpdateTime = currentTime
}
