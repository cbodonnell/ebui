package ebui

import (
	"image"
	"image/color"
	"math"
	"time"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.design/x/clipboard"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

var _ FocusableComponent = &TextInput{}

var clipboardDisabled bool

func init() {
	err := clipboard.Init()
	if err != nil {
		println("Warning: Failed to initialize clipboard:", err)
		clipboardDisabled = true
	}
}

type TextInput struct {
	*BaseFocusable
	*BaseContainer
	text             []rune
	cursorPos        int
	selectionStart   int
	selectionEnd     int
	scrollOffset     float64 // Tracks horizontal scroll position
	font             font.Face
	textColor        color.Color
	backgroundColor  color.Color
	cursorColor      color.Color
	selectionColor   color.Color
	focusBorderColor color.Color
	isFocused        bool
	lastBlink        time.Time
	showCursor       bool
	repeatKey        ebiten.Key
	repeatStart      time.Time
	lastRepeat       time.Time
	onChange         func(string)
	onSubmit         func(string)
	isPassword       bool
	maskChar         rune
	focusable        bool
	tabIndex         int
}

type TextInputColors struct {
	Text        color.Color
	Background  color.Color
	Cursor      color.Color
	Selection   color.Color
	FocusBorder color.Color
}

func DefaultTextInputColors() TextInputColors {
	return TextInputColors{
		Text:        color.Black,
		Background:  color.White,
		Cursor:      color.Black,
		Selection:   color.RGBA{100, 149, 237, 127}, // Dodger Blue
		FocusBorder: color.Black,
	}
}

func WithTextInputColors(colors TextInputColors) ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*TextInput); ok {
			t.textColor = colors.Text
			t.backgroundColor = colors.Background
			t.cursorColor = colors.Cursor
			t.selectionColor = colors.Selection
			t.focusBorderColor = colors.FocusBorder
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

func WithChangeHandler(handler func(string)) ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*TextInput); ok {
			t.onChange = handler
		}
	}
}

func WithSubmitHandler(handler func(string)) ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*TextInput); ok {
			t.onSubmit = handler
		}
	}
}

func WithPasswordMasking() ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*TextInput); ok {
			t.isPassword = true
		}
	}
}

func WithTabIndex(index int) ComponentOpt {
	return func(c Component) {
		if t, ok := c.(*TextInput); ok {
			t.tabIndex = index
		}
	}
}

func NewTextInput(opts ...ComponentOpt) *TextInput {
	colors := DefaultTextInputColors()
	t := &TextInput{
		BaseFocusable:   NewBaseFocusable(),
		BaseContainer:   NewBaseContainer(opts...),
		text:            make([]rune, 0),
		cursorPos:       0,
		selectionStart:  -1,
		selectionEnd:    -1,
		scrollOffset:    0,
		repeatKey:       -1,
		font:            basicfont.Face7x13,
		textColor:       colors.Text,
		backgroundColor: colors.Background,
		cursorColor:     colors.Cursor,
		selectionColor:  colors.Selection,
		lastBlink:       time.Now(),
		onChange:        func(string) {},
		onSubmit:        func(string) {},
		isPassword:      false,
		maskChar:        '*', // Default mask character
		focusable:       true,
		tabIndex:        0,
	}

	for _, opt := range opts {
		opt(t)
	}

	t.registerEventListeners()
	return t
}

func (t *TextInput) registerEventListeners() {
	t.AddEventListener(MouseDown, func(e *Event) {
		t.Focus()

		// Calculate cursor position from click, accounting for scroll
		clickX := e.MouseX - t.GetAbsolutePosition().X + t.scrollOffset
		t.cursorPos = t.getCharIndexAtX(clickX)
		t.selectionStart = t.cursorPos
		t.selectionEnd = t.cursorPos
	})

	t.AddEventListener(MouseUp, func(e *Event) {
		if t.isFocused {
			clickX := e.MouseX - t.GetAbsolutePosition().X + t.scrollOffset
			endPos := t.getCharIndexAtX(clickX)

			if endPos != t.selectionStart {
				t.selectionEnd = endPos
				t.cursorPos = endPos
			}
		}
	})

	t.AddEventListener(Drag, func(e *Event) {
		if t.isFocused {
			clickX := e.MouseX - t.GetAbsolutePosition().X + t.scrollOffset
			t.selectionEnd = t.getCharIndexAtX(clickX)
			t.cursorPos = t.selectionEnd
			t.ensureCursorVisible()
		}
	})

	t.AddEventListener(Focus, func(e *Event) {
		t.Focus()
	})

	t.AddEventListener(Blur, func(e *Event) {
		t.Blur()
	})
}

func (t *TextInput) Update() error {
	if t.isFocused {
		t.handleKeyboardInput()
		if time.Since(t.lastBlink) > 530*time.Millisecond {
			t.showCursor = !t.showCursor
			t.lastBlink = time.Now()
		}
	}
	return t.BaseContainer.Update()
}

func (t *TextInput) handleKeyboardInput() {
	t.handleCharacterInput()
	t.handleSpecialKeys()
	t.handleKeyRepeat()
}

func (t *TextInput) handleCharacterInput() {
	inputChars := ebiten.AppendInputChars(nil)
	if len(inputChars) > 0 {
		if t.hasSelection() {
			t.deleteSelection()
		}

		for _, ch := range inputChars {
			if unicode.IsPrint(ch) {
				newText := make([]rune, len(t.text)+1)
				copy(newText, t.text[:t.cursorPos])
				newText[t.cursorPos] = ch
				copy(newText[t.cursorPos+1:], t.text[t.cursorPos:])
				t.text = newText
				t.cursorPos++
				t.ensureCursorVisible()

				if t.onChange != nil {
					t.onChange(string(t.text))
				}
			}
		}
	}
}

func (t *TextInput) handleSpecialKeys() {
	ctrlPressed := ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyMeta)
	shiftPressed := ebiten.IsKeyPressed(ebiten.KeyShift)

	// Handle keyboard shortcuts only if it's not an arrow key
	if ctrlPressed {
		shortcuts := map[ebiten.Key]func(){
			ebiten.KeyA: t.selectAll,
			ebiten.KeyX: t.handleCut,
			ebiten.KeyC: t.handleCopy,
			ebiten.KeyV: t.handlePaste,
		}

		for key, handler := range shortcuts {
			if ebiten.IsKeyPressed(key) {
				if t.repeatKey != key {
					handler()
					t.repeatKey = key
					t.repeatStart = time.Now()
					t.lastRepeat = time.Now()
				}
				return
			} else if t.repeatKey == key {
				t.repeatKey = -1
			}
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
			if t.repeatKey != key {
				t.repeatKey = key
				t.repeatStart = time.Now()
				t.lastRepeat = time.Now()
				t.handleKey(key, ctrlPressed, shiftPressed)
			}
		} else if t.repeatKey == key {
			t.repeatKey = -1
		}
	}
}

func (t *TextInput) handleKeyRepeat() {
	if t.repeatKey == -1 {
		return
	}

	now := time.Now()
	initialDelay := 500 * time.Millisecond
	repeatDelay := 50 * time.Millisecond

	shouldRepeat := time.Since(t.repeatStart) >= initialDelay &&
		time.Since(t.lastRepeat) >= repeatDelay

	if !shouldRepeat {
		return
	}

	ctrlPressed := ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyMeta)
	shiftPressed := ebiten.IsKeyPressed(ebiten.KeyShift)

	// Handle repeatable shortcuts when ctrl is pressed
	if ctrlPressed {
		switch t.repeatKey {
		case ebiten.KeyV:
			t.handlePaste()
		case ebiten.KeyLeft, ebiten.KeyRight:
			// Allow ctrl+arrow keys to repeat for word-by-word movement
			t.handleKey(t.repeatKey, ctrlPressed, shiftPressed)
		}
	} else {
		// Handle regular key repeats
		t.handleKey(t.repeatKey, ctrlPressed, shiftPressed)
	}

	t.lastRepeat = now
}

func (t *TextInput) handleKey(key ebiten.Key, ctrlPressed, shiftPressed bool) {
	switch key {
	case ebiten.KeyLeft:
		t.handleLeftKey(ctrlPressed, shiftPressed)
	case ebiten.KeyRight:
		t.handleRightKey(ctrlPressed, shiftPressed)
	case ebiten.KeyBackspace:
		t.handleBackspace()
	case ebiten.KeyDelete:
		t.handleDelete()
	case ebiten.KeyEnter:
		if t.onSubmit != nil {
			t.onSubmit(string(t.text))
		}
	case ebiten.KeyHome:
		t.handleHome(shiftPressed)
	case ebiten.KeyEnd:
		t.handleEnd(shiftPressed)
	}

	t.ensureCursorVisible()
	t.showCursor = true
	t.lastBlink = time.Now()
}

func (t *TextInput) Draw(screen *ebiten.Image) {
	pos := t.GetAbsolutePosition()
	size := t.GetSize()
	padding := t.GetPadding()

	if t.isFocused {
		// Draw the focus border 1px
		bg := ebiten.NewImage(int(size.Width+2), int(size.Height+2))
		bg.Fill(t.focusBorderColor)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pos.X-1, pos.Y-1)
		screen.DrawImage(bg, op)
	}

	// Draw background first on the main screen
	bg := ebiten.NewImage(int(size.Width-padding.Left-padding.Right), int(size.Height-padding.Top-padding.Bottom))
	bg.Fill(t.backgroundColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos.X+padding.Left, pos.Y+padding.Top)
	screen.DrawImage(bg, op)

	// Create clip bounds for text content
	clipBounds := image.Rect(
		int(pos.X+padding.Left),
		int(pos.Y+padding.Top),
		int(pos.X+size.Width-padding.Right),
		int(pos.Y+size.Height-padding.Bottom),
	)
	clippedScreen := screen.SubImage(clipBounds).(*ebiten.Image)

	// Draw selection if exists
	if t.hasSelection() {
		t.drawSelection(clippedScreen)
	}

	// Draw text or password mask
	if len(t.text) > 0 {
		displayText := string(t.text)
		if t.isPassword {
			// Create a string of mask characters the same length as the text
			masked := make([]rune, len(t.text))
			for i := range masked {
				masked[i] = t.maskChar
			}
			displayText = string(masked)
		}
		text.Draw(
			clippedScreen,
			displayText,
			t.font,
			int(pos.X+padding.Left-t.scrollOffset),
			int(pos.Y+padding.Top+t.getTextBaseline()),
			t.textColor,
		)
	}

	// Draw cursor on the main screen (not clipped)
	if t.isFocused && t.showCursor {
		t.drawCursor(screen)
	}

	t.BaseContainer.Draw(screen)
}

func (t *TextInput) getTextBaseline() float64 {
	metrics := t.font.Metrics()
	return (t.GetSize().Height-float64(metrics.Height.Ceil()))/2 + float64(metrics.Ascent.Ceil())
}

func (t *TextInput) getCharIndexAtX(x float64) int {
	if x <= 0 {
		return 0
	}

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

	// If in password mode, measure using mask characters
	if t.isPassword {
		masked := make([]rune, index)
		for i := range masked {
			masked[i] = t.maskChar
		}
		return float64(font.MeasureString(t.font, string(masked)).Ceil())
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

func (t *TextInput) drawSelection(screen *ebiten.Image) {
	if !t.hasSelection() {
		return
	}

	start, end := t.getOrderedSelection()
	startX := t.getXPositionForIndex(start) - t.scrollOffset
	endX := t.getXPositionForIndex(end) - t.scrollOffset

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

	cursorX := t.getXPositionForIndex(t.cursorPos) - t.scrollOffset

	cursor := ebiten.NewImage(1, int(size.Height-padding.Top-padding.Bottom))
	cursor.Fill(t.cursorColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		pos.X+padding.Left+cursorX,
		pos.Y+padding.Top,
	)
	screen.DrawImage(cursor, op)
}

func (t *TextInput) ensureCursorVisible() {
	padding := t.GetPadding()
	availableWidth := t.GetSize().Width - padding.Left - padding.Right
	cursorX := t.getXPositionForIndex(t.cursorPos)

	// If cursor is to the left of the visible area
	if cursorX < t.scrollOffset {
		t.scrollOffset = cursorX
	}

	// If cursor is beyond the right edge of the visible area
	if cursorX > t.scrollOffset+availableWidth {
		t.scrollOffset = cursorX - availableWidth
	}

	// Clamp scroll offset to valid range
	maxScroll := math.Max(0, t.getXPositionForIndex(len(t.text))-availableWidth)
	t.scrollOffset = clamp(t.scrollOffset, 0, maxScroll)
}

func (t *TextInput) handleLeftKey(ctrlPressed, shiftPressed bool) {
	prevPos := t.cursorPos
	if ctrlPressed {
		t.cursorPos = t.findPreviousWordBoundary()
	} else if t.cursorPos > 0 {
		t.cursorPos--
	}

	t.updateSelection(prevPos, shiftPressed)
}

func (t *TextInput) handleRightKey(ctrlPressed, shiftPressed bool) {
	prevPos := t.cursorPos
	if ctrlPressed {
		t.cursorPos = t.findNextWordBoundary()
	} else if t.cursorPos < len(t.text) {
		t.cursorPos++
	}

	t.updateSelection(prevPos, shiftPressed)
}

func (t *TextInput) updateSelection(prevPos int, shiftPressed bool) {
	if !shiftPressed {
		t.ClearSelection()
	} else {
		if t.selectionStart == -1 {
			t.selectionStart = prevPos
		}
		t.selectionEnd = t.cursorPos
	}
}

func (t *TextInput) handleHome(shiftPressed bool) {
	prevPos := t.cursorPos
	t.cursorPos = 0
	t.updateSelection(prevPos, shiftPressed)
}

func (t *TextInput) handleEnd(shiftPressed bool) {
	prevPos := t.cursorPos
	t.cursorPos = len(t.text)
	t.updateSelection(prevPos, shiftPressed)
}

func (t *TextInput) handleBackspace() {
	if t.hasSelection() {
		t.deleteSelection()
	} else if t.cursorPos > 0 {
		t.text = append(t.text[:t.cursorPos-1], t.text[t.cursorPos:]...)
		t.cursorPos--
		if t.onChange != nil {
			t.onChange(string(t.text))
		}
	}
}

func (t *TextInput) handleDelete() {
	if t.hasSelection() {
		t.deleteSelection()
	} else if t.cursorPos < len(t.text) {
		t.text = append(t.text[:t.cursorPos], t.text[t.cursorPos+1:]...)
		if t.onChange != nil {
			t.onChange(string(t.text))
		}
	}
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
	t.ClearSelection()

	if t.onChange != nil {
		t.onChange(string(t.text))
	}
}

// Public API methods

func (t *TextInput) SetText(text string) {
	t.text = []rune(text)
	t.cursorPos = len(t.text)
	t.ClearSelection()
	t.scrollOffset = 0
	t.ensureCursorVisible()
	if t.onChange != nil {
		t.onChange(text)
	}
}

func (t *TextInput) GetText() string {
	return string(t.text)
}

func (t *TextInput) Focus() {
	t.isFocused = true
	t.showCursor = true
	t.lastBlink = time.Now()
	t.ensureCursorVisible()
}

func (t *TextInput) Blur() {
	t.isFocused = false
	t.showCursor = false
	t.ClearSelection()
}

func (t *TextInput) IsFocused() bool {
	return t.isFocused
}

func (t *TextInput) SetFont(font font.Face) {
	t.font = font
}

func (t *TextInput) selectAll() {
	if len(t.text) > 0 {
		t.selectionStart = 0
		t.selectionEnd = len(t.text)
		t.cursorPos = t.selectionEnd
		t.showCursor = true
		t.lastBlink = time.Now()
		t.ensureCursorVisible()
	}
}

func (t *TextInput) handleCopy() {
	if clipboardDisabled {
		return
	}
	if !t.hasSelection() {
		return
	}
	start, end := t.getOrderedSelection()
	text := string(t.text[start:end])
	clipboard.Write(clipboard.FmtText, []byte(text))
}

func (t *TextInput) handleCut() {
	if clipboardDisabled {
		return
	}
	if !t.hasSelection() {
		return
	}
	t.handleCopy()
	t.deleteSelection()
}

func (t *TextInput) handlePaste() {
	if clipboardDisabled {
		return
	}

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
	t.ensureCursorVisible()

	if t.onChange != nil {
		t.onChange(string(t.text))
	}
}

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
	t.ensureCursorVisible()
}

func (t *TextInput) ClearSelection() {
	t.selectionStart = -1
	t.selectionEnd = -1
}
