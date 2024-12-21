# EBUI - Ebitengine UI Framework

EBUI is a comprehensive UI framework for [Ebitengine](https://ebitengine.org/), providing a robust set of components and layouts for building game interfaces and tools.

## Features

- **Component-Based Architecture**: Build UIs using reusable, composable components
- **Flexible Layout System**: Arrange components using various layout strategies
- **Event System**: Handle user interactions with a comprehensive event system
- **Rich Component Library**:
  - Buttons with customizable colors and states
  - Labels with text alignment options
  - Text inputs with selection and clipboard support
  - Scrollable containers
  - Windows with drag-and-drop functionality

## Installation

```bash
go get github.com/cbodonnell/ebui
```

## Quick Start

```go
package main

import (
    "log"
    "github.com/cbodonnell/ebui"
    "github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
    ui *ebui.Manager
}

func NewGame() *Game {
    // Create a root container
    root := ebui.NewBaseContainer(
        ebui.WithSize(800, 600),
    )

    // Create a button
    button := ebui.NewButton(
        ebui.WithSize(120, 40),
        ebui.WithLabelText("Click Me"),
    )
    button.SetClickHandler(func() {
        println("Button clicked!")
    })

    // Add button to root
    root.AddChild(button)

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

Custom 

```go
// Basic container
container := ebui.NewBaseContainer(
    ebui.WithSize(400, 300),
    ebui.WithPadding(10, 10, 10, 10),
)

// Layout container with vertical stack
layoutContainer := ebui.NewLayoutContainer(
    ebui.WithSize(400, 300),
    ebui.WithLayout(ebui.NewVerticalStackLayout(10, ebui.AlignStart)),
)
```

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

EBUI includes a debug mode that visualizes component bounds and layout information:

```go
ebui.Debug = true
```

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
