package element

import (
	"errors"
	"github.com/faiface/pixel"
	"github.com/orfby/ui/pkg/ui/util"
	"image"
	"image/color"
)

//Interface for a background
type Bkg interface {
	//Function to get the element's
	//background XML field
	GetBkgField() string
	//Function to get the element's
	//background sprite
	GetBkgSprite() *pixel.Sprite
	//Function to set the element's
	//background sprite
	SetBkgSprite(*pixel.Sprite)
	//Function to determine whether
	//the background should repeat.
	//Should only return true if
	//the background isn't a colour
	ShouldRepeat() bool
}

//Type for an element's background
type BkgImpl struct {
	//The element's background
	//from xml
	BackgroundField string `uixml:"http://github.com/orfby/ui/api/schema background,optional"`
	//The element's background
	//sprite
	sprite *pixel.Sprite
	//Whether the background should
	//repeat
	Repeat bool `uixml:"http://github.com/orfby/ui/api/schema bkg-repeat,optional"`
}

//Function to get the background's
//XML field
func (e *BkgImpl) GetBkgField() string { return e.BackgroundField }

//Function to get the background's
//sprite
func (e *BkgImpl) GetBkgSprite() *pixel.Sprite { return e.sprite }

//Function to set the background's
//sprite
func (e *BkgImpl) SetBkgSprite(s *pixel.Sprite) { e.sprite = s }

//Function to determine whether
//the background should repeat.
//Should only return true if
//the background isn't a colour
func (e *BkgImpl) ShouldRepeat() bool {
	return e.GetBkgField() != "" && e.GetBkgField()[0] != '#' && e.Repeat
}

//Function to determine whether the
//background has been initialised, by
//whether its background sprite has been
//set (assuming it's meant to be set).
//This function doesn't call
//element.IsInitialised
func (e *BkgImpl) IsInitialised() bool {
	//If the element doesn't have a background
	return e.GetBkgField() == "" ||
		//Or the background has been initialised
		e.GetBkgSprite() != nil
}

//Function to create a sprite
//from an XML string
//todo move me somewhere more applicable
func CreateSpriteFromField(field string) (*pixel.Sprite, error) {
	if field != "" {
		//If the first character is a hash
		if field[0] == '#' {
			//Convert the field to a colour type
			colour, err := util.ParseColor(field)
			if err != nil {
				return nil, errors.New("invalid colour attribute value '" + field + "'")
			}
			//Create a 1x1 image
			img := image.NewRGBA(image.Rect(0, 0, 2, 2))
			for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
				for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
					img.SetRGBA(x, y, colour)
				}
			}
			//Convert it to a pixel picture
			pic := pixel.PictureDataFromImage(img)
			//Create a sprite from the picture
			return pixel.NewSprite(pic, pic.Bounds()), nil
		} else {
			//Load the picture
			pic, err := util.LoadPicture(field)
			if err != nil {
				return nil, err
			}
			//Create a sprite from the picture
			return pixel.NewSprite(pic, pic.Bounds()), nil
		}
	}
	return nil, nil
}

//Function to initialise an element's
//background. Doesn't call element.Init.
//Should be called last (as the size of
//the background depends on the actual
//size of the element)
func InitBkg(e Bkg, bounds *pixel.Rect) error {
	//If the bounds are known
	if bounds != nil {
		//If the sprite doesn't exist
		if e.GetBkgSprite() == nil {
			//Create the background
			sprite, err := CreateSpriteFromField(e.GetBkgField())
			if err != nil {
				return err
			}
			e.SetBkgSprite(sprite)
		}
	}

	return nil
}

//Function to draw the background
//of an element. This function
//should be called first
func DrawBkg(e Element) {
	//Draw the background sprite, if it exists
	if e.GetBkgSprite() != nil {
		//If the background should repeat
		if e.ShouldRepeat() {
			//Get the size of the background sprite
			spriteSize := e.GetBkgSprite().Picture().Bounds().Size()
			//Iterate over the y coords of each tile
			for y := spriteSize.Y / 2; y < e.GetCanvas().Bounds().Max.Y; y += spriteSize.Y {
				//Iterate over the x coords of each tile
				for x := spriteSize.X / 2; x < e.GetCanvas().Bounds().Max.X; x += spriteSize.X {
					mat := pixel.IM
					//Move the tile to the position
					mat = mat.Moved(pixel.V(x, y))
					//Draw it
					e.GetBkgSprite().Draw(e.GetCanvas(), mat)
				}
			}
		} else {
			mat := pixel.IM
			//Scale the background up to the size of the element
			mat = mat.ScaledXY(pixel.ZV, pixel.V(
				e.GetCanvas().Bounds().Size().X/e.GetBkgSprite().Picture().Bounds().Size().X,
				e.GetCanvas().Bounds().Size().Y/e.GetBkgSprite().Picture().Bounds().Size().Y))
			//Move it to the center of the canvas
			mat = mat.Moved(e.GetCanvas().Bounds().Center())
			//Draw the background
			e.GetBkgSprite().Draw(e.GetCanvas(), mat)
		}

		//If there isn't a background
	} else {
		//Clear the canvas (with a transparent background)
		//(this is why it's important the background is
		//drawn first)
		e.GetCanvas().Clear(color.Transparent)
	}
}
