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
	text        string
	justify     Justify
	color       color.Color
	font        font.Face
	wrap        bool
	lines       []string
	lineSpacing int
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
		}
	}
}

func WithLineSpacing(spacing int) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Label); ok {
			b.lineSpacing = spacing
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
	}

	for _, opt := range opts {
		opt(b)
	}

	b.calculateWrappedLines()

	return b
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
}

func (b *Label) GetText() string {
	return b.text
}

func (b *Label) SetText(text string) {
	if b.text != text {
		b.text = text
		b.calculateWrappedLines()
	}
}

func (b *Label) GetColor() color.Color {
	return b.color
}

func (b *Label) SetColor(color color.Color) {
	b.color = color
}

func (b *Label) GetNumberOfLines() int {
	return len(b.lines)
}

func (b *Label) GetLineHeight() int {
	return b.font.Metrics().Height.Ceil()
}

func (b *Label) GetTextHeight() int {
	lineCount := b.GetNumberOfLines()
	if lineCount <= 0 {
		return 0
	}

	lineHeight := b.GetLineHeight()

	// For a single line, just return line height
	if lineCount == 1 {
		return lineHeight
	}

	// For multiple lines, add spacing between lines
	return lineHeight + (lineCount-1)*(lineHeight+b.lineSpacing)
}

func (b *Label) Draw(screen *ebiten.Image) {
	if !b.size.IsDrawable() {
		panic("Label must have a size")
	}
	if b.hidden {
		return
	}
	b.BaseComponent.drawBackground(screen)
	b.draw(screen)
	b.BaseComponent.drawDebug(screen)
}

// draw renders the label to the screen
func (b Label) draw(screen *ebiten.Image) {
	pos := b.GetAbsolutePosition()
	size := b.GetSize()
	padding := b.GetPadding()

	// Calculate total height of text with line spacing
	lineHeight := b.font.Metrics().Height.Ceil()
	totalLineHeight := lineHeight + b.lineSpacing
	totalTextHeight := totalLineHeight*(len(b.lines)-1) + lineHeight
	if len(b.lines) == 0 {
		totalTextHeight = 0
	}

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

		// Use line spacing when calculating Y position
		lineY := startY + float64(totalLineHeight*i) + float64(b.font.Metrics().Ascent.Ceil())
		text.Draw(screen, line, b.font, int(textX), int(lineY), b.color)
	}
}
