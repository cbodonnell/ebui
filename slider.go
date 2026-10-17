package ebui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var _ FocusableComponent = &Slider{}

// SliderColors represents the color scheme for a slider
type SliderColors struct {
	Track        color.Color
	TrackFilled  color.Color
	Thumb        color.Color
	ThumbHovered color.Color
	ThumbDragged color.Color
	FocusBorder  color.Color
}

// DefaultSliderColors returns a default color scheme for sliders
func DefaultSliderColors() SliderColors {
	return SliderColors{
		Track:        color.RGBA{200, 200, 200, 255},
		TrackFilled:  color.RGBA{100, 149, 237, 255}, // Cornflower blue
		Thumb:        color.RGBA{255, 255, 255, 255},
		ThumbHovered: color.RGBA{230, 230, 230, 255},
		ThumbDragged: color.RGBA{220, 220, 220, 255},
		FocusBorder:  color.Black,
	}
}

// Slider is a component that allows selecting a value within a range
type Slider struct {
	*BaseFocusable
	*LayoutContainer
	min         float64
	max         float64
	value       float64
	stepSize    float64
	thumbWidth  float64
	thumbHeight float64
	trackHeight float64
	isDragging  bool
	isHovered   bool
	colors      SliderColors
	onChange    func(value float64)
	valueLabel  *Label
	showValue   bool
	valueSuffix string
}

// SliderOpt is a function that configures a Slider
type SliderOpt func(s *Slider)

// WithMinValue sets the minimum value for the slider
func WithMinValue(min float64) ComponentOpt {
	return func(c Component) {
		if s, ok := c.(*Slider); ok {
			s.min = min
		}
	}
}

// WithMaxValue sets the maximum value for the slider
func WithMaxValue(max float64) ComponentOpt {
	return func(c Component) {
		if s, ok := c.(*Slider); ok {
			s.max = max
		}
	}
}

// WithValue sets the initial value for the slider
func WithValue(value float64) ComponentOpt {
	return func(c Component) {
		if s, ok := c.(*Slider); ok {
			s.SetValue(value)
		}
	}
}

// WithStepSize sets the step size for the slider
func WithStepSize(step float64) ComponentOpt {
	return func(c Component) {
		if s, ok := c.(*Slider); ok {
			s.stepSize = step
		}
	}
}

// WithTrackHeight sets the height of the slider track
func WithTrackHeight(height float64) ComponentOpt {
	return func(c Component) {
		if s, ok := c.(*Slider); ok {
			s.trackHeight = height
		}
	}
}

// WithThumbSize sets the size of the slider thumb
func WithThumbSize(width, height float64) ComponentOpt {
	return func(c Component) {
		if s, ok := c.(*Slider); ok {
			s.thumbWidth = width
			s.thumbHeight = height
		}
	}
}

// WithSliderColors sets the colors for the slider
func WithSliderColors(colors SliderColors) ComponentOpt {
	return func(c Component) {
		if s, ok := c.(*Slider); ok {
			s.colors = colors
		}
	}
}

// WithOnChangeHandler sets the handler for value changes
func WithOnChangeHandler(handler func(value float64)) ComponentOpt {
	return func(c Component) {
		if s, ok := c.(*Slider); ok {
			s.onChange = handler
		}
	}
}

// WithShowValue makes the slider display its current value
func WithShowValue() ComponentOpt {
	return func(c Component) {
		if s, ok := c.(*Slider); ok {
			s.showValue = true
		}
	}
}

// WithValueSuffix sets a suffix for the displayed value (e.g., "%", "px")
func WithValueSuffix(suffix string) ComponentOpt {
	return func(c Component) {
		if s, ok := c.(*Slider); ok {
			s.valueSuffix = suffix
		}
	}
}

// NewSlider creates a new slider component
func NewSlider(opts ...ComponentOpt) *Slider {
	// Default layout is a horizontal container
	withLayout := WithLayout(NewHorizontalStackLayout(10, AlignCenter))
	s := &Slider{
		LayoutContainer: NewLayoutContainer(
			append([]ComponentOpt{withLayout}, opts...)...,
		),
		BaseFocusable: NewBaseFocusable(),
		min:           0,
		max:           100,
		value:         50,
		stepSize:      1,
		thumbWidth:    20,
		thumbHeight:   20,
		trackHeight:   6,
		colors:        DefaultSliderColors(),
		onChange:      func(value float64) {},
		showValue:     false,
		valueSuffix:   "",
	}

	// Value label (optional, shown if WithShowValue is used)
	s.valueLabel = NewLabel(
		s.getValueText(),
		WithSize(60, s.size.Height),
		WithJustify(JustifyCenter),
	)

	// Process options
	for _, opt := range opts {
		opt(s)
	}

	// We'll handle the value label differently - we won't add it as a child
	// to avoid it being drawn behind the slider

	s.registerEventListeners()

	return s
}

func (s *Slider) registerEventListeners() {
	s.AddEventListener(MouseEnter, func(e *Event) {
		s.isHovered = true
	})

	s.AddEventListener(MouseLeave, func(e *Event) {
		s.isHovered = false
	})

	s.AddEventListener(DragStart, func(e *Event) {
		// Check if user clicked on the track or thumb
		if s.isPointOverThumb(e.MouseX, e.MouseY) || s.isPointOverTrack(e.MouseX, e.MouseY) {
			s.isDragging = true
			// Update slider value based on click position
			s.updateValueFromPosition(e.MouseX)
		}
	})

	s.AddEventListener(DragEnd, func(e *Event) {
		s.isDragging = false
	})

	s.AddEventListener(Drag, func(e *Event) {
		if s.isDragging {
			s.updateValueFromPosition(e.MouseX)
		}
	})

	s.AddEventListener(Nav, func(e *Event) {
		if e.Phase != PhaseTarget {
			return
		}
		switch e.Nav {
		case NavLeft:
			s.decrementValue()
			e.Handled = true
		case NavRight:
			s.incrementValue()
			e.Handled = true
		}
	})
}

func (s *Slider) Update() error {
	s.handleInput()
	return s.LayoutContainer.Update()
}

func (s *Slider) handleInput() {
	if !s.IsFocused() {
		return
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyHome) {
		s.SetValue(s.min)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnd) {
		s.SetValue(s.max)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		s.SetValue(s.value - (s.max-s.min)/10)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		s.SetValue(s.value + (s.max-s.min)/10)
	}
}

func (s *Slider) Draw(screen *ebiten.Image) {
	if s.IsFocusVisible() {
		paint(screen, Paint{Part: PartFocusRing, Role: s.GetRole(), State: StateFocused, Rect: componentRect(s), Border: s.colors.FocusBorder})
	}

	// Let the container handle drawing of any children (like the value label)
	s.LayoutContainer.Draw(screen)

	// Draw the slider track and thumb
	s.drawTrack(screen)
	s.drawThumb(screen)

	// Draw the value label if needed
	if s.showValue {
		s.drawValueLabel(screen)
	}
}

// state is the slider's interaction state, for the painter.
func (s *Slider) state() State {
	var st State
	if s.isHovered {
		st |= StateHovered
	}
	if s.isDragging {
		st |= StatePressed
	}
	if s.IsFocusVisible() {
		st |= StateFocused
	}
	if s.IsDisabled() {
		st |= StateDisabled
	}
	return st
}

func (s *Slider) drawTrack(screen *ebiten.Image) {
	pos := s.GetAbsolutePosition()
	size := s.GetSize()
	trackY := pos.Y + (size.Height-s.trackHeight)/2
	track := rectAt(s.getTrackX(), trackY, s.getTrackWidth(), s.trackHeight)
	paint(screen, Paint{Part: PartSliderTrack, Role: s.GetRole(), State: s.state(), Rect: track, Fill: s.colors.Track})

	// Filled track (from left to thumb)
	filled := track
	filled.Max.X = int(s.getThumbPosition())
	if filled.Dx() > 0 {
		paint(screen, Paint{Part: PartSliderFill, Role: s.GetRole(), State: s.state(), Rect: filled, Fill: s.colors.TrackFilled})
	}
}

func (s *Slider) drawThumb(screen *ebiten.Image) {
	thumbColor := s.colors.Thumb
	if s.isDragging {
		thumbColor = s.colors.ThumbDragged
	} else if s.isHovered {
		thumbColor = s.colors.ThumbHovered
	}

	pos := s.GetAbsolutePosition()
	size := s.GetSize()
	thumbX := s.getThumbPosition() - s.thumbWidth/2
	thumbY := pos.Y + (size.Height-s.thumbHeight)/2
	thumb := rectAt(thumbX, thumbY, s.thumbWidth, s.thumbHeight)
	paint(screen, Paint{Part: PartSliderThumb, Role: s.GetRole(), State: s.state(), Rect: thumb, Fill: thumbColor})
}

// drawValueLabel draws the value label on the right side of the slider
func (s *Slider) drawValueLabel(screen *ebiten.Image) {
	pos := s.GetAbsolutePosition()
	size := s.GetSize()

	// Draw value text directly instead of using a Label component
	// Position it on the right side of the slider
	valueX := pos.X + size.Width - 45
	valueY := pos.Y + size.Height/2 + 5 // Center vertically

	text.Draw(
		screen,
		s.getValueText(),
		s.valueLabel.font,
		int(valueX),
		int(valueY),
		s.valueLabel.GetColor(),
	)
}

func (s *Slider) getTrackX() float64 {
	pos := s.GetAbsolutePosition()
	return pos.X + s.thumbWidth/2
}

func (s *Slider) getTrackWidth() float64 {
	size := s.GetSize()
	// We'll handle the value label separately, so no need to adjust the track width
	return size.Width - s.thumbWidth
}

func (s *Slider) getThumbPosition() float64 {
	// Calculate the x position of the thumb center based on current value
	valueRatio := (s.value - s.min) / (s.max - s.min)
	trackWidth := s.getTrackWidth()
	return s.getTrackX() + valueRatio*trackWidth
}

func (s *Slider) updateValueFromPosition(mouseX float64) {
	// Convert mouse position to slider value
	trackX := s.getTrackX()
	trackWidth := s.getTrackWidth()

	// Calculate relative position (0.0 to 1.0)
	relativePos := clamp((mouseX-trackX)/trackWidth, 0, 1)

	// Convert to value range and apply step size
	newValue := s.min + relativePos*(s.max-s.min)

	// Apply step size if set
	if s.stepSize > 0 {
		newValue = math.Round(newValue/s.stepSize) * s.stepSize
	}

	// Set the new value (this will call onChange if value changed)
	s.SetValue(newValue)
}

func (s *Slider) isPointOverThumb(x, y float64) bool {
	pos := s.GetAbsolutePosition()
	size := s.GetSize()

	thumbX := s.getThumbPosition() - s.thumbWidth/2
	thumbY := pos.Y + (size.Height-s.thumbHeight)/2

	return x >= thumbX && x <= thumbX+s.thumbWidth &&
		y >= thumbY && y <= thumbY+s.thumbHeight
}

func (s *Slider) isPointOverTrack(x, y float64) bool {
	pos := s.GetAbsolutePosition()
	size := s.GetSize()

	trackX := s.getTrackX()
	trackY := pos.Y + (size.Height-s.trackHeight)/2
	trackWidth := s.getTrackWidth()

	return x >= trackX && x <= trackX+trackWidth &&
		y >= trackY && y <= trackY+s.trackHeight
}

func (s *Slider) incrementValue() {
	s.SetValue(s.value + s.stepSize)
}

func (s *Slider) decrementValue() {
	s.SetValue(s.value - s.stepSize)
}

func (s *Slider) getValueText() string {
	// Format value text with suffix if provided
	if s.valueSuffix != "" {
		return fmt.Sprintf("%.0f%s", s.value, s.valueSuffix)
	}
	return fmt.Sprintf("%.0f", s.value)
}

// SetValue sets the slider's value, respecting min/max bounds
func (s *Slider) SetValue(value float64) {
	// Clamp value to min/max range
	newValue := clamp(value, s.min, s.max)

	// Always update the value (needed for reset button functionality)
	s.value = newValue

	// Update the value label if visible
	if s.showValue {
		s.valueLabel.SetText(s.getValueText())
	}

	// Call onChange handler
	if s.onChange != nil {
		s.onChange(s.value)
	}
}

// GetValue returns the current slider value
func (s *Slider) GetValue() float64 {
	return s.value
}

// SetColors sets the color scheme for the slider
func (s *Slider) SetColors(colors SliderColors) {
	s.colors = colors
}

// SetOnChange sets the handler for value changes
func (s *Slider) SetOnChange(handler func(value float64)) {
	s.onChange = handler
}
