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
}

// DefaultWindowColors returns a default color scheme for windows
func DefaultWindowColors() WindowColors {
	return WindowColors{
		Background: color.RGBA{240, 240, 240, 255},
		Header:     color.RGBA{200, 200, 200, 255},
	}
}

type Window struct {
	*BaseInteractive
	*LayoutContainer
	manager       *WindowManager
	header        *BaseComponent
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

// Show makes the window visible
func (w *Window) Show() {
	w.state = WindowStateNormal
	w.manager.SetActiveWindow(w)
}

// Hide makes the window invisible
func (w *Window) Hide() {
	w.state = WindowStateHidden
	w.closeCallback()
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
	w.LayoutContainer.Draw(screen)
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
	headerHeight := 30.0

	// Create the window
	window := &Window{
		BaseInteractive: NewBaseInteractive(),
		LayoutContainer: NewLayoutContainer(
			WithSize(width, height),
			WithLayout(NewVerticalStackLayout(0, AlignStart)),
		),
		manager:      wm,
		headerHeight: headerHeight,
		colors:       DefaultWindowColors(),
		state:        WindowStateNormal,
	}

	for _, opt := range opts {
		opt(window)
	}

	// Create header and content area that fills the window
	window.header = NewBaseComponent(
		WithSize(width, headerHeight),
		WithBackground(window.colors.Header),
	)
	window.content = NewLayoutContainer(
		WithSize(width, height-headerHeight),
		WithBackground(window.colors.Background),
		WithPadding(10, 10, 10, 10),
		WithLayout(NewVerticalStackLayout(10, AlignStart)),
	)
	window.LayoutContainer.AddChild(window.header)
	window.LayoutContainer.AddChild(window.content)

	window.registerEventListeners()
	window.SetPosition(Position{
		X:        100,
		Y:        100,
		ZIndex:   wm.nextZIndex,
		Relative: false,
	})
	wm.nextZIndex++

	wm.AddChild(window)
	wm.SetActiveWindow(window)

	return window
}

func (w *Window) registerEventListeners() {
	w.eventDispatcher.AddEventListener(EventMouseDown, func(e Event) {
		// Activate window on any mouse down
		w.manager.SetActiveWindow(w)

		// Start dragging if over header
		if w.isOverHeader(e.X, e.Y) {
			w.isDragging = true
			w.dragStartX = e.X
			w.dragStartY = e.Y
			absPos := w.GetAbsolutePosition()
			w.windowStartX = absPos.X
			w.windowStartY = absPos.Y
		}
	})

	w.eventDispatcher.AddEventListener(EventMouseUp, func(e Event) {
		w.isDragging = false
	})

	w.eventDispatcher.AddEventListener(EventMouseMove, func(e Event) {
		if w.isDragging {
			deltaX := e.X - w.dragStartX
			deltaY := e.Y - w.dragStartY
			newPos := w.GetPosition()
			newPos.X = w.windowStartX + deltaX
			newPos.Y = w.windowStartY + deltaY
			w.SetPosition(newPos)
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
