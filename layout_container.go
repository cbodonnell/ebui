// layout_container.go
package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type LayoutContainer struct {
	*BaseContainer
	layout     Layout
	background color.Color
}

// Convenience constructors for common layouts
func NewVStackContainer(spacing float64, alignment Alignment) *LayoutContainer {
	return NewLayoutContainer(NewVerticalStack(StackConfig{
		Spacing:   spacing,
		Alignment: alignment,
	}))
}

func NewHStackContainer(spacing float64, alignment Alignment) *LayoutContainer {
	return NewLayoutContainer(NewHorizontalStack(StackConfig{
		Spacing:   spacing,
		Alignment: alignment,
	}))
}

func NewLayoutContainer(layout Layout) *LayoutContainer {
	return &LayoutContainer{
		BaseContainer: NewBaseContainer(),
		layout:        layout,
		background:    color.RGBA{0, 0, 0, 0}, // Transparent by default
	}
}

func (c *LayoutContainer) SetLayout(layout Layout) {
	c.layout = layout
}

func (c *LayoutContainer) SetBackground(color color.Color) {
	c.background = color
}

func (c *LayoutContainer) Draw(screen *ebiten.Image) {
	// Draw background if set
	if c.background != nil {
		pos := c.GetAbsolutePosition()
		size := c.GetSize()
		bg := ebiten.NewImage(int(size.Width), int(size.Height))
		bg.Fill(c.background)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(bg, op)
	}

	c.BaseContainer.Draw(screen)
}

func (c *LayoutContainer) Update() error {
	if c.layout != nil {
		c.layout.ArrangeChildren(c)
	}
	return c.BaseContainer.Update()
}
