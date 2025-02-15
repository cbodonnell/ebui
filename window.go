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
	Border     color.Color
}

// DefaultWindowColors returns a default color scheme for windows
func DefaultWindowColors() WindowColors {
	return WindowColors{
		Background: color.RGBA{240, 240, 240, 255},
		Header:     color.RGBA{200, 200, 200, 255},
		Border:     color.RGBA{0, 0, 0, 255},
	}
}

type Window struct {
	*BaseFocusable
	*LayoutContainer
	manager       *WindowManager
	header        *LayoutContainer
	content       *LayoutContainer
	title         string
	state         WindowState
	isDragging    bool
	dragStartX    float64
	dragStartY    float64
	windowStartX  float64
	windowStartY  float64
	headerHeight  float64
	closeCallback func()
	colors        WindowColors
	borderWidth   float64
	isStatic      bool
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
	// Get the game window bounds
	width, height := ebiten.WindowSize()
	screenWidth, screenHeight := float64(width), float64(height)

	pos := w.GetPosition()
	size := w.GetSize()

	// Keep the window title bar on screen
	minX := -size.Width + size.Width/2    // Keep half the window on screen
	minY := float64(0)                    // Prevent dragging above screen
	maxX := screenWidth - size.Width/2    // Keep half the window on screen
	maxY := screenHeight - w.headerHeight // Keep header on screen

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
		manager:      wm,
		headerHeight: 30,
		colors:       DefaultWindowColors(),
		borderWidth:  1,
		state:        WindowStateNormal,
	}

	for _, opt := range opts {
		opt(window)
	}

	window.header = NewLayoutContainer(
		WithSize(width, window.headerHeight),
		WithBackground(window.colors.Header),
		WithLayout(NewVerticalStackLayout(0, AlignStart)),
	)

	titleLabel := NewLabel(
		window.title,
		WithSize(width, window.headerHeight),
		WithJustify(JustifyCenter),
	)
	window.header.AddChild(titleLabel)

	window.content = NewLayoutContainer(
		WithSize(width, height-window.headerHeight),
		WithBackground(window.colors.Background),
		WithLayout(NewVerticalStackLayout(10, AlignStart)),
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
