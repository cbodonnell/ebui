package main

import (
	"image/color"
	"log"

	"github.com/cbodonnell/ebui"
	"github.com/hajimehoshi/ebiten/v2"
)

type LoginGame struct {
	ui            *ebui.Manager
	usernameInput *ebui.TextInput
	passwordInput *ebui.TextInput
	statusLabel   *ebui.Label
}

func NewLoginGame() *LoginGame {
	game := &LoginGame{}

	// Create root container
	root := ebui.NewBaseContainer(
		ebui.WithSize(400, 400),
		ebui.WithBackground(color.RGBA{240, 240, 240, 255}),
	)

	// Create centered container for login form
	formContainer := ebui.NewLayoutContainer(
		ebui.WithSize(300, 350),
		ebui.WithPosition(ebui.Position{X: 50, Y: 25}),
		// ebui.WithBackground(color.RGBA{255, 255, 255, 255}),
		// use a darker gray instead of white for the form background
		ebui.WithBackground(color.RGBA{220, 220, 220, 255}),
		ebui.WithPadding(20, 20, 20, 20),
		ebui.WithLayout(ebui.NewVerticalStackLayout(15, ebui.AlignStart)),
	)

	// Title
	titleLabel := ebui.NewLabel(
		"Login",
		ebui.WithSize(260, 30),
		ebui.WithJustify(ebui.JustifyCenter),
	)

	// Username input
	usernameLabel := ebui.NewLabel(
		"Username:",
		ebui.WithSize(260, 20),
		ebui.WithJustify(ebui.JustifyLeft),
	)

	game.usernameInput = ebui.NewTextInput(
		ebui.WithSize(260, 30),
		ebui.WithTextInputColors(ebui.TextInputColors{
			Text:       color.Black,
			Background: color.White,
			Cursor:     color.Black,
			Selection:  color.RGBA{100, 149, 237, 127},
		}),
	)

	// Password input
	passwordLabel := ebui.NewLabel(
		"Password:",
		ebui.WithSize(260, 20),
		ebui.WithJustify(ebui.JustifyLeft),
	)

	game.passwordInput = ebui.NewTextInput(
		ebui.WithSize(260, 30),
		ebui.WithPasswordMasking(),
		ebui.WithTextInputColors(ebui.TextInputColors{
			Text:       color.Black,
			Background: color.White,
			Cursor:     color.Black,
			Selection:  color.RGBA{100, 149, 237, 127},
		}),
	)

	// Login button
	loginBtn := ebui.NewButton(
		ebui.WithSize(260, 40),
		ebui.WithLabelText("Login"),
		ebui.WithButtonColors(ebui.ButtonColors{
			Default: color.RGBA{46, 139, 87, 255},  // Sea Green
			Hovered: color.RGBA{60, 179, 113, 255}, // Medium Sea Green
			Pressed: color.RGBA{32, 97, 61, 255},   // Dark Sea Green
		}),
	)

	// Status label for showing login results
	game.statusLabel = ebui.NewLabel(
		"",
		ebui.WithSize(260, 20),
		ebui.WithJustify(ebui.JustifyCenter),
		ebui.WithColor(color.RGBA{255, 0, 0, 255}), // Red for errors
	)

	loginBtn.SetClickHandler(game.handleLogin)

	// Add all components to form container
	formContainer.AddChild(titleLabel)
	formContainer.AddChild(usernameLabel)
	formContainer.AddChild(game.usernameInput)
	formContainer.AddChild(passwordLabel)
	formContainer.AddChild(game.passwordInput)
	formContainer.AddChild(loginBtn)
	formContainer.AddChild(game.statusLabel)

	// Add form container to root
	root.AddChild(formContainer)

	game.ui = ebui.NewManager(root)
	return game
}

func (g *LoginGame) handleLogin() {
	username := g.usernameInput.GetText()
	password := g.passwordInput.GetText()

	// Simple validation
	if username == "" || password == "" {
		g.statusLabel.SetText("Please enter both username and password")
		return
	}

	// In a real app, you would validate credentials here
	if username == "admin" && password == "password" {
		g.statusLabel.SetText("Login successful!")
		g.statusLabel.SetColor(color.RGBA{0, 128, 0, 255}) // Green for success
	} else {
		g.statusLabel.SetText("Invalid credentials")
		g.statusLabel.SetColor(color.RGBA{255, 0, 0, 255}) // Red for error
	}
}

func (g *LoginGame) Update() error {
	return g.ui.Update()
}

func (g *LoginGame) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

func (g *LoginGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 400, 400
}

func main() {
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("EBUI Login Example")

	if err := ebiten.RunGame(NewLoginGame()); err != nil {
		log.Fatal(err)
	}
}
