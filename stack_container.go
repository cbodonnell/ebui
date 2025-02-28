package ebui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// StackContainer manages a stack of views/screens with transitions
type StackContainer struct {
	*BaseContainer
	stack           []Component
	transitioning   bool
	transitionPos   float64
	pushing         bool
	transitionSpeed float64
}

func WithTransitionSpeed(speed float64) ComponentOpt {
	return func(c Component) {
		if sc, ok := c.(*StackContainer); ok {
			sc.transitionSpeed = speed
		}
	}
}

func NewStackContainer(opts ...ComponentOpt) *StackContainer {
	sc := &StackContainer{
		BaseContainer:   NewBaseContainer(opts...),
		stack:           make([]Component, 0),
		transitionSpeed: 0.05, // Default transition speed
	}
	for _, opt := range opts {
		opt(sc)
	}
	return sc
}

// Push adds a new view to the stack with a transition
func (sc *StackContainer) Push(view Container) {
	if sc.transitioning {
		return
	}

	// For the very first view, just add it directly without animation
	if len(sc.stack) == 0 {
		view.SetSize(sc.GetSize())
		view.SetPosition(Position{
			X:        0, // Start in normal position
			Y:        0,
			Relative: true,
		})
		sc.AddChild(view)
		sc.stack = append(sc.stack, view)
		return
	}

	// For subsequent views, do the sliding animation
	sc.transitioning = true
	sc.pushing = true
	sc.transitionPos = 0

	viewSize := sc.GetSize()
	view.SetSize(viewSize)
	view.SetPosition(Position{
		X:        viewSize.Width, // Start offscreen to the right
		Y:        0,
		Relative: true,
	})

	sc.AddChild(view)
	sc.stack = append(sc.stack, view)
}

// Update handles the transition animation
func (sc *StackContainer) Update() error {
	if len(sc.stack) == 0 {
		return nil
	}

	if sc.transitioning {
		// Update transition position
		sc.transitionPos += sc.transitionSpeed
		if sc.transitionPos >= 1 {
			sc.finishTransition()
		}

		// Calculate positions for views
		width := sc.GetSize().Width
		if sc.pushing && len(sc.stack) >= 2 {
			// Move previous view left and new view in from right
			prevView := sc.stack[len(sc.stack)-2]
			currentView := sc.stack[len(sc.stack)-1]

			prevView.SetPosition(Position{
				X:        -width * sc.transitionPos,
				Y:        0,
				Relative: true,
			})

			currentView.SetPosition(Position{
				X:        width * (1 - sc.transitionPos),
				Y:        0,
				Relative: true,
			})
		} else if !sc.pushing && len(sc.stack) >= 2 {
			// Move current view right and previous view in from left
			currentView := sc.stack[len(sc.stack)-1]
			previousView := sc.stack[len(sc.stack)-2]

			currentView.SetPosition(Position{
				X:        width * sc.transitionPos,
				Y:        0,
				Relative: true,
			})
			previousView.SetPosition(Position{
				X:        -width * (1 - sc.transitionPos),
				Y:        0,
				Relative: true,
			})
		}
	}

	return sc.BaseContainer.Update()
}

// Pop removes the top view from the stack with a transition
func (sc *StackContainer) Pop() {
	if len(sc.stack) <= 1 || sc.transitioning {
		return
	}

	sc.transitioning = true
	sc.pushing = false
	sc.transitionPos = 0

	// Previous view should already be in the visual hierarchy
	// Just need to ensure it's positioned correctly
	previousView := sc.stack[len(sc.stack)-2]
	previousView.SetPosition(Position{
		X:        -sc.GetSize().Width,
		Y:        0,
		Relative: true,
	})
}

// finishTransition completes the current transition
func (sc *StackContainer) finishTransition() {
	sc.transitionPos = 1

	if !sc.pushing {
		// When popping, remove the top view after transition
		currentView := sc.stack[len(sc.stack)-1]
		sc.RemoveChild(currentView)
		sc.stack = sc.stack[:len(sc.stack)-1]

		// Ensure the new top view is positioned correctly
		if len(sc.stack) > 0 {
			topView := sc.stack[len(sc.stack)-1]
			topView.SetPosition(Position{
				X:        0,
				Y:        0,
				Relative: true,
			})
		}
	}

	sc.transitioning = false
}

// Draw overrides the BaseContainer.Draw method to implement clipping
func (sc *StackContainer) Draw(screen *ebiten.Image) {
	// Draw the container's background and debug info
	sc.BaseComponent.Draw(screen)

	// Create a sub-image for clipping to the container's bounds
	bounds := sc.getVisibleBounds()
	subScreen := screen.SubImage(bounds).(*ebiten.Image)

	// Draw all children to the clipped sub-image
	for _, child := range sc.children {
		child.Draw(subScreen)
	}
}

// getVisibleBounds returns the visible rectangle of the container
func (sc *StackContainer) getVisibleBounds() image.Rectangle {
	pos := sc.GetAbsolutePosition()
	size := sc.GetSize()
	padding := sc.GetPadding()

	return image.Rectangle{
		Min: image.Point{
			X: int(pos.X + padding.Left),
			Y: int(pos.Y + padding.Top),
		},
		Max: image.Point{
			X: int(pos.X + size.Width - padding.Right),
			Y: int(pos.Y + size.Height - padding.Bottom),
		},
	}
}

// GetActiveView returns the currently active view
func (sc *StackContainer) GetActiveView() Component {
	if len(sc.stack) == 0 {
		return nil
	}
	return sc.stack[len(sc.stack)-1]
}

// Clear removes all views from the stack except the root view
func (sc *StackContainer) Clear() {
	if len(sc.stack) <= 1 {
		return
	}

	rootView := sc.stack[0]
	sc.stack = []Component{rootView}

	// Remove all children except root
	sc.children = []Component{rootView}

	// Ensure root view is positioned correctly
	rootView.SetPosition(Position{
		X:        0,
		Y:        0,
		Relative: true,
	})
}
