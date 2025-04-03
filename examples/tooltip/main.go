package main

import (
	"flag"
	"image/color"
	"log"

	"github.com/cbodonnell/ebui"
	"github.com/hajimehoshi/ebiten/v2"
)

type TooltipDemoGame struct {
	ui *ebui.Manager
}

func NewTooltipDemoGame() *TooltipDemoGame {
	game := &TooltipDemoGame{}

	// Create root container
	root := ebui.NewBaseContainer(
		ebui.WithSize(800, 600),
	)

	// Create content container
	content := ebui.NewLayoutContainer(
		ebui.WithSize(800, 600),
		ebui.WithBackground(color.RGBA{240, 240, 240, 255}),
		ebui.WithPadding(20, 20, 20, 20),
		ebui.WithLayout(ebui.NewVerticalStackLayout(15, ebui.AlignStart)),
	)
	root.AddChild(content)

	// Create the tooltip manager
	tm := ebui.NewTooltipManager(
		ebui.WithSize(800, 600),
	)

	// Title label
	titleLabel := ebui.NewLabel(
		"Enhanced Tooltip Demo",
		ebui.WithSize(760, 40),
		ebui.WithJustify(ebui.JustifyCenter),
		ebui.WithColor(color.RGBA{50, 50, 50, 255}),
	)
	content.AddChild(titleLabel)

	// Instructions label
	instructionsLabel := ebui.NewLabel(
		"Hover over elements to see tooltips that follow your mouse cursor",
		ebui.WithSize(760, 30),
		ebui.WithJustify(ebui.JustifyCenter),
		ebui.WithColor(color.RGBA{100, 100, 100, 255}),
	)
	content.AddChild(instructionsLabel)

	// Create a section for position types
	positionTypesTitle := ebui.NewLabel(
		"Tooltip Positions",
		ebui.WithSize(760, 30),
		ebui.WithJustify(ebui.JustifyLeft),
		ebui.WithColor(color.RGBA{80, 80, 80, 255}),
	)
	content.AddChild(positionTypesTitle)

	// Create a container for position types
	positionContainer := ebui.NewLayoutContainer(
		ebui.WithSize(760, 50),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignCenter)),
	)
	content.AddChild(positionContainer)

	// Create buttons with different tooltip positions
	createPositionButton(tm, positionContainer, "Top-Right", ebui.TooltipPositionTopRight)
	createPositionButton(tm, positionContainer, "Top-Left", ebui.TooltipPositionTopLeft)
	createPositionButton(tm, positionContainer, "Bottom-Right", ebui.TooltipPositionBottomRight)
	createPositionButton(tm, positionContainer, "Bottom-Left", ebui.TooltipPositionBottomLeft)

	// Create a second row of position buttons
	positionContainer2 := ebui.NewLayoutContainer(
		ebui.WithSize(760, 50),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignCenter)),
	)
	content.AddChild(positionContainer2)

	createPositionButton(tm, positionContainer2, "Top", ebui.TooltipPositionTop)
	createPositionButton(tm, positionContainer2, "Right", ebui.TooltipPositionRight)
	createPositionButton(tm, positionContainer2, "Bottom", ebui.TooltipPositionBottom)
	createPositionButton(tm, positionContainer2, "Left", ebui.TooltipPositionLeft)

	// Create a section for auto-positioning demo
	autoTitle := ebui.NewLabel(
		"Auto-Positioning Demo",
		ebui.WithSize(760, 30),
		ebui.WithJustify(ebui.JustifyLeft),
		ebui.WithColor(color.RGBA{80, 80, 80, 255}),
	)
	content.AddChild(autoTitle)

	autoDesc := ebui.NewLabel(
		"Try moving your cursor to the edges of the screen",
		ebui.WithSize(760, 30),
		ebui.WithJustify(ebui.JustifyCenter),
		ebui.WithColor(color.RGBA{100, 100, 100, 255}),
	)
	content.AddChild(autoDesc)

	// Create a container for edge buttons
	edgeContainer := ebui.NewLayoutContainer(
		ebui.WithSize(760, 80),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(600, ebui.AlignCenter)),
	)
	content.AddChild(edgeContainer)

	// Left edge button
	leftBtn := ebui.NewButton(
		ebui.WithSize(100, 40),
		ebui.WithLabelText("Left Edge"),
		ebui.WithButtonColors(ebui.ButtonColors{
			Default:     color.RGBA{46, 139, 87, 255},
			Hovered:     color.RGBA{66, 159, 107, 255},
			Pressed:     color.RGBA{26, 119, 67, 255},
			FocusBorder: color.Black,
		}),
	)

	leftTooltip := ebui.NewTooltip(
		ebui.WithTooltipPosition(ebui.TooltipPositionAuto),
	)
	leftTooltipContent := ebui.NewLabel(
		"This tooltip will automatically reposition to stay on screen",
		ebui.WithSize(250, 40),
		ebui.WithJustify(ebui.JustifyCenter),
		ebui.WithColor(color.Black),
		ebui.WithTextWrap(),
	)
	leftTooltip.SetContent(leftTooltipContent)
	tm.RegisterTooltip(leftBtn, leftTooltip)

	// Right edge button
	rightBtn := ebui.NewButton(
		ebui.WithSize(100, 40),
		ebui.WithLabelText("Right Edge"),
		ebui.WithButtonColors(ebui.ButtonColors{
			Default:     color.RGBA{46, 139, 87, 255},
			Hovered:     color.RGBA{66, 159, 107, 255},
			Pressed:     color.RGBA{26, 119, 67, 255},
			FocusBorder: color.Black,
		}),
	)

	rightTooltip := ebui.NewTooltip(
		ebui.WithTooltipPosition(ebui.TooltipPositionAuto),
	)
	rightTooltipContent := ebui.NewLabel(
		"This tooltip will automatically reposition to stay on screen",
		ebui.WithSize(250, 40),
		ebui.WithJustify(ebui.JustifyCenter),
		ebui.WithColor(color.Black),
		ebui.WithTextWrap(),
	)
	rightTooltip.SetContent(rightTooltipContent)
	tm.RegisterTooltip(rightBtn, rightTooltip)

	edgeContainer.AddChild(leftBtn)
	edgeContainer.AddChild(rightBtn)

	// Create a section for styled tooltips
	styledTitle := ebui.NewLabel(
		"Styled Tooltips",
		ebui.WithSize(760, 30),
		ebui.WithJustify(ebui.JustifyLeft),
		ebui.WithColor(color.RGBA{80, 80, 80, 255}),
	)
	content.AddChild(styledTitle)

	// Create a container for styled buttons
	styledContainer := ebui.NewLayoutContainer(
		ebui.WithSize(760, 50),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignCenter)),
	)
	content.AddChild(styledContainer)

	// Blue tooltip
	blueBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabelText("Blue Tooltip"),
		ebui.WithButtonColors(ebui.ButtonColors{
			Default:     color.RGBA{100, 149, 237, 255},
			Hovered:     color.RGBA{120, 169, 255, 255},
			Pressed:     color.RGBA{80, 129, 217, 255},
			FocusBorder: color.Black,
		}),
	)

	blueTooltip := ebui.NewTooltip(
		ebui.WithTooltipColors(ebui.TooltipColors{
			Background: color.RGBA{100, 149, 237, 230},
			Border:     color.RGBA{70, 119, 207, 255},
		}),
	)
	blueTooltipContent := ebui.NewLabel(
		"Blue styled tooltip",
		ebui.WithSize(150, 30),
		ebui.WithJustify(ebui.JustifyCenter),
		ebui.WithColor(color.White),
	)
	blueTooltip.SetContent(blueTooltipContent)
	tm.RegisterTooltip(blueBtn, blueTooltip)
	styledContainer.AddChild(blueBtn)

	// Green tooltip
	greenBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabelText("Green Tooltip"),
		ebui.WithButtonColors(ebui.ButtonColors{
			Default:     color.RGBA{46, 139, 87, 255},
			Hovered:     color.RGBA{66, 159, 107, 255},
			Pressed:     color.RGBA{26, 119, 67, 255},
			FocusBorder: color.Black,
		}),
	)

	greenTooltip := ebui.NewTooltip(
		ebui.WithTooltipColors(ebui.TooltipColors{
			Background: color.RGBA{46, 139, 87, 230},
			Border:     color.RGBA{26, 119, 67, 255},
		}),
	)
	greenTooltipContent := ebui.NewLabel(
		"Green styled tooltip",
		ebui.WithSize(150, 30),
		ebui.WithJustify(ebui.JustifyCenter),
		ebui.WithColor(color.White),
	)
	greenTooltip.SetContent(greenTooltipContent)
	tm.RegisterTooltip(greenBtn, greenTooltip)
	styledContainer.AddChild(greenBtn)

	// Red tooltip
	redBtn := ebui.NewButton(
		ebui.WithSize(120, 40),
		ebui.WithLabelText("Red Tooltip"),
		ebui.WithButtonColors(ebui.ButtonColors{
			Default:     color.RGBA{220, 53, 69, 255},
			Hovered:     color.RGBA{240, 73, 89, 255},
			Pressed:     color.RGBA{200, 33, 49, 255},
			FocusBorder: color.Black,
		}),
	)

	redTooltip := ebui.NewTooltip(
		ebui.WithTooltipColors(ebui.TooltipColors{
			Background: color.RGBA{220, 53, 69, 230},
			Border:     color.RGBA{200, 33, 49, 255},
		}),
	)
	redTooltipContent := ebui.NewLabel(
		"Red styled tooltip",
		ebui.WithSize(150, 30),
		ebui.WithJustify(ebui.JustifyCenter),
		ebui.WithColor(color.White),
	)
	redTooltip.SetContent(redTooltipContent)
	tm.RegisterTooltip(redBtn, redTooltip)
	styledContainer.AddChild(redBtn)

	// Create a section for controls
	controlsTitle := ebui.NewLabel(
		"Controls",
		ebui.WithSize(760, 30),
		ebui.WithJustify(ebui.JustifyLeft),
		ebui.WithColor(color.RGBA{80, 80, 80, 255}),
	)
	content.AddChild(controlsTitle)

	// Create container for control buttons
	controlsContainer := ebui.NewLayoutContainer(
		ebui.WithSize(760, 50),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(20, ebui.AlignCenter)),
	)
	content.AddChild(controlsContainer)

	// Enable tooltip button
	enableBtn := ebui.NewButton(
		ebui.WithSize(150, 40),
		ebui.WithLabelText("Enable Tooltips"),
		ebui.WithButtonColors(ebui.ButtonColors{
			Default:     color.RGBA{46, 139, 87, 255},
			Hovered:     color.RGBA{66, 159, 107, 255},
			Pressed:     color.RGBA{26, 119, 67, 255},
			FocusBorder: color.Black,
		}),
	)
	enableBtn.SetClickHandler(func() {
		tm.Enable()
	})
	controlsContainer.AddChild(enableBtn)

	// Disable tooltip button
	disableBtn := ebui.NewButton(
		ebui.WithSize(150, 40),
		ebui.WithLabelText("Disable Tooltips"),
		ebui.WithButtonColors(ebui.ButtonColors{
			Default:     color.RGBA{220, 53, 69, 255},
			Hovered:     color.RGBA{240, 73, 89, 255},
			Pressed:     color.RGBA{200, 33, 49, 255},
			FocusBorder: color.Black,
		}),
	)
	disableBtn.SetClickHandler(func() {
		tm.Disable()
	})
	controlsContainer.AddChild(disableBtn)

	// Add tooltip manager and create UI manager
	root.AddChild(tm)
	game.ui = ebui.NewManager(root)

	return game
}

// Helper function to create a button with a tooltip positioned relative to mouse
func createPositionButton(tm *ebui.TooltipManager, container ebui.Container, labelText string, position ebui.TooltipPosition) {
	btn := ebui.NewButton(
		ebui.WithSize(100, 40),
		ebui.WithLabelText(labelText),
		ebui.WithButtonColors(ebui.ButtonColors{
			Default:     color.RGBA{100, 149, 237, 255},
			Hovered:     color.RGBA{120, 169, 255, 255},
			Pressed:     color.RGBA{80, 129, 217, 255},
			FocusBorder: color.Black,
		}),
	)

	tooltip := ebui.NewTooltip(
		ebui.WithTooltipPosition(position),
	)
	tooltipContent := ebui.NewLabel(
		"Position: "+labelText,
		ebui.WithSize(150, 30),
		ebui.WithJustify(ebui.JustifyCenter),
		ebui.WithColor(color.Black),
	)
	tooltip.SetContent(tooltipContent)
	tm.RegisterTooltip(btn, tooltip)

	container.AddChild(btn)
}

func (g *TooltipDemoGame) Update() error {
	return g.ui.Update()
}

func (g *TooltipDemoGame) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

func (g *TooltipDemoGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("EBUI Tooltip Demo")

	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *debug {
		ebui.Debug = true
	}

	if err := ebiten.RunGame(NewTooltipDemoGame()); err != nil {
		log.Fatal(err)
	}
}
