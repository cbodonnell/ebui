package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var _ EbitenLifecycle = &Manager{}

type Manager struct {
	root  Component
	input *InputManager
}

type ManagerOpt func(m *Manager)

func WithInputManager(im *InputManager) ManagerOpt {
	return func(m *Manager) {
		m.input = im
	}
}

// NewManager creates a new UI Manager with the given root container.
func NewManager(root Container, opts ...ManagerOpt) *Manager {
	m := &Manager{
		root:  root,
		input: NewInputManager(),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// Update updates the UI Manager.
func (u *Manager) Update() error {
	u.input.Update(u.root)
	return u.root.Update()
}

func (u *Manager) Draw(screen *ebiten.Image) {
	u.root.Draw(screen)
}

// DisableFocus disables focus management for the UI.
func (u *Manager) DisableFocus() {
	u.input.DisableFocusManagement()
}

// EnableFocus enables focus management for the UI.
func (u *Manager) EnableFocus() {
	u.input.EnableFocusManagement()
}

// IsFocusEnabled returns whether focus management is enabled.
func (u *Manager) IsFocusEnabled() bool {
	return u.input.IsFocusManagementEnabled()
}
