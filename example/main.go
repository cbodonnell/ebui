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

	root := ebui.NewVStackContainer(20, ebui.AlignCenter)
	root.SetSize(ebui.Size{
		Width:  800,
		Height: 600,
	})
	root.SetPadding(ebui.Padding{
		Top:    20,
		Right:  20,
		Bottom: 20,
		Left:   20,
	})
	root.SetBackground(color.RGBA{240, 240, 240, 255})

	// Header section
	header := ebui.NewHStackContainer(10, ebui.AlignCenter)
	header.SetSize(ebui.Size{
		Width:  760,
		Height: 60,
	})
	header.SetBackground(color.RGBA{220, 220, 220, 255})

	// Add Item button
	addBtn := ebui.NewButton("Add Item")
	addBtn.SetSize(ebui.Size{
		Width:  120,
		Height: 40,
	})
	addBtn.OnClick(func() {
		game.addItem()
	})

	// Add 5 Items button
	addMultiBtn := ebui.NewButton("Add 5 Items")
	addMultiBtn.SetSize(ebui.Size{
		Width:  120,
		Height: 40,
	})
	addMultiBtn.OnClick(func() {
		for i := 0; i < 5; i++ {
			game.addItem()
		}
	})

	// Clear All button
	clearBtn := ebui.NewButton("Clear All")
	clearBtn.SetSize(ebui.Size{
		Width:  120,
		Height: 40,
	})
	clearBtn.OnClick(func() {
		game.clearItems()
	})

	header.AddChild(addBtn)
	header.AddChild(addMultiBtn)
	header.AddChild(clearBtn)

	// Create scrollable content area
	scrollable := ebui.NewScrollableContainer(ebui.NewVerticalStack(ebui.StackConfig{
		Spacing:   10,
		Alignment: ebui.AlignStart,
	}))
	scrollable.SetSize(ebui.Size{
		Width:  760,
		Height: 480,
	})
	scrollable.SetBackground(color.RGBA{255, 255, 255, 255})
	scrollable.SetPadding(ebui.Padding{
		Top:    10,
		Right:  10,
		Bottom: 10,
		Left:   10,
	})

	game.scrollable = scrollable

	// Add some initial items
	for i := 0; i < 5; i++ {
		game.addItem()
	}

	root.AddChild(header)
	root.AddChild(scrollable)

	game.ui = ebui.NewManager(root)
	return game
}

func (g *Game) addItem() {
	priority := []string{"Low", "Medium", "High"}[rand.Intn(3)]
	status := []string{"New", "In Progress", "Done"}[rand.Intn(3)]

	// Create row container
	row := ebui.NewHStackContainer(10, ebui.AlignStart)
	row.SetSize(ebui.Size{
		Width:  740,
		Height: 50,
	})
	row.SetBackground(color.RGBA{245, 245, 245, 255})

	// ID label
	idLabel := ebui.NewButton(fmt.Sprintf("Item %d", g.nextID))
	idLabel.SetSize(ebui.Size{
		Width:  100,
		Height: 40,
	})

	// Priority label
	priorityLabel := ebui.NewButton(priority)
	priorityLabel.SetSize(ebui.Size{
		Width:  100,
		Height: 40,
	})
	// Set different colors for different priorities
	switch priority {
	case "Low":
		priorityLabel.SetBackground(color.RGBA{144, 238, 144, 255}) // Light green
	case "Medium":
		priorityLabel.SetBackground(color.RGBA{255, 218, 185, 255}) // Peach
	case "High":
		priorityLabel.SetBackground(color.RGBA{255, 182, 193, 255}) // Light pink
	}

	// Status button that cycles through states
	statusBtn := ebui.NewButton(status)
	statusBtn.SetSize(ebui.Size{
		Width:  120,
		Height: 40,
	})
	statusBtn.OnClick(func() {
		currentStatus := statusBtn.GetLabel()
		var newStatus string
		switch currentStatus {
		case "New":
			newStatus = "In Progress"
			statusBtn.SetBackground(color.RGBA{255, 218, 185, 255}) // Peach
		case "In Progress":
			newStatus = "Done"
			statusBtn.SetBackground(color.RGBA{144, 238, 144, 255}) // Light green
		case "Done":
			newStatus = "New"
			statusBtn.SetBackground(color.RGBA{200, 200, 200, 255}) // Gray
		}
		statusBtn.SetLabel(newStatus)
	})

	// Delete button
	deleteBtn := ebui.NewButton("Delete")
	deleteBtn.SetSize(ebui.Size{
		Width:  100,
		Height: 40,
	})
	deleteBtn.SetBackground(color.RGBA{255, 192, 192, 255}) // Light red
	deleteBtn.OnClick(func() {
		g.scrollable.RemoveChild(row)
	})

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
