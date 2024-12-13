package ebui

import (
	"math"
)

type ScrollableContainer struct {
	*BaseInteractive
	*LayoutContainer
	scrollOffset Position
}

func NewScrollableContainer(layout Layout) *ScrollableContainer {
	sc := &ScrollableContainer{
		BaseInteractive: NewBaseInteractive(),
		LayoutContainer: NewLayoutContainer(layout),
		scrollOffset:    Position{X: 0, Y: 0},
	}

	// handle mouse wheel events
	sc.eventDispatcher.AddEventListener(EventMouseWheel, func(e Event) {
		wheelY := e.Y
		sc.scrollOffset.Y -= wheelY * 10

		// Clamp scroll position
		contentSize := sc.layout.GetMinSize(sc)
		viewportSize := sc.GetSize()
		maxScroll := math.Max(0, contentSize.Height-viewportSize.Height)
		sc.scrollOffset.Y = clamp(sc.scrollOffset.Y, 0, maxScroll)
	})

	return sc
}

func (sc *ScrollableContainer) Update() error {
	if sc.layout != nil {
		sc.layout.ArrangeChildren(sc)
	}

	// Update children positions
	for _, child := range sc.children {
		child.SetPosition(Position{
			X: child.GetPosition().X,
			Y: child.GetPosition().Y - sc.scrollOffset.Y,
		})
	}

	return sc.BaseContainer.Update()
}
