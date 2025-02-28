package ebui

import (
	"image/color"
)

var _ FocusableComponent = &Button{}

type Button struct {
	*ButtonContainer

	label *Label
}

func WithLabelText(text string) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Button); ok {
			b.label.SetText(text)
		}
	}
}

func WithLabelColor(color color.Color) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Button); ok {
			b.label.SetColor(color)
		}
	}
}

func NewButton(opts ...ComponentOpt) *Button {
	b := &Button{
		ButtonContainer: NewButtonContainer(opts...),
	}

	b.label = NewLabel(
		"",
		WithSize(b.size.Width, b.size.Height),
		WithJustify(JustifyCenter),
	)
	b.AddChild(b.label)

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *Button) GetLabel() string {
	return b.label.GetText()
}

func (b *Button) SetLabel(text string) {
	b.label.SetText(text)
}
