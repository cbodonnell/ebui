package ebui

import "testing"

func button(name string, x, y float64, clicks *[]string) *Button {
	return NewButton(
		WithSize(20, 20),
		WithPosition(Position{X: x, Y: y, Relative: true}),
		WithLabelText(name),
		WithClickHandler(func() { *clicks = append(*clicks, name) }),
	)
}

func newTestManager(root Container) *Manager {
	return NewManager(root, WithNavInput(nil))
}

func focusedLabel(t *testing.T, m *Manager) string {
	t.Helper()
	b, ok := m.Focused().(*Button)
	if !ok {
		return ""
	}
	return b.GetLabel()
}

func grid(root Container, clicks *[]string) {
	for i, name := range []string{"a", "b", "c", "d", "e", "f"} {
		root.AddChild(button(name, float64(i%3)*30, float64(i/3)*30, clicks))
	}
}

func TestArrowsMoveThroughAGrid(t *testing.T) {
	var clicks []string
	root := NewBaseContainer(WithSize(200, 200))
	grid(root, &clicks)
	m := newTestManager(root)

	steps := []struct {
		action NavAction
		want   string
		used   bool
	}{
		{NavRight, "a", true},
		{NavRight, "b", true},
		{NavDown, "e", true},
		{NavDown, "e", false},
		{NavLeft, "d", true},
		{NavUp, "a", true},
		{NavRight, "b", true},
		{NavRight, "c", true},
		{NavRight, "c", false},
	}
	for i, s := range steps {
		if used := m.Navigate(s.action); used != s.used {
			t.Fatalf("step %d %v: used = %v, want %v", i, s.action, used, s.used)
		}
		if got := focusedLabel(t, m); got != s.want {
			t.Fatalf("step %d %v: focus %q, want %q", i, s.action, got, s.want)
		}
	}

	if !m.Navigate(NavActivate) || len(clicks) != 1 || clicks[0] != "c" {
		t.Fatalf("activate clicked %v, want [c]", clicks)
	}
}

func TestTabCyclesOneScopeAndWraps(t *testing.T) {
	var clicks []string
	root := NewBaseContainer(WithSize(200, 200))
	grid(root, &clicks)
	m := newTestManager(root)

	m.Navigate(NavNext)
	for _, want := range []string{"b", "c", "d", "e", "f", "a"} {
		m.Navigate(NavNext)
		if got := focusedLabel(t, m); got != want {
			t.Fatalf("tab: focus %q, want %q", got, want)
		}
	}
	m.Navigate(NavPrev)
	if got := focusedLabel(t, m); got != "f" {
		t.Fatalf("shift+tab: focus %q, want f", got)
	}
}

func TestUnusedActionsReportUnused(t *testing.T) {
	root := NewBaseContainer(WithSize(200, 200))
	m := newTestManager(root)
	for _, a := range []NavAction{NavUp, NavNext, NavActivate, NavCancel} {
		if m.Navigate(a) {
			t.Errorf("%v used with nothing to focus", a)
		}
	}

	var clicks []string
	root.AddChild(button("a", 0, 0, &clicks))
	m.Navigate(NavDown)
	if m.Navigate(NavCancel) {
		t.Errorf("cancel used outside a window")
	}
}

func windowWith(wm *WindowManager, title string, x float64, names []string, clicks *[]string) *Window {
	w := wm.CreateWindow(100, 100, WithWindowTitle(title), WithWindowPosition(x, 0))
	content := NewBaseContainer(WithSize(100, 70), WithPosition(Position{Relative: true}))
	for i, name := range names {
		content.AddChild(button(name, float64(i)*30, 0, clicks))
	}
	w.AddChild(content)
	return w
}

func TestWindowsAreScopes(t *testing.T) {
	var clicks []string
	root := NewBaseContainer(WithSize(400, 200))
	wm := NewWindowManager(WithSize(400, 200))
	root.AddChild(wm)
	inv := windowWith(wm, "inventory", 0, []string{"i1", "i2"}, &clicks)
	eq := windowWith(wm, "equips", 200, []string{"e1", "e2"}, &clicks)
	inv.Hide()
	eq.Hide()
	m := newTestManager(root)
	if err := m.Update(); err != nil {
		t.Fatal(err)
	}

	inv.Show()
	m.Update()
	if got := focusedLabel(t, m); got != "i1" {
		t.Fatalf("after opening inventory: focus %q, want i1", got)
	}
	eq.Show()
	m.Update()
	if got := focusedLabel(t, m); got != "e1" {
		t.Fatalf("after opening equips: focus %q, want e1", got)
	}

	m.Navigate(NavRight)
	if m.Navigate(NavRight) {
		t.Fatalf("right from the last in equips left the window")
	}
	if got := focusedLabel(t, m); got != "e2" {
		t.Fatalf("focus %q, want e2", got)
	}

	m.Navigate(NavNext)
	if got := focusedLabel(t, m); got != "i1" {
		t.Fatalf("tab to inventory: focus %q, want i1", got)
	}
	m.Navigate(NavNext)
	if got := focusedLabel(t, m); got != "e2" {
		t.Fatalf("tab back to equips: focus %q, want e2", got)
	}

	if !m.Navigate(NavCancel) || eq.IsVisible() {
		t.Fatalf("cancel didn't close equips")
	}
	m.Update()
	if got := focusedLabel(t, m); got != "i1" {
		t.Fatalf("after closing equips: focus %q, want i1", got)
	}
	m.Navigate(NavCancel)
	m.Update()
	if m.Focused() != nil {
		t.Fatalf("focus %q with every window closed", focusedLabel(t, m))
	}
}

func TestModalKeepsTheFocus(t *testing.T) {
	var clicks []string
	root := NewBaseContainer(WithSize(400, 200))
	wm := NewWindowManager(WithSize(400, 200))
	root.AddChild(wm)
	windowWith(wm, "inventory", 0, []string{"i1"}, &clicks)
	prompt := wm.CreateWindow(100, 100, WithModal())
	content := NewBaseContainer(WithSize(100, 70), WithPosition(Position{Relative: true}))
	content.AddChild(button("ok", 0, 0, &clicks))
	content.AddChild(button("cancel", 30, 0, &clicks))
	prompt.AddChild(content)
	m := newTestManager(root)
	m.Update()

	if got := focusedLabel(t, m); got != "ok" {
		t.Fatalf("focus %q, want ok", got)
	}
	m.Navigate(NavNext)
	if got := focusedLabel(t, m); got != "cancel" {
		t.Fatalf("tab in the modal: focus %q, want cancel", got)
	}
	m.Navigate(NavNext)
	if got := focusedLabel(t, m); got != "ok" {
		t.Fatalf("tab wrapped out of the modal to %q", got)
	}
}

func TestTextInputCapturesKeysButNotWhenAWindowOpens(t *testing.T) {
	root := NewBaseContainer(WithSize(400, 200))
	wm := NewWindowManager(WithSize(400, 200))
	root.AddChild(wm)
	w := wm.CreateWindow(100, 100)
	content := NewBaseContainer(WithSize(100, 70), WithPosition(Position{Relative: true}))
	input := NewTextInput(WithSize(60, 20), WithPosition(Position{Relative: true}))
	var clicks []string
	content.AddChild(input)
	content.AddChild(button("ok", 0, 30, &clicks))
	w.AddChild(content)
	m := newTestManager(root)
	m.Update()

	if got := focusedLabel(t, m); got != "ok" {
		t.Fatalf("opening the window focused %T, want the button", m.Focused())
	}
	if m.IsCapturingKeys() {
		t.Fatalf("capturing keys with a button focused")
	}

	input.Focus()
	m.Update()
	if m.Focused() != input || !m.IsCapturingKeys() {
		t.Fatalf("focus %T after Focus(), want the input capturing keys", m.Focused())
	}
	if !m.Navigate(NavLeft) || !m.Navigate(NavActivate) {
		t.Fatalf("the input let left or enter through")
	}
	input.Blur()
	m.Update()
	if m.Focused() != nil {
		t.Fatalf("focus %T after Blur()", m.Focused())
	}
}

func TestRebuiltContentKeepsTheFocusNearby(t *testing.T) {
	var clicks []string
	root := NewBaseContainer(WithSize(200, 200))
	list := NewBaseContainer(WithSize(200, 200))
	root.AddChild(list)
	grid(list, &clicks)
	m := newTestManager(root)
	m.Navigate(NavDown)
	m.Navigate(NavRight)
	m.Navigate(NavDown)
	if got := focusedLabel(t, m); got != "e" {
		t.Fatalf("focus %q, want e", got)
	}

	list.ClearChildren()
	grid(list, &clicks)
	m.Update()
	if got := focusedLabel(t, m); got != "e" {
		t.Fatalf("after rebuilding: focus %q, want the new e", got)
	}
}

func TestPointerHidesTheFocusRing(t *testing.T) {
	var clicks []string
	root := NewBaseContainer(WithSize(200, 200))
	grid(root, &clicks)
	m := newTestManager(root)
	m.Navigate(NavDown)
	if !m.Focused().IsFocusVisible() {
		t.Fatalf("navigated focus isn't visible")
	}
	b := root.GetChildren()[4].(*Button)
	m.input.focusManager.focusFromPointer(b)
	if m.Focused() != b || b.IsFocusVisible() {
		t.Fatalf("pressed button: focused %v, visible %v", m.Focused() == b, b.IsFocusVisible())
	}
	m.Navigate(NavLeft)
	if got := focusedLabel(t, m); got != "d" || !m.Focused().IsFocusVisible() {
		t.Fatalf("navigating after a press: focus %q, visible %v", got, m.Focused().IsFocusVisible())
	}
}

func TestKeyRepeat(t *testing.T) {
	var fired []int
	for d := 0; d <= 40; d++ {
		if repeats(d, 60) {
			fired = append(fired, d)
		}
	}
	want := []int{1, 28, 32, 36, 40}
	if len(fired) != len(want) {
		t.Fatalf("fired on %v, want %v", fired, want)
	}
	for i := range want {
		if fired[i] != want[i] {
			t.Fatalf("fired on %v, want %v", fired, want)
		}
	}
}
