package ebui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

var (
	Debug = false

	debugColors = []color.Color{
		color.RGBA{255, 0, 0, 25},     // Red, more transparent
		color.RGBA{0, 255, 0, 25},     // Green, more transparent
		color.RGBA{0, 0, 255, 25},     // Blue, more transparent
		color.RGBA{255, 0, 255, 25},   // Magenta, more transparent
		color.RGBA{0, 255, 255, 25},   // Cyan, more transparent
		color.RGBA{255, 255, 255, 25}, // White, more transparent
		color.RGBA{0, 0, 0, 25},       // Black, more transparent
		color.RGBA{128, 128, 128, 25}, // Gray, more transparent
		color.RGBA{128, 0, 0, 25},     // Maroon, more transparent
		color.RGBA{128, 128, 0, 25},   // Olive, more transparent
		color.RGBA{0, 128, 0, 25},     // Green, more transparent
		color.RGBA{128, 0, 128, 25},   // Purple, more transparent
		color.RGBA{0, 128, 128, 25},   // Teal, more transparent
		color.RGBA{0, 0, 128, 25},     // Navy, more transparent
	}

	debugDepth = 0
	colorMap   = make(map[uint64]color.Color)
)

// debugDraw draws debugging visualization for a component
func debugDraw(screen *ebiten.Image, component Component) {
	if !Debug {
		return
	}

	// Get a color for this component
	debugColor, ok := colorMap[component.GetID()]
	if !ok {
		debugColor = debugColors[debugDepth%len(debugColors)]
		colorMap[component.GetID()] = debugColor
	}

	pos := component.GetAbsolutePosition()
	size := component.GetSize()
	padding := component.GetPadding()

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
