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

	root := ebui.NewVStackContainer(10, ebui.AlignCenter)

	// Set up the root container with full window size
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
	root.SetBackground(color.RGBA{245, 245, 245, 255})

	// Create a horizontal stack for buttons - centered in root container
	buttonContainer := ebui.NewHStackContainer(10, ebui.AlignCenter)
	buttonContainer.SetSize(ebui.Size{
		Width:  760, // Container width - padding
		Height: 100, // Fixed height
	})
	buttonContainer.SetBackground(color.RGBA{235, 235, 235, 255})

	// Add buttons
	for i := 1; i <= 3; i++ {
		button := ebui.NewButton(fmt.Sprintf("Button %d", i))
		button.SetSize(ebui.Size{
			Width:  150,
			Height: 50,
		})
		button.OnClick(func() {
			log.Printf("Button %d clicked!\n", i)
		})
		buttonContainer.AddChild(button)
	}

	root.AddChild(buttonContainer)

	ui := ebui.NewManager(root)

	return &Game{
		ui: ui,
	}
}

func (g *Game) Update() error {
	g.ui.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("EBUI Example")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
