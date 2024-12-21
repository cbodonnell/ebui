package ebui

import (
	"image/color"
	"time"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.design/x/clipboard"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

var _ InteractiveComponent = &TextInput{}

func init() {
	err := clipboard.Init()
	if err != nil {
		// Note: In a real application, you might want to handle this error differently
		// For now, we'll just print it and continue without clipboard support
		println("Warning: Failed to initialize clipboard:", err)
	}
}

// TextInput represents a single-line text input field
type TextInput struct {
	*BaseInteractive
	*LayoutContainer
	text            []rune
	cursorPos       int
	selectionStart  int
	selectionEnd    int
	font            font.Face
	textColor       color.Color
	backgroundColor color.Color
	cursorColor     color.Color
	selectionColor  color.Color
	isFocused       bool
	lastBlink       time.Time
	showCursor      bool
	repeatKey       ebiten.Key
	repeatStart     time.Time
	lastRepeat      time.Time
	onChange        func(string)
	onSubmit        func(string)
}

// TextInputColors holds the color scheme for a text input
type TextInputColors struct {
	Text       color.Color
	Background color.Color
	Cursor     color.Color
	Selection  color.Color
}

// DefaultTextInputColors returns the default color scheme
func DefaultTextInputColors() TextInputColors {
	return TextInputColors{
		Text:       color.Black,
		Background: color.White,
		Cursor:     color.Black,
		Selection:  color.RGBA{100, 149, 237, 127}, // Semi-transparent cornflower blue
	}
}

func WithTextInputColors(colors TextInputColors) ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*TextInput); ok {
			t.textColor = colors.Text
			t.backgroundColor = colors.Background
			t.cursorColor = colors.Cursor
			t.selectionColor = colors.Selection
		}
	}
}

func WithInitialText(text string) ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*TextInput); ok {
			t.SetText(text)
		}
	}
}

func WithOnChange(handler func(string)) ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*TextInput); ok {
			t.onChange = handler
		}
	}
}

func WithOnSubmit(handler func(string)) ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*TextInput); ok {
			t.onSubmit = handler
		}
	}
}

// NewTextInput creates a new text input component
func NewTextInput(opts ...ComponentOpt) *TextInput {
	colors := DefaultTextInputColors()
	t := &TextInput{
		BaseInteractive: NewBaseInteractive(),
		LayoutContainer: NewLayoutContainer(opts...),
		text:            make([]rune, 0),
		cursorPos:       0,
		selectionStart:  -1,
		selectionEnd:    -1,
		repeatKey:       -1,
		font:            basicfont.Face7x13,
		textColor:       colors.Text,
		backgroundColor: colors.Background,
		cursorColor:     colors.Cursor,
		selectionColor:  colors.Selection,
		lastBlink:       time.Now(),
		onChange:        func(string) {},
		onSubmit:        func(string) {},
	}

	for _, opt := range opts {
		opt(t)
	}

	t.registerEventListeners()
	return t
}

func (t *TextInput) registerEventListeners() {
	t.eventDispatcher.AddEventListener(MouseDown, func(e *Event) {
		t.isFocused = true
		t.showCursor = true
		t.lastBlink = time.Now()

		// Calculate cursor position from click
		clickX := e.MouseX - t.GetAbsolutePosition().X
		t.cursorPos = t.getCharIndexAtX(clickX)
		t.selectionStart = t.cursorPos
		t.selectionEnd = t.cursorPos
	})

	t.eventDispatcher.AddEventListener(MouseUp, func(e *Event) {
		if t.isFocused {
			clickX := e.MouseX - t.GetAbsolutePosition().X
			endPos := t.getCharIndexAtX(clickX)

			if endPos != t.selectionStart {
				t.selectionEnd = endPos
				t.cursorPos = endPos
			}
		}
	})

	t.eventDispatcher.AddEventListener(Drag, func(e *Event) {
		if t.isFocused {
			clickX := e.MouseX - t.GetAbsolutePosition().X
			t.selectionEnd = t.getCharIndexAtX(clickX)
			t.cursorPos = t.selectionEnd
		}
	})
}

func (t *TextInput) Update() error {
	// Handle keyboard input if focused
	if t.isFocused {
		t.handleKeyboardInput()
	}

	// Update cursor blink
	if t.isFocused {
		if time.Since(t.lastBlink) > 530*time.Millisecond {
			t.showCursor = !t.showCursor
			t.lastBlink = time.Now()
		}
	}

	return t.LayoutContainer.Update()
}

func (t *TextInput) Draw(screen *ebiten.Image) {
	t.LayoutContainer.Draw(screen)

	pos := t.GetAbsolutePosition()
	size := t.GetSize()
	padding := t.GetPadding()

	// Draw background
	bg := ebiten.NewImage(int(size.Width), int(size.Height))
	bg.Fill(t.backgroundColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X, pos.Y)
	screen.DrawImage(bg, op)

	// Draw selection if exists
	if t.hasSelection() {
		t.drawSelection(screen)
	}

	// Draw text
	if len(t.text) > 0 {
		text.Draw(
			screen,
			string(t.text),
			t.font,
			int(pos.X+padding.Left),
			int(pos.Y+padding.Top+t.getTextBaseline()),
			t.textColor,
		)
	}

	// Draw cursor
	if t.isFocused && t.showCursor {
		t.drawCursor(screen)
	}
}

func (t *TextInput) handleKeyboardInput() {
	// Handle character input
	t.handleCharacterInput()

	// Handle special keys
	t.handleSpecialKeys()

	// Handle key repeating
	t.handleKeyRepeat()
}

func (t *TextInput) handleCharacterInput() {
	// Get input string from Ebitengine
	inputChars := ebiten.AppendInputChars(nil)
	if len(inputChars) > 0 {
		if t.hasSelection() {
			t.deleteSelection()
		}

		for _, ch := range inputChars {
			if unicode.IsPrint(ch) {
				// Insert character at cursor position
				newText := make([]rune, len(t.text)+1)
				copy(newText, t.text[:t.cursorPos])
				newText[t.cursorPos] = ch
				copy(newText[t.cursorPos+1:], t.text[t.cursorPos:])
				t.text = newText
				t.cursorPos++

				if t.onChange != nil {
					t.onChange(string(t.text))
				}
			}
		}
	}
}

func (t *TextInput) handleSpecialKeys() {
	// Handle keyboard shortcuts first
	ctrlPressed := ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyMeta)
	shiftPressed := ebiten.IsKeyPressed(ebiten.KeyShift)
	if ctrlPressed {
		// Handle Ctrl+A (Select All)
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			if t.repeatKey != ebiten.KeyA {
				t.selectAll()
				t.repeatKey = ebiten.KeyA
				t.repeatStart = time.Now()
				t.lastRepeat = time.Now()
			}
			return
		} else if t.repeatKey == ebiten.KeyA {
			t.repeatKey = -1
		}

		// Handle Ctrl+X (Cut)
		if ebiten.IsKeyPressed(ebiten.KeyX) {
			if t.repeatKey != ebiten.KeyX {
				t.handleCut()
				t.repeatKey = ebiten.KeyX
				t.repeatStart = time.Now()
				t.lastRepeat = time.Now()
			}
			return
		} else if t.repeatKey == ebiten.KeyX {
			t.repeatKey = -1
		}

		// Handle Ctrl+C (Copy)
		if ebiten.IsKeyPressed(ebiten.KeyC) {
			if t.repeatKey != ebiten.KeyC {
				t.handleCopy()
				t.repeatKey = ebiten.KeyC
				t.repeatStart = time.Now()
				t.lastRepeat = time.Now()
			}
			return
		} else if t.repeatKey == ebiten.KeyC {
			t.repeatKey = -1
		}

		// Handle Ctrl+V (Paste)
		if ebiten.IsKeyPressed(ebiten.KeyV) {
			if t.repeatKey != ebiten.KeyV {
				t.handlePaste()
				t.repeatKey = ebiten.KeyV
				t.repeatStart = time.Now()
				t.lastRepeat = time.Now()
			}
			return
		} else if t.repeatKey == ebiten.KeyV {
			t.repeatKey = -1
		}
	}

	// Define the keys we want to handle
	keys := []ebiten.Key{
		ebiten.KeyLeft,
		ebiten.KeyRight,
		ebiten.KeyBackspace,
		ebiten.KeyDelete,
		ebiten.KeyEnter,
		ebiten.KeyHome,
		ebiten.KeyEnd,
	}

	for _, key := range keys {
		if ebiten.IsKeyPressed(key) {
			// Start key repeat when key is first pressed
			if t.repeatKey != key {
				t.repeatKey = key
				t.repeatStart = time.Now()
				t.lastRepeat = time.Now()
				t.handleKey(key, ctrlPressed, shiftPressed)
			}
		} else if t.repeatKey == key {
			// Key was released
			t.repeatKey = -1
		}
	}
}

func (t *TextInput) handleKeyRepeat() {
	if t.repeatKey == -1 {
		return
	}

	now := time.Now()
	// Initial delay of 500ms, then repeat every 50ms
	var shouldRepeat bool
	if time.Since(t.repeatStart) < 500*time.Millisecond {
		shouldRepeat = false
	} else if time.Since(t.lastRepeat) >= 50*time.Millisecond {
		shouldRepeat = true
	}

	if !shouldRepeat {
		return
	}

	ctrlPressed := ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyMeta)
	shiftPressed := ebiten.IsKeyPressed(ebiten.KeyShift)

	// Handle ctrl+key combinations
	if ctrlPressed {
		switch t.repeatKey {
		case ebiten.KeyV:
			t.handlePaste()
		case ebiten.KeyC:
			t.handleCopy()
		case ebiten.KeyX:
			t.handleCut()
		case ebiten.KeyA:
			t.selectAll()
		default:
			t.handleKey(t.repeatKey, ctrlPressed, shiftPressed)
		}
	} else {
		t.handleKey(t.repeatKey, ctrlPressed, shiftPressed)
	}

	t.lastRepeat = now
}

func (t *TextInput) handleKey(key ebiten.Key, ctrlPressed, shiftPressed bool) {
	switch key {
	case ebiten.KeyLeft:
		prevPos := t.cursorPos
		if ctrlPressed {
			// Move to previous word
			t.cursorPos = t.findPreviousWordBoundary()
		} else if t.cursorPos > 0 {
			t.cursorPos--
		}
		if !shiftPressed {
			t.ClearSelection()
		} else {
			if t.selectionStart == -1 {
				t.selectionStart = prevPos
			}
			t.selectionEnd = t.cursorPos
		}

	case ebiten.KeyRight:
		prevPos := t.cursorPos
		if ctrlPressed {
			// Move to next word
			t.cursorPos = t.findNextWordBoundary()
		} else if t.cursorPos < len(t.text) {
			t.cursorPos++
		}
		if !shiftPressed {
			t.ClearSelection()
		} else {
			if t.selectionStart == -1 {
				t.selectionStart = prevPos
			}
			t.selectionEnd = t.cursorPos
		}

	case ebiten.KeyBackspace:
		if t.hasSelection() {
			t.deleteSelection()
		} else if t.cursorPos > 0 {
			t.text = append(t.text[:t.cursorPos-1], t.text[t.cursorPos:]...)
			t.cursorPos--
			if t.onChange != nil {
				t.onChange(string(t.text))
			}
		}

	case ebiten.KeyDelete:
		if t.hasSelection() {
			t.deleteSelection()
		} else if t.cursorPos < len(t.text) {
			t.text = append(t.text[:t.cursorPos], t.text[t.cursorPos+1:]...)
			if t.onChange != nil {
				t.onChange(string(t.text))
			}
		}

	case ebiten.KeyEnter:
		if t.onSubmit != nil {
			t.onSubmit(string(t.text))
		}

	case ebiten.KeyHome:
		t.cursorPos = 0

	case ebiten.KeyEnd:
		t.cursorPos = len(t.text)
	}

	t.showCursor = true
	t.lastBlink = time.Now()
}

func (t *TextInput) drawSelection(screen *ebiten.Image) {
	if !t.hasSelection() {
		return
	}

	start, end := t.getOrderedSelection()
	startX := t.getXPositionForIndex(start)
	endX := t.getXPositionForIndex(end)

	pos := t.GetAbsolutePosition()
	size := t.GetSize()
	padding := t.GetPadding()

	selectionWidth := endX - startX
	selectionHeight := size.Height - padding.Top - padding.Bottom

	selection := ebiten.NewImage(int(selectionWidth), int(selectionHeight))
	selection.Fill(t.selectionColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		pos.X+padding.Left+startX,
		pos.Y+padding.Top,
	)
	screen.DrawImage(selection, op)
}

func (t *TextInput) drawCursor(screen *ebiten.Image) {
	pos := t.GetAbsolutePosition()
	size := t.GetSize()
	padding := t.GetPadding()

	cursorX := t.getXPositionForIndex(t.cursorPos)

	cursor := ebiten.NewImage(1, int(size.Height-padding.Top-padding.Bottom))
	cursor.Fill(t.cursorColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		pos.X+padding.Left+cursorX,
		pos.Y+padding.Top,
	)
	screen.DrawImage(cursor, op)
}

// Helper methods
func (t *TextInput) getTextBaseline() float64 {
	metrics := t.font.Metrics()
	return (t.GetSize().Height-float64(metrics.Height.Ceil()))/2 + float64(metrics.Ascent.Ceil())
}

func (t *TextInput) getCharIndexAtX(x float64) int {
	for i := 0; i <= len(t.text); i++ {
		charX := t.getXPositionForIndex(i)
		if x < charX {
			return i - 1
		}
	}
	return len(t.text)
}

func (t *TextInput) getXPositionForIndex(index int) float64 {
	if index <= 0 {
		return 0
	}

	if index > len(t.text) {
		index = len(t.text)
	}

	return float64(font.MeasureString(t.font, string(t.text[:index])).Ceil())
}

func (t *TextInput) hasSelection() bool {
	return t.selectionStart != -1 && t.selectionEnd != -1 && t.selectionStart != t.selectionEnd
}

func (t *TextInput) getOrderedSelection() (int, int) {
	if t.selectionStart < t.selectionEnd {
		return t.selectionStart, t.selectionEnd
	}
	return t.selectionEnd, t.selectionStart
}

func (t *TextInput) findPreviousWordBoundary() int {
	pos := t.cursorPos
	// Skip spaces before cursor
	for pos > 0 && unicode.IsSpace(t.text[pos-1]) {
		pos--
	}
	// Skip word characters
	for pos > 0 && !unicode.IsSpace(t.text[pos-1]) {
		pos--
	}
	return pos
}

func (t *TextInput) findNextWordBoundary() int {
	pos := t.cursorPos
	// Skip current word
	for pos < len(t.text) && !unicode.IsSpace(t.text[pos]) {
		pos++
	}
	// Skip spaces
	for pos < len(t.text) && unicode.IsSpace(t.text[pos]) {
		pos++
	}
	return pos
}

func (t *TextInput) deleteSelection() {
	if !t.hasSelection() {
		return
	}

	start, end := t.getOrderedSelection()
	t.text = append(t.text[:start], t.text[end:]...)
	t.cursorPos = start
	t.selectionStart = -1
	t.selectionEnd = -1

	if t.onChange != nil {
		t.onChange(string(t.text))
	}
}

// Public API methods

// SetText sets the text content of the input field
func (t *TextInput) SetText(text string) {
	t.text = []rune(text)
	t.cursorPos = len(t.text)
	t.selectionStart = -1
	t.selectionEnd = -1
	if t.onChange != nil {
		t.onChange(text)
	}
}

// GetText returns the current text content
func (t *TextInput) GetText() string {
	return string(t.text)
}

// Focus sets the input field as focused
func (t *TextInput) Focus() {
	t.isFocused = true
	t.showCursor = true
	t.lastBlink = time.Now()
}

// Blur removes focus from the input field
func (t *TextInput) Blur() {
	t.isFocused = false
	t.showCursor = false
	t.selectionStart = -1
	t.selectionEnd = -1
}

// IsFocused returns whether the input field is currently focused
func (t *TextInput) IsFocused() bool {
	return t.isFocused
}

// SetFont sets the font used for rendering text
func (t *TextInput) SetFont(font font.Face) {
	t.font = font
}

// selectAll selects all text in the input
func (t *TextInput) selectAll() {
	if len(t.text) > 0 {
		t.selectionStart = 0
		t.selectionEnd = len(t.text)
		t.cursorPos = t.selectionEnd
		t.showCursor = true
		t.lastBlink = time.Now()
	}
}

// handleCopy copies selected text to clipboard
func (t *TextInput) handleCopy() {
	if !t.hasSelection() {
		return
	}
	start, end := t.getOrderedSelection()
	text := string(t.text[start:end])
	clipboard.Write(clipboard.FmtText, []byte(text))
}

// handleCut cuts selected text to clipboard
func (t *TextInput) handleCut() {
	if !t.hasSelection() {
		return
	}
	t.handleCopy()
	t.deleteSelection()
}

// handlePaste pastes text from clipboard
func (t *TextInput) handlePaste() {
	bytes := clipboard.Read(clipboard.FmtText)
	if len(bytes) == 0 {
		return
	}
	text := string(bytes)

	if t.hasSelection() {
		t.deleteSelection()
	}

	// Filter out any newlines or tabs from the pasted text
	filtered := []rune{}
	for _, ch := range text {
		if ch != '\n' && ch != '\r' && ch != '\t' {
			filtered = append(filtered, ch)
		}
	}

	// Insert the filtered text at cursor position
	newText := make([]rune, 0, len(t.text)+len(filtered))
	newText = append(newText, t.text[:t.cursorPos]...)
	newText = append(newText, filtered...)
	newText = append(newText, t.text[t.cursorPos:]...)
	t.text = newText
	t.cursorPos += len(filtered)

	if t.onChange != nil {
		t.onChange(string(t.text))
	}
}

// Select selects a range of text
func (t *TextInput) Select(start, end int) {
	if start < 0 {
		start = 0
	}
	if end > len(t.text) {
		end = len(t.text)
	}
	if start > end {
		start, end = end, start
	}

	t.selectionStart = start
	t.selectionEnd = end
	t.cursorPos = end
	t.showCursor = true
	t.lastBlink = time.Now()
}

// ClearSelection clears the current selection
func (t *TextInput) ClearSelection() {
	t.selectionStart = -1
	t.selectionEnd = -1
}
