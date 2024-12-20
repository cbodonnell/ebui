package main

import (
	"flag"
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
	windows    *ebui.WindowManager
}

func NewGame() *Game {
	game := &Game{nextID: 1}

	// Create root container
	root := ebui.NewBaseContainer(
		ebui.WithSize(800, 600),
	)

	// Main content area
	vstack := ebui.NewLayoutContainer(
		ebui.WithSize(800, 600),
		ebui.WithPadding(20, 20, 20, 20),
		ebui.WithBackground(color.RGBA{240, 240, 240, 255}),
		ebui.WithLayout(ebui.NewVerticalStackLayout(20, ebui.AlignStart)),
	)

	// Header section with buttons
	header := ebui.NewLayoutContainer(
		ebui.WithSize(760, 60),
		ebui.WithBackground(color.RGBA{220, 220, 220, 255}),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignCenter)),
	)

	// Task management buttons
	addBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabel("Add Task"),
	)
	addBtn.SetClickHandler(func() { game.addItem() })

	addMultiBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabel("Add 5 Tasks"),
	)
	addMultiBtn.SetClickHandler(func() {
		for i := 0; i < 5; i++ {
			game.addItem()
		}
	})

	clearBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabel("Clear Tasks"),
	)
	clearBtn.SetClickHandler(func() { game.clearItems() })

	// Window management buttons
	newWindowBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabel("New Window"),
	)
	newWindowBtn.SetClickHandler(func() { game.createRandomWindow() })

	header.AddChild(addBtn)
	header.AddChild(addMultiBtn)
	header.AddChild(clearBtn)
	header.AddChild(newWindowBtn)

	// Create scrollable task list
	scrollable := ebui.NewScrollableContainer(
		ebui.WithSize(760, 480),
		ebui.WithPadding(10, 10, 10, 10),
		ebui.WithBackground(color.RGBA{255, 255, 255, 255}),
		ebui.WithLayout(ebui.NewVerticalStackLayout(10, ebui.AlignStart)),
	)
	game.scrollable = scrollable

	vstack.AddChild(header)
	vstack.AddChild(scrollable)

	// Create window manager
	game.windows = ebui.NewWindowManager()

	// Create initial windows
	game.createStatsWindow()
	game.createInfoWindow()

	// Add everything to root
	root.AddChild(vstack)
	root.AddChild(game.windows)

	game.ui = ebui.NewManager(root)

	// Add some initial tasks
	for i := 0; i < 5; i++ {
		game.addItem()
	}

	return game
}

func (g *Game) createInfoWindow() {
	window := g.windows.CreateWindow(300, 200,
		ebui.WithWindowTitle("Welcome"),
		ebui.WithWindowColors(ebui.WindowColors{
			Background: color.RGBA{230, 230, 230, 255},
			Header:     color.RGBA{100, 149, 237, 255}, // Cornflower blue
			Border:     color.RGBA{100, 149, 237, 255},
		}),
	)

	infoLbl := ebui.NewLabel(
		"Try dragging this window!",
		ebui.WithSize(300, 40),
	)

	window.AddChild(infoLbl)
}

func (g *Game) createStatsWindow() {
	window := g.windows.CreateWindow(250, 180,
		ebui.WithWindowTitle("Statistics"),
		ebui.WithWindowColors(ebui.WindowColors{
			Background: color.RGBA{230, 230, 230, 255},
			Header:     color.RGBA{46, 139, 87, 255},
			Border:     color.RGBA{46, 139, 87, 255},
		}),
	)

	updateStatsBtn := ebui.NewButton(
		ebui.WithSize(250, 40),
		ebui.WithLabel(fmt.Sprintf("Tasks: %d", len(g.scrollable.GetChildren()))),
	)
	updateStatsBtn.SetClickHandler(func() {
		updateStatsBtn.SetLabel(fmt.Sprintf("Tasks: %d", len(g.scrollable.GetChildren())))
	})

	window.AddChild(updateStatsBtn)
}

func (g *Game) createRandomWindow() {
	titles := []string{"Notes", "Settings", "Help", "About", "Debug"}
	colors := []ebui.WindowColors{
		{
			Background: color.RGBA{230, 230, 230, 255},
			Header:     color.RGBA{218, 112, 214, 255}, // Orchid
			Border:     color.RGBA{218, 112, 214, 255},
		},
		{
			Background: color.RGBA{230, 230, 230, 255},
			Header:     color.RGBA{210, 105, 30, 255}, // Chocolate
			Border:     color.RGBA{210, 105, 30, 255},
		},
	}

	title := titles[rand.Intn(len(titles))]
	colorScheme := colors[rand.Intn(len(colors))]

	window := g.windows.CreateWindow(
		200+rand.Float64()*100,
		150+rand.Float64()*100,
		ebui.WithWindowTitle(title),
		ebui.WithWindowColors(colorScheme),
	)

	content := ebui.NewScrollableContainer(
		ebui.WithSize(window.GetSize().Width, window.GetSize().Height-30), // Subtract header height
		ebui.WithPadding(0, 4, 0, 4),
		ebui.WithLayout(ebui.NewVerticalStackLayout(0, ebui.AlignStart)),
	)

	for i := 0; i < 20; i++ {
		lbl := ebui.NewLabel(
			fmt.Sprintf("Item %d", i),
			ebui.WithSize(60, 20),
		)
		content.AddChild(lbl)
	}

	window.AddChild(content)
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
	)
	statusBtn.SetClickHandler(func() {
		currentStatus := statusBtn.GetLabel()
		var newStatus string
		switch currentStatus {
		case "New":
			newStatus = "In Progress"
		case "In Progress":
			newStatus = "Done"
		case "Done":
			newStatus = "New"
		}
		statusBtn.SetLabel(newStatus)
	})

	// Delete button
	deleteBtn := ebui.NewButton(
		ebui.WithSize(100, 40),
		ebui.WithLabel("Delete"),
		ebui.WithBackground(color.RGBA{255, 192, 192, 255}), // Light red
	)
	deleteBtn.SetClickHandler(func() {
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
	ebiten.SetWindowTitle("EBUI Framework Demo")

	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *debug {
		ebui.Debug = true
	}

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
