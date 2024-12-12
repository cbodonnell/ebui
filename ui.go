package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var _ EbitenLifecycle = &Manager{}

type Manager struct {
	root  Component
	input *InputManager
}

// NewManager creates a new UI Manager with the given root container.
func NewManager(root Container) *Manager {
	return &Manager{
		root:  root,
		input: NewInputManager(),
	}
}

// Update updates the UI Manager.
func (u *Manager) Update() error {
	u.input.Update(u.root)
	return u.root.Update()
}

func (u *Manager) Draw(screen *ebiten.Image) {
	u.root.Draw(screen)
}
