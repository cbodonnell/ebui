package ebui

import "github.com/hajimehoshi/ebiten/v2"

type EbitenLifecycle interface {
	Update() error
	Draw(screen *ebiten.Image)
}
