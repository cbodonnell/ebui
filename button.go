package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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

// Draw draws the button's face, then its label, unless the painter draws the label itself.
func (b *Button) Draw(screen *ebiten.Image) {
	if b.IsHidden() {
		return
	}
	b.drawFace(screen)
	drawn := paintCustom(screen, Paint{
		Part:  PartButtonLabel,
		Role:  b.GetRole(),
		State: b.state(),
		Rect:  componentRect(b),
		Fill:  b.label.GetColor(),
		Text:  b.label.GetText(),
	})
	for _, child := range b.children {
		if drawn && child == b.label {
			continue
		}
		child.Draw(screen)
	}
}
