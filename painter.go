package ebui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Part is a piece of a component that a Painter draws.
type Part int

const (
	PartWindow       Part = iota // a window's frame and ground
	PartWindowHeader             // the strip across a window's top
	PartWindowTitle              // a window's title; Text is set
	PartButton                   // a button's face
	PartButtonLabel              // a button's label; Text is set
	PartFocusRing                // the ring around a focused component
	PartTooltip                  // a tooltip's frame and ground
	PartSliderTrack              // a slider's groove
	PartSliderFill               // the groove up to the thumb
	PartSliderThumb              // a slider's thumb
	PartScrollTrack              // a scroll bar's groove
	PartScrollThumb              // a scroll bar's thumb
	PartTextInput                // a text input's well
	PartCaret                    // a text input's caret
)

// State is the interaction state a part is painted in.
type State int

const (
	StateHovered State = 1 << iota
	StatePressed
	StateFocused
	StateDisabled
)

func (s State) Has(f State) bool { return s&f != 0 }

// Paint is one part to draw. Fill and Border are the component's flat colors, which the flat
// look draws and a painter may ignore.
type Paint struct {
	Part   Part
	Role   string // the component's role, set with WithRole or WithWindowRole
	State  State
	Rect   image.Rectangle
	Fill   color.Color
	Border color.Color
	Text   string
}

// Painter draws parts. Paint reports whether it drew p; a part it declines is drawn flat. For
// PartWindowTitle and PartButtonLabel, drawing the part replaces the component's own label.
type Painter interface {
	Paint(dst *ebiten.Image, p Paint) bool
}

// DefaultPainter draws every component; nil draws them all flat. Set it before drawing.
var DefaultPainter Painter

// WithRole names what a component is for (primary, danger…), for the painter.
func WithRole(role string) ComponentOpt {
	return func(c Component) {
		if r, ok := c.(interface{ SetRole(string) }); ok {
			r.SetRole(role)
		}
	}
}

// paint draws p with the default painter, or flat if it declines.
func paint(dst *ebiten.Image, p Paint) {
	if !paintCustom(dst, p) {
		paintFlat(dst, p)
	}
}

// paintCustom draws p with the default painter, reporting whether it did.
func paintCustom(dst *ebiten.Image, p Paint) bool {
	return DefaultPainter != nil && DefaultPainter.Paint(dst, p)
}

// paintFlat draws p in its flat colors: Fill over the rect, and Border as a 1px line, outside
// the rect for windows and focus rings and inside it for tooltips.
func paintFlat(dst *ebiten.Image, p Paint) {
	r := p.Rect
	switch p.Part {
	case PartWindow:
		fillRect(dst, r.Inset(-1), p.Border)
		fillRect(dst, r, p.Fill)
	case PartFocusRing:
		borderRect(dst, r.Inset(-1), p.Border)
	case PartTooltip:
		fillRect(dst, r, p.Fill)
		borderRect(dst, r, p.Border)
	case PartWindowTitle, PartButtonLabel:
		// the component's label draws itself
	default:
		fillRect(dst, r, p.Fill)
	}
}

func fillRect(dst *ebiten.Image, r image.Rectangle, c color.Color) {
	if c == nil || r.Dx() <= 0 || r.Dy() <= 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
	dst.DrawImage(GetCache().ImageWithColor(r.Dx(), r.Dy(), c), op)
}

func borderRect(dst *ebiten.Image, r image.Rectangle, c color.Color) {
	if c == nil || r.Dx() <= 0 || r.Dy() <= 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
	dst.DrawImage(GetCache().BorderImageWithColor(r.Dx(), r.Dy(), c), op)
}

// rectAt is the rectangle at pos of size, in whole pixels.
func rectAt(x, y, w, h float64) image.Rectangle {
	return image.Rect(int(x), int(y), int(x)+int(w), int(y)+int(h))
}

// componentRect is c's on-screen rectangle.
func componentRect(c Component) image.Rectangle {
	pos := c.GetAbsolutePosition()
	size := c.GetSize()
	return rectAt(pos.X, pos.Y, size.Width, size.Height)
}
