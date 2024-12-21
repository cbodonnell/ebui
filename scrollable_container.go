package ebui

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ EventBoundary = &ScrollableContainer{}

type ScrollableContainer struct {
	*BaseInteractive
	*LayoutContainer
	scrollOffset    Position
	isDraggingThumb bool
	dragStartY      float64
	dragStartOffset float64
	scrollBarWidth  float64
}

func NewScrollableContainer(opts ...ComponentOpt) *ScrollableContainer {
	sc := &ScrollableContainer{
		BaseInteractive: NewBaseInteractive(),
		LayoutContainer: NewLayoutContainer(opts...),
		scrollBarWidth:  12, // Width of scroll bar in pixels
	}

	for _, opt := range opts {
		opt(sc)
	}

	sc.registerEventListeners()

	return sc
}

func (sc *ScrollableContainer) registerEventListeners() {
	// Handle mouse wheel events
	sc.AddEventListener(Wheel, func(e *Event) {
		wheelY := e.WheelDeltaY
		sc.scrollOffset.Y -= wheelY * 10
		sc.clampScrollOffset()
	})

	// Handle scroll bar dragging
	sc.AddEventListener(DragStart, func(e *Event) {
		if sc.isOverScrollBar(e.MouseX, e.MouseY) {
			sc.isDraggingThumb = true
			sc.dragStartY = e.MouseY
			sc.dragStartOffset = sc.scrollOffset.Y
		}
	})

	sc.AddEventListener(DragEnd, func(e *Event) {
		sc.isDraggingThumb = false
	})

	sc.AddEventListener(Drag, func(e *Event) {
		if sc.isDraggingThumb {
			deltaY := e.MouseY - sc.dragStartY
			contentSize := sc.layout.GetMinSize(sc)
			viewportSize := sc.GetSize()
			scrollRatio := (contentSize.Height - viewportSize.Height) / (viewportSize.Height - sc.getScrollThumbHeight())

			sc.scrollOffset.Y = sc.dragStartOffset + deltaY*scrollRatio
			sc.clampScrollOffset()
		}
	})
}

func (sc *ScrollableContainer) AddChild(child Component) {
	sc.LayoutContainer.AddChild(child)
	sc.clampScrollOffset()
}

func (sc *ScrollableContainer) RemoveChild(child Component) {
	sc.LayoutContainer.RemoveChild(child)
	sc.clampScrollOffset()
}

func (sc *ScrollableContainer) Update() error {
	if sc.layout != nil {
		sc.layout.ArrangeChildren(sc)
	}

	// Update children positions
	for _, child := range sc.children {
		pos := child.GetPosition()
		pos.Y -= sc.scrollOffset.Y
		child.SetPosition(pos)
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

func (sc *ScrollableContainer) IsWithinBounds(x, y float64) bool {
	return sc.Contains(x, y)
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
