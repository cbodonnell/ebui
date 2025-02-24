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

type TaskData struct {
	ID       int
	Priority string
	Status   string
}

type Game struct {
	ui         *ebui.Manager
	navStack   *ebui.StackContainer
	scrollable *ebui.ScrollableContainer
	nextID     int
	windows    *ebui.WindowManager
	tasks      map[int]*TaskData
}

func NewGame() *Game {
	game := &Game{
		nextID: 1,
		tasks:  make(map[int]*TaskData),
	}

	// Create root container
	root := ebui.NewBaseContainer(
		ebui.WithSize(800, 600),
	)

	// Create navigation stack with same size as root
	game.navStack = ebui.NewStackContainer(
		ebui.WithSize(800, 600),
	)

	// Push main view first - this will place it without animation
	mainView := game.createMainView()
	game.navStack.Push(mainView)

	// Add navigation stack to root
	root.AddChild(game.navStack)

	// Create window manager
	game.windows = ebui.NewWindowManager(
		ebui.WithSize(800, 600),
	)
	root.AddChild(game.windows)

	// Create initial windows
	game.createStatsWindow()
	game.createInfoWindow()

	game.ui = ebui.NewManager(root)

	// Add some initial tasks
	for i := 0; i < 5; i++ {
		game.addItem()
	}

	return game
}

func (g *Game) createMainView() *ebui.LayoutContainer {
	mainView := ebui.NewLayoutContainer(
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

	input := ebui.NewTextInput(
		ebui.WithSize(200, 30),
		ebui.WithTextInputColors(ebui.DefaultTextInputColors()),
		ebui.WithInitialText("Hello"),
		ebui.WithChangeHandler(func(text string) {
			fmt.Printf("Text changed: %s\n", text)
		}),
		ebui.WithSubmitHandler(func(text string) {
			fmt.Printf("Text submitted: %s\n", text)
		}),
	)

	// Task management buttons
	addBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabelText("Add Task"),
	)
	addBtn.SetClickHandler(func() { g.addItem() })

	addMultiBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabelText("Add 5 Tasks"),
	)
	addMultiBtn.SetClickHandler(func() {
		for i := 0; i < 5; i++ {
			g.addItem()
		}
	})

	clearBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabelText("Clear Tasks"),
	)
	clearBtn.SetClickHandler(func() { g.clearItems() })

	header.AddChild(input)
	header.AddChild(addBtn)
	header.AddChild(addMultiBtn)
	header.AddChild(clearBtn)

	// Create scrollable task list
	scrollable := ebui.NewScrollableContainer(
		ebui.WithSize(760, 480),
		ebui.WithPadding(10, 10, 10, 10),
		ebui.WithBackground(color.RGBA{255, 255, 255, 255}),
		ebui.WithLayout(ebui.NewVerticalStackLayout(10, ebui.AlignStart)),
	)
	g.scrollable = scrollable

	vstack.AddChild(header)
	vstack.AddChild(scrollable)
	mainView.AddChild(vstack)

	return mainView
}

func (g *Game) createTaskDetailView(task *TaskData) *ebui.LayoutContainer {
	detailView := ebui.NewLayoutContainer(
		ebui.WithSize(800, 600),
		ebui.WithBackground(color.RGBA{240, 240, 240, 255}),
		ebui.WithPadding(20, 20, 20, 20),
		ebui.WithLayout(ebui.NewVerticalStackLayout(20, ebui.AlignStart)),
	)

	// Header
	header := ebui.NewLayoutContainer(
		ebui.WithSize(760, 60),
		ebui.WithBackground(color.RGBA{220, 220, 220, 255}),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignCenter)),
	)

	// Back button
	backBtn := ebui.NewButton(
		ebui.WithSize(100, 40),
		ebui.WithLabelText("< Back"),
		ebui.WithButtonColors(ebui.ButtonColors{
			Default:     color.RGBA{100, 149, 237, 255}, // Cornflower blue
			Hovered:     color.RGBA{120, 169, 255, 255},
			Pressed:     color.RGBA{80, 129, 217, 255},
			FocusBorder: color.Black,
		}),
	)
	backBtn.SetClickHandler(func() {
		g.navStack.Pop()
	})

	titleLabel := ebui.NewLabel(
		fmt.Sprintf("Task %d Details", task.ID),
		ebui.WithSize(200, 40),
		ebui.WithJustify(ebui.JustifyCenter),
	)

	header.AddChild(backBtn)
	header.AddChild(titleLabel)

	// Task details form
	form := ebui.NewLayoutContainer(
		ebui.WithSize(760, 400),
		ebui.WithBackground(color.RGBA{255, 255, 255, 255}),
		ebui.WithPadding(20, 20, 20, 20),
		ebui.WithLayout(ebui.NewVerticalStackLayout(20, ebui.AlignStart)),
	)

	// Priority selection
	priorityLabel := ebui.NewLabel(
		"Priority:",
		ebui.WithSize(720, 30),
	)

	priorityContainer := ebui.NewLayoutContainer(
		ebui.WithSize(720, 40),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignStart)),
	)

	priorities := []struct {
		label string
		color ebui.ButtonColors
	}{
		{
			label: "Low",
			color: ebui.ButtonColors{
				Default:     color.RGBA{144, 238, 144, 255}, // Light green
				Hovered:     color.RGBA{152, 251, 152, 255},
				Pressed:     color.RGBA{50, 205, 50, 255},
				FocusBorder: color.Black,
			},
		},
		{
			label: "Medium",
			color: ebui.ButtonColors{
				Default:     color.RGBA{255, 218, 185, 255}, // Light orange
				Hovered:     color.RGBA{255, 228, 196, 255},
				Pressed:     color.RGBA{255, 165, 0, 255},
				FocusBorder: color.Black,
			},
		},
		{
			label: "High",
			color: ebui.ButtonColors{
				Default:     color.RGBA{255, 182, 193, 255}, // Light pink
				Hovered:     color.RGBA{255, 192, 203, 255},
				Pressed:     color.RGBA{255, 105, 180, 255},
				FocusBorder: color.Black,
			},
		},
	}

	for _, p := range priorities {
		btn := ebui.NewButton(
			ebui.WithSize(100, 40),
			ebui.WithLabelText(p.label),
			ebui.WithButtonColors(p.color),
		)

		// Highlight if selected
		if p.label == task.Priority {
			// Make the button darker when selected
			btn.SetBackground(p.color.Pressed)
		}

		p := p // Capture for closure
		btn.SetClickHandler(func() {
			task.Priority = p.label
			g.updateTaskRow(task)
			g.navStack.Pop()
		})
		priorityContainer.AddChild(btn)
	}

	// Status selection
	statusLabel := ebui.NewLabel(
		"Status:",
		ebui.WithSize(720, 30),
	)

	statusContainer := ebui.NewLayoutContainer(
		ebui.WithSize(720, 40),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignStart)),
	)

	statuses := []struct {
		label string
		color ebui.ButtonColors
	}{
		{
			label: "New",
			color: ebui.ButtonColors{
				Default:     color.RGBA{173, 216, 230, 255}, // Light blue
				Hovered:     color.RGBA{187, 222, 251, 255},
				Pressed:     color.RGBA{100, 149, 237, 255},
				FocusBorder: color.Black,
			},
		},
		{
			label: "In Progress",
			color: ebui.ButtonColors{
				Default:     color.RGBA{255, 218, 185, 255}, // Peach
				Hovered:     color.RGBA{255, 228, 196, 255},
				Pressed:     color.RGBA{255, 165, 0, 255},
				FocusBorder: color.Black,
			},
		},
		{
			label: "Done",
			color: ebui.ButtonColors{
				Default:     color.RGBA{152, 251, 152, 255}, // Pale green
				Hovered:     color.RGBA{162, 255, 162, 255},
				Pressed:     color.RGBA{50, 205, 50, 255},
				FocusBorder: color.Black,
			},
		},
	}

	for _, s := range statuses {
		btn := ebui.NewButton(
			ebui.WithSize(120, 40),
			ebui.WithLabelText(s.label),
			ebui.WithButtonColors(s.color),
		)

		// Highlight if selected
		if s.label == task.Status {
			// Make the button darker when selected
			btn.SetBackground(s.color.Pressed)
		}

		s := s // Capture for closure
		btn.SetClickHandler(func() {
			task.Status = s.label
			g.updateTaskRow(task)
			g.navStack.Pop()
		})
		statusContainer.AddChild(btn)
	}

	form.AddChild(priorityLabel)
	form.AddChild(priorityContainer)
	form.AddChild(statusLabel)
	form.AddChild(statusContainer)

	detailView.AddChild(header)
	detailView.AddChild(form)

	return detailView
}

func (g *Game) updateTaskRow(task *TaskData) {
	// Find and update the task row
	for _, child := range g.scrollable.GetChildren() {
		if row, ok := child.(*ebui.LayoutContainer); ok {
			// Check if this is the right row by looking at the ID label
			if idLabel := row.GetChildren()[0].(*ebui.Label); idLabel != nil {
				if idLabel.GetText() == fmt.Sprintf("Item %d", task.ID) {
					// Update priority and status labels
					priorityContainer := row.GetChildren()[1].(*ebui.LayoutContainer)
					priorityLabel := priorityContainer.GetChildren()[0].(*ebui.Label)
					statusContainer := row.GetChildren()[2].(*ebui.LayoutContainer)
					statusLabel := statusContainer.GetChildren()[0].(*ebui.Label)

					// Update priority text and color
					priorityLabel.SetText(task.Priority)
					var priorityColor color.Color
					switch task.Priority {
					case "Low":
						priorityColor = color.RGBA{144, 238, 144, 255} // Light green
					case "Medium":
						priorityColor = color.RGBA{255, 218, 185, 255} // Light orange
					case "High":
						priorityColor = color.RGBA{255, 182, 193, 255} // Light pink
					}
					priorityContainer.SetBackground(priorityColor)

					// Update status text and color
					statusLabel.SetText(task.Status)
					var statusColor color.Color
					switch task.Status {
					case "New":
						statusColor = color.RGBA{173, 216, 230, 255} // Light blue
					case "In Progress":
						statusColor = color.RGBA{255, 218, 185, 255} // Peach
					case "Done":
						statusColor = color.RGBA{152, 251, 152, 255} // Pale green
					}
					statusContainer.SetBackground(statusColor)

					return
				}
			}
		}
	}
}

func (g *Game) addItem() {
	i := rand.Intn(3)
	priority := []string{"Low", "Medium", "High"}[i]

	statusIndex := rand.Intn(3)
	status := []string{"New", "In Progress", "Done"}[statusIndex]

	taskData := &TaskData{
		ID:       g.nextID,
		Priority: priority,
		Status:   status,
	}
	g.tasks[g.nextID] = taskData

	// Priority tag colors
	priorityColors := []color.Color{
		color.RGBA{144, 238, 144, 255}, // Light green for Low
		color.RGBA{255, 218, 185, 255}, // Light orange for Medium
		color.RGBA{255, 182, 193, 255}, // Light pink for High
	}[i]

	// Status tag colors
	statusColors := []color.Color{
		color.RGBA{173, 216, 230, 255}, // Light blue for New
		color.RGBA{255, 218, 185, 255}, // Peach for In Progress
		color.RGBA{152, 251, 152, 255}, // Pale green for Done
	}[statusIndex]

	// Create row container
	row := ebui.NewLayoutContainer(
		ebui.WithSize(740, 50),
		ebui.WithBackground(color.RGBA{245, 245, 245, 255}),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignStart)),
	)

	// ID label
	idLabel := ebui.NewLabel(
		fmt.Sprintf("Item %d", taskData.ID),
		ebui.WithSize(100, 40),
		ebui.WithJustify(ebui.JustifyCenter),
	)

	// Priority label with background
	priorityContainer := ebui.NewLayoutContainer(
		ebui.WithSize(100, 40),
		ebui.WithBackground(priorityColors),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(0, ebui.AlignCenter)),
	)
	priorityLabel := ebui.NewLabel(
		priority,
		ebui.WithSize(100, 40),
		ebui.WithJustify(ebui.JustifyCenter),
	)
	priorityContainer.AddChild(priorityLabel)

	// Status label with background
	statusContainer := ebui.NewLayoutContainer(
		ebui.WithSize(120, 40),
		ebui.WithBackground(statusColors),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(0, ebui.AlignCenter)),
	)
	statusLabel := ebui.NewLabel(
		status,
		ebui.WithSize(120, 40),
		ebui.WithJustify(ebui.JustifyCenter),
	)
	statusContainer.AddChild(statusLabel)

	// Edit button
	editBtn := ebui.NewButton(
		ebui.WithSize(100, 40),
		ebui.WithLabelText("Edit"),
		ebui.WithButtonColors(ebui.ButtonColors{
			Default:     color.RGBA{100, 149, 237, 255}, // Cornflower blue
			Hovered:     color.RGBA{120, 169, 255, 255},
			Pressed:     color.RGBA{80, 129, 217, 255},
			FocusBorder: color.Black,
		}),
	)
	editBtn.SetClickHandler(func() {
		detailView := g.createTaskDetailView(taskData)
		g.navStack.Push(detailView)
	})

	// Delete button
	deleteBtn := ebui.NewButton(
		ebui.WithSize(100, 40),
		ebui.WithLabelText("Delete"),
		ebui.WithButtonColors(ebui.ButtonColors{
			Default:     color.RGBA{255, 99, 71, 255},
			Hovered:     color.RGBA{255, 69, 0, 255},
			Pressed:     color.RGBA{178, 34, 34, 255},
			FocusBorder: color.Black,
		}),
	)
	deleteBtn.SetClickHandler(func() {
		g.scrollable.RemoveChild(row)
		delete(g.tasks, taskData.ID)
	})

	// Add all elements to row
	row.AddChild(idLabel)
	row.AddChild(priorityContainer)
	row.AddChild(statusContainer)
	row.AddChild(editBtn)
	row.AddChild(deleteBtn)

	// Add row to scrollable container
	g.scrollable.AddChild(row)
	g.nextID++
}

func (g *Game) clearItems() {
	for len(g.scrollable.GetChildren()) > 0 {
		g.scrollable.RemoveChild(g.scrollable.GetChildren()[0])
	}
	g.tasks = make(map[int]*TaskData)
}

func (g *Game) createStatsWindow() {
	window := g.windows.CreateWindow(250, 180,
		ebui.WithWindowPosition(500, 100),
		ebui.WithWindowTitle("Statistics"),
		ebui.WithWindowColors(ebui.WindowColors{
			Background: color.RGBA{230, 230, 230, 255},
			Header:     color.RGBA{46, 139, 87, 255},
			HeaderText: color.Black,
			Border:     color.RGBA{46, 139, 87, 255},
		}),
	)

	updateStatsBtn := ebui.NewButton(
		ebui.WithSize(250, 40),
		ebui.WithLabelText(fmt.Sprintf("Tasks: %d", len(g.scrollable.GetChildren()))),
	)
	updateStatsBtn.SetClickHandler(func() {
		updateStatsBtn.SetLabel(fmt.Sprintf("Tasks: %d", len(g.scrollable.GetChildren())))
	})

	window.AddChild(updateStatsBtn)
}

func (g *Game) createInfoWindow() {
	window := g.windows.CreateWindow(300, 200,
		ebui.WithWindowPosition(100, 100),
		ebui.WithWindowTitle("Welcome"),
		ebui.WithWindowColors(ebui.WindowColors{
			Background: color.RGBA{230, 230, 230, 255},
			Header:     color.RGBA{100, 149, 237, 255}, // Cornflower blue
			HeaderText: color.Black,
			Border:     color.RGBA{100, 149, 237, 255},
		}),
	)

	infoLbl := ebui.NewLabel(
		"Click any task to view details!",
		ebui.WithSize(300, 40),
		ebui.WithJustify(ebui.JustifyCenter),
	)

	window.AddChild(infoLbl)
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
	ebiten.SetWindowTitle("EBUI Tasks Example")

	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *debug {
		ebui.Debug = true
	}

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
