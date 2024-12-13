package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type InputManager struct {
	lastMousePressed bool
	hovered          Interactive
}

func NewInputManager() *InputManager {
	return &InputManager{}
}

func (u *InputManager) Update(root Component) {
	// Handle hover and click events
	x, y := ebiten.CursorPosition()
	fx, fy := float64(x), float64(y)
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	target := findInteractableAt(fx, fy, root)
	if target != u.hovered {
		if u.hovered != nil {
			u.hovered.HandleEvent(Event{
				Type: EventMouseLeave,
				X:    fx,
				Y:    fy,
			})
		}
		if target != nil {
			target.HandleEvent(Event{
				Type: EventMouseEnter,
				X:    fx,
				Y:    fy,
			})
		}
		u.hovered = target
	}

	if u.hovered != nil {
		if mousePressed && !u.lastMousePressed {
			u.hovered.HandleEvent(Event{
				Type: EventMouseDown,
				X:    fx,
				Y:    fy,
			})
		}
		if !mousePressed && u.lastMousePressed {
			u.hovered.HandleEvent(Event{
				Type: EventMouseUp,
				X:    fx,
				Y:    fy,
			})
		}
	}

	u.lastMousePressed = mousePressed

	// Handle mouse wheel events
	wheelX, wheelY := ebiten.Wheel()
	scrollable := findScrollableAt(fx, fy, root)
	if scrollable != nil {
		scrollable.HandleEvent(Event{
			Type: EventMouseWheel,
			X:    wheelX,
			Y:    wheelY,
		})
	}
}

func findInteractableAt(x, y float64, c Component) Interactive {
	contains := c.Contains(x, y)
	if !contains {
		return nil
	}

	if container, ok := c.(Container); ok {
		// First check children if this component is a container
		children := container.GetChildren()
		// Iterate backwards to check the top-most first
		for i := len(children) - 1; i >= 0; i-- {
			child := children[i]
			if interactive := findInteractableAt(x, y, child); interactive != nil {
				return interactive
			}
		}
	}
	// Check if this component is interactive
	interactive, ok := c.(Interactive)
	if !ok {
		return nil
	}
	// If it is, and the point is within the bounds, return it
	if contains {
		return interactive
	}
	return nil
}

func findScrollableAt(x, y float64, c Component) *ScrollableContainer {
	contains := c.Contains(x, y)
	if !contains {
		return nil
	}

	if container, ok := c.(Container); ok {
		// First check children if this component is a container
		children := container.GetChildren()
		// Iterate backwards to check the top-most first
		for i := len(children) - 1; i >= 0; i-- {
			child := children[i]
			if interactive := findScrollableAt(x, y, child); interactive != nil {
				return interactive
			}
		}
	}
	// Check if this component is scrollable
	scrollable, ok := c.(*ScrollableContainer)
	if !ok {
		return nil
	}
	// If it is, and the point is within the bounds, return it
	if contains {
		return scrollable
	}
	return nil
}
