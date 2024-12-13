package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/cbodonnell/ebui"
	"github.com/hajimehoshi/ebiten/v2"
)

var _ ebiten.Game = &Game{}

type Game struct {
	ui *ebui.Manager
}

func NewGame() *Game {
	ebui.Debug = true

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

	// Add header buttons
	for i := 1; i <= 3; i++ {
		btn := ebui.NewButton(fmt.Sprintf("Action %d", i))
		btn.SetSize(ebui.Size{
			Width:  120,
			Height: 40,
		})
		btn.OnClick(func(i int) func() {
			return func() {
				log.Printf("Action %d clicked", i)
			}
		}(i))
		header.AddChild(btn)
	}

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

	// Add items to scrollable container
	for i := 1; i <= 10; i++ {
		// Create row container
		row := ebui.NewHStackContainer(10, ebui.AlignStart)
		row.SetSize(ebui.Size{
			Width:  740,
			Height: 50,
		})
		row.SetBackground(color.RGBA{245, 245, 245, 255})

		// Add item label
		label := ebui.NewButton(fmt.Sprintf("Item %d", i))
		label.SetSize(ebui.Size{
			Width:  200,
			Height: 40,
		})
		label.OnClick(func(i int) func() {
			return func() {
				log.Printf("Item %d selected", i)
			}
		}(i))

		// Add action buttons
		editBtn := ebui.NewButton("Edit")
		editBtn.SetSize(ebui.Size{
			Width:  100,
			Height: 40,
		})
		editBtn.OnClick(func(i int) func() {
			return func() {
				log.Printf("Edit item %d", i)
			}
		}(i))

		deleteBtn := ebui.NewButton("Delete")
		deleteBtn.SetSize(ebui.Size{
			Width:  100,
			Height: 40,
		})
		deleteBtn.OnClick(func(i int) func() {
			return func() {
				log.Printf("Delete item %d", i)
			}
		}(i))

		// Add all buttons to row
		row.AddChild(label)
		row.AddChild(editBtn)
		row.AddChild(deleteBtn)

		// Add row to scrollable container
		scrollable.AddChild(row)
	}

	root.AddChild(header)
	root.AddChild(scrollable)

	ui := ebui.NewManager(root)

	return &Game{ui: ui}
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
	ebiten.SetWindowTitle("EBUI Demo")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
