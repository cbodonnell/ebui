package ebui

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
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

	// Handle mouse wheel events
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

// Draw overrides BaseContainer's Draw to implement clipping
func (sc *ScrollableContainer) Draw(screen *ebiten.Image) {
	// Draw the container's background and debug info
	sc.BaseComponent.Draw(screen)

	// Create a sub-image for clipping
	bounds := sc.getVisibleBounds()
	subScreen := screen.SubImage(bounds).(*ebiten.Image)

	// Draw all children to the clipped sub-image
	for _, child := range sc.children {
		child.Draw(subScreen)
	}
}

// Contains overrides BaseContainer's Contains to implement proper event handling
func (sc *ScrollableContainer) Contains(x, y float64) bool {
	// First check if the point is within the container's bounds
	if !sc.BaseContainer.Contains(x, y) {
		return false
	}

	// Get the container's visible bounds
	bounds := sc.getVisibleBounds()

	// Check if the point is within the visible area
	return y >= float64(bounds.Min.Y) && y <= float64(bounds.Max.Y)
}

// getVisibleBounds returns the visible rectangle of the container
func (sc *ScrollableContainer) getVisibleBounds() image.Rectangle {
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
