package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

var _ Component = &Label{}

type Label struct {
	*BaseComponent
	text    string
	justify Justify
	color   color.Color
	font    font.Face
}

type Justify int

const (
	JustifyLeft Justify = iota
	JustifyCenter
	JustifyRight
)

func WithColor(color color.Color) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Label); ok {
			b.color = color
		}
	}
}

func WithFont(font font.Face) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Label); ok {
			b.font = font
		}
	}
}

func WithJustify(justify Justify) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Label); ok {
			b.justify = justify
		}
	}
}

func NewLabel(text string, opts ...ComponentOpt) *Label {
	b := &Label{
		BaseComponent: NewBaseComponent(opts...),
		text:          text,
		color:         color.Black,        // Default text color
		font:          basicfont.Face7x13, // Default font
		justify:       JustifyCenter,      // Default text justification
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *Label) GetText() string {
	return b.text
}

func (b *Label) SetText(text string) {
	b.text = text
}

func (b *Label) Draw(screen *ebiten.Image) {
	if !b.size.IsDrawable() {
		panic("Label must have a size")
	}
	b.BaseComponent.drawBackground(screen)
	b.draw(screen)
	b.BaseComponent.drawDebug(screen)
}

// draw renders the button to the screen
func (b Label) draw(screen *ebiten.Image) {
	pos := b.GetAbsolutePosition()
	size := b.GetSize()
	padding := b.GetPadding()

	// Draw text
	bounds, _ := font.BoundString(b.font, b.text)
	textWidth := (bounds.Max.X - bounds.Min.X).Ceil()
	textHeight := (bounds.Max.Y - bounds.Min.Y).Ceil()
	var textX float64
	switch b.justify {
	case JustifyLeft:
		textX = pos.X + padding.Left
	case JustifyCenter:
		textX = pos.X + padding.Left + (size.Width-float64(textWidth))/2
	case JustifyRight:
		textX = pos.X + size.Width - padding.Right - float64(textWidth)
	}
	textY := pos.Y + padding.Top + (size.Height-float64(textHeight))/2 + float64(textHeight)
	text.Draw(screen, b.text, b.font, int(textX), int(textY), b.color)
}
