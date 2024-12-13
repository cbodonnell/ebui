package ebui

import "math"

func clamp(value, min, max float64) float64 {
	return math.Min(math.Max(value, min), max)
}
