package button

import (
	"encoding/xml"
	"errors"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/orfby/ui/pkg/ui/element"
)

//Element type for an image button
type ImageButton struct {
	//An image button is a button
	element.ButtonImpl

	//It has an image
	element.ImageImpl
}

//Function to create a new image button
func NewImageButton(name xml.Name, parent element.Layout) element.Element {
	return &ImageButton{ButtonImpl: element.NewButton(name, parent)}
}

//The XML name of the element
var ImageButtonTypeName = xml.Name{Space: "http://github.com/orfby/ui/api/schema", Local: "ImageButton"}

//Function to unmarshal an XML element into
//an element. This function is usually only
//called by xml.Unmarshal
func (e *ImageButton) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	//Unmarshal the button
	err = e.ButtonImpl.UnmarshalXML(d, start)
	if err != nil {
		return err
	}

	//Set the element's attributes
	err = element.SetAttrs(e, start.Attr)
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

	return d.Skip()
}

//Function to determine whether
//the element is initialised
func (e *ImageButton) IsInitialised() bool {
	return e.ButtonImpl.IsInitialised() &&
		e.ImageImpl.IsInitialised()
}

//Function to initialise the element
func (e *ImageButton) Init(window *pixelgl.Window, bounds *pixel.Rect) error {
	//Initialise the button
	err := e.ButtonImpl.Init(window, bounds)
	if err != nil {
		return err
	}

	//Initialise the button's image
	err = element.InitImage(e)
	if err != nil {
		return err
	}
	return nil
}

//Function that is called when there
//is a new event
func (e *ImageButton) NewEvent(window *pixelgl.Window) {
	//Call the button's new event
	element.ButtonNewEvent(e, window)
}

//Function to draw the element
func (e *ImageButton) Draw() {
	//Draw the button
	e.ButtonImpl.Draw()
	//Draw the image
	element.DrawImage(e)
}
