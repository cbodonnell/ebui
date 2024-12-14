package ebui

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type ZIndexedContainer struct {
	*BaseContainer
}

func NewZIndexedContainer(opts ...ComponentOpt) *ZIndexedContainer {
	z := &ZIndexedContainer{
		BaseContainer: NewBaseContainer(opts...),
	}
	for _, opt := range opts {
		opt(z)
	}
	return z
}

// GetChildren returns the children of the container sorted by ZIndex.
func (z *ZIndexedContainer) GetChildren() []Component {
	sorted := make([]Component, len(z.children))
	copy(sorted, z.children)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].GetPosition().ZIndex < sorted[j].GetPosition().ZIndex
	})
	return sorted
}

func (z *ZIndexedContainer) Update() error {
	for _, child := range z.GetChildren() {
		if err := child.Update(); err != nil {
			return err
		}
	}
	return z.BaseComponent.Update()
}

func (z *ZIndexedContainer) Draw(screen *ebiten.Image) {
	z.BaseComponent.Draw(screen)
	for _, child := range z.GetChildren() {
		child.Draw(screen)
	}
}
