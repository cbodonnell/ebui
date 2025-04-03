package ebui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// TooltipPosition defines the preferred position of a tooltip relative to the mouse cursor
type TooltipPosition int

const (
	// Tooltip positions
	TooltipPositionAuto        TooltipPosition = iota // Automatically determine best position
	TooltipPositionTop                                // Above the cursor
	TooltipPositionRight                              // To the right of the cursor
	TooltipPositionBottom                             // Below the cursor
	TooltipPositionLeft                               // To the left of the cursor
	TooltipPositionTopLeft                            // Above and to the left of the cursor
	TooltipPositionTopRight                           // Above and to the right of the cursor
	TooltipPositionBottomLeft                         // Below and to the left of the cursor
	TooltipPositionBottomRight                        // Below and to the right of the cursor
)

// Tooltipable is an interface for components that can have tooltips attached
type Tooltipable interface {
	SetTooltip(tooltip *Tooltip)
	GetTooltip() *Tooltip
	ClearTooltip()
}

type TooltipableComponent interface {
	InteractiveComponent
	Tooltipable
}

var _ Tooltipable = &BaseTooltipable{}

// BaseTooltipable is a base struct that implements the Tooltipable interface
type BaseTooltipable struct {
	tooltip *Tooltip
}

// NewBaseTooltipable creates a new base for tooltipable components
func NewBaseTooltipable() *BaseTooltipable {
	return &BaseTooltipable{}
}

// SetTooltip sets the tooltip for the component
func (b *BaseTooltipable) SetTooltip(tooltip *Tooltip) {
	b.tooltip = tooltip
}

// GetTooltip returns the tooltip for the component
func (b *BaseTooltipable) GetTooltip() *Tooltip {
	return b.tooltip
}

// ClearTooltip clears the tooltip for the component
func (b *BaseTooltipable) ClearTooltip() {
	b.tooltip = nil
}

// Tooltip represents a tooltip component that can be attached to other components
type Tooltip struct {
	*LayoutContainer
	backgroundColor color.Color
	borderColor     color.Color
	borderWidth     float64
	position        TooltipPosition
	offsetX         float64   // Additional horizontal offset
	offsetY         float64   // Additional vertical offset
	mouseX          float64   // Current mouse X position
	mouseY          float64   // Current mouse Y position
	target          Component // Target component this tooltip belongs to
}

// TooltipColors contains all the colors used by a tooltip
type TooltipColors struct {
	Background color.Color
	Border     color.Color
}

// DefaultTooltipColors returns the default color scheme for tooltips
func DefaultTooltipColors() TooltipColors {
	return TooltipColors{
		Background: color.RGBA{250, 250, 250, 240},
		Border:     color.RGBA{200, 200, 200, 255},
	}
}

func WithBackgroundColor(col color.Color) ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*Tooltip); ok {
			t.backgroundColor = col
		}
	}
}

func WithBorderColor(col color.Color) ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*Tooltip); ok {
			t.borderColor = col
		}
	}
}

func WithTooltipColors(colors TooltipColors) ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*Tooltip); ok {
			t.backgroundColor = colors.Background
			t.borderColor = colors.Border
		}
	}
}

func WithTooltipPosition(position TooltipPosition) ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*Tooltip); ok {
			t.position = position
		}
	}
}

func WithTooltipOffset(x, y float64) ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*Tooltip); ok {
			t.offsetX = x
			t.offsetY = y
		}
	}
}

// NewTooltip creates a new tooltip component
func NewTooltip(opts ...ComponentOpt) *Tooltip {
	colors := DefaultTooltipColors()

	tooltip := &Tooltip{
		LayoutContainer: NewLayoutContainer(
			WithLayout(NewVerticalStackLayout(0, AlignStart)),
		),
		backgroundColor: colors.Background,
		borderColor:     colors.Border,
		borderWidth:     1,
		position:        TooltipPositionTopRight, // Default position - top right of cursor
		offsetX:         10,                      // Default horizontal offset
		offsetY:         10,                      // Default vertical offset
	}

	// Process options
	for _, opt := range opts {
		opt(tooltip)
	}

	// Apply background
	tooltip.SetBackground(tooltip.backgroundColor)

	return tooltip
}

// Draw renders the tooltip
func (t *Tooltip) Draw(screen *ebiten.Image) {
	if t.IsHidden() {
		return
	}

	pos := t.GetAbsolutePosition()
	size := t.GetSize()

	// Draw border if border width > 0
	if t.borderWidth > 0 {
		border := GetCache().BorderImageWithColor(
			int(size.Width),
			int(size.Height),
			t.borderColor,
		)
		borderOp := &ebiten.DrawImageOptions{}
		borderOp.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(border, borderOp)
	}

	// Draw the component
	t.LayoutContainer.Draw(screen)
}

// SetContent sets the content of the tooltip
func (t *Tooltip) SetContent(content Component) {
	// Clear existing children
	t.ClearChildren()

	// Add new content
	t.AddChild(content)

	// Update sizing
	t.updateSize()
}

// SetTarget sets the target component this tooltip belongs to
func (t *Tooltip) SetTarget(target Component) {
	t.target = target
}

// GetTarget returns the target component this tooltip belongs to
func (t *Tooltip) GetTarget() Component {
	return t.target
}

// updateSize updates the tooltip size based on its content
func (t *Tooltip) updateSize() {
	if len(t.GetChildren()) == 0 {
		return
	}

	// Get the size from the layout
	minSize := t.layout.GetMinSize(t)

	t.SetSize(Size{
		Width:  minSize.Width,
		Height: minSize.Height,
	})
}

// UpdateMousePosition updates the current mouse position for the tooltip
func (t *Tooltip) UpdateMousePosition(x, y float64) {
	t.mouseX = x
	t.mouseY = y
}

// GetTooltipPosition returns the tooltip's preferred position
func (t *Tooltip) GetTooltipPosition() TooltipPosition {
	return t.position
}

// GetOffset returns the tooltip's offset
func (t *Tooltip) GetOffset() (float64, float64) {
	return t.offsetX, t.offsetY
}

// GetMousePosition returns the current mouse position
func (t *Tooltip) GetMousePosition() (float64, float64) {
	return t.mouseX, t.mouseY
}

// TooltipManager handles the display and positioning of tooltips
type TooltipManager struct {
	*ZIndexedContainer
	activeTooltip   *Tooltip
	lastHoverTarget Component
	disabled        bool
	currentMouseX   float64
	currentMouseY   float64
}

// NewTooltipManager creates a new TooltipManager
func NewTooltipManager(opts ...ComponentOpt) *TooltipManager {
	tm := &TooltipManager{
		ZIndexedContainer: NewZIndexedContainer(opts...),
	}
	return tm
}

// Update updates the tooltip manager state
func (tm *TooltipManager) Update() error {
	// Update tooltip position when mouse moves
	if tm.activeTooltip != nil {
		tm.updateTooltipPosition()
	}
	return tm.ZIndexedContainer.Update()
}

// HandleTargetHover handles when a component with a tooltip is hovered
func (tm *TooltipManager) HandleTargetHover(component Component, mouseX, mouseY float64) {
	// Store current mouse position
	tm.currentMouseX = mouseX
	tm.currentMouseY = mouseY

	// Skip if tooltips are disabled
	if tm.disabled {
		return
	}

	// Check if component is tooltipable
	tooltipable, ok := component.(Tooltipable)
	if !ok {
		return
	}

	// Get tooltip
	tooltip := tooltipable.GetTooltip()
	if tooltip == nil {
		return
	}

	// Update hover state
	tm.lastHoverTarget = component

	// Update mouse position in tooltip
	tooltip.UpdateMousePosition(mouseX, mouseY)

	// Set the target component
	tooltip.SetTarget(component)

	// Show tooltip immediately
	tm.ShowTooltip(tooltip)
}

// HandleTargetLeave handles when the mouse leaves a component with a tooltip
func (tm *TooltipManager) HandleTargetLeave(component Component) {
	// Only process if this is the current hover target
	if tm.lastHoverTarget != component {
		return
	}

	// Clear hover state
	tm.lastHoverTarget = nil

	// Hide tooltip immediately
	tm.HideTooltip()
}

// HandleMouseMove handles mouse movement to update tooltip position
func (tm *TooltipManager) HandleMouseMove(mouseX, mouseY float64) {
	tm.currentMouseX = mouseX
	tm.currentMouseY = mouseY

	// Update tooltip position if active
	if tm.activeTooltip != nil {
		tm.activeTooltip.UpdateMousePosition(mouseX, mouseY)
		tm.updateTooltipPosition()
	}
}

// ShowTooltip displays a tooltip at the appropriate position
func (tm *TooltipManager) ShowTooltip(tooltip *Tooltip) {
	// Hide any existing tooltip
	tm.HideTooltip()

	// Position the tooltip based on preference
	tm.positionTooltip(tooltip)

	// Show the tooltip
	tooltip.Show()
	tm.AddChild(tooltip)
	tm.activeTooltip = tooltip
}

// HideTooltip hides the currently visible tooltip
func (tm *TooltipManager) HideTooltip() {
	if tm.activeTooltip != nil {
		tm.RemoveChild(tm.activeTooltip)
		tm.activeTooltip = nil
	}
}

// positionTooltip calculates and sets the position for a tooltip relative to cursor
func (tm *TooltipManager) positionTooltip(tooltip *Tooltip) {
	tooltipSize := tooltip.GetSize()
	parentSize := tm.GetSize()
	mouseX, mouseY := tooltip.GetMousePosition()
	offsetX, offsetY := tooltip.GetOffset()
	var posX, posY float64

	// Calculate position based on preference
	switch tooltip.GetTooltipPosition() {
	case TooltipPositionTop:
		posX = mouseX - tooltipSize.Width/2
		posY = mouseY - tooltipSize.Height - offsetY

	case TooltipPositionRight:
		posX = mouseX + offsetX
		posY = mouseY - tooltipSize.Height/2

	case TooltipPositionBottom:
		posX = mouseX - tooltipSize.Width/2
		posY = mouseY + offsetY

	case TooltipPositionLeft:
		posX = mouseX - tooltipSize.Width - offsetX
		posY = mouseY - tooltipSize.Height/2

	case TooltipPositionTopLeft:
		posX = mouseX - tooltipSize.Width - offsetX
		posY = mouseY - tooltipSize.Height - offsetY

	case TooltipPositionTopRight: // Default
		posX = mouseX + offsetX
		posY = mouseY - tooltipSize.Height - offsetY

	case TooltipPositionBottomLeft:
		posX = mouseX - tooltipSize.Width - offsetX
		posY = mouseY + offsetY

	case TooltipPositionBottomRight:
		posX = mouseX + offsetX
		posY = mouseY + offsetY

	case TooltipPositionAuto:
		// Start with top-right position
		posX = mouseX + offsetX
		posY = mouseY - tooltipSize.Height - offsetY
	}

	// Apply smart positioning to keep tooltip on screen
	tm.applySmartPositioning(&posX, &posY, mouseX, mouseY, tooltipSize, parentSize, tooltip.GetTooltipPosition())

	// Set the calculated position
	tooltip.SetPosition(Position{
		X:      posX,
		Y:      posY,
		ZIndex: 9999, // Ensure tooltips are on top
	})
}

// applySmartPositioning implements intelligent repositioning to keep tooltips on screen
func (tm *TooltipManager) applySmartPositioning(posX, posY *float64, mouseX, mouseY float64, tooltipSize, parentSize Size, position TooltipPosition) {
	// First, check if the tooltip is out of bounds
	outTop := *posY < 0
	outBottom := *posY+tooltipSize.Height > parentSize.Height
	outLeft := *posX < 0
	outRight := *posX+tooltipSize.Width > parentSize.Width

	// Apply the appropriate adjustment based on current position and out-of-bounds status
	if position == TooltipPositionAuto {
		// For auto position, try all positions until one fits
		// Priority order: top-right, bottom-right, top-left, bottom-left

		// First try: top-right
		*posX = mouseX + 10
		*posY = mouseY - tooltipSize.Height - 10
		if *posX < 0 || *posX+tooltipSize.Width > parentSize.Width ||
			*posY < 0 || *posY+tooltipSize.Height > parentSize.Height {

			// Second try: bottom-right
			*posX = mouseX + 10
			*posY = mouseY + 10
			if *posX < 0 || *posX+tooltipSize.Width > parentSize.Width ||
				*posY < 0 || *posY+tooltipSize.Height > parentSize.Height {

				// Third try: top-left
				*posX = mouseX - tooltipSize.Width - 10
				*posY = mouseY - tooltipSize.Height - 10
				if *posX < 0 || *posX+tooltipSize.Width > parentSize.Width ||
					*posY < 0 || *posY+tooltipSize.Height > parentSize.Height {

					// Last try: bottom-left
					*posX = mouseX - tooltipSize.Width - 10
					*posY = mouseY + 10
				}
			}
		}
	} else {
		// For non-auto positions, handle specific out-of-bounds conditions
		offsetX, offsetY := 10.0, 10.0 // Default offsets

		// Handle vertical out-of-bounds
		if outTop {
			// If tooltip is above and goes out of top, flip to below
			switch position {
			case TooltipPositionTop:
				*posY = mouseY + offsetY
			case TooltipPositionTopLeft:
				*posY = mouseY + offsetY
			case TooltipPositionTopRight:
				*posY = mouseY + offsetY
			}
		}

		if outBottom {
			// If tooltip is below and goes out of bottom, flip to above
			switch position {
			case TooltipPositionBottom:
				*posY = mouseY - tooltipSize.Height - offsetY
			case TooltipPositionBottomLeft:
				*posY = mouseY - tooltipSize.Height - offsetY
			case TooltipPositionBottomRight:
				*posY = mouseY - tooltipSize.Height - offsetY
			}
		}

		// Handle horizontal out-of-bounds
		if outLeft {
			// If tooltip is left and goes out of left, flip to right
			switch position {
			case TooltipPositionLeft:
				*posX = mouseX + offsetX
			case TooltipPositionTopLeft:
				*posX = mouseX + offsetX
			case TooltipPositionBottomLeft:
				*posX = mouseX + offsetX
			}
		}

		if outRight {
			// If tooltip is right and goes out of right, flip to left
			switch position {
			case TooltipPositionRight:
				*posX = mouseX - tooltipSize.Width - offsetX
			case TooltipPositionTopRight:
				*posX = mouseX - tooltipSize.Width - offsetX
			case TooltipPositionBottomRight:
				*posX = mouseX - tooltipSize.Width - offsetX
			}
		}
	}

	// As a final safety measure, ensure the tooltip is completely on screen
	tm.adjustPositionToFitScreen(posX, posY, tooltipSize, parentSize)
}

// adjustPositionToFitScreen adjusts the tooltip position to fit within screen bounds
func (tm *TooltipManager) adjustPositionToFitScreen(posX, posY *float64, tooltipSize, parentSize Size) {
	// Adjust horizontal position if needed
	if *posX < 0 {
		*posX = 0
	} else if *posX+tooltipSize.Width > parentSize.Width {
		*posX = parentSize.Width - tooltipSize.Width
	}

	// Adjust vertical position if needed
	if *posY < 0 {
		*posY = 0
	} else if *posY+tooltipSize.Height > parentSize.Height {
		*posY = parentSize.Height - tooltipSize.Height
	}
}

// updateTooltipPosition updates the position of a tooltip when the mouse moves
func (tm *TooltipManager) updateTooltipPosition() {
	if tm.activeTooltip == nil {
		return
	}

	// Update the tooltip's position based on current mouse position
	tm.activeTooltip.UpdateMousePosition(tm.currentMouseX, tm.currentMouseY)
	tm.positionTooltip(tm.activeTooltip)
}

// RegisterTooltip registers a component to show tooltips
// It adds event listeners for mouse enter, leave, and move events
// to manage tooltip visibility and positioning
func (tm *TooltipManager) RegisterTooltip(component TooltipableComponent, tooltip *Tooltip) {
	component.SetTooltip(tooltip)

	// Add hover and leave event listeners
	component.AddEventListener(MouseEnter, func(e *Event) {
		tm.HandleTargetHover(component, e.MouseX, e.MouseY)
	})

	component.AddEventListener(MouseLeave, func(e *Event) {
		tm.HandleTargetLeave(component)
	})

	// Add mouse move listener to update tooltip position
	component.AddEventListener(MouseMove, func(e *Event) {
		tm.HandleMouseMove(e.MouseX, e.MouseY)
	})
}

// Enable enables the tooltip manager
func (tm *TooltipManager) Enable() {
	tm.disabled = false
}

// Disable disables the tooltip manager and hides any active tooltips
func (tm *TooltipManager) Disable() {
	tm.disabled = true
	if tm.activeTooltip != nil {
		tm.HideTooltip()
	}
}

// IsDisabled returns whether the tooltip manager is disabled
func (tm *TooltipManager) IsDisabled() bool {
	return tm.disabled
}
