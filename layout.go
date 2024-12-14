package ebui

// Alignment represents horizontal or vertical alignment options
type Alignment int

const (
	AlignStart Alignment = iota
	AlignCenter
	AlignEnd
)

// Layout defines how a container should arrange its children
type Layout interface {
	ArrangeChildren(container Container)
	GetMinSize(container Container) Size
}

// StackConfig holds configuration for stack layouts
type StackConfig struct {
	Spacing   float64
	Alignment Alignment
}

// StackLayout implements vertical or horizontal stacking of components
type StackLayout struct {
	Horizontal bool
	Config     StackConfig
}

type StackLayoutOpt func(l *StackLayout)

func WithHorizontal() StackLayoutOpt {
	return func(l *StackLayout) {
		l.Horizontal = true
	}
}

func WithSpacing(spacing float64) StackLayoutOpt {
	return func(l *StackLayout) {
		l.Config.Spacing = spacing
	}
}

func WithAlignment(align Alignment) StackLayoutOpt {
	return func(l *StackLayout) {
		l.Config.Alignment = align
	}
}

func NewStackLayout(opts ...StackLayoutOpt) *StackLayout {
	l := &StackLayout{}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// NewVerticalStackLayout is a helper function to create a vertical stack layout
func NewVerticalStackLayout(spacing float64, align Alignment) *StackLayout {
	return NewStackLayout(
		WithSpacing(spacing),
		WithAlignment(align),
	)
}

// NewHorizontalStackLayout is a helper function to create a horizontal stack layout
func NewHorizontalStackLayout(spacing float64, align Alignment) *StackLayout {
	return NewStackLayout(
		WithHorizontal(),
		WithSpacing(spacing),
		WithAlignment(align),
	)
}

// ArrangeChildren positions all children in a vertical or horizontal stack
func (l *StackLayout) ArrangeChildren(container Container) {
	children := container.GetChildren()
	if len(children) == 0 {
		return
	}

	containerSize := container.GetSize()
	containerPadding := container.GetPadding()

	// Calculate available space
	availableWidth := containerSize.Width - containerPadding.Left - containerPadding.Right
	availableHeight := containerSize.Height - containerPadding.Top - containerPadding.Bottom

	// First pass: calculate total size and count flexible children
	var totalFixedMainAxis float64
	childSizes := make([]Size, len(children))

	for i, child := range children {
		size := child.GetSize()
		childSizes[i] = size

		if !l.Horizontal {
			totalFixedMainAxis += size.Height
		} else {
			totalFixedMainAxis += size.Width
		}
	}

	// Calculate total spacing
	totalSpacing := l.Config.Spacing * float64(len(children)-1)

	// Calculate total width/height including spacing
	totalSize := totalFixedMainAxis + totalSpacing

	// Calculate starting position based on alignment
	startX := containerPadding.Left
	startY := containerPadding.Top

	if l.Horizontal {
		// For horizontal stack, calculate starting X for the entire group
		switch l.Config.Alignment {
		case AlignStart:
			startX = containerPadding.Left
		case AlignCenter:
			startX = containerPadding.Left + (availableWidth-totalSize)/2
		case AlignEnd:
			startX = containerPadding.Left + availableWidth - totalSize
		}
	} else {
		// For vertical stack, calculate starting Y for the entire group
		switch l.Config.Alignment {
		case AlignStart:
			startY = containerPadding.Top
		case AlignCenter:
			startY = containerPadding.Top + (availableHeight-totalSize)/2
		case AlignEnd:
			startY = containerPadding.Top + availableHeight - totalSize
		}
	}

	// Position all children
	currentX := startX
	currentY := startY

	for i, child := range children {
		size := childSizes[i]
		pos := Position{Relative: true}

		if !l.Horizontal {
			// For vertical stack
			pos.Y = currentY
			switch l.Config.Alignment {
			case AlignStart:
				pos.X = containerPadding.Left
			case AlignCenter:
				pos.X = containerPadding.Left + (availableWidth-size.Width)/2
			case AlignEnd:
				pos.X = containerPadding.Left + availableWidth - size.Width
			}
			currentY += size.Height + l.Config.Spacing
		} else {
			// For horizontal stack
			pos.X = currentX
			// Center vertically within the container
			pos.Y = containerPadding.Top + (availableHeight-size.Height)/2
			currentX += size.Width + l.Config.Spacing
		}

		child.SetPosition(pos)
	}
}

// GetMinSize returns the minimum size required to fit all children
func (l *StackLayout) GetMinSize(container Container) Size {
	children := container.GetChildren()
	if len(children) == 0 {
		return Size{}
	}

	padding := container.GetPadding()
	totalWidth, totalHeight := 0.0, 0.0
	maxWidth, maxHeight := 0.0, 0.0

	for i, child := range children {
		size := child.GetSize()

		if !l.Horizontal {
			totalHeight += size.Height
			if i > 0 {
				totalHeight += l.Config.Spacing
			}
			if size.Width > maxWidth {
				maxWidth = size.Width
			}
		} else {
			totalWidth += size.Width
			if i > 0 {
				totalWidth += l.Config.Spacing
			}
			if size.Height > maxHeight {
				maxHeight = size.Height
			}
		}
	}

	if !l.Horizontal {
		totalWidth = maxWidth
	} else {
		totalHeight = maxHeight
	}

	return Size{
		Width:  totalWidth + padding.Left + padding.Right,
		Height: totalHeight + padding.Top + padding.Bottom,
	}
}
