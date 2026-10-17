package ebui

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// DefaultFont is the face labels and text inputs are created with when no font is given.
// Set it before building components.
var DefaultFont font.Face = basicfont.Face7x13
