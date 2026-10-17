package ebui

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// NavAction is a keyboard navigation input, delivered to the focused component as a Nav event.
type NavAction int

const (
	NavNone NavAction = iota
	NavUp
	NavDown
	NavLeft
	NavRight
	NavNext
	NavPrev
	NavActivate
	NavCancel
)

// NavInput reports this tick's navigation actions.
type NavInput func() []NavAction

// KeyboardNavInput reads the arrows, Tab, Enter and Escape.
func KeyboardNavInput() NavInput {
	return func() []NavAction {
		var actions []NavAction
		for _, k := range []struct {
			key    ebiten.Key
			action NavAction
		}{
			{ebiten.KeyArrowUp, NavUp},
			{ebiten.KeyArrowDown, NavDown},
			{ebiten.KeyArrowLeft, NavLeft},
			{ebiten.KeyArrowRight, NavRight},
		} {
			if KeyRepeated(k.key) {
				actions = append(actions, k.action)
			}
		}
		if KeyRepeated(ebiten.KeyTab) {
			if ebiten.IsKeyPressed(ebiten.KeyShift) {
				actions = append(actions, NavPrev)
			} else {
				actions = append(actions, NavNext)
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) {
			actions = append(actions, NavActivate)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			actions = append(actions, NavCancel)
		}
		return actions
	}
}

// KeyRepeated reports whether the key was just pressed or is repeating.
func KeyRepeated(key ebiten.Key) bool {
	return repeats(inpututil.KeyPressDuration(key), ebiten.TPS())
}

func repeats(d, tps int) bool {
	delay, interval := tps*2/5, max(tps/15, 1)
	return d == 1 || d > delay && (d-delay)%interval == 0
}

type rect struct {
	x0, y0, x1, y1 float64
}

func rectOf(c Component) rect {
	pos := c.GetAbsolutePosition()
	size := c.GetSize()
	return rect{pos.X, pos.Y, pos.X + size.Width, pos.Y + size.Height}
}

func (r rect) center() (float64, float64) {
	return (r.x0 + r.x1) / 2, (r.y0 + r.y1) / 2
}

func gap(a0, a1, b0, b1 float64) float64 {
	return math.Max(0, math.Max(b0-a1, a0-b1))
}

// directionalScore rates a move from one rect to another; lower is better.
func directionalScore(from, to rect, action NavAction) (float64, bool) {
	fx, fy := from.center()
	tx, ty := to.center()
	var along, across, offset float64
	switch action {
	case NavRight:
		if tx <= fx || to.x0 < from.x0+1 {
			return 0, false
		}
		along, across, offset = math.Max(0, to.x0-from.x1), gap(from.y0, from.y1, to.y0, to.y1), math.Abs(ty-fy)
	case NavLeft:
		if tx >= fx || to.x1 > from.x1-1 {
			return 0, false
		}
		along, across, offset = math.Max(0, from.x0-to.x1), gap(from.y0, from.y1, to.y0, to.y1), math.Abs(ty-fy)
	case NavDown:
		if ty <= fy || to.y0 < from.y0+1 {
			return 0, false
		}
		along, across, offset = math.Max(0, to.y0-from.y1), gap(from.x0, from.x1, to.x0, to.x1), math.Abs(tx-fx)
	case NavUp:
		if ty >= fy || to.y1 > from.y1-1 {
			return 0, false
		}
		along, across, offset = math.Max(0, from.y0-to.y1), gap(from.x0, from.x1, to.x0, to.x1), math.Abs(tx-fx)
	default:
		return 0, false
	}
	return along + across*4 + offset*0.1, true
}

func (fm *FocusManager) nearest(action NavAction) FocusableComponent {
	from := rectOf(fm.currentFocus)
	var best FocusableComponent
	bestScore := math.Inf(1)
	for _, c := range fm.inScope(fm.currentScope) {
		if c == fm.currentFocus {
			continue
		}
		if score, ok := directionalScore(from, rectOf(c), action); ok && score < bestScore {
			best, bestScore = c, score
		}
	}
	return best
}

func (fm *FocusManager) nearestTo(scope Component, r rect) FocusableComponent {
	rx, ry := r.center()
	var best FocusableComponent
	bestDist := math.Inf(1)
	for _, c := range fm.inScope(scope) {
		cx, cy := rectOf(c).center()
		if d := math.Hypot(cx-rx, cy-ry); d < bestDist {
			best, bestDist = c, d
		}
	}
	return best
}
