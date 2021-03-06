package element

import (
	"encoding/xml"
	"errors"
	"net/http"
)

// Type for the root of a UI XML document
type Root struct {
	// The root element's parent
	// (if it has one)
	parent Layout
	// The filesystem to use
	fs http.FileSystem
	// The root element itself
	Element
}

// Function to create a new design from
// an XML string
func NewRoot(fs http.FileSystem, parent Layout, path string) (e *Root, err error) {
	// Create a new root struct
	e = new(Root)
	e.parent = parent
	e.fs = fs

	// Open the file
	file, err := fs.Open(path)
	if err != nil {
		return nil, err
	}

	// Create an xml decoder
	d := xml.NewDecoder(file)
	// Decode into this element
	err = d.Decode(e)
	if err != nil {
		return nil, err
	}

	return
}

// Function to unmarshal an XML element into
// a root element. This function is usually
// only called by xml.Unmarshal
func (e *Root) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	// Create an element of the type
	e.Element = New(e.fs, start.Name, e.parent)
	// If the element was created
	if e.Element != nil {
		// Decode the XML element into it
		err := decoder.DecodeElement(e.Element, &start)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("unknown element type '" +
			XMLNameToString(start.Name) + "'")
	}
}
