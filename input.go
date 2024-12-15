package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type InputManager struct {
	lastMousePressed bool
	lastX, lastY     float64
	hovered          []InteractiveComponent
}

func NewInputManager() *InputManager {
	return &InputManager{
		hovered: make([]InteractiveComponent, 0),
	}
}

// Update handles all input events for the frame
func (u *InputManager) Update(root Component) {
	x, y := ebiten.CursorPosition()
	fx, fy := float64(x), float64(y)
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	wheelX, wheelY := ebiten.Wheel()

	var events []*Event
	if !u.lastMousePressed && mousePressed {
		events = append(events, &Event{Type: EventMouseDown, X: fx, Y: fy})
	}
	if u.lastMousePressed && !mousePressed {
		events = append(events, &Event{Type: EventMouseUp, X: fx, Y: fy})
	}

	if u.lastX != fx || u.lastY != fy {
		events = append(events, &Event{Type: EventMouseMove, X: fx, Y: fy})
	}

	if wheelX != 0 || wheelY != 0 {
		data := &MouseWheelEvent{WheelX: wheelX, WheelY: wheelY}
		events = append(events, &Event{Type: EventMouseWheel, X: fx, Y: fy, Data: data})
	}

	// Find all interactive components at the current mouse position
	components := u.findInteractiveComponentsAt(root, fx, fy)

	// Emit all events
	for _, event := range events {
		u.emitEvent(components, event)
	}

	// Detect hover changes and emit enter/leave events
	entered, left := u.detectHoverChanges(components)
	u.emitEvent(entered, &Event{Type: EventMouseEnter, X: fx, Y: fy})
	u.emitEvent(left, &Event{Type: EventMouseLeave, X: fx, Y: fy})

	u.lastMousePressed = mousePressed
	u.lastX, u.lastY = fx, fy
	u.hovered = components
}

// detectHoverChanges compares the currently hovered components with the previously hovered
// components and returns the components that were entered and left.
func (u *InputManager) detectHoverChanges(currentHovered []InteractiveComponent) (entered, left []InteractiveComponent) {
	// Track which components were previously hovered but aren't anymore
	for _, prev := range u.hovered {
		found := false
		for _, curr := range currentHovered {
			if prev == curr {
				found = true
				break
			}
		}
		if !found {
			left = append(left, prev)
		}
	}

	// Track which components are newly hovered
	for _, curr := range currentHovered {
		found := false
		for _, prev := range u.hovered {
			if curr == prev {
				found = true
				break
			}
		}
		if !found {
			entered = append(entered, curr)
		}
	}

	return entered, left
}

func (u *InputManager) emitEvent(components []InteractiveComponent, event *Event) {
	for _, comp := range components {
		comp.HandleEvent(event)
		if event.PropagationStopped() {
			break
		}
	}
}

// findInteractiveComponentsAt returns a stack of components that are interactive
// and are at the given coordinates. The stack is ordered from topmost to bottommost.
func (m *InputManager) findInteractiveComponentsAt(root Component, x, y float64) []InteractiveComponent {
	stack := make([]InteractiveComponent, 0)
	m.buildInteractiveStack(root, x, y, &stack)
	return stack
}

// TODO: should this be event-specific?
func (m *InputManager) buildInteractiveStack(c Component, x, y float64, stack *[]InteractiveComponent) {
	// Check if component has an event boundary
	if boundary, ok := c.(EventBoundary); ok {
		if !boundary.IsWithinBounds(x, y) {
			// Stop checking this component and its children
			return
		}
	}

	// Check children first (in reverse order for proper z-index handling)
	if container, ok := c.(Container); ok {
		children := container.GetChildren()
		for i := len(children) - 1; i >= 0; i-- {
			m.buildInteractiveStack(children[i], x, y, stack)
		}
	}

	// Finally check this component itself
	if interactive, ok := c.(InteractiveComponent); ok && c.Contains(x, y) {
		*stack = append(*stack, interactive)
	}
}
