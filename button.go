package ebui

import (
	"image/color"
)

var _ InteractiveComponent = &Button{}

type Button struct {
	*LayoutContainer
	*BaseInteractive
	label     *Label
	colors    ButtonColors
	isHovered bool
	isPressed bool
	onClick   func()
}

type ButtonColors struct {
	Default color.Color
	Hovered color.Color
	Pressed color.Color
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
		BaseInteractive: NewBaseInteractive(),
		colors: ButtonColors{ // Default colors
			Default: color.RGBA{200, 200, 200, 255},
			Hovered: color.RGBA{220, 220, 220, 255},
			Pressed: color.RGBA{170, 170, 170, 255},
		},
		onClick: func() {}, // Default click handler
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
	return b.LayoutContainer.Update()
}
