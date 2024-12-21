package main

import (
	"image/color"
	"log"

	"github.com/cbodonnell/ebui"
	"github.com/hajimehoshi/ebiten/v2"
)

type InventorySlot struct {
	*ebui.LayoutContainer
	*ebui.BaseInteractive
	item      *Item
	label     *ebui.Label
	isHovered bool
	game      *InventoryGame
}

type Item struct {
	Name  string
	Color color.Color
}

func NewInventorySlot(game *InventoryGame) *InventorySlot {
	slot := &InventorySlot{
		LayoutContainer: ebui.NewLayoutContainer(
			ebui.WithSize(60, 60),
			ebui.WithBackground(color.RGBA{200, 200, 200, 255}),
			ebui.WithLayout(ebui.NewVerticalStackLayout(0, ebui.AlignCenter)),
		),
		BaseInteractive: ebui.NewBaseInteractive(),
		game:            game,
	}

	slot.label = ebui.NewLabel(
		"",
		ebui.WithSize(60, 20),
		ebui.WithJustify(ebui.JustifyCenter),
	)
	slot.AddChild(slot.label)

	slot.registerEventListeners()
	return slot
}

func (s *InventorySlot) registerEventListeners() {
	s.AddEventListener(ebui.MouseEnter, func(e *ebui.Event) {
		s.isHovered = true
	})

	s.AddEventListener(ebui.MouseLeave, func(e *ebui.Event) {
		s.isHovered = false
	})

	s.AddEventListener(ebui.DragStart, func(e *ebui.Event) {
		if s.item != nil {
			s.game.startDragging(s, e.MouseX, e.MouseY)
		}
	})

	s.AddEventListener(ebui.Drag, func(e *ebui.Event) {
		s.game.updateDragPosition(e.MouseX, e.MouseY)
	})

	s.AddEventListener(ebui.Drop, func(e *ebui.Event) {
		if sourceSlot, ok := e.RelatedTarget.(*InventorySlot); ok {
			// Swap items between slots
			s.item, sourceSlot.item = sourceSlot.item, s.item
			s.updateDisplay()
			sourceSlot.updateDisplay()
			s.game.endDragging()
		}
		s.isHovered = false
	})

	s.AddEventListener(ebui.DragEnd, func(e *ebui.Event) {
		if s.game.draggedItem != nil {
			// Check if the mouse is within the inventory bounds
			if s.game.isWithinInventory(e.MouseX, e.MouseY) {
				// Do nothing, item will return to original slot
				s.game.endDragging()
			} else {
				// Remove the item from the source slot
				s.game.dragStartSlot.item = nil
				s.game.dragStartSlot.updateDisplay()
				s.game.endDragging()
			}
		}
		s.isHovered = false
	})
}

func (s *InventorySlot) HandleEvent(event *ebui.Event) {
	s.BaseInteractive.HandleEvent(event)
}

func (s *InventorySlot) Update() error {
	if s.item != nil {
		if s.isHovered {
			// Lighten the item's color when hovered
			r, g, b, _ := s.item.Color.RGBA()
			s.SetBackground(color.RGBA{
				uint8(min(255, (r>>8)+20)),
				uint8(min(255, (g>>8)+20)),
				uint8(min(255, (b>>8)+20)),
				255,
			})
		} else {
			s.SetBackground(s.item.Color)
		}
	} else {
		if s.isHovered {
			s.SetBackground(color.RGBA{220, 220, 220, 255})
		} else {
			s.SetBackground(color.RGBA{200, 200, 200, 255})
		}
	}
	return s.LayoutContainer.Update()
}

func (s *InventorySlot) SetItem(item *Item) {
	s.item = item
	s.updateDisplay()
}

func (s *InventorySlot) updateDisplay() {
	if s.item != nil {
		s.label.SetText(s.item.Name)
		s.SetBackground(s.item.Color)
	} else {
		s.label.SetText("")
		s.SetBackground(color.RGBA{200, 200, 200, 255})
	}
}

type InventoryGame struct {
	ui            *ebui.Manager
	slots         []*InventorySlot
	draggedItem   *Item
	dragStartSlot *InventorySlot
	dragX, dragY  float64
	gridContainer *ebui.LayoutContainer
}

func NewInventoryGame() *InventoryGame {
	game := &InventoryGame{}

	// Create root container
	root := ebui.NewBaseContainer(
		ebui.WithSize(310, 310),
		ebui.WithBackground(color.RGBA{240, 240, 240, 255}),
	)

	// Create inventory grid container
	gridContainer := ebui.NewLayoutContainer(
		ebui.WithSize(290, 290),
		ebui.WithPosition(ebui.Position{X: 10, Y: 10}),
		ebui.WithBackground(color.RGBA{255, 255, 255, 255}),
		ebui.WithPadding(10, 10, 10, 10),
		ebui.WithLayout(ebui.NewVerticalStackLayout(10, ebui.AlignStart)),
	)
	game.gridContainer = gridContainer

	// Create rows of inventory slots
	game.slots = make([]*InventorySlot, 0)
	for row := 0; row < 4; row++ {
		rowContainer := ebui.NewLayoutContainer(
			ebui.WithSize(260, 60),
			ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignStart)),
		)

		for col := 0; col < 4; col++ {
			slot := NewInventorySlot(game)
			game.slots = append(game.slots, slot)
			rowContainer.AddChild(slot)
		}

		gridContainer.AddChild(rowContainer)
	}

	// Add some sample items
	sampleItems := []Item{
		{Name: "Sword", Color: color.RGBA{255, 0, 0, 255}},
		{Name: "Shield", Color: color.RGBA{0, 0, 255, 255}},
		{Name: "Potion", Color: color.RGBA{0, 255, 0, 255}},
		{Name: "Bow", Color: color.RGBA{139, 69, 19, 255}},
		{Name: "Arrow", Color: color.RGBA{128, 128, 128, 255}},
	}

	// Place items in first row
	for i := 0; i < len(sampleItems); i++ {
		item := sampleItems[i]
		game.slots[i].SetItem(&item)
	}

	root.AddChild(gridContainer)
	game.ui = ebui.NewManager(root)
	return game
}

func (g *InventoryGame) startDragging(slot *InventorySlot, mouseX, mouseY float64) {
	g.draggedItem = slot.item
	g.dragStartSlot = slot
	g.dragX = mouseX - 30 // Center the item on cursor
	g.dragY = mouseY - 30
}

func (g *InventoryGame) updateDragPosition(mouseX, mouseY float64) {
	g.dragX = mouseX - 30
	g.dragY = mouseY - 30
}

func (g *InventoryGame) endDragging() {
	g.draggedItem = nil
	g.dragStartSlot = nil
}

func (g *InventoryGame) isWithinInventory(x, y float64) bool {
	pos := g.gridContainer.GetAbsolutePosition()
	size := g.gridContainer.GetSize()

	// Account for padding in bounds check
	return x >= pos.X &&
		x <= pos.X+size.Width &&
		y >= pos.Y &&
		y <= pos.Y+size.Height
}

func (g *InventoryGame) Update() error {
	return g.ui.Update()
}

func (g *InventoryGame) Draw(screen *ebiten.Image) {
	// Draw the main UI
	g.ui.Draw(screen)

	// Draw the dragged item if there is one
	if g.draggedItem != nil {
		// Create the dragged item visual
		dragImg := ebiten.NewImage(60, 60)
		dragImg.Fill(g.draggedItem.Color)

		// Set up drawing options
		op := &ebiten.DrawImageOptions{}

		// Make it more transparent when outside inventory bounds
		if !g.isWithinInventory(g.dragX+30, g.dragY+30) {
			op.ColorM.Scale(1, 1, 1, 0.4) // More transparent when outside
		} else {
			op.ColorM.Scale(1, 1, 1, 0.7) // Normal transparency when inside
		}

		op.GeoM.Translate(g.dragX, g.dragY)

		// Draw the dragged item
		screen.DrawImage(dragImg, op)

		// Draw the item name
		ebui.NewLabel(
			g.draggedItem.Name,
			ebui.WithSize(60, 20),
			ebui.WithPosition(ebui.Position{X: g.dragX, Y: g.dragY + 20}),
			ebui.WithJustify(ebui.JustifyCenter),
		).Draw(screen)
	}
}

func (g *InventoryGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 310, 310
}

func main() {
	ebiten.SetWindowSize(310, 310)
	ebiten.SetWindowTitle("EBUI Inventory Example")

	if err := ebiten.RunGame(NewInventoryGame()); err != nil {
		log.Fatal(err)
	}
}
