package ebui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Container interface {
	Component
	AddChild(child Component)
	RemoveChild(child Component)
	GetChildren() []Component
}

var _ Container = &BaseContainer{}

type BaseContainer struct {
	*BaseComponent
	children []Component
}

func NewBaseContainer() *BaseContainer {
	return &BaseContainer{
		BaseComponent: NewBaseComponent(),
	}
}

func (c *BaseContainer) AddChild(child Component) {
	child.SetParent(c)
	c.children = append(c.children, child)
}

func (c *BaseContainer) RemoveChild(child Component) {
	for i, ch := range c.children {
		if ch == child {
			c.children = append(c.children[:i], c.children[i+1:]...)
			return
		}
	}
}

func (c *BaseContainer) GetChildren() []Component {
	return c.children
}

func (c *BaseContainer) Update() error {
	for _, child := range c.children {
		if err := child.Update(); err != nil {
			return err
		}
	}
	return c.BaseComponent.Update()
}

func (c *BaseContainer) Draw(screen *ebiten.Image) {
	c.BaseComponent.Draw(screen)
	for _, child := range c.children {
		child.Draw(screen)
	}
}
