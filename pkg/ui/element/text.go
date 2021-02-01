package element

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"io/ioutil"
)

// Interface for some text
type Text interface {
	// Function to get the text
	// XML field
	GetField() string

	// Function to get the text
	// font
	GetFont() string

	// Function to get the text
	// currently being displayed
	GetText() string
	// Function to set the text
	// currently being displayed.
	// This function doesn't alter
	// the text's width and height,
	// for this element.SetText
	// should be used instead
	SetText(string) error

	// Function to get the text size
	GetTextSize() float64

	// Function to get the text
	// "sprite"
	GetSprite() *text.Text
	// Function to set the text
	// "sprite"
	SetSprite(*text.Text)
}

// Implementation of the text interface
type TextImpl struct {
	// The element's text from xml
	Text string `uixml:"http://github.com/bhollier/ui/api/schema text,optional"`
	// The element's font
	Font string `uixml:"http://github.com/bhollier/ui/api/schema font,optional"`
	// The element's text size
	Size float64 `uixml:"http://github.com/bhollier/ui/api/schema text-size,optional"`
	// The element's text string
	textStr string
	// The element's text "sprite"
	textSprite *text.Text
}

// Function to get the text
// XML field
func (t *TextImpl) GetField() string { return t.Text }

// Function to get the text
// font
func (t *TextImpl) GetFont() string { return t.Font }

// Function to get the text
// currently being displayed
func (t *TextImpl) GetText() string { return t.textStr }

// Function to set the text
// currently being displayed.
// This function doesn't alter
// the text's width and height,
// for this element.SetText
// should be used instead
func (t *TextImpl) SetText(s string) error {
	t.textStr = s
	t.GetSprite().Clear()
	_, err := fmt.Fprintf(t.GetSprite(), t.textStr)
	return err
}

// Function to get the text
// size
func (t *TextImpl) GetTextSize() float64 { return t.Size }

// Function to get the text
// sprite
func (t *TextImpl) GetSprite() *text.Text { return t.textSprite }

// Function to set the text
// sprite
func (t *TextImpl) SetSprite(s *text.Text) { t.textSprite = s }

// Function to reset the text
func (t *TextImpl) Reset() {}

// Function to determine whether the
// text has been initialised, by
// whether its text sprite has been
// set (assuming it's meant to be set).
// This function doesn't call
// element.IsInitialised
func (t *TextImpl) IsInitialised() bool {
	// If the element doesn't have any text
	// or the text has been set
	return t.GetField() == "" || t.GetSprite() != nil
}

// Function to initialise an element's
// text. Doesn't call element.Init
func InitText(e Element, t Text) error {
	// If the text hasn't been made yet
	if t.GetSprite() == nil {
		// Get the font
		var ttf *truetype.Font
		var err error
		if t.GetFont() == "" {
			ttf, err = truetype.Parse(goregular.TTF)
			if err != nil {
				return err
			}
		} else {
			// Open the file
			file, err := e.GetFS().Open(t.GetFont())
			if err != nil {
				return err
			}
			defer file.Close()

			// Read the file
			fontData, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}

			// Parse the data as a font
			ttf, err = truetype.Parse(fontData)
			if err != nil {
				return err
			}
		}

		// Get the text size
		textSize := t.GetTextSize()
		if textSize == 0 {
			textSize = 24
		}

		// Create a new font face
		face := truetype.NewFace(ttf, &truetype.Options{Size: textSize})
		// Make a new text object
		t.SetSprite(text.New(pixel.V(0, 0), text.NewAtlas(face, text.ASCII)))
		// Set the text colour
		t.GetSprite().Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}
		// Set the text
		err = t.SetText(t.GetField())
		if err != nil {
			return err
		}
	}

	// If the sprite has been created
	if t.GetSprite() != nil {
		// If the element's width isn't known and
		// the width is meant to match the content size
		if e.GetActualWidth() == nil && e.GetRelWidth().MatchContent {
			// Set the actual width as the size of the text sprite
			newWidth := t.GetSprite().Bounds().Size().X
			e.SetActualWidth(&newWidth)
		}

		// If the element's height isn't known and
		// the height is meant to match the content size
		if e.GetActualHeight() == nil && e.GetRelHeight().MatchContent {
			// Set the actual height as the size of the text sprite
			newHeight := t.GetSprite().Bounds().Size().Y
			e.SetActualHeight(&newHeight)
		}
	}

	return nil
}

// Function to draw an element's
// text
func DrawText(e Element, t Text) {
	// Draw the text sprite, if it exists
	if t.GetSprite() != nil {
		mat := pixel.IM
		// Move it to the center of the canvas
		mat = mat.Moved(pixel.V(e.GetCanvas().Bounds().Center().X-(t.GetSprite().Bounds().Size().X/2),
			e.GetCanvas().Bounds().Center().Y-(t.GetSprite().Bounds().Size().Y/2)))
		// todo why is it not centered properly?
		// Draw the text
		t.GetSprite().Draw(e.GetCanvas(), mat)
	}
}
