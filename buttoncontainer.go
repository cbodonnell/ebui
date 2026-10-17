package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ FocusableComponent = &ButtonContainer{}

type ButtonContainer struct {
	*LayoutContainer
	*BaseFocusable
	colors    ButtonColors
	isHovered bool
	isPressed bool
	onClick   func()
}

type ButtonColors struct {
	Default     color.Color
	Hovered     color.Color
	Pressed     color.Color
	FocusBorder color.Color
}

func DefaultButtonColors() ButtonColors {
	return ButtonColors{
		Default:     color.RGBA{200, 200, 200, 255},
		Hovered:     color.RGBA{220, 220, 220, 255},
		Pressed:     color.RGBA{170, 170, 170, 255},
		FocusBorder: color.Black,
	}
}

func WithButtonColors(colors ButtonColors) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*ButtonContainer); ok {
			b.colors = colors
		}
	}
}

func WithClickHandler(handler func()) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*ButtonContainer); ok {
			b.onClick = handler
		}
	}
}

func NewButtonContainer(opts ...ComponentOpt) *ButtonContainer {
	b := &ButtonContainer{
		LayoutContainer: NewLayoutContainer(
			WithLayout(NewVerticalStackLayout(0, AlignStart)),
		),
		BaseFocusable: NewBaseFocusable(),
		colors:        DefaultButtonColors(),
		onClick:       func() {},
	}

	for _, opt := range opts {
		opt(b)
	}

	b.registerEventListeners()

	return b
}

func (b *ButtonContainer) registerEventListeners() {
	b.AddEventListener(MouseEnter, func(e *Event) {
		b.isHovered = true
	})

	b.AddEventListener(MouseLeave, func(e *Event) {
		b.isHovered = false
		b.isPressed = false
	})

	b.AddEventListener(MouseDown, func(e *Event) {
		b.isPressed = true
	})

	b.AddEventListener(MouseUp, func(e *Event) {
		if b.isPressed && b.isHovered {
			b.onClick()
		}
		b.isPressed = false
	})

	b.AddEventListener(Nav, func(e *Event) {
		if e.Phase == PhaseTarget && e.Nav == NavActivate {
			b.onClick()
			e.Handled = true
		}
	})
}

func (b *ButtonContainer) SetClickHandler(handler func()) {
	b.onClick = handler
}

func (b *ButtonContainer) Update() error {
	b.updateAppearance()
	return b.BaseContainer.Update()
}

func (b *ButtonContainer) updateAppearance() {
	var bgColor color.Color
	switch {
	case b.isPressed:
		bgColor = b.colors.Pressed
	case b.isHovered:
		bgColor = b.colors.Hovered
	default:
		bgColor = b.colors.Default
	}
	b.BaseContainer.SetBackground(bgColor)
}

func (b *ButtonContainer) Draw(screen *ebiten.Image) {
	if b.IsHidden() {
		return
	}

	b.drawFace(screen)
	for _, child := range b.children {
		child.Draw(screen)
	}
}

// drawFace draws the button's focus ring and face.
func (b *ButtonContainer) drawFace(screen *ebiten.Image) {
	if b.IsFocusVisible() {
		paint(screen, Paint{Part: PartFocusRing, Role: b.GetRole(), State: StateFocused, Rect: componentRect(b), Border: b.colors.FocusBorder})
	}
	paint(screen, Paint{Part: PartButton, Role: b.GetRole(), State: b.state(), Rect: componentRect(b), Fill: b.GetBackground()})
	b.drawDebug(screen)
}

// state is the button's interaction state, for the painter.
func (b *ButtonContainer) state() State {
	var s State
	if b.isHovered {
		s |= StateHovered
	}
	if b.isPressed {
		s |= StatePressed
	}
	if b.IsFocusVisible() {
		s |= StateFocused
	}
	if b.IsDisabled() {
		s |= StateDisabled
	}
	return s
}

// DrawFocusRing draws the focus ring around a component.
func DrawFocusRing(screen *ebiten.Image, c Component, col color.Color) {
	if col == nil {
		return
	}
	role := ""
	if r, ok := c.(interface{ GetRole() string }); ok {
		role = r.GetRole()
	}
	paint(screen, Paint{Part: PartFocusRing, Role: role, State: StateFocused, Rect: componentRect(c), Border: col})
}
