package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var _ FocusableComponent = &ButtonContainer{}

type ButtonContainer struct {
	*LayoutContainer
	*BaseFocusable
	colors    ButtonColors
	isHovered bool
	isPressed bool
	isFocused bool
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

	b.AddEventListener(Focus, func(e *Event) {
		b.isFocused = true
	})

	b.AddEventListener(Blur, func(e *Event) {
		b.isFocused = false
	})
}

func (b *ButtonContainer) SetClickHandler(handler func()) {
	b.onClick = handler
}

func (b *ButtonContainer) Update() error {
	b.handleInput()
	b.updateAppearance()
	return b.BaseContainer.Update()
}

func (b *ButtonContainer) handleInput() {
	if !b.isFocused {
		return
	}

	enterPressed := inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	spacePressed := inpututil.IsKeyJustPressed(ebiten.KeySpace)
	if enterPressed || spacePressed {
		b.onClick()
	}
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

	if b.isFocused {
		// Draw the focus border 1px
		pos := b.GetAbsolutePosition()
		size := b.GetSize()
		focusBorder := GetCache().BorderImageWithColor(int(size.Width+2), int(size.Height+2), b.colors.FocusBorder)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pos.X-1, pos.Y-1)
		screen.DrawImage(focusBorder, op)
	}
	b.BaseContainer.Draw(screen)
}
