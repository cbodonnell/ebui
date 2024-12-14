package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

var _ InteractiveComponent = &Button{}

type Button struct {
	*BaseComponent
	*BaseInteractive
	label     string
	textColor color.Color
	font      font.Face
	colors    ButtonColors
	isHovered bool
	isPressed bool
}

type ButtonColors struct {
	Default color.Color
	Hovered color.Color
	Pressed color.Color
}

func WithLabel(label string) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Button); ok {
			b.label = label
		}
	}
}

func WithTextColor(color color.Color) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Button); ok {
			b.textColor = color
		}
	}
}

func WithFont(font font.Face) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Button); ok {
			b.font = font
		}
	}
}

func WithClickHandler(handler EventHandler) ComponentOpt {
	return func(c Component) {
		if b, ok := c.(*Button); ok {
			b.eventDispatcher.AddEventListener(EventClick, handler)
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

func NewButton(opts ...ComponentOpt) *Button {
	b := &Button{
		BaseComponent:   NewBaseComponent(opts...),
		BaseInteractive: NewBaseInteractive(),
		textColor:       color.Black,        // Default text color
		font:            basicfont.Face7x13, // Default font
		colors: ButtonColors{ // Default colors
			Default: color.RGBA{200, 200, 200, 255},
			Hovered: color.RGBA{220, 220, 220, 255},
			Pressed: color.RGBA{170, 170, 170, 255},
		},
	}

	for _, opt := range opts {
		opt(b)
	}

	b.registerEventListeners()

	return b
}

func (b *Button) registerEventListeners() {
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
}

func (b *Button) GetLabel() string {
	return b.label
}

func (b *Button) SetLabel(label string) {
	b.label = label
}

func (b *Button) Draw(screen *ebiten.Image) {
	if !b.size.IsDrawable() {
		panic("Button must have a size")
	}
	b.BaseComponent.drawBackground(screen)
	b.draw(screen)
	b.BaseComponent.drawDebug(screen)
}

// draw renders the button to the screen
func (b Button) draw(screen *ebiten.Image) {
	pos := b.GetAbsolutePosition()
	size := b.GetSize()

	// Draw background
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X, pos.Y)

	bgColor := b.colors.Default
	if b.isPressed {
		bgColor = b.colors.Pressed
	} else if b.isHovered {
		bgColor = b.colors.Hovered
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
}
