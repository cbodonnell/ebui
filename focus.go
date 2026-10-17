package ebui

import "sort"

type FocusableComponent interface {
	InteractiveComponent
	IsFocusable() bool
	SetFocusable(focusable bool)
	GetTabIndex() int
	SetTabIndex(index int)
	IsFocused() bool
	IsFocusVisible() bool
	RequestFocus()
	setFocusState(focused, visible bool)
	takeFocusRequest() bool
}

// KeyCapturer is a component that takes the keyboard while focused.
type KeyCapturer interface {
	CapturesKeys() bool
}

// FocusScope groups focusables for navigation, as windows do.
type FocusScope interface {
	Component
	IsFocusScope() bool
}

// ModalScope is a scope that keeps the focus while open.
type ModalScope interface {
	IsModal() bool
}

type BaseFocusable struct {
	*BaseInteractive
	focusable      bool
	tabIndex       int
	focused        bool
	focusVisible   bool
	focusRequested bool
}

func NewBaseFocusable() *BaseFocusable {
	return &BaseFocusable{
		BaseInteractive: NewBaseInteractive(),
		focusable:       true,
		tabIndex:        0,
	}
}

func (b *BaseFocusable) IsFocusable() bool {
	return b.focusable
}

func (b *BaseFocusable) SetFocusable(focusable bool) {
	b.focusable = focusable
}

func (b *BaseFocusable) GetTabIndex() int {
	return b.tabIndex
}

func (b *BaseFocusable) SetTabIndex(index int) {
	b.tabIndex = index
}

func (b *BaseFocusable) IsFocused() bool {
	return b.focused
}

func (b *BaseFocusable) IsFocusVisible() bool {
	return b.focused && b.focusVisible
}

func (b *BaseFocusable) RequestFocus() {
	b.focusRequested = true
}

func (b *BaseFocusable) setFocusState(focused, visible bool) {
	b.focused = focused
	b.focusVisible = focused && visible
}

func (b *BaseFocusable) takeFocusRequest() bool {
	requested := b.focusRequested
	b.focusRequested = false
	return requested
}

func WithFocusable(focusable bool) ComponentOpt {
	return func(c Component) {
		if f, ok := c.(FocusableComponent); ok {
			f.SetFocusable(focusable)
		}
	}
}

func WithTabIndex(index int) ComponentOpt {
	return func(c Component) {
		if f, ok := c.(FocusableComponent); ok {
			f.SetTabIndex(index)
		}
	}
}

type focusEntry struct {
	component FocusableComponent
	scope     Component
}

type FocusManager struct {
	entries []focusEntry
	scopes  []Component
	// parents maps components to their outer containers; GetParent returns embedded bases
	parents        map[Component]Component
	currentFocus   FocusableComponent
	currentScope   Component
	lastRect       rect
	lostFromWindow bool
	remembered     map[Component]FocusableComponent
	topmost        Component
	keyboard       bool
	enabled        bool
}

func NewFocusManager() *FocusManager {
	return &FocusManager{
		remembered: make(map[Component]FocusableComponent),
		parents:    make(map[Component]Component),
		keyboard:   true,
		enabled:    true,
	}
}

// RefreshFocusableComponents finds all focusable components in the component tree
func (fm *FocusManager) RefreshFocusableComponents(root Component) {
	if !fm.enabled {
		return
	}

	fm.collect(root)

	for _, e := range fm.entries {
		if e.component.takeFocusRequest() {
			fm.setFocus(e.component, fm.keyboard)
		}
	}

	if fm.currentFocus != nil {
		switch {
		case !fm.currentFocus.IsFocused():
			fm.currentFocus = nil
			fm.currentScope = nil
		case fm.scopeOf(fm.currentFocus) == nil:
			scope := fm.currentScope
			fm.blur()
			delete(fm.remembered, scope)
			if fm.isOpen(scope) {
				if near := fm.nearestTo(scope, fm.lastRect); near != nil {
					fm.setFocus(near, fm.keyboard)
				}
			}
		}
	}

	topmost := fm.topmostWindow()
	switch {
	case topmost != fm.topmost:
		fm.topmost = topmost
		if topmost != nil && fm.currentScope != topmost {
			fm.enterScope(topmost, true)
		}
	case fm.currentFocus == nil && fm.lostFromWindow && topmost != nil:
		fm.enterScope(topmost, true)
	}
	fm.lostFromWindow = false
}

func (fm *FocusManager) collect(root Component) {
	fm.entries = fm.entries[:0]
	fm.scopes = fm.scopes[:0]
	clear(fm.parents)
	var windows []Component

	var visit func(c Component, scope Component)
	visit = func(c Component, scope Component) {
		if c.IsDisabled() || c.IsHidden() {
			return
		}
		if s, ok := c.(FocusScope); ok && s.IsFocusScope() && c != root {
			scope = c
			windows = append(windows, c)
		}
		focusable, isFocusable := c.(FocusableComponent)
		isFocusable = isFocusable && focusable.IsFocusable()
		at := len(fm.entries)
		if container, ok := c.(Container); ok {
			for _, child := range container.GetChildren() {
				fm.parents[child] = c
				visit(child, scope)
			}
		}
		// containers only take focus when nothing inside them can
		if isFocusable && len(fm.entries) == at {
			fm.entries = append(fm.entries, focusEntry{component: focusable, scope: scope})
		}
	}
	visit(root, root)

	sort.SliceStable(fm.entries, func(i, j int) bool {
		return fm.entries[i].component.GetTabIndex() < fm.entries[j].component.GetTabIndex()
	})

	for _, e := range fm.entries {
		if e.scope == root {
			fm.scopes = append(fm.scopes, root)
			break
		}
	}
	sort.SliceStable(windows, func(i, j int) bool { return windows[i].GetID() < windows[j].GetID() })
	fm.scopes = append(fm.scopes, windows...)

	var modal Component
	for _, w := range windows {
		if m, ok := w.(ModalScope); ok && m.IsModal() {
			if modal == nil || w.GetPosition().ZIndex > modal.GetPosition().ZIndex {
				modal = w
			}
		}
	}
	if modal != nil {
		fm.scopes = append(fm.scopes[:0], modal)
		entries := fm.entries[:0]
		for _, e := range fm.entries {
			if e.scope == modal {
				entries = append(entries, e)
			}
		}
		fm.entries = entries
	}
}

func (fm *FocusManager) scopeOf(c FocusableComponent) Component {
	for _, e := range fm.entries {
		if e.component == c {
			return e.scope
		}
	}
	return nil
}

func (fm *FocusManager) isOpen(scope Component) bool {
	for _, s := range fm.scopes {
		if s == scope {
			return true
		}
	}
	return false
}

func (fm *FocusManager) inScope(scope Component) []FocusableComponent {
	var cs []FocusableComponent
	for _, e := range fm.entries {
		if e.scope == scope {
			cs = append(cs, e.component)
		}
	}
	return cs
}

func (fm *FocusManager) topmostWindow() Component {
	var top Component
	for _, s := range fm.scopes {
		if _, ok := s.(FocusScope); !ok {
			continue
		}
		if top == nil || s.GetPosition().ZIndex > top.GetPosition().ZIndex {
			top = s
		}
	}
	return top
}

// enterScope restores the scope's last focus, or focuses its first component.
func (fm *FocusManager) enterScope(scope Component, auto bool) bool {
	if r := fm.remembered[scope]; r != nil && fm.scopeOf(r) == scope {
		fm.setFocus(r, fm.keyboard)
		return true
	}
	for _, c := range fm.inScope(scope) {
		if kc, ok := c.(KeyCapturer); auto && ok && kc.CapturesKeys() {
			continue
		}
		fm.setFocus(c, fm.keyboard)
		return true
	}
	return false
}

func (fm *FocusManager) enterDefault() bool {
	if top := fm.topmostWindow(); top != nil {
		return fm.enterScope(top, false)
	}
	for _, s := range fm.scopes {
		if fm.enterScope(s, false) {
			return true
		}
	}
	return false
}

func (fm *FocusManager) SetFocus(component FocusableComponent) {
	if !fm.enabled {
		return
	}
	if component == nil {
		fm.blur()
		return
	}
	fm.setFocus(component, fm.keyboard)
}

func (fm *FocusManager) setFocus(component FocusableComponent, visible bool) {
	if fm.currentFocus == component {
		component.setFocusState(true, visible)
		return
	}

	previous := fm.currentFocus
	fm.blur()

	fm.currentFocus = component
	fm.currentScope = fm.scopeOf(component)
	fm.lastRect = rectOf(component)
	if fm.currentScope != nil {
		fm.remembered[fm.currentScope] = component
	}
	component.setFocusState(true, visible)
	component.HandleEvent(&Event{
		Type:          Focus,
		Target:        component,
		RelatedTarget: previous,
	})

	for p := fm.parents[component]; p != nil; p = fm.parents[p] {
		if s, ok := p.(interface{ ScrollIntoView(Component) }); ok {
			s.ScrollIntoView(component)
		}
	}
	if w, ok := fm.currentScope.(interface{ Raise() }); ok {
		w.Raise()
		fm.topmost = fm.currentScope
	}
}

func (fm *FocusManager) blur() {
	if fm.currentFocus == nil {
		return
	}
	previous := fm.currentFocus
	_, fm.lostFromWindow = fm.currentScope.(FocusScope)
	fm.lastRect = rectOf(previous)
	fm.currentFocus = nil
	fm.currentScope = nil
	previous.setFocusState(false, false)
	previous.HandleEvent(&Event{
		Type:   Blur,
		Target: previous,
	})
}

func (fm *FocusManager) GetCurrentFocus() FocusableComponent {
	if !fm.enabled {
		return nil
	}
	return fm.currentFocus
}

func (fm *FocusManager) focusFromPointer(target InteractiveComponent) {
	fm.keyboard = false
	if target != nil {
		for c := Component(target); c != nil; c = fm.parents[c] {
			if f, ok := c.(FocusableComponent); ok && fm.scopeOf(f) != nil {
				fm.setFocus(f, false)
				return
			}
		}
	}
	fm.blur()
	fm.lostFromWindow = false
}

// Navigate applies a navigation action and reports whether it was handled.
func (fm *FocusManager) Navigate(root Component, action NavAction) bool {
	if !fm.enabled {
		return false
	}
	fm.RefreshFocusableComponents(root)
	fm.keyboard = true

	if fm.currentFocus != nil {
		fm.currentFocus.setFocusState(true, true)
		if fm.dispatchNav(action) {
			return true
		}
	}

	switch action {
	case NavUp, NavDown, NavLeft, NavRight:
		if fm.currentFocus == nil {
			return fm.enterDefault()
		}
		next := fm.nearest(action)
		if next == nil {
			return false
		}
		fm.setFocus(next, true)
		return true
	case NavNext, NavPrev:
		return fm.cycle(action == NavNext)
	}
	return false
}

func (fm *FocusManager) dispatchNav(action NavAction) bool {
	event := &Event{
		Type:          Nav,
		Nav:           action,
		Target:        fm.currentFocus,
		CurrentTarget: fm.currentFocus,
		Bubbles:       true,
		Phase:         PhaseTarget,
	}
	fm.currentFocus.HandleEvent(event)
	event.Phase = PhaseBubble
	for p := fm.parents[fm.currentFocus]; p != nil && !event.Handled; p = fm.parents[p] {
		if i, ok := p.(InteractiveComponent); ok {
			event.CurrentTarget = i
			i.HandleEvent(event)
		}
	}
	return event.Handled
}

// cycle moves between scopes, or within the only open scope.
func (fm *FocusManager) cycle(forward bool) bool {
	if fm.currentFocus == nil {
		return fm.enterDefault()
	}
	step := func(i, n int) int {
		if forward {
			return (i + 1) % n
		}
		return (i - 1 + n) % n
	}
	if len(fm.scopes) > 1 {
		at := 0
		for i, s := range fm.scopes {
			if s == fm.currentScope {
				at = i
			}
		}
		for n, i := len(fm.scopes), step(at, len(fm.scopes)); i != at; i = step(i, n) {
			if fm.enterScope(fm.scopes[i], false) {
				return true
			}
		}
		return false
	}
	cs := fm.inScope(fm.currentScope)
	if len(cs) < 2 {
		return false
	}
	at := 0
	for i, c := range cs {
		if c == fm.currentFocus {
			at = i
		}
	}
	fm.setFocus(cs[step(at, len(cs))], true)
	return true
}

// Enable turns on focus management
func (fm *FocusManager) Enable() {
	fm.enabled = true
}

// Disable turns off focus management
func (fm *FocusManager) Disable() {
	fm.blur()
	fm.enabled = false
}

// IsEnabled returns whether focus management is enabled
func (fm *FocusManager) IsEnabled() bool {
	return fm.enabled
}
