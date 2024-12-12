package ebui

type LayoutContainer struct {
	*BaseContainer
	layout Layout
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
	}
}

func (c *LayoutContainer) SetLayout(layout Layout) {
	c.layout = layout
}

func (c *LayoutContainer) Update() error {
	if c.layout != nil {
		c.layout.ArrangeChildren(c)
	}
	return c.BaseContainer.Update()
}
