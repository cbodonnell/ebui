package ebui

import (
	"image/color"
)

var (
	Debug = false

	debugColors = []color.Color{
		color.RGBA{255, 0, 0, 25},     // Red, more transparent
		color.RGBA{0, 255, 0, 25},     // Green, more transparent
		color.RGBA{0, 0, 255, 25},     // Blue, more transparent
		color.RGBA{255, 0, 255, 25},   // Magenta, more transparent
		color.RGBA{0, 255, 255, 25},   // Cyan, more transparent
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
