# EBUI - Ebitengine UI Framework

EBUI is a UI framework for [Ebitengine](https://ebitengine.org/), providing a set of components and layouts for building game interfaces and tools.

## Features

- **Component-Based Architecture**: Build UIs using reusable, composable components
- **Flexible Layout System**: Arrange components using various layout strategies
- **Event System**: Handle user interactions with a flexible event system
- **Component Library**:
  - Labels with text alignment options
  - Buttons with customizable colors and states
  - Text inputs with selection and clipboard support
  - Scrollable content containers
  - Windows with drag-and-drop functionality

## Installation

```bash
go get github.com/cbodonnell/ebui
```

## Quick Start

```go
package main

import (
	"image/color"
	"log"

	"github.com/cbodonnell/ebui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	ui *ebui.Manager
}

func NewGame() *Game {
	// Create a root container
	root := ebui.NewLayoutContainer(
		ebui.WithSize(800, 600),
		ebui.WithLayout(ebui.NewVerticalStackLayout(0, ebui.AlignCenter)),
	)

	// Create a label
	label := ebui.NewLabel(
		"Welcome to EBUI!",
		ebui.WithSize(800, 40),
		ebui.WithColor(color.White),
	)

	// Create a button
	button := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabelText("Click Me"),
		ebui.WithClickHandler(func() {
			log.Println("Button clicked!")
		}),
	)

	// Create a container to center the button
	buttonContainer := ebui.NewLayoutContainer(
		ebui.WithSize(800, 40),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(0, ebui.AlignCenter)),
	)

	// Add button to button container
	buttonContainer.AddChild(button)

	// Add label and button container to root
	root.AddChild(label)
	root.AddChild(buttonContainer)

	return &Game{
		// Create a UI manager using the root component
		ui: ebui.NewManager(root),
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
	ebiten.SetWindowTitle("EBUI Example")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
```

More comprehensive examples can be found in the [examples](examples) directory.

## Core Concepts

### Components

Components are the building blocks of EBUI. Every UI element implements the `Component` interface:

```go
type Component interface {
    Identifiable
    EbitenLifecycle
    SetPosition(pos Position)
    GetPosition() Position
    SetSize(size Size)
    GetSize() Size
    SetParent(parent Container)
    GetParent() Container
    SetPadding(padding Padding)
    GetPadding() Padding
    Contains(x, y float64) bool
    GetAbsolutePosition() Position
}
```

### Containers

Containers are components that can hold other components. EBUI provides several types of containers:

- **BaseContainer**: Basic container with no layout management
- **LayoutContainer**: Container that arranges children using a layout strategy
- **ScrollableContainer**: Container with scrolling capability
- **ZIndexedContainer**: Container that manages component layering
- **WindowManager**: Special container for managing multiple windows

### Layouts

EBUI provides flexible layout options:

- **Vertical Stack Layout**: Arrange components vertically
- **Horizontal Stack Layout**: Arrange components horizontally
- Custom layouts can be implemented by implementing the `Layout` interface

### Event System

The event system supports:
- Mouse events (click, hover, drag)
- Keyboard input
- Focus management
- Event bubbling and capturing

## Components

### Label

```go
label := ebui.NewLabel(
    "Hello World",
    ebui.WithSize(200, 40),
    ebui.WithJustify(ebui.JustifyCenter),
    ebui.WithColor(color.Black),
    ebui.WithFont(basicfont.Face7x13), // Default font
)
```

### Button

```go
button := ebui.NewButton(
    ebui.WithSize(120, 40),
    ebui.WithLabelText("Click Me"),
    ebui.WithButtonColors(ebui.ButtonColors{
        Default: color.RGBA{200, 200, 200, 255},
        Hovered: color.RGBA{220, 220, 220, 255},
        Pressed: color.RGBA{170, 170, 170, 255},
    }),
)
```

### Text Input

```go
input := ebui.NewTextInput(
    ebui.WithSize(200, 40),
    ebui.WithInitialText("Hello"),
    ebui.WithOnChange(func(text string) {
        println("Text changed:", text)
    }),
)
```

### Scrollable Container

```go
scrollable := ebui.NewScrollableContainer(
    ebui.WithSize(300, 400),
    ebui.WithLayout(ebui.NewVerticalStackLayout(10, ebui.AlignStart)),
)
```

### Window

```go
windowManager := ebui.NewWindowManager(
	ebui.WithSize(800, 600),
)
window := windowManager.CreateWindow(400, 300,
    ebui.WithWindowTitle("My Window"),
    ebui.WithWindowColors(ebui.WindowColors{
        Background: color.RGBA{240, 240, 240, 255},
        Header:     color.RGBA{200, 200, 200, 255},
        Border:     color.RGBA{0, 0, 0, 255},
    }),
)
```

## Debugging

EBUI includes a debug mode that visualizes component bounds and layout information. Set the global `Debug` variable to `true` to enable debug mode:

```go
ebui.Debug = true
```

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
