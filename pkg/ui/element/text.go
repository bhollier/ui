package element

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
)

//Interface for an element with text
type HasText interface {
	//An element with text is an
	//element
	IsElement

	//Function to get the element's
	//text XML field
	GetTextField() string
	//Function to get the element's
	//text "sprite"
	GetTextSprite() *text.Text
	//Function to set the element's
	//text "sprite"
	SetTextSprite(*text.Text)
	//Function to get the element's
	//text size
	GetTextSize() float64
}

//Type for an element's text
type Text struct {
	//The element's text from xml
	Text string `uixml:"http://github.com/orfby/ui/api/schema text,optional"`
	//The element's text size
	Size float64 `uixml:"http://github.com/orfby/ui/api/schema text-size,optional"`
	//The element's text "sprite"
	textSprite *text.Text
}

//Function to get the element's
//text XML field
func (e *Text) GetTextField() string {return e.Text}
//Function to get the element's
//text sprite
func (e *Text) GetTextSprite() *text.Text {return e.textSprite}
//Function to set the element's
//text sprite
func (e *Text) SetTextSprite(s *text.Text) {e.textSprite = s}
//Function to get the element's
//text size
func (e *Text) GetTextSize() float64 {return e.Size}

//Function to determine whether the
//text has been initialised, by
//whether its text sprite has been
//set (assuming it's meant to be set).
//This function doesn't call
//element.IsInitialised
func (e *Text) IsInitialised() bool {
	//If the element doesn't have any text
	//or the text has been set
	return e.GetTextField() == "" || e.GetTextSprite() != nil
}

//Function to initialise an element's
//text. Doesn't call element.Init
func InitText(e HasText, _ *pixel.Rect) error {
	//If there should be text
	//but it hasn't been made yet
	if e.GetTextField() != "" && e.GetTextSprite() == nil {
		//Get the font
		ttf, err := truetype.Parse(goregular.TTF)
		if err != nil {return err}

		//Get the text size
		textSize := e.GetTextSize()
		if textSize == 0 {textSize = 24}

		//Create a new font face
		face := truetype.NewFace(ttf, &truetype.Options{Size: textSize})
		//Make a new text object
		e.SetTextSprite(text.New(pixel.V(0, 0), text.NewAtlas(face, text.ASCII)))
		//Set the text colour
		e.GetTextSprite().Color = color.RGBA{R: 0, G: 0, B: 0, A:255}
		//Add the text
		_, err = fmt.Fprintf(e.GetTextSprite(), e.GetTextField())
		if err != nil {return err}
	}

	//If the sprite has been created
	if e.GetTextSprite() != nil {
		//If the element's width isn't known and
		//the width is meant to match the content size
		if e.GetActualWidth() == nil && e.GetRelWidth().MatchContent {
			//Set the actual width as the size of the text sprite
			newWidth := e.GetTextSprite().Bounds().Size().X
			e.SetActualWidth(&newWidth)
		}

		//If the element's height isn't known and
		//the height is meant to match the content size
		if e.GetActualHeight() == nil && e.GetRelHeight().MatchContent {
			//Set the actual height as the size of the text sprite
			newHeight := e.GetTextSprite().Bounds().Size().Y
			e.SetActualHeight(&newHeight)
		}
	}

	return nil
}

//Function to draw an element's
//text
func DrawText(e HasText) {
	//Draw the text sprite, if it exists
	if e.GetTextSprite() != nil {
		mat := pixel.IM
		//Move it to the center of the canvas
		mat = mat.Moved(pixel.V(e.GetCanvas().Bounds().Center().X - (e.GetTextSprite().Bounds().Size().X / 2),
			e.GetCanvas().Bounds().Center().Y - (e.GetTextSprite().Bounds().Size().Y / 2)))
		//todo why is it not centered properly?
		//Draw the image
		e.GetTextSprite().Draw(e.GetCanvas(), mat)
	}
}