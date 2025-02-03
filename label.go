package ebui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

var _ Component = &Label{}

type Label struct {
	*BaseComponent
	text       string
	justify    Justify
	color      color.Color
	font       font.Face
	wrap       bool
	lines      []string
	linesDirty bool
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

func WithTextWrap() ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Label); ok {
			b.wrap = true
			b.linesDirty = true
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
		wrap:          false,              // Default to no wrapping
		linesDirty:    true,
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
	if b.text != text {
		b.text = text
		b.linesDirty = true
	}
}

func (b *Label) GetColor() color.Color {
	return b.color
}

func (b *Label) SetColor(color color.Color) {
	b.color = color
}

func (b *Label) Draw(screen *ebiten.Image) {
	if !b.size.IsDrawable() {
		panic("Label must have a size")
	}
	b.BaseComponent.drawBackground(screen)
	b.draw(screen)
	b.BaseComponent.drawDebug(screen)
}

// calculateWrappedLines splits the text into lines that fit within the label width
func (b *Label) calculateWrappedLines() {
	if !b.wrap {
		b.lines = []string{b.text}
		return
	}

	padding := b.GetPadding()
	maxWidth := b.size.Width - padding.Left - padding.Right

	// Reset lines
	b.lines = []string{}

	// Split text into words
	words := strings.Fields(b.text)
	if len(words) == 0 {
		return
	}

	currentLine := words[0]

	for _, word := range words[1:] {
		// Try adding the next word
		testLine := currentLine + " " + word
		bounds, _ := font.BoundString(b.font, testLine)
		width := (bounds.Max.X - bounds.Min.X).Ceil()

		if float64(width) <= maxWidth {
			// Word fits, add it to the current line
			currentLine = testLine
		} else {
			// Word doesn't fit, start a new line
			b.lines = append(b.lines, currentLine)
			currentLine = word
		}
	}

	// Add the last line
	if currentLine != "" {
		b.lines = append(b.lines, currentLine)
	}

	b.linesDirty = false
}

// draw renders the label to the screen
func (b Label) draw(screen *ebiten.Image) {
	pos := b.GetAbsolutePosition()
	size := b.GetSize()
	padding := b.GetPadding()

	// Calculate wrapped lines if needed
	if b.linesDirty {
		b.calculateWrappedLines()
	}

	// Calculate total height of text
	lineHeight := b.font.Metrics().Height.Ceil()
	totalTextHeight := lineHeight * len(b.lines)

	// Calculate starting Y position based on total text height
	startY := pos.Y + padding.Top + (size.Height-float64(totalTextHeight))/2

	// Draw each line
	for i, line := range b.lines {
		bounds, _ := font.BoundString(b.font, line)
		textWidth := (bounds.Max.X - bounds.Min.X).Ceil()

		var textX float64
		switch b.justify {
		case JustifyLeft:
			textX = pos.X + padding.Left
		case JustifyCenter:
			textX = pos.X + padding.Left + (size.Width-float64(textWidth))/2
		case JustifyRight:
			textX = pos.X + size.Width - padding.Right - float64(textWidth)
		}

		lineY := startY + float64(lineHeight*i) + float64(b.font.Metrics().Ascent.Ceil())
		text.Draw(screen, line, b.font, int(textX), int(lineY), b.color)
	}
}
