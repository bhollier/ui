package element

import (
	"github.com/faiface/pixel"
	"github.com/orfby/ui/pkg/ui/util"
	"image/color"
	"net/http"
)

//Interface for a background
type Bkg interface {
	//Function to get the element's
	//background XML field
	GetBkgField() string
	//Function to get the element's
	//background scale options
	GetBkgScale() util.ScaleOption
	//Function to get the element's
	//background sprite
	GetBkgSprite() *pixel.Sprite
	//Function to set the element's
	//background sprite
	SetBkgSprite(*pixel.Sprite)
}

//Type for an element's background
type BkgImpl struct {
	//The element's background
	//from xml
	BackgroundField string `uixml:"http://github.com/orfby/ui/api/schema background,optional"`
	//The element's scale option
	BackgroundScale util.ScaleOption `uixml:"http://github.com/orfby/ui/api/schema bkg-scale,optional"`
	//The element's background
	//sprite
	sprite *pixel.Sprite
}

//Function to get the background's
//XML field
func (e *BkgImpl) GetBkgField() string { return e.BackgroundField }

//Function to get the background's
//scale option
func (e *BkgImpl) GetBkgScale() util.ScaleOption { return e.BackgroundScale }

//Function to get the background's
//sprite
func (e *BkgImpl) GetBkgSprite() *pixel.Sprite { return e.sprite }

//Function to set the background's
//sprite
func (e *BkgImpl) SetBkgSprite(s *pixel.Sprite) { e.sprite = s }

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

//Function to initialise an element's
//background. Doesn't call element.Init.
//Should be called last (as the size of
//the background depends on the actual
//size of the element)
func InitBkg(e Bkg, fs http.FileSystem, bounds *pixel.Rect) error {
	//If the bounds are known
	if bounds != nil {
		//If the sprite doesn't exist
		if e.GetBkgSprite() == nil {
			//Load the background
			picture, err := util.CreatePictureFromField(fs, e.GetBkgField())
			if err != nil {
				return err
			}
			if picture != nil {
				//Create the sprite
				sprite := pixel.NewSprite(picture, picture.Bounds())
				e.SetBkgSprite(sprite)
			}
		}
	}

	return nil
}

//Function to draw the background
//of an element. This function
//should be called first
func DrawBkg(e Element) {
	//Clear the canvas (with a transparent background)
	//(this is why it's important the background is
	//drawn first)
	e.GetCanvas().Clear(color.Transparent)

	scale := e.GetBkgScale()
	//If the scale is zero
	if scale == util.ZeroScaleOption {
		//Set the default as scale to fill
		scale = util.ScaleToFill
	}

	//Draw the background sprite
	util.DrawSprite(e.GetCanvas(), e.GetBkgSprite(),
		scale, e.GetGravity())
}
