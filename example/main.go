package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/cbodonnell/ebui"
	"github.com/hajimehoshi/ebiten/v2"
)

var _ ebiten.Game = &Game{}

type Game struct {
	ui         *ebui.Manager
	scrollable *ebui.ScrollableContainer
	nextID     int
}

func NewGame() *Game {
	// ebui.Debug = true
	game := &Game{nextID: 1}

	root := ebui.NewBaseContainer(
		ebui.WithSize(800, 600),
	)

	vstack := ebui.NewLayoutContainer(
		ebui.WithSize(800, 600),
		ebui.WithPadding(20, 20, 20, 20),
		ebui.WithBackground(color.RGBA{240, 240, 240, 255}),
		ebui.WithLayout(ebui.NewVerticalStackLayout(20, ebui.AlignStart)),
	)

	// Header section
	header := ebui.NewLayoutContainer(
		ebui.WithSize(760, 60),
		ebui.WithBackground(color.RGBA{220, 220, 220, 255}),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignCenter)),
	)

	// Add Item button
	addBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabel("Add Item"),
		ebui.WithClickHandler(func(e ebui.Event) { game.addItem() }),
	)

	// Add 5 Items button
	addMultiBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabel("Add 5 Items"),
		ebui.WithClickHandler(func(e ebui.Event) {
			for i := 0; i < 5; i++ {
				game.addItem()
			}
		}),
	)

	// Clear All button
	clearBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabel("Clear All"),
		ebui.WithClickHandler(func(e ebui.Event) { game.clearItems() }),
	)

	header.AddChild(addBtn)
	header.AddChild(addMultiBtn)
	header.AddChild(clearBtn)

	// Create scrollable content area
	scrollable := ebui.NewScrollableContainer(
		ebui.WithSize(760, 480),
		ebui.WithPadding(10, 10, 10, 10),
		ebui.WithBackground(color.RGBA{255, 255, 255, 255}),
		ebui.WithLayout(ebui.NewVerticalStackLayout(10, ebui.AlignStart)),
	)

	game.scrollable = scrollable

	vstack.AddChild(header)
	vstack.AddChild(scrollable)

	floaters := ebui.NewBaseContainer()

	floatBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithPosition(60, 60),
		ebui.WithLabel("Floater"),
		ebui.WithClickHandler(func(e ebui.Event) {
			fmt.Println("Floating button clicked!")
		}),
	)

	floaters.AddChild(floatBtn)

	root.AddChild(vstack)
	root.AddChild(floaters)

	game.ui = ebui.NewManager(root)

	// Add some initial items
	for i := 0; i < 5; i++ {
		game.addItem()
	}

	return game
}

func (g *Game) addItem() {
	i := rand.Intn(3)
	priority := []string{"Low", "Medium", "High"}[i]
	status := []string{"New", "In Progress", "Done"}[rand.Intn(3)]
	background := []color.RGBA{
		{144, 238, 144, 255}, // Light green
		{255, 218, 185, 255}, // Peach
		{255, 182, 193, 255}, // Light pink
	}[i]

	// Create row container
	row := ebui.NewLayoutContainer(
		ebui.WithSize(740, 50),
		ebui.WithBackground(color.RGBA{245, 245, 245, 255}),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignStart)),
	)

	// ID label
	idLabel := ebui.NewButton(
		ebui.WithSize(100, 40),
		ebui.WithLabel(fmt.Sprintf("Item %d", g.nextID)),
	)

	// Priority label
	priorityLabel := ebui.NewButton(
		ebui.WithSize(100, 40),
		ebui.WithLabel(priority),
		ebui.WithBackground(background),
	)

	// Status button that cycles through states
	statusBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabel(status),
		ebui.WithClickHandler(func(e ebui.Event) {
			b := e.Component.(*ebui.Button)
			currentStatus := b.GetLabel()
			var newStatus string
			switch currentStatus {
			case "New":
				newStatus = "In Progress"
			case "In Progress":
				newStatus = "Done"
			case "Done":
				newStatus = "New"
			}
			b.SetLabel(newStatus)
		}),
	)

	// Delete button
	deleteBtn := ebui.NewButton(
		ebui.WithSize(100, 40),
		ebui.WithLabel("Delete"),
		ebui.WithBackground(color.RGBA{255, 192, 192, 255}), // Light red
		ebui.WithClickHandler(func(e ebui.Event) {
			g.scrollable.RemoveChild(row)
		}),
	)

	// Add all elements to row
	row.AddChild(idLabel)
	row.AddChild(priorityLabel)
	row.AddChild(statusBtn)
	row.AddChild(deleteBtn)

	// Add row to scrollable container
	g.scrollable.AddChild(row)
	g.nextID++
}

func (g *Game) clearItems() {
	for len(g.scrollable.GetChildren()) > 0 {
		g.scrollable.RemoveChild(g.scrollable.GetChildren()[0])
	}
}

func (g *Game) Update() error {
	return g.ui.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("EBUI Task Manager Demo")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
