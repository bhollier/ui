package element

import (
	"encoding/xml"
	"errors"
	"os"
)

//Type for the root of a UI XML document
type Root struct {
	//The root element's parent
	//(if it has one)
	parent IsLayout
	//The root element itself
	IsElement
}

//Function to create a new design from
//an XML string
func NewRoot(parent IsLayout, path string) (e *Root, err error) {
	//Create a new root struct
	e = new(Root)
	e.parent = parent

	//Open the file
	file, err := os.Open(path)
	if err != nil {return nil, err}

	//Create an xml decoder
	d := xml.NewDecoder(file)
	//Decode into this element
	err = d.Decode(e)
	if err != nil {return nil, err}

	return
}

//Function to unmarshal an XML element into
//a root element. This function is usually
//only called by xml.Unmarshal
func (e *Root) UnmarshalXML(decoder* xml.Decoder, start xml.StartElement) error {
	//Create an element of the type
	e.IsElement = New(start.Name, e.parent)
	//If the element was created
	if e.IsElement != nil {
		//Decode the XML element into it
		err := decoder.DecodeElement(e.IsElement, &start)
		if err != nil {return err}
		return nil
	} else {
		return errors.New("unknown element type '" +
			XMLNameToString(start.Name) + "'")
	}
}


