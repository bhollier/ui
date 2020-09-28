package builtin

import (
	"encoding/xml"
	"errors"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/orfby/ui/pkg/ui/element"
	"net/http"
)

//Element type for an image
type Image struct {
	//An image element is an
	//element
	element.Impl

	//It also has an image
	element.ImageImpl
}

//Function to create a new image
func NewImage(fs http.FileSystem, name xml.Name, parent element.Layout) element.Element {
	return &Image{Impl: element.NewElement(fs, name, parent)}
}

//The XML name of the element
var ImageTypeName = xml.Name{Space: "http://github.com/orfby/ui/api/schema", Local: "Image"}

//Function to unmarshal an XML element into
//an element. This function is usually only
//called by xml.Unmarshal
func (e *Image) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	//Unmarshal the element
	err = e.Impl.UnmarshalXML(d, start)
	if err != nil {
		return err
	}
	//Set the element's attributes
	err = element.SetAttrs(e, start.Attr)
	if err != nil {
		return err
	}
	return d.Skip()
}

//Function to determine whether
//the element is initialised
func (e *Image) IsInitialised() bool {
	//If the element is initialised
	return e.Impl.IsInitialised() &&
		//And the image has been initialised
		e.ImageImpl.IsInitialised()
}

//Function to initialise the element
//(load textures, create sprites, set
//sprite locations, etc.)
func (e *Image) Init(window *pixelgl.Window, bounds *pixel.Rect) error {
	//Initialise the element
	err := e.Impl.Init(window, bounds)
	if err != nil {
		return err
	}

	//Initialise the image
	err = element.InitImage(e)
	if err != nil {
		return err
	}

	//If no image was given
	if e.GetImageField() == "" {
		//If it wants to match the content
		if e.GetRelWidth().MatchContent {
			return errors.New("invalid width attribute value 'match_content' on XML element '" +
				element.FullName(e, ".", false) +
				"': no content to match")
		} else if e.GetRelHeight().MatchContent {
			return errors.New("invalid height attribute value 'match_content' on XML element '" +
				element.FullName(e, ".", false) +
				"': no content to match")
		}
	}

	return nil
}

//Function to draw the element
func (e *Image) Draw() {
	//Draw the element
	e.Impl.Draw()
	//Draw the image
	element.DrawImage(e)
}
