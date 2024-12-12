package ebui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// Position represents the position of a component
type Position struct {
	X, Y             float64
	RelativeToParent bool
}

// Size represents the dimensions of a component
type Size struct {
	Width, Height float64
	AutoWidth     bool
	AutoHeight    bool
}

// Padding represents padding around a component
type Padding struct {
	Top, Right, Bottom, Left float64
}

// Component is the base interface that all UI elements must implement
type Component interface {
	Identifiable
	EbitenLifecycle
	SetPosition(pos Position)
	GetPosition() Position
	SetSize(size Size)
	GetSize() Size
	SetParent(parent Container)
	GetParent() Container
	SetPadding(padding Padding)
	GetPadding() Padding
	Contains(x, y float64) bool
	GetAbsolutePosition() Position
}

var _ Component = &BaseComponent{}

// BaseComponent provides common functionality for all components
type BaseComponent struct {
	id         uint64
	position   Position
	size       Size
	padding    Padding
	background color.Color
	parent     Container
}

func NewBaseComponent() *BaseComponent {
	return &BaseComponent{
		id:         GenerateID(),
		background: color.RGBA{0, 0, 0, 0}, // Transparent by default
	}
}

func (b *BaseComponent) GetID() uint64 {
	return b.id
}

func (b *BaseComponent) Update() error {
	return nil
}

func (b *BaseComponent) Draw(screen *ebiten.Image) {
	b.drawBackground(screen)
	b.drawDebug(screen)
}

func (b *BaseComponent) SetPosition(pos Position) {
	b.position = pos
}

func (b *BaseComponent) GetPosition() Position {
	return b.position
}

func (b *BaseComponent) SetSize(size Size) {
	b.size = size
}

func (b *BaseComponent) GetSize() Size {
	return b.size
}

func (b *BaseComponent) SetParent(parent Container) {
	b.parent = parent
}

func (b *BaseComponent) GetParent() Container {
	return b.parent
}

func (b *BaseComponent) SetPadding(padding Padding) {
	b.padding = padding
}

func (b *BaseComponent) GetPadding() Padding {
	return b.padding
}

func (b *BaseComponent) SetBackground(color color.Color) {
	b.background = color
}

func (b *BaseComponent) drawBackground(screen *ebiten.Image) {
	if b.background == nil {
		return
	}
	pos := b.GetAbsolutePosition()
	size := b.GetSize()
	bg := ebiten.NewImage(int(size.Width), int(size.Height))
	bg.Fill(b.background)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(bg, op)
}

func (b *BaseComponent) Contains(x, y float64) bool {
	pos := b.GetAbsolutePosition()
	size := b.GetSize()
	return x >= pos.X && x <= pos.X+size.Width &&
		y >= pos.Y && y <= pos.Y+size.Height
}

func (b *BaseComponent) GetAbsolutePosition() Position {
	pos := b.position
	if b.parent != nil && b.position.RelativeToParent {
		parentPos := b.parent.GetAbsolutePosition()
		pos.X += parentPos.X
		pos.Y += parentPos.Y
	}
	return pos
}

func (b *BaseComponent) drawDebug(screen *ebiten.Image) {
	if !Debug {
		return
	}

	// Get a color for this component
	debugColor, ok := colorMap[b.GetID()]
	if !ok {
		debugColor = debugColors[debugDepth%len(debugColors)]
		colorMap[b.GetID()] = debugColor
	}

	pos := b.GetAbsolutePosition()
	size := b.GetSize()
	padding := b.GetPadding()

	// Draw component bounds
	debugRect := ebiten.NewImage(int(size.Width), int(size.Height))
	debugRect.Fill(debugColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(debugRect, op)

	// Draw padding bounds with even more transparent color
	if padding.Top > 0 || padding.Right > 0 || padding.Bottom > 0 || padding.Left > 0 {
		paddingRect := ebiten.NewImage(
			int(size.Width-padding.Left-padding.Right),
			int(size.Height-padding.Top-padding.Bottom),
		)
		paddingRect.Fill(color.RGBA{255, 255, 255, 15}) // Very transparent white

		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pos.X+padding.Left, pos.Y+padding.Top)
		screen.DrawImage(paddingRect, op)
	}

	// Draw component info with a slight shadow for better visibility
	info := fmt.Sprintf("Pos: (%.0f, %.0f)\nSize: %.0f x %.0f\nPadding: %.0f, %.0f, %.0f, %.0f",
		pos.X, pos.Y, size.Width, size.Height,
		padding.Top, padding.Right, padding.Bottom, padding.Left)

	// Draw text shadow
	text.Draw(screen, info, basicfont.Face7x13,
		int(pos.X)+5, int(pos.Y)+14, color.RGBA{0, 0, 0, 40})
	// Draw text
	text.Draw(screen, info, basicfont.Face7x13,
		int(pos.X)+4, int(pos.Y)+13, color.RGBA{0, 0, 0, 180})

	debugDepth++
}
