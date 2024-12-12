package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type Button struct {
	BaseInteractive
	label           string
	backgroundColor color.Color
	textColor       color.Color
	font            font.Face
	isHovered       bool
	isPressed       bool
}

func NewButton(label string) *Button {
	b := &Button{
		BaseInteractive: NewBaseInteractive(),
		label:           label,
		backgroundColor: color.RGBA{200, 200, 200, 255},
		textColor:       color.Black,
		font:            basicfont.Face7x13,
	}

	// Set up event handlers
	b.eventDispatcher.AddEventListener(EventMouseEnter, func(e Event) {
		b.isHovered = true
	})

	b.eventDispatcher.AddEventListener(EventMouseLeave, func(e Event) {
		b.isHovered = false
		b.isPressed = false
	})

	b.eventDispatcher.AddEventListener(EventMouseDown, func(e Event) {
		b.isPressed = true
	})

	b.eventDispatcher.AddEventListener(EventMouseUp, func(e Event) {
		if b.isPressed && b.isHovered {
			b.eventDispatcher.DispatchEvent(Event{
				Type:      EventClick,
				X:         e.X,
				Y:         e.Y,
				Component: b,
			})
		}
		b.isPressed = false
	})

	return b
}

func (b *Button) OnClick(handler func()) {
	b.eventDispatcher.AddEventListener(EventClick, func(e Event) {
		handler()
	})
}

func (b *Button) Draw(screen *ebiten.Image) {
	pos := b.GetAbsolutePosition()
	size := b.GetSize()

	// Draw background
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X, pos.Y)

	bgColor := b.backgroundColor
	if b.isPressed {
		bgColor = color.RGBA{170, 170, 170, 255}
	} else if b.isHovered {
		bgColor = color.RGBA{220, 220, 220, 255}
	}

	buttonImage := ebiten.NewImage(int(size.Width), int(size.Height))
	buttonImage.Fill(bgColor)
	screen.DrawImage(buttonImage, op)

	// Draw text
	bounds, _ := font.BoundString(b.font, b.label)
	textWidth := (bounds.Max.X - bounds.Min.X).Ceil()
	textHeight := (bounds.Max.Y - bounds.Min.Y).Ceil()
	textX := pos.X + (size.Width-float64(textWidth))/2
	textY := pos.Y + (size.Height-float64(textHeight))/2 + float64(textHeight)
	text.Draw(screen, b.label, b.font, int(textX), int(textY), b.textColor)

	b.BaseComponent.Draw(screen)
}
