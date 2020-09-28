package element

import (
	"github.com/faiface/pixel"
	"github.com/orfby/ui/pkg/ui/util"
)

//Interface for an element with an image
type HasImage interface {
	//An element with an image is an
	//element
	Element

	//Function to get the element's
	//image XML field
	GetImageField() string
	//Function to get the element's
	//image scale options
	GetImageScale() util.ScaleOption
	//Function to get the element's
	//image sprite
	GetImageSprite() *pixel.Sprite
	//Function to set the element's
	//image sprite
	SetImageSprite(*pixel.Sprite)
}

//Type for an element's image
type ImageImpl struct {
	//The element's image
	//from xml
	ImageField string `uixml:"http://github.com/orfby/ui/api/schema source,optional"`
	//The element's scale
	//option
	ImageScale util.ScaleOption `uixml:"http://github.com/orfby/ui/api/schema scale,optional"`
	//The element's image
	//sprite
	sprite *pixel.Sprite
}

//Function to get the background's
//XML field
func (e *ImageImpl) GetImageField() string { return e.ImageField }

//Function to get the element's
//image scale options
func (e *ImageImpl) GetImageScale() util.ScaleOption { return e.ImageScale }

//Function to get the background's
//sprite
func (e *ImageImpl) GetImageSprite() *pixel.Sprite { return e.sprite }

//Function to set the background's
//sprite
func (e *ImageImpl) SetImageSprite(s *pixel.Sprite) { e.sprite = s }

//Function to determine whether the
//image has been initialised, by whether
//its image sprite has been set (assuming
//it's meant to be set). This function
//doesn't call element.IsInitialised
func (e *ImageImpl) IsInitialised() bool {
	//If the element doesn't have an image
	return e.GetImageField() == "" ||
		//Or the image has been initialised
		e.GetImageSprite() != nil
}

//Function to initialise an element's
//image. Doesn't call element.Init
func InitImage(e HasImage) error {
	//If the image hasn't been made yet
	if e.GetImageSprite() == nil {
		//Load the image
		picture, err := util.CreatePictureFromField(e.GetFS(), e.GetImageField())
		if err != nil {
			return err
		}
		if picture != nil {
			e.SetImageSprite(pixel.NewSprite(picture, picture.Bounds()))
		}
	}

	//If the sprite has been created
	if e.GetImageSprite() != nil {
		//If the element's width isn't known
		//and the width should be the content
		if e.GetActualWidth() == nil && e.GetRelWidth().MatchContent {
			//Set the actual width as the size of the image
			newWidth := e.GetImageSprite().Frame().Size().X
			e.SetActualWidth(&newWidth)
		}

		//If the element's height isn't known
		//and the height should be the content
		if e.GetActualHeight() == nil && e.GetRelHeight().MatchContent {
			//Set the actual height as the size of the image
			newHeight := e.GetImageSprite().Frame().Size().Y
			e.SetActualHeight(&newHeight)
		}

		//todo factor in scale options
	}

	return nil
}

//Function to draw an element's
//image
func DrawImage(e HasImage) {
	//Draw the image
	util.DrawSprite(e.GetCanvas(), e.GetImageSprite(),
		e.GetImageScale(), e.GetGravity())
}
