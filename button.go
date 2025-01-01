package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var _ FocusableComponent = &Button{}

type Button struct {
	*LayoutContainer
	*BaseFocusable
	label     *Label
	colors    ButtonColors
	isHovered bool
	isPressed bool
	isFocused bool
	onClick   func()
	focusable bool
	tabIndex  int
}

type ButtonColors struct {
	Default     color.Color
	Hovered     color.Color
	Pressed     color.Color
	FocusBorder color.Color
}

func WithLabelText(text string) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Button); ok {
			b.label.SetText(text)
		}
	}
}

func WithButtonColors(colors ButtonColors) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Button); ok {
			b.colors = colors
		}
	}
}

func WithClickHandler(handler func()) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Button); ok {
			b.onClick = handler
		}
	}
}

func NewButton(opts ...ComponentOpt) *Button {
	// Button layout is a vertical stack layout with no spacing and items aligned to the start
	withLayout := WithLayout(NewVerticalStackLayout(0, AlignStart))
	b := &Button{
		LayoutContainer: NewLayoutContainer(
			append([]ComponentOpt{withLayout}, opts...)...,
		),
		BaseFocusable: NewBaseFocusable(),
		colors: ButtonColors{
			Default:     color.RGBA{200, 200, 200, 255}, // Light gray
			Hovered:     color.RGBA{220, 220, 220, 255}, // Light gray
			Pressed:     color.RGBA{170, 170, 170, 255}, // Dark gray
			FocusBorder: color.Black,
		},
		onClick:   func() {},
		focusable: true,
		tabIndex:  0,
	}

	// Button label fills the button's width and is centered
	b.label = NewLabel(
		"",
		WithSize(b.size.Width, b.size.Height),
		WithJustify(JustifyCenter),
	)
	b.AddChild(b.label)

	for _, opt := range opts {
		opt(b)
	}

	b.registerEventListeners()

	return b
}

func (b *Button) registerEventListeners() {
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

func (b *Button) SetClickHandler(handler func()) {
	b.onClick = handler
}

func (b *Button) GetLabel() string {
	return b.label.GetText()
}

func (b *Button) SetLabel(text string) {
	b.label.SetText(text)
}

func (b *Button) Update() error {
	b.handleInput()
	b.updateAppearance()
	return b.LayoutContainer.Update()
}

func (b *Button) handleInput() {
	if !b.isFocused {
		return
	}

	enterPressed := inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	spacePressed := inpututil.IsKeyJustPressed(ebiten.KeySpace)
	if enterPressed || spacePressed {
		b.onClick()
	}
}

func (b *Button) updateAppearance() {
	var bgColor color.Color
	switch {
	case b.isPressed:
		bgColor = b.colors.Pressed
	case b.isHovered:
		bgColor = b.colors.Hovered
	default:
		bgColor = b.colors.Default
	}
	b.LayoutContainer.SetBackground(bgColor)
}

func (b *Button) Draw(screen *ebiten.Image) {
	if b.isFocused {
		// Draw the focus border 1px
		pos := b.GetAbsolutePosition()
		size := b.GetSize()
		bg := ebiten.NewImage(int(size.Width+2), int(size.Height+2))
		bg.Fill(b.colors.FocusBorder)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pos.X-1, pos.Y-1)
		screen.DrawImage(bg, op)
	}
	b.LayoutContainer.Draw(screen)
}
