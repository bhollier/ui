package button

import (
	"encoding/xml"
	"errors"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/orfby/ui/pkg/ui/element"
	"net/http"
)

// Element type for a text button
type TextButton struct {
	// A text button is a button
	element.ButtonImpl

	// It has text
	element.TextImpl
}

// Function to create a new text button
func NewTextButton(fs http.FileSystem, name xml.Name, parent element.Layout) element.Element {
	return &TextButton{ButtonImpl: element.NewButton(fs, name, parent)}
}

// The XML name of the element
var TextButtonTypeName = xml.Name{Space: "http://github.com/orfby/ui/api/schema", Local: "TextButton"}

// Function to unmarshal an XML element into
// an element. This function is usually only
// called by xml.Unmarshal
func (e *TextButton) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	// Unmarshal the button
	err = e.ButtonImpl.UnmarshalXML(d, start)
	if err != nil {
		return err
	}

	// Set the element's attributes
	err = element.SetAttrs(e, start.Attr)
	if err != nil {
		return err
	}

	// If no text was given
	if e.TextImpl.GetField() == "" {
		// If it wants to match the content
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

// Function to reset the element
func (e *TextButton) Reset() {
	e.ButtonImpl.Reset()
	e.TextImpl.Reset()
}

// Function to determine whether
// the element is initialised
func (e *TextButton) IsInitialised() bool {
	return e.ButtonImpl.IsInitialised() &&
		e.TextImpl.IsInitialised()
}

// Function to initialise the element
func (e *TextButton) Init(window *pixelgl.Window, bounds *pixel.Rect) error {
	// Initialise the button
	err := e.ButtonImpl.Init(window, bounds)
	if err != nil {
		return err
	}

	// Initialise the button's text
	err = element.InitText(e, &e.TextImpl, bounds)
	if err != nil {
		return err
	}
	return nil
}

// Function that is called when there
// is a new event
func (e *TextButton) NewEvent(window *pixelgl.Window) {
	// Call the button's new event
	element.ButtonNewEvent(e, window)
}

// Function to draw the element
func (e *TextButton) Draw() {
	// Draw the button
	e.ButtonImpl.Draw()
	// Draw the text
	element.DrawText(e, &e.TextImpl)
}
