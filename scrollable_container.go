package ebui

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type ScrollableContainer struct {
	*BaseInteractive
	*LayoutContainer
	scrollOffset    Position
	isDraggingThumb bool
	dragStartY      float64
	dragStartOffset float64
	scrollBarWidth  float64
}

func NewScrollableContainer(layout Layout) *ScrollableContainer {
	sc := &ScrollableContainer{
		BaseInteractive: NewBaseInteractive(),
		LayoutContainer: NewLayoutContainer(layout),
		scrollOffset:    Position{X: 0, Y: 0},
		scrollBarWidth:  12, // Width of scroll bar in pixels
	}

	// Handle mouse wheel events
	sc.eventDispatcher.AddEventListener(EventMouseWheel, func(e Event) {
		wheelY := e.Y
		sc.scrollOffset.Y -= wheelY * 10
		sc.clampScrollOffset()
	})

	// Handle scroll bar dragging
	sc.eventDispatcher.AddEventListener(EventMouseDown, func(e Event) {
		if sc.isOverScrollBar(e.X, e.Y) {
			sc.isDraggingThumb = true
			sc.dragStartY = e.Y
			sc.dragStartOffset = sc.scrollOffset.Y
		}
	})

	sc.eventDispatcher.AddEventListener(EventMouseUp, func(e Event) {
		sc.isDraggingThumb = false
	})

	sc.eventDispatcher.AddEventListener(EventMouseMove, func(e Event) {
		if sc.isDraggingThumb {
			deltaY := e.Y - sc.dragStartY

			contentSize := sc.layout.GetMinSize(sc)
			viewportSize := sc.GetSize()
			scrollRatio := (contentSize.Height - viewportSize.Height) / (viewportSize.Height - sc.getScrollThumbHeight())

			sc.scrollOffset.Y = sc.dragStartOffset + deltaY*scrollRatio
			sc.clampScrollOffset()
		}
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

	// Draw scroll bar if needed
	if sc.needsScrollBar() {
		sc.drawScrollBar(screen)
	}
}

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

// Scroll bar related methods
func (sc *ScrollableContainer) needsScrollBar() bool {
	contentSize := sc.layout.GetMinSize(sc)
	viewportSize := sc.GetSize()
	return contentSize.Height > viewportSize.Height
}

func (sc *ScrollableContainer) getScrollThumbHeight() float64 {
	contentSize := sc.layout.GetMinSize(sc)
	viewportSize := sc.GetSize()
	ratio := viewportSize.Height / contentSize.Height
	return math.Max(viewportSize.Height*ratio, 20) // Minimum thumb size of 20px
}

func (sc *ScrollableContainer) getScrollThumbPosition() float64 {
	contentSize := sc.layout.GetMinSize(sc)
	viewportSize := sc.GetSize()
	scrollableHeight := viewportSize.Height - sc.getScrollThumbHeight()
	scrollRatio := sc.scrollOffset.Y / (contentSize.Height - viewportSize.Height)
	return scrollableHeight * scrollRatio
}

func (sc *ScrollableContainer) isOverScrollBar(x, y float64) bool {
	if !sc.needsScrollBar() {
		return false
	}

	pos := sc.GetAbsolutePosition()
	size := sc.GetSize()
	thumbY := pos.Y + sc.getScrollThumbPosition()

	return x >= pos.X+size.Width-sc.scrollBarWidth &&
		x <= pos.X+size.Width &&
		y >= thumbY &&
		y <= thumbY+sc.getScrollThumbHeight()
}

func (sc *ScrollableContainer) drawScrollBar(screen *ebiten.Image) {
	pos := sc.GetAbsolutePosition()
	size := sc.GetSize()

	// Draw track
	trackImg := ebiten.NewImage(int(sc.scrollBarWidth), int(size.Height))
	trackImg.Fill(color.RGBA{200, 200, 200, 255})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X+size.Width-sc.scrollBarWidth, pos.Y)
	screen.DrawImage(trackImg, op)

	// Draw thumb
	thumbHeight := sc.getScrollThumbHeight()
	thumbY := sc.getScrollThumbPosition()

	thumbImg := ebiten.NewImage(int(sc.scrollBarWidth), int(thumbHeight))
	if sc.isDraggingThumb {
		thumbImg.Fill(color.RGBA{120, 120, 120, 255}) // Darker when dragging
	} else {
		thumbImg.Fill(color.RGBA{160, 160, 160, 255})
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X+size.Width-sc.scrollBarWidth, pos.Y+thumbY)
	screen.DrawImage(thumbImg, op)
}

func (sc *ScrollableContainer) clampScrollOffset() {
	contentSize := sc.layout.GetMinSize(sc)
	viewportSize := sc.GetSize()
	maxScroll := math.Max(0, contentSize.Height-viewportSize.Height)
	sc.scrollOffset.Y = clamp(sc.scrollOffset.Y, 0, maxScroll)
}

// getVisibleBounds returns the visible rectangle of the container
func (sc *ScrollableContainer) getVisibleBounds() image.Rectangle {
	pos := sc.GetAbsolutePosition()
	size := sc.GetSize()
	padding := sc.GetPadding()

	// Adjust width to account for scroll bar if needed
	scrollBarAdjustment := float64(0)
	if sc.needsScrollBar() {
		scrollBarAdjustment = sc.scrollBarWidth
	}

	return image.Rectangle{
		Min: image.Point{
			X: int(pos.X + padding.Left),
			Y: int(pos.Y + padding.Top),
		},
		Max: image.Point{
			X: int(pos.X + size.Width - padding.Right - scrollBarAdjustment),
			Y: int(pos.Y + size.Height - padding.Bottom),
		},
	}
}
