package ebui

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// recorder records what it is asked to paint and draws only the parts in draws.
type recorder struct {
	paints []Paint
	draws  map[Part]bool
}

func (r *recorder) Paint(dst *ebiten.Image, p Paint) bool {
	r.paints = append(r.paints, p)
	return r.draws[p.Part]
}

func (r *recorder) find(part Part, role string) (Paint, bool) {
	for _, p := range r.paints {
		if p.Part == part && p.Role == role {
			return p, true
		}
	}
	return Paint{}, false
}

func withPainter(t *testing.T, p Painter) {
	t.Helper()
	was := DefaultPainter
	DefaultPainter = p
	t.Cleanup(func() { DefaultPainter = was })
}

func TestButtonsPaintTheirRoleAndState(t *testing.T) {
	rec := &recorder{}
	withPainter(t, rec)
	b := NewButton(WithSize(40, 20), WithRole("primary"), WithLabelText("Buy"))
	b.isHovered = true
	b.Disable()
	b.Draw(ebiten.NewImage(50, 30))

	face, ok := rec.find(PartButton, "primary")
	if !ok {
		t.Fatal("no button face painted")
	}
	if !face.State.Has(StateHovered) || !face.State.Has(StateDisabled) || face.State.Has(StatePressed) {
		t.Errorf("state = %b, want hovered and disabled", face.State)
	}
	if face.Rect.Dx() != 40 || face.Rect.Dy() != 20 {
		t.Errorf("rect = %v, want 40×20", face.Rect)
	}
	label, ok := rec.find(PartButtonLabel, "primary")
	if !ok || label.Text != "Buy" {
		t.Errorf("label paint = %+v, want Buy", label)
	}
}

func TestPaintedTitlesReplaceTheLabel(t *testing.T) {
	for _, drawn := range []bool{false, true} {
		rec := &recorder{draws: map[Part]bool{PartWindowTitle: drawn}}
		withPainter(t, rec)
		wm := NewWindowManager(WithSize(200, 200))
		w := wm.CreateWindow(100, 80, WithWindowTitle("Bag"), WithWindowRole("company"))
		w.Draw(ebiten.NewImage(200, 200))

		if _, ok := rec.find(PartWindow, "company"); !ok {
			t.Fatal("no window frame painted in its role")
		}
		title, _ := rec.find(PartWindowTitle, "company")
		if title.Text != "Bag" {
			t.Errorf("title text = %q", title.Text)
		}
		if w.titleLabel.IsHidden() != drawn {
			t.Errorf("painter drew title %v, label hidden %v", drawn, w.titleLabel.IsHidden())
		}
		if _, ok := rec.find(PartButton, "close"); !ok {
			t.Error("close button not painted as close")
		}
	}
}

func TestDeclinedPartsDrawFlat(t *testing.T) {
	rec := &recorder{}
	withPainter(t, rec)
	s := NewSlider(WithSize(100, 20))
	s.Draw(ebiten.NewImage(120, 30))
	for _, part := range []Part{PartSliderTrack, PartSliderFill, PartSliderThumb} {
		if _, ok := rec.find(part, ""); !ok {
			t.Errorf("part %d not offered to the painter", part)
		}
	}
}
