package ebui

import "sort"

type FocusableComponent interface {
	InteractiveComponent
	IsFocusable() bool
	SetFocusable(focusable bool)
	GetTabIndex() int
	SetTabIndex(index int)
}

type BaseFocusable struct {
	*BaseInteractive
	focusable bool
	tabIndex  int
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

type FocusManager struct {
	focusableComponents []FocusableComponent
	currentFocus        FocusableComponent
}

func NewFocusManager() *FocusManager {
	return &FocusManager{}
}

// RefreshFocusableComponents finds all focusable components in the component tree
func (fm *FocusManager) RefreshFocusableComponents(root Component) {
	fm.focusableComponents = nil

	// Find all focusable components
	var findFocusables func(Component)
	findFocusables = func(c Component) {
		if focusable, ok := c.(FocusableComponent); ok {
			fm.focusableComponents = append(fm.focusableComponents, focusable)
		}

		if container, ok := c.(Container); ok {
			for _, child := range container.GetChildren() {
				findFocusables(child)
			}
		}
	}

	findFocusables(root)

	// Sort focusables by tab index
	sort.SliceStable(fm.focusableComponents, func(i, j int) bool {
		return fm.focusableComponents[i].GetTabIndex() < fm.focusableComponents[j].GetTabIndex()
	})
}

func (fm *FocusManager) SetFocus(component FocusableComponent) {
	if fm.currentFocus == component {
		return
	}

	// Handle blur for previous focus
	if fm.currentFocus != nil {
		blurEvent := &Event{
			Type:          Blur,
			Target:        fm.currentFocus,
			RelatedTarget: component,
		}
		fm.currentFocus.HandleEvent(blurEvent)
	}

	fm.currentFocus = component

	// Handle focus for new component
	if component != nil {
		focusEvent := &Event{
			Type:          Focus,
			Target:        component,
			RelatedTarget: fm.currentFocus,
		}
		component.HandleEvent(focusEvent)
	}
}

func (fm *FocusManager) GetCurrentFocus() FocusableComponent {
	return fm.currentFocus
}

func (fm *FocusManager) HandleTab(shiftPressed bool) {
	if len(fm.focusableComponents) == 0 {
		return
	}

	// Find current focus index
	currentIndex := -1
	for i, c := range fm.focusableComponents {
		if c == fm.currentFocus {
			currentIndex = i
			break
		}
	}

	// Calculate next focus index
	var nextIndex int
	if shiftPressed {
		if currentIndex <= 0 {
			nextIndex = len(fm.focusableComponents) - 1
		} else {
			nextIndex = currentIndex - 1
		}
	} else {
		if currentIndex >= len(fm.focusableComponents)-1 {
			nextIndex = 0
		} else {
			nextIndex = currentIndex + 1
		}
	}

	fm.SetFocus(fm.focusableComponents[nextIndex])
}
