package ebui

type LayoutContainer struct {
	*BaseContainer
	layout Layout
}

func WithLayout(layout Layout) ComponentOpt {
	return func(c Component) {
		if lc, ok := c.(*LayoutContainer); ok {
			lc.layout = layout
		}
	}
}

func NewLayoutContainer(opts ...ComponentOpt) *LayoutContainer {
	l := &LayoutContainer{
		BaseContainer: NewBaseContainer(opts...),
		layout:        NewStackLayout(), // Default layout
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (c *LayoutContainer) Update() error {
	if c.layout != nil {
		c.layout.ArrangeChildren(c)
	}
	return c.BaseContainer.Update()
}
