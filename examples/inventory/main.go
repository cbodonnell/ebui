package main

import (
	"flag"
	"image/color"
	"log"

	"github.com/cbodonnell/ebui"
	"github.com/hajimehoshi/ebiten/v2"
)

type InventoryGame struct {
	ui *ebui.Manager
}

func NewInventoryGame() *InventoryGame {
	game := &InventoryGame{}

	// Create root container
	root := ebui.NewBaseContainer()

	// Create window manager
	wm := ebui.NewWindowManager()

	// Create the window
	window := wm.CreateWindow(318, 336,
		ebui.WithWindowPosition(160, 160),
		ebui.WithWindowTitle("Inventory"),
		ebui.WithWindowColors(ebui.WindowColors{
			Background: color.RGBA{230, 230, 230, 255},
			Header:     color.RGBA{46, 139, 87, 255},
			Border:     color.RGBA{46, 139, 87, 255},
		}),
	)

	// Create the inventory component
	inv := NewInventory(
		ebui.WithSize(318, 306),
		ebui.WithPosition(ebui.Position{X: 10, Y: 10}),
		ebui.WithBackground(color.RGBA{255, 255, 255, 255}),
		ebui.WithPadding(10, 10, 10, 10),
	)

	// Add items to the inventory
	sampleItems := []Item{
		{Name: "Sword", Color: color.RGBA{255, 0, 0, 255}},
		{Name: "Shield", Color: color.RGBA{0, 0, 255, 255}},
		{Name: "Potion", Color: color.RGBA{0, 255, 0, 255}},
		{Name: "Bow", Color: color.RGBA{139, 69, 19, 255}},
		{Name: "Arrow", Color: color.RGBA{128, 128, 128, 255}},
		{Name: "Gem", Color: color.RGBA{147, 112, 219, 255}},
		{Name: "Ring", Color: color.RGBA{255, 215, 0, 255}},
		{Name: "Staff", Color: color.RGBA{65, 105, 225, 255}},
	}
	for i := 0; i < len(sampleItems); i++ {
		item := sampleItems[i]
		inv.slots[i].SetItem(&item)
	}

	// Add inventory component to window
	window.AddChild(inv)

	// Add window manager to root
	root.AddChild(wm)

	// Create the UI manager
	game.ui = ebui.NewManager(root)

	return game
}

// Inventory is a custom component implementing a scrollable inventory.
// It embeds the ScrollableContainer component and adds inventory-specific
// functionality.
type Inventory struct {
	*ebui.ScrollableContainer
	slots         []*InventorySlot
	draggedItem   *Item
	dragStartSlot *InventorySlot
	dragX, dragY  float64
}

func WithNumSlots(n int) ebui.ComponentOpt {
	return func(c ebui.Component) {
		if inv, ok := c.(*Inventory); ok {
			inv.slots = make([]*InventorySlot, n)
		}
	}
}

func NewInventory(opts ...ebui.ComponentOpt) *Inventory {
	// Create the inventory as a scrollable container
	inv := &Inventory{
		ScrollableContainer: ebui.NewScrollableContainer(opts...),
		slots:               make([]*InventorySlot, 32), // Default to 32 slots
	}

	for _, opt := range opts {
		opt(inv)
	}

	// Create a vertical container to hold the inventory rows
	gridContainer := ebui.NewLayoutContainer(
		ebui.WithSize(286, 582), // 4x8 grid of 64x64 slots with 10px padding
		ebui.WithBackground(color.RGBA{255, 255, 255, 255}),
		ebui.WithLayout(ebui.NewVerticalStackLayout(10, ebui.AlignStart)),
	)

	var rows []*ebui.LayoutContainer
	var row *ebui.LayoutContainer
	for i := 0; i < len(inv.slots); i++ {
		if i%4 == 0 {
			row = ebui.NewLayoutContainer(
				ebui.WithSize(286, 64), // each slot is 64x64 with 10px spacing
				ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignStart)),
			)
			rows = append(rows, row)
		}

		inv.slots[i] = NewInventorySlot(inv)
		row.AddChild(inv.slots[i])
	}
	for _, r := range rows {
		gridContainer.AddChild(r)
	}

	// Add grid container to scrollable container
	inv.AddChild(gridContainer)

	return inv
}

// Update isWithinInventory to account for scrolling
func (inv *Inventory) isWithinInventory(x, y float64) bool {
	return inv.BaseComponent.Contains(x, y)
}

func (inv *Inventory) startDragging(slot *InventorySlot, mouseX, mouseY float64) {
	inv.draggedItem = slot.item
	inv.dragStartSlot = slot
	inv.dragX = mouseX - 32 // Center the item on cursor
	inv.dragY = mouseY - 32
}

func (inv *Inventory) updateDragPosition(mouseX, mouseY float64) {
	inv.dragX = mouseX - 32
	inv.dragY = mouseY - 32
}

func (inv *Inventory) endDragging() {
	inv.draggedItem = nil
	inv.dragStartSlot = nil
}

func (inv *Inventory) Draw(screen *ebiten.Image) {
	// Draw the scrollable container
	inv.ScrollableContainer.Draw(screen)

	// Draw the dragged item if there is one
	if inv.draggedItem != nil {
		inv.drawDraggedItem(screen)
	}
}

func (inv *Inventory) drawDraggedItem(screen *ebiten.Image) {
	// Draw the dragged item with a semi-transparent background
	transparency := 0.7
	if !inv.isWithinInventory(inv.dragX+32, inv.dragY+32) {
		// If the item is outside the inventory bounds, increase transparency
		transparency = 0.4
	}

	r, g, b, a := inv.draggedItem.Color.RGBA()
	background := color.RGBA{
		uint8(float64(r>>8) * transparency),
		uint8(float64(g>>8) * transparency),
		uint8(float64(b>>8) * transparency),
		uint8(float64(a>>8) * transparency),
	}

	draggedSlot := ebui.NewLayoutContainer(
		ebui.WithSize(64, 64),
		ebui.WithPosition(ebui.Position{X: inv.dragX, Y: inv.dragY}),
		ebui.WithBackground(background),
	)

	// Add the item name label
	nameLbl := ebui.NewLabel(
		inv.draggedItem.Name,
		ebui.WithSize(64, 64),
		ebui.WithPosition(ebui.Position{X: inv.dragX, Y: inv.dragY}),
		ebui.WithJustify(ebui.JustifyCenter),
	)
	draggedSlot.AddChild(nameLbl)

	// Draw the dragged item
	draggedSlot.Draw(screen)
}

// InventorySlot is a custom component representing a slot in an inventory.
// It embeds the LayoutContainer and BaseFocusable components and adds
// inventory slot-specific functionality.
type InventorySlot struct {
	*ebui.LayoutContainer
	*ebui.BaseFocusable
	item      *Item
	label     *ebui.Label
	isHovered bool
	inv       *Inventory
}

type Item struct {
	Name  string
	Color color.Color
}

func NewInventorySlot(inv *Inventory) *InventorySlot {
	slot := &InventorySlot{
		LayoutContainer: ebui.NewLayoutContainer(
			ebui.WithSize(64, 64),
			ebui.WithBackground(color.RGBA{200, 200, 200, 255}),
			ebui.WithLayout(ebui.NewVerticalStackLayout(0, ebui.AlignCenter)),
		),
		BaseFocusable: ebui.NewBaseFocusable(),
		inv:           inv,
	}

	slot.label = ebui.NewLabel(
		"",
		ebui.WithSize(64, 64),
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
			s.inv.startDragging(s, e.MouseX, e.MouseY)
		}
	})

	s.AddEventListener(ebui.Drag, func(e *ebui.Event) {
		s.inv.updateDragPosition(e.MouseX, e.MouseY)
	})

	s.AddEventListener(ebui.Drop, func(e *ebui.Event) {
		if sourceSlot, ok := e.RelatedTarget.(*InventorySlot); ok {
			// Swap items between slots
			s.item, sourceSlot.item = sourceSlot.item, s.item
			s.updateDisplay()
			sourceSlot.updateDisplay()
			s.inv.endDragging()
		}
	})

	s.AddEventListener(ebui.DragEnd, func(e *ebui.Event) {
		if s.inv.draggedItem != nil {
			// Check if the mouse is within the inventory bounds
			if s.inv.isWithinInventory(e.MouseX, e.MouseY) {
				// Do nothing, item will return to original slot
				s.inv.endDragging()
			} else {
				// Remove the item from the source slot
				s.inv.dragStartSlot.item = nil
				s.inv.dragStartSlot.updateDisplay()
				s.inv.endDragging()
			}
		}
	})
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

func (g *InventoryGame) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

func (g *InventoryGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 638, 656
}

func (g *InventoryGame) Update() error {
	return g.ui.Update()
}

func main() {
	ebiten.SetWindowSize(638, 656)
	ebiten.SetWindowTitle("EBUI Inventory Example")

	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *debug {
		ebui.Debug = true
	}

	if err := ebiten.RunGame(NewInventoryGame()); err != nil {
		log.Fatal(err)
	}
}
