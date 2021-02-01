package builtin

import (
	"encoding/xml"
	"errors"
	"github.com/bhollier/ui/pkg/ui/element"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"net/http"
)

// Element type for a text
type Text struct {
	// A text element is an
	// element
	element.Impl

	// It also has text
	element.TextImpl

	// The window the element
	// is in
	window *pixelgl.Window
	// The parent bounds
	parentBounds *pixel.Rect
}

// Function to create a new text
func NewText(fs http.FileSystem, name xml.Name, parent element.Layout) element.Element {
	return &Text{Impl: element.NewElement(fs, name, parent)}
}

// The XML name of the element
var TextTypeName = xml.Name{Space: "http://github.com/bhollier/ui/api/schema", Local: "Text"}

// Function to set the text content
func (e *Text) SetText(s string) error {
	// Set the text
	err := e.TextImpl.SetText(s)
	if err != nil {
		return err
	}

	// If the width depends on the content
	if e.GetRelWidth().MatchContent ||
		e.GetRelHeight().MatchContent {
		// Reset the element
		e.Reset()
		// Update the width and/or height
		if e.GetRelWidth().MatchContent {
			newWidth := e.GetSprite().Bounds().Size().X
			e.SetActualWidth(&newWidth)
		}
		if e.GetRelHeight().MatchContent {
			newHeight := e.GetSprite().Bounds().Size().Y
			e.SetActualHeight(&newHeight)
		}
		// Re-initialise the element
		err = e.Init(e.window, e.parentBounds)
		if err != nil {
			return err
		}
	}
	return nil
}

// Function to unmarshal an XML element into
// an element. This function is usually only
// called by xml.Unmarshal
func (e *Text) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	// Unmarshal the element
	err = e.Impl.UnmarshalXML(d, start)
	if err != nil {
		return err
	}
	// Set the element's attributes
	err = element.SetAttrs(e, start.Attr)
	if err != nil {
		return err
	}
	return d.Skip()
}

// Function to reset the element
func (e *Text) Reset() {
	e.Impl.Reset()
	e.TextImpl.Reset()
}

// Function to determine whether
// the element is initialised
func (e *Text) IsInitialised() bool {
	// If the element is initialised
	return e.Impl.IsInitialised() &&
		// And the image has been initialised
		e.TextImpl.IsInitialised()
}

// Function to initialise the element
// (load textures, create sprites, set
// sprite locations, etc.)
func (e *Text) Init(window *pixelgl.Window, bounds *pixel.Rect) error {
	// Save the window and parent bounds
	e.window = window
	e.parentBounds = bounds

	// Initialise the element
	err := e.Impl.Init(window, bounds)
	if err != nil {
		return err
	}

	// Initialise the text
	err = element.InitText(e, &e.TextImpl)
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

	return nil
}

// Function to draw the element
func (e *Text) Draw() {
	// Draw the element
	e.Impl.Draw()
	// Draw the text
	element.DrawText(e, &e.TextImpl)
}
