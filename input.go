package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type InputManager struct {
	lastMousePressed bool
	hovered          InteractiveComponent
	activeScroller   *ScrollableContainer
	activeWindow     *Window
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

	// Handle mouse move events for active components
	if mousePressed {
		u.handleMouseMove(fx, fy)
	}

	// Handle active component release
	if u.lastMousePressed && !mousePressed {
		u.handleMouseRelease(fx, fy)
	}

	// Only handle hover and new mouse press events if we're not dragging anything
	if u.activeWindow == nil && u.activeScroller == nil {
		u.handleHoverEvents(fx, fy, root)
		if mousePressed && !u.lastMousePressed {
			u.handleMousePress(fx, fy)
		}
	}

	u.handleWheelEvents(fx, fy, root)

	u.lastMousePressed = mousePressed
	u.lastX, u.lastY = fx, fy
}

func (u *InputManager) handleMouseMove(fx, fy float64) {
	if u.activeWindow != nil {
		u.activeWindow.HandleEvent(Event{Type: EventMouseMove, X: fx, Y: fy, Component: u.activeWindow})
	} else if u.activeScroller != nil {
		u.activeScroller.HandleEvent(Event{Type: EventMouseMove, X: fx, Y: fy, Component: u.activeScroller})
	}
}

func (u *InputManager) handleMouseRelease(fx, fy float64) {
	if u.activeWindow != nil {
		u.activeWindow.HandleEvent(Event{Type: EventMouseUp, X: fx, Y: fy, Component: u.activeWindow})
		u.activeWindow = nil
	} else if u.activeScroller != nil {
		u.activeScroller.HandleEvent(Event{Type: EventMouseUp, X: fx, Y: fy, Component: u.activeScroller})
		u.activeScroller = nil
	} else if u.hovered != nil {
		u.hovered.HandleEvent(Event{Type: EventMouseUp, X: fx, Y: fy, Component: u.hovered})
	}
}

func (u *InputManager) handleMousePress(fx, fy float64) {
	if u.hovered == nil {
		return
	}

	u.hovered.HandleEvent(Event{Type: EventMouseDown, X: fx, Y: fy, Component: u.hovered})

	// Check if we're starting a window drag
	if window, ok := u.hovered.(*Window); ok && window.isOverHeader(fx, fy) {
		u.activeWindow = window
		return
	}

	// Check if we're starting a scroll
	if scrollable, ok := u.hovered.(*ScrollableContainer); ok && scrollable.isOverScrollBar(fx, fy) {
		u.activeScroller = scrollable
		return
	}
}

func (u *InputManager) handleHoverEvents(fx, fy float64, root Component) {
	target := findInteractableAt(fx, fy, root)
	if target == u.hovered {
		return
	}

	if u.hovered != nil {
		u.hovered.HandleEvent(Event{Type: EventMouseLeave, X: fx, Y: fy, Component: u.hovered})
	}
	if target != nil {
		target.HandleEvent(Event{Type: EventMouseEnter, X: fx, Y: fy, Component: target})
	}
	u.hovered = target
}

func (u *InputManager) handleWheelEvents(fx, fy float64, root Component) {
	wheelX, wheelY := ebiten.Wheel()
	if scrollable := findScrollableAt(fx, fy, root); scrollable != nil {
		scrollable.HandleEvent(Event{Type: EventMouseWheel, X: wheelX, Y: wheelY, Component: scrollable})
	}
}

// findComponentAt is a generic component finder
func findComponentAt[T any](x, y float64, c Component, check func(Component) (T, bool)) (T, bool) {
	var zero T

	// Check if component controls its own event boundary
	if boundary, ok := c.(EventBoundary); ok {
		if !boundary.ShouldPropagateEvent(Event{}, x, y) {
			return zero, false
		}
	}

	// Check children first (in reverse order for proper z-index handling)
	if container, ok := c.(Container); ok {
		children := container.GetChildren()
		for i := len(children) - 1; i >= 0; i-- {
			if found, ok := findComponentAt(x, y, children[i], check); ok {
				return found, true
			}
		}
	}

	// Finally check this component itself
	if !c.Contains(x, y) {
		return zero, false
	}

	return check(c)
}

// findInteractableAt finds the topmost interactive component at the given coordinates
func findInteractableAt(x, y float64, c Component) InteractiveComponent {
	found, _ := findComponentAt(x, y, c, func(c Component) (InteractiveComponent, bool) {
		if i, ok := c.(InteractiveComponent); ok {
			return i, true
		}
		return nil, false
	})
	return found
}

// findScrollableAt finds the topmost scrollable container at the given coordinates
func findScrollableAt(x, y float64, c Component) *ScrollableContainer {
	found, _ := findComponentAt(x, y, c, func(c Component) (*ScrollableContainer, bool) {
		if s, ok := c.(*ScrollableContainer); ok {
			return s, true
		}
		return nil, false
	})
	return found
}
