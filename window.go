package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type WindowState int

const (
	WindowStateHidden WindowState = iota
	WindowStateNormal
)

// WindowColors represents the color scheme for a window
type WindowColors struct {
	Background color.Color
	Header     color.Color
	HeaderText color.Color
	Border     color.Color
	// Close button colors (optional, will be derived from header color if nil)
	CloseButton  color.Color
	CloseHovered color.Color
	ClosePressed color.Color
	CloseCross   color.Color
}

// DefaultWindowColors returns a default color scheme for windows
func DefaultWindowColors() WindowColors {
	return WindowColors{
		Background: color.RGBA{240, 240, 240, 255},
		Header:     color.RGBA{200, 200, 200, 255},
		HeaderText: color.Black,
		Border:     color.RGBA{0, 0, 0, 255},
		// Close button colors will be derived from header color
	}
}

type Window struct {
	*BaseFocusable
	*LayoutContainer
	manager         *WindowManager
	header          *BaseContainer
	content         *LayoutContainer
	title           string
	titleLabel      *Label
	closeButton     *Button
	state           WindowState
	isDragging      bool
	dragStartX      float64
	dragStartY      float64
	windowStartX    float64
	windowStartY    float64
	headerHeight    float64
	closeCallback   func()
	colors          WindowColors
	borderWidth     float64
	isStatic        bool
	closeButtonSize Size
}

type WindowOpt func(w *Window)

// WithWindowTitle sets the window title
func WithWindowTitle(title string) WindowOpt {
	return func(w *Window) {
		w.title = title
	}
}

// WithCloseCallback sets the callback for when the window is closed
func WithCloseCallback(callback func()) WindowOpt {
	return func(w *Window) {
		w.closeCallback = callback
	}
}

// WithWindowColors sets custom colors for the window
func WithWindowColors(colors WindowColors) WindowOpt {
	return func(w *Window) {
		w.colors = colors
	}
}

// WithBorderWidth sets the width of the window border
func WithBorderWidth(width float64) WindowOpt {
	return func(w *Window) {
		w.borderWidth = width
	}
}

// WithHeaderHeight sets the height of the window header
func WithHeaderHeight(height float64) WindowOpt {
	return func(w *Window) {
		w.headerHeight = height
	}
}

// WithPosition sets the window position
func WithWindowPosition(x, y float64) WindowOpt {
	return func(w *Window) {
		w.SetPosition(Position{X: x, Y: y})
	}
}

// WithStatic makes the window non-draggable
func WithStatic() WindowOpt {
	return func(w *Window) {
		w.isStatic = true
	}
}

// WithCloseButtonSize sets the size of the close button
func WithCloseButtonSize(size Size) WindowOpt {
	return func(w *Window) {
		w.closeButtonSize = size
	}
}

// Show makes the window visible
func (w *Window) Show() {
	w.state = WindowStateNormal
	w.manager.SetActiveWindow(w)
	w.Enable()
}

// Hide makes the window invisible
func (w *Window) Hide() {
	w.state = WindowStateHidden
	w.Disable()
	if w.closeCallback != nil {
		w.closeCallback()
	}
}

// Toggle shows the window if hidden, hides it if visible
func (w *Window) Toggle() {
	if w.IsVisible() {
		w.Hide()
	} else {
		w.Show()
	}
}

func (w *Window) AddChild(child Component) {
	w.content.AddChild(child)
}

func (w *Window) RemoveChild(child Component) {
	w.content.RemoveChild(child)
}

// IsVisible returns whether the window is currently visible
func (w *Window) IsVisible() bool {
	return w.state == WindowStateNormal
}

func (w *Window) SetSize(size Size) {
	w.LayoutContainer.SetSize(size)
	// Update header and content sizes
	w.header.SetSize(Size{Width: size.Width, Height: w.headerHeight})
	w.content.SetSize(Size{Width: size.Width, Height: size.Height - w.headerHeight})
}

func (w *Window) SetTitle(title string) {
	w.title = title
	w.titleLabel.SetText(title)
}

func (w *Window) Draw(screen *ebiten.Image) {
	if !w.IsVisible() {
		return
	}

	// Draw the window border 1px
	pos := w.GetAbsolutePosition()
	size := w.GetSize()
	bg := GetCache().ImageWithColor(int(size.Width+2), int(size.Height+2), w.colors.Border)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X-1, pos.Y-1)
	screen.DrawImage(bg, op)

	w.LayoutContainer.Draw(screen)
}

func (w *Window) clampToScreen() {
	// Get the window manager bounds
	bounds := w.manager.GetSize()

	pos := w.GetPosition()
	size := w.GetSize()

	// Keep the window title bar within bounds
	minX := -size.Width + size.Width/2     // Keep half the window within bounds
	minY := float64(0)                     // Prevent dragging above bounds
	maxX := bounds.Width - size.Width/2    // Keep half the window within bounds
	maxY := bounds.Height - w.headerHeight // Keep header within bounds

	pos.X = clamp(pos.X, minX, maxX)
	pos.Y = clamp(pos.Y, minY, maxY)

	w.SetPosition(pos)
}

func (w *Window) registerEventListeners() {
	w.AddEventListener(DragStart, func(e *Event) {
		// Always activate window on any mouse down within the window
		w.manager.SetActiveWindow(w)

		// Don't drag if window is static
		if w.isStatic {
			return
		}

		// Start dragging only if over header
		if w.isOverHeader(e.MouseX, e.MouseY) {
			w.isDragging = true
			w.dragStartX = e.MouseX
			w.dragStartY = e.MouseY
			absPos := w.GetAbsolutePosition()
			w.windowStartX = absPos.X
			w.windowStartY = absPos.Y
		}
	})

	w.AddEventListener(DragEnd, func(e *Event) {
		if w.isStatic {
			return
		}

		w.isDragging = false
	})

	w.AddEventListener(Drag, func(e *Event) {
		if w.isStatic {
			return
		}

		if w.isDragging {
			deltaX := e.MouseX - w.dragStartX
			deltaY := e.MouseY - w.dragStartY
			newPos := w.GetPosition()
			newPos.X = w.windowStartX + deltaX
			newPos.Y = w.windowStartY + deltaY
			w.SetPosition(newPos)
			w.clampToScreen() // Clamp after setting new position
		}
	})
}

func (w *Window) isOverHeader(x, y float64) bool {
	absPos := w.GetAbsolutePosition()

	// Get position of close button
	cbPos := w.closeButton.GetAbsolutePosition()
	cbSize := w.closeButton.GetSize()

	// Skip dragging if over close button
	if x >= cbPos.X && x <= cbPos.X+cbSize.Width &&
		y >= cbPos.Y && y <= cbPos.Y+cbSize.Height {
		return false
	}

	return x >= absPos.X &&
		x <= absPos.X+w.GetSize().Width &&
		y >= absPos.Y &&
		y <= absPos.Y+w.headerHeight
}

type WindowManager struct {
	*ZIndexedContainer
	activeWindow *Window
	nextZIndex   int
}

func NewWindowManager(opts ...ComponentOpt) *WindowManager {
	wm := &WindowManager{
		ZIndexedContainer: NewZIndexedContainer(opts...),
		nextZIndex:        1,
	}
	return wm
}

func (wm *WindowManager) CreateWindow(width, height float64, opts ...WindowOpt) *Window {
	window := &Window{
		BaseFocusable: NewBaseFocusable(),
		LayoutContainer: NewLayoutContainer(
			WithSize(width, height),
			WithLayout(NewVerticalStackLayout(0, AlignStart)),
		),
		manager:         wm,
		headerHeight:    30,
		colors:          DefaultWindowColors(),
		borderWidth:     1,
		state:           WindowStateNormal,
		closeButtonSize: Size{Width: 20, Height: 20},
	}

	for _, opt := range opts {
		opt(window)
	}

	// Initialize header with layout that allows absolute positioning of children
	window.header = NewBaseContainer(
		WithSize(width, window.headerHeight),
		WithBackground(window.colors.Header),
	)

	// Create title label centered in header
	window.titleLabel = NewLabel(
		window.title,
		WithSize(width, window.headerHeight),
		WithPosition(Position{X: 0, Y: 0, Relative: true}),
		WithColor(window.colors.HeaderText),
		WithJustify(JustifyCenter),
	)
	window.header.AddChild(window.titleLabel)

	// Derive close button colors from window colors if not explicitly set
	closeButtonColor := window.colors.CloseButton
	closeHoveredColor := window.colors.CloseHovered
	closePressedColor := window.colors.ClosePressed
	closeCrossColor := window.colors.CloseCross

	// Set default colors based on header if not specified
	if closeButtonColor == nil {
		// close button color should match the header color
		closeButtonColor = window.colors.Header
	}

	if closeHoveredColor == nil {
		// Use a slightly brighter version of the close button color
		r, g, b, a := closeButtonColor.RGBA()
		r = uint32(float64(r>>8) * 1.4) // 140% brightness
		g = uint32(float64(g>>8) * 1.4)
		b = uint32(float64(b>>8) * 1.4)
		// Clamp to 255
		r = min(r, 255)
		g = min(g, 255)
		b = min(b, 255)
		closeHoveredColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a >> 8)}
	}

	if closePressedColor == nil {
		// Use a darker version of the close button color
		r, g, b, a := closeButtonColor.RGBA()
		r = uint32(float64(r>>8) * 0.6) // 60% brightness
		g = uint32(float64(g>>8) * 0.6)
		b = uint32(float64(b>>8) * 0.6)
		closePressedColor = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a >> 8)}
	}

	if closeCrossColor == nil {
		// cross color should match the header text color
		closeCrossColor = window.colors.HeaderText
	}

	// Create close button in top left corner
	window.closeButton = NewButton(
		WithSize(window.closeButtonSize.Width, window.closeButtonSize.Height),
		WithPosition(Position{
			X:        5,                                                         // 5px padding from left edge
			Y:        (window.headerHeight - window.closeButtonSize.Height) / 2, // Centered vertically
			Relative: true,
		}),
		WithLabelText("X"), // Unicode multiplication sign as "X"
		WithButtonColors(ButtonColors{
			Default:     closeButtonColor,
			Hovered:     closeHoveredColor,
			Pressed:     closePressedColor,
			FocusBorder: color.Transparent, // No focus border for close button
		}),
		WithLabelColor(closeCrossColor),
		WithClickHandler(func() {
			window.Hide()
		}),
	)
	window.header.AddChild(window.closeButton)

	// Create content container
	window.content = NewLayoutContainer(
		WithSize(width, height-window.headerHeight),
		WithBackground(window.colors.Background),
		WithLayout(NewVerticalStackLayout(0, AlignStart)),
	)

	window.LayoutContainer.AddChild(window.header)
	window.LayoutContainer.AddChild(window.content)

	window.registerEventListeners()

	window.SetPosition(Position{
		X:        window.position.X,
		Y:        window.position.Y,
		ZIndex:   wm.nextZIndex,
		Relative: false,
	})
	wm.nextZIndex++

	wm.AddChild(window)
	wm.SetActiveWindow(window)

	return window
}

func (wm *WindowManager) SetActiveWindow(window *Window) {
	if wm.activeWindow == window {
		return
	}

	wm.activeWindow = window
	maxZ := 0
	for _, child := range wm.GetChildren() {
		if z := child.GetPosition().ZIndex; z > maxZ {
			maxZ = z
		}
	}
	pos := window.GetPosition()
	pos.ZIndex = maxZ + 1
	window.SetPosition(pos)
	wm.nextZIndex = maxZ + 2
}
