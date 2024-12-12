package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
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
	id       uint64
	position Position
	size     Size
	padding  Padding
	parent   Container
}

func NewBaseComponent() *BaseComponent {
	return &BaseComponent{
		id: generateID(),
	}
}

var nextID uint64

func generateID() uint64 {
	nextID++
	if nextID == 0 {
		panic("ID overflow")
	}
	return nextID
}

func (b *BaseComponent) GetID() uint64 {
	return b.id
}

func (b *BaseComponent) Update() error {
	return nil
}

func (b *BaseComponent) Draw(screen *ebiten.Image) {
	debugDraw(screen, b)
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
