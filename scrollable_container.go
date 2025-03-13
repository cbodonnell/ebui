package ebui

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scrollable interface {
	Container
	Interactive
	EventBoundary
	GetScrollOffset() Position
	SetScrollOffset(offset Position)
	HideScrollBar()
	ShowScrollBar()
	IsScrollBarHidden() bool
	ScrollToTop()
	ScrollToBottom()
}

var _ Scrollable = &ScrollableContainer{}

type ScrollableColors struct {
	Track     color.Color
	Thumb     color.Color
	ThumbDrag color.Color
}

func DefaultScrollableColors() ScrollableColors {
	return ScrollableColors{
		Track:     color.RGBA{200, 200, 200, 255},
		Thumb:     color.RGBA{160, 160, 160, 255},
		ThumbDrag: color.RGBA{120, 120, 120, 255},
	}
}

type ScrollableContainer struct {
	*BaseFocusable
	*LayoutContainer
	colors            ScrollableColors
	scrollOffset      Position
	isDraggingThumb   bool
	dragStartY        float64
	dragStartOffset   float64
	scrollBarWidth    float64
	isFocused         bool
	isScrollBarHidden bool
}

func WithScrollableColors(colors ScrollableColors) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*ScrollableContainer); ok {
			b.colors = colors
		}
	}
}

func NewScrollableContainer(opts ...ComponentOpt) *ScrollableContainer {
	sc := &ScrollableContainer{
		BaseFocusable:   NewBaseFocusable(),
		LayoutContainer: NewLayoutContainer(opts...),
		colors:          DefaultScrollableColors(),
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
		scrollOffset := sc.GetScrollOffset()
		scrollOffset.Y -= wheelY * 10
		sc.SetScrollOffset(scrollOffset)
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

			scrollOffset := sc.GetScrollOffset()
			scrollOffset.Y = sc.dragStartOffset + deltaY*scrollRatio
			sc.SetScrollOffset(scrollOffset)
		}
	})

	sc.AddEventListener(Focus, func(e *Event) {
		sc.isFocused = true
	})

	sc.AddEventListener(Blur, func(e *Event) {
		sc.isFocused = false
	})
}

func (sc *ScrollableContainer) GetScrollOffset() Position {
	return sc.scrollOffset
}

func (sc *ScrollableContainer) SetScrollOffset(offset Position) {
	sc.scrollOffset = offset
	sc.clampScrollOffset()
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
	sc.handleInput()

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

func (sc *ScrollableContainer) handleInput() {
	if !sc.isFocused {
		return
	}

	scrollAmount := float64(20)
	shiftPressed := ebiten.IsKeyPressed(ebiten.KeyShift)
	if shiftPressed {
		scrollAmount = 100
	}

	switch {
	case ebiten.IsKeyPressed(ebiten.KeyArrowUp):
		scrollOffset := sc.GetScrollOffset()
		scrollOffset.Y -= scrollAmount
		sc.SetScrollOffset(scrollOffset)
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown):
		scrollOffset := sc.GetScrollOffset()
		scrollOffset.Y += scrollAmount
		sc.SetScrollOffset(scrollOffset)
	case ebiten.IsKeyPressed(ebiten.KeyPageUp):
		scrollOffset := sc.GetScrollOffset()
		scrollOffset.Y -= sc.GetSize().Height
		sc.SetScrollOffset(scrollOffset)
	case ebiten.IsKeyPressed(ebiten.KeyPageDown):
		scrollOffset := sc.GetScrollOffset()
		scrollOffset.Y += sc.GetSize().Height
		sc.SetScrollOffset(scrollOffset)
	case ebiten.IsKeyPressed(ebiten.KeyHome):
		scrollOffset := sc.GetScrollOffset()
		scrollOffset.Y = 0
		sc.SetScrollOffset(scrollOffset)
	case ebiten.IsKeyPressed(ebiten.KeyEnd):
		scrollOffset := sc.GetScrollOffset()
		contentSize := sc.layout.GetMinSize(sc)
		viewportSize := sc.GetSize()
		scrollOffset.Y = contentSize.Height - viewportSize.Height
		sc.SetScrollOffset(scrollOffset)
	}
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
	if sc.isScrollBarHidden {
		return false
	}
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
	trackWidth := int(sc.scrollBarWidth)
	trackHeight := int(size.Height)

	trackImg := GetCache().ImageWithColor(trackWidth, trackHeight, sc.colors.Track)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X+size.Width-sc.scrollBarWidth, pos.Y)
	screen.DrawImage(trackImg, op)

	// Draw thumb
	thumbHeight := int(sc.getScrollThumbHeight())

	// Calculate max position where thumb would touch bottom exactly
	maxThumbY := trackHeight - thumbHeight

	// Calculate normal position using ratio
	contentSize := sc.layout.GetMinSize(sc)
	viewportSize := sc.GetSize()
	maxScroll := math.Max(0, contentSize.Height-viewportSize.Height)

	// Calculate thumb position, ensuring it doesn't exceed the track
	var thumbY int
	if maxScroll <= 0 {
		thumbY = 0
	} else {
		ratio := float64(sc.scrollOffset.Y) / maxScroll
		// Use integer math to avoid fractional positioning
		thumbY = int(math.Min(float64(maxThumbY), float64(maxThumbY)*ratio))
	}

	thumbColor := sc.colors.Thumb
	if sc.isDraggingThumb {
		thumbColor = sc.colors.ThumbDrag
	}

	thumbImg := GetCache().ImageWithColor(trackWidth, thumbHeight, thumbColor)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X+size.Width-sc.scrollBarWidth, pos.Y+float64(thumbY))
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

func (sc *ScrollableContainer) ScrollToTop() {
	scrollOffset := sc.GetScrollOffset()
	scrollOffset.Y = 0
	sc.SetScrollOffset(scrollOffset)
}

func (sc *ScrollableContainer) ScrollToBottom() {
	contentSize := sc.layout.GetMinSize(sc)
	viewportSize := sc.GetSize()
	maxScroll := math.Max(0, contentSize.Height-viewportSize.Height)

	scrollOffset := sc.GetScrollOffset()
	scrollOffset.Y = maxScroll
	sc.SetScrollOffset(scrollOffset)
}

func (sc *ScrollableContainer) HideScrollBar() {
	sc.isScrollBarHidden = true
}

func (sc *ScrollableContainer) ShowScrollBar() {
	sc.isScrollBarHidden = false
}

func (sc *ScrollableContainer) IsScrollBarHidden() bool {
	return sc.isScrollBarHidden
}
