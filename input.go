package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type InputManager struct {
	lastMousePressed bool
	hovered          Interactive
	activeScroller   *ScrollableContainer
	lastX, lastY     float64
}

func NewInputManager() *InputManager {
	return &InputManager{}
}

// Update handles all input events for the frame
func (u *InputManager) Update(root Component) {
	x, y := ebiten.CursorPosition()
	fx, fy := float64(x), float64(y)
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	u.handleActiveScroller(fx, fy, mousePressed)
	u.handleHoverEvents(fx, fy, root)
	u.handleMouseEvents(fx, fy, mousePressed)
	u.handleWheelEvents(fx, fy, root)

	u.lastMousePressed = mousePressed
	u.lastX, u.lastY = fx, fy
}

// handleActiveScroller manages the currently active scrolling component
func (u *InputManager) handleActiveScroller(fx, fy float64, mousePressed bool) {
	if u.lastMousePressed && !mousePressed && u.activeScroller != nil {
		u.activeScroller.HandleEvent(newEvent(EventMouseUp, fx, fy))
		u.activeScroller = nil
	}

	if u.activeScroller != nil {
		u.activeScroller.HandleEvent(newEvent(EventMouseMove, fx, fy))
	}
}

// handleHoverEvents manages hover state changes
func (u *InputManager) handleHoverEvents(fx, fy float64, root Component) {
	target := findInteractableAt(fx, fy, root)
	if target == u.hovered {
		return
	}

	if u.hovered != nil {
		u.hovered.HandleEvent(newEvent(EventMouseLeave, fx, fy))
	}
	if target != nil {
		target.HandleEvent(newEvent(EventMouseEnter, fx, fy))
	}
	u.hovered = target
}

// handleMouseEvents manages mouse press and release events
func (u *InputManager) handleMouseEvents(fx, fy float64, mousePressed bool) {
	if u.hovered == nil {
		return
	}

	if mousePressed && !u.lastMousePressed {
		u.hovered.HandleEvent(newEvent(EventMouseDown, fx, fy))
		if scrollable, ok := u.hovered.(*ScrollableContainer); ok {
			if scrollable.isOverScrollBar(fx, fy) {
				u.activeScroller = scrollable
			}
		}
	}

	if !mousePressed && u.lastMousePressed {
		u.hovered.HandleEvent(newEvent(EventMouseUp, fx, fy))
	}
}

// handleWheelEvents manages mouse wheel scrolling
func (u *InputManager) handleWheelEvents(fx, fy float64, root Component) {
	wheelX, wheelY := ebiten.Wheel()
	if scrollable := findScrollableAt(fx, fy, root); scrollable != nil {
		scrollable.HandleEvent(Event{Type: EventMouseWheel, X: wheelX, Y: wheelY})
	}
}

// newEvent creates a new event with the given type and coordinates
func newEvent(eventType EventType, x, y float64) Event {
	return Event{Type: eventType, X: x, Y: y}
}

// findComponentAt is a generic component finder
func findComponentAt[T any](x, y float64, c Component, check func(Component) (T, bool)) (T, bool) {
	var zero T
	if !c.Contains(x, y) {
		return zero, false
	}

	if container, ok := c.(Container); ok {
		children := container.GetChildren()
		for i := len(children) - 1; i >= 0; i-- {
			if found, ok := findComponentAt(x, y, children[i], check); ok {
				return found, true
			}
		}
	}

	return check(c)
}

// findInteractableAt finds the topmost interactive component at the given coordinates
func findInteractableAt(x, y float64, c Component) Interactive {
	found, _ := findComponentAt(x, y, c, func(c Component) (Interactive, bool) {
		if i, ok := c.(Interactive); ok && c.Contains(x, y) {
			return i, true
		}
		return nil, false
	})
	return found
}

// findScrollableAt finds the topmost scrollable container at the given coordinates
func findScrollableAt(x, y float64, c Component) *ScrollableContainer {
	found, _ := findComponentAt(x, y, c, func(c Component) (*ScrollableContainer, bool) {
		if s, ok := c.(*ScrollableContainer); ok && c.Contains(x, y) {
			return s, true
		}
		return nil, false
	})
	return found
}
