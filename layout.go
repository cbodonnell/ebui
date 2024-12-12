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
	Vertical bool
	Config   StackConfig
}

func NewVerticalStack(config StackConfig) *StackLayout {
	return &StackLayout{
		Vertical: true,
		Config:   config,
	}
}

func NewHorizontalStack(config StackConfig) *StackLayout {
	return &StackLayout{
		Vertical: false,
		Config:   config,
	}
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
	var flexibleChildren int
	childSizes := make([]Size, len(children))

	for i, child := range children {
		size := child.GetSize()
		childSizes[i] = size

		if l.Vertical {
			if size.AutoHeight {
				flexibleChildren++
			} else {
				totalFixedMainAxis += size.Height
			}
		} else {
			if size.AutoWidth {
				flexibleChildren++
			} else {
				totalFixedMainAxis += size.Width
			}
		}
	}

	// Calculate total spacing
	totalSpacing := l.Config.Spacing * float64(len(children)-1)

	// Calculate total width/height including spacing
	totalSize := totalFixedMainAxis + totalSpacing

	// Calculate starting position based on alignment
	startX := containerPadding.Left
	startY := containerPadding.Top

	if !l.Vertical {
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
		pos := Position{RelativeToParent: true}

		if l.Vertical {
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
		child.SetSize(size)
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

		if l.Vertical {
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

	if l.Vertical {
		totalWidth = maxWidth
	} else {
		totalHeight = maxHeight
	}

	return Size{
		Width:  totalWidth + padding.Left + padding.Right,
		Height: totalHeight + padding.Top + padding.Bottom,
	}
}
