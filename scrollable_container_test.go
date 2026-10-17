package ebui

import "testing"

func TestScrollableContainerOverflows(t *testing.T) {
	sc := NewScrollableContainer(
		WithSize(100, 100),
		WithLayout(NewVerticalStackLayout(0, AlignStart)),
	)
	sc.AddChild(NewLabel("a", WithSize(100, 60)))
	if sc.Overflows() {
		t.Fatal("60px of content overflows a 100px container")
	}
	sc.AddChild(NewLabel("b", WithSize(100, 60)))
	if !sc.Overflows() {
		t.Fatal("120px of content doesn't overflow a 100px container")
	}
}
