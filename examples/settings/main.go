package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"

	"github.com/cbodonnell/ebui"
	"github.com/hajimehoshi/ebiten/v2"
)

type SettingsGame struct {
	ui               *ebui.Manager
	brightnessVal    float64
	volumeVal        float64
	fontSizeVal      float64
	speedVal         float64
	brightnessLbl    *ebui.Label
	volumeLbl        *ebui.Label
	fontSizeLbl      *ebui.Label
	speedLbl         *ebui.Label
	brightnessSlider *ebui.Slider
	volumeSlider     *ebui.Slider
	fontSizeSlider   *ebui.Slider
	speedSlider      *ebui.Slider
}

func NewSettingsGame() *SettingsGame {
	game := &SettingsGame{
		brightnessVal: 75,
		volumeVal:     50,
		fontSizeVal:   16,
		speedVal:      1,
	}

	// Create root container
	root := ebui.NewLayoutContainer(
		ebui.WithSize(500, 400),
		ebui.WithBackground(color.RGBA{240, 240, 240, 255}),
		ebui.WithPadding(20, 20, 20, 20),
		ebui.WithLayout(ebui.NewVerticalStackLayout(20, ebui.AlignStart)),
	)

	// Title label
	titleLabel := ebui.NewLabel(
		"Settings",
		ebui.WithSize(460, 40),
		ebui.WithJustify(ebui.JustifyCenter),
		ebui.WithColor(color.RGBA{50, 50, 50, 255}),
	)
	root.AddChild(titleLabel)

	// Create sliders and store references to them
	game.brightnessSlider = ebui.NewSlider(
		ebui.WithSize(290, 40),
		ebui.WithMinValue(0),
		ebui.WithMaxValue(100),
		ebui.WithValue(game.brightnessVal),
		ebui.WithStepSize(1),
		ebui.WithTrackHeight(8),
		ebui.WithThumbSize(20, 20),
		ebui.WithOnChangeHandler(func(val float64) {
			game.brightnessVal = val
			game.brightnessLbl.SetText(fmt.Sprintf("Brightness: %.0f%%", val))
		}),
	)

	brightnessSlider := game.createSettingRow("Brightness", &game.brightnessLbl, game.brightnessSlider, "%")
	root.AddChild(brightnessSlider)

	game.volumeSlider = ebui.NewSlider(
		ebui.WithSize(290, 40),
		ebui.WithMinValue(0),
		ebui.WithMaxValue(100),
		ebui.WithValue(game.volumeVal),
		ebui.WithStepSize(5),
		ebui.WithTrackHeight(8),
		ebui.WithThumbSize(20, 20),
		ebui.WithOnChangeHandler(func(val float64) {
			game.volumeVal = val
			game.volumeLbl.SetText(fmt.Sprintf("Volume: %.0f%%", val))
		}),
	)

	volumeSlider := game.createSettingRow("Volume", &game.volumeLbl, game.volumeSlider, "%")
	root.AddChild(volumeSlider)

	game.fontSizeSlider = ebui.NewSlider(
		ebui.WithSize(290, 40),
		ebui.WithMinValue(8),
		ebui.WithMaxValue(32),
		ebui.WithValue(game.fontSizeVal),
		ebui.WithStepSize(1),
		ebui.WithTrackHeight(8),
		ebui.WithThumbSize(20, 20),
		ebui.WithOnChangeHandler(func(val float64) {
			game.fontSizeVal = val
			game.fontSizeLbl.SetText(fmt.Sprintf("Font Size: %.0fpx", val))
		}),
	)

	fontSizeSlider := game.createSettingRow("Font Size", &game.fontSizeLbl, game.fontSizeSlider, "px")
	root.AddChild(fontSizeSlider)

	// Create a slider with custom colors
	speedColors := ebui.SliderColors{
		Track:        color.RGBA{200, 200, 200, 255},
		TrackFilled:  color.RGBA{76, 175, 80, 255}, // Green
		Thumb:        color.RGBA{255, 255, 255, 255},
		ThumbHovered: color.RGBA{240, 240, 240, 255},
		ThumbDragged: color.RGBA{230, 230, 230, 255},
		FocusBorder:  color.RGBA{76, 175, 80, 255},
	}

	game.speedSlider = ebui.NewSlider(
		ebui.WithSize(290, 40),
		ebui.WithMinValue(0.5),
		ebui.WithMaxValue(2.0),
		ebui.WithValue(game.speedVal),
		ebui.WithStepSize(0.1),
		ebui.WithTrackHeight(8),
		ebui.WithThumbSize(20, 20),
		ebui.WithSliderColors(speedColors),
		ebui.WithOnChangeHandler(func(val float64) {
			game.speedVal = val
			game.speedLbl.SetText(fmt.Sprintf("Game Speed: %.1fx", val))
		}),
	)

	speedSlider := game.createCustomSettingRow("Game Speed", &game.speedLbl, game.speedSlider, "x")
	root.AddChild(speedSlider)

	// Reset button
	resetBtn := ebui.NewButton(
		ebui.WithSize(150, 40),
		ebui.WithLabelText("Reset to Defaults"),
		ebui.WithButtonColors(ebui.ButtonColors{
			Default:     color.RGBA{220, 53, 69, 255}, // Red
			Hovered:     color.RGBA{240, 73, 89, 255},
			Pressed:     color.RGBA{200, 33, 49, 255},
			FocusBorder: color.Black,
		}),
	)

	resetBtn.SetClickHandler(func() {
		// Reset all values to defaults
		game.brightnessVal = 75
		game.volumeVal = 50
		game.fontSizeVal = 16
		game.speedVal = 1

		// Update labels
		game.brightnessLbl.SetText(fmt.Sprintf("Brightness: %.0f%%", game.brightnessVal))
		game.volumeLbl.SetText(fmt.Sprintf("Volume: %.0f%%", game.volumeVal))
		game.fontSizeLbl.SetText(fmt.Sprintf("Font Size: %.0fpx", game.fontSizeVal))
		game.speedLbl.SetText(fmt.Sprintf("Game Speed: %.1fx", game.speedVal))

		// Directly update sliders with their reference
		game.brightnessSlider.SetValue(game.brightnessVal)
		game.volumeSlider.SetValue(game.volumeVal)
		game.fontSizeSlider.SetValue(game.fontSizeVal)
		game.speedSlider.SetValue(game.speedVal)
	})

	// Button container for centering
	btnContainer := ebui.NewLayoutContainer(
		ebui.WithSize(460, 40),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(0, ebui.AlignCenter)),
	)
	btnContainer.AddChild(resetBtn)
	root.AddChild(btnContainer)

	game.ui = ebui.NewManager(root)
	return game
}

func (g *SettingsGame) createSettingRow(
	label string,
	labelRef **ebui.Label,
	slider *ebui.Slider,
	suffix string,
) *ebui.LayoutContainer {
	container := ebui.NewLayoutContainer(
		ebui.WithSize(460, 40),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignCenter)),
	)

	// Create label
	*labelRef = ebui.NewLabel(
		fmt.Sprintf("%s: %.0f%s", label, slider.GetValue(), suffix),
		ebui.WithSize(150, 40),
		ebui.WithJustify(ebui.JustifyLeft),
	)
	container.AddChild(*labelRef)

	// Add the slider to the container
	container.AddChild(slider)

	return container
}

func (g *SettingsGame) createCustomSettingRow(
	label string,
	labelRef **ebui.Label,
	slider *ebui.Slider,
	suffix string,
) *ebui.LayoutContainer {
	container := ebui.NewLayoutContainer(
		ebui.WithSize(460, 40),
		ebui.WithLayout(ebui.NewHorizontalStackLayout(10, ebui.AlignCenter)),
	)

	// Create label
	*labelRef = ebui.NewLabel(
		fmt.Sprintf("%s: %.1f%s", label, slider.GetValue(), suffix),
		ebui.WithSize(150, 40),
		ebui.WithJustify(ebui.JustifyLeft),
	)
	container.AddChild(*labelRef)

	// Add the slider to the container
	container.AddChild(slider)

	return container
}

func (g *SettingsGame) Update() error {
	return g.ui.Update()
}

func (g *SettingsGame) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

func (g *SettingsGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 500, 400
}

func main() {
	ebiten.SetWindowSize(500, 400)
	ebiten.SetWindowTitle("EBUI Settings Example")

	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *debug {
		ebui.Debug = true
	}

	if err := ebiten.RunGame(NewSettingsGame()); err != nil {
		log.Fatal(err)
	}
}
