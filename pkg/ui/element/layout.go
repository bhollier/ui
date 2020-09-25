package element

import (
	"encoding/xml"
	"errors"
	"github.com/faiface/pixel/pixelgl"
)

//Interface type for a layout element
type IsLayout interface {
	//A layout is also an element
	IsElement

	//Function to get one of a layout's child
	//elements
	GetChild(n int) IsElement
	//Function to get one of a layout's child
	//elements by its ID
	GetChildByID(id string) IsElement
	//Function to get the number of child
	//elements this layout has
	NumChildren() int
}

//Type for a layout
type Layout struct {
	//The layout's child elements
	//(in order)
	Children []IsElement
}

//Function to get one of a layout's
//child elements
func (e *Layout) GetChild(n int) IsElement {return e.Children[n]}
//Function to get the number element
//elements a layout has
func (e *Layout) NumChildren() int {return len(e.Children)}
//Function to get one of a layout's child
//elements by its ID. Returns nil if no
//child could be found
func (e *Layout) GetChildByID(id string) IsElement {
	for _, child := range e.Children {
		if child.GetID() != nil && *child.GetID() == id {return child}
	}
	return nil
}

//Function to reset the element
func (e *Layout) Reset() {
	for _, child := range e.Children {
		child.Reset()
	}
}

//Function to unmarshal an XML element into
//a number of child elements. This function
//is usually only called by
//Layout.UnmarshalXML
func ChildrenUnmarshalXML(parent IsLayout, d* xml.Decoder,
	start xml.StartElement) ([]IsElement, error) {
	//Iterate over the XML tokens
	children := make([]IsElement, 0)
	for {
		//Get the next token
		t, err := d.Token()
		if err != nil {
			return nil, err
		}
		var elem IsElement
		switch tt := t.(type) {
		//If this is the start of an element
		case xml.StartElement:
			//Create an element of the type
			elem = New(tt.Name, parent)
			//If the element was created
			if elem != nil {
				//Decode the XML element into it
				err = d.DecodeElement(elem, &tt)
				if err != nil {
					return nil, err
				}
				//Add it to the children array
				children = append(children, elem)
				elem = nil
			} else {
				return nil, errors.New("unknown element type '" +
					XMLNameToString(tt.Name) + "'")
			}
			//If this is the end of the element
		case xml.EndElement:
			if tt == start.End() {
				return children, nil
			}
		}
	}
}

//Function to determine whether a layout's
//children have been initialised. This
//function doesn't call element.IsInitialised
func ChildrenAreInitialised(e IsLayout) bool {
	//Iterate over the children
	for i := 0; i < e.NumChildren(); i++ {
		//If the child hasn't been initialised
		if !e.GetChild(i).IsInitialised() {
			//The layout hasn't been initialised,
			//return false
			return false
		}
	}
	//Return whether the layout itself
	//has been initialised
	return true
}

//Function that is called when there
//is a new event. This function only
//calls NewEvent on the child elements
func (e *Layout) NewEvent(window *pixelgl.Window) {
	for _, child := range e.Children {
		child.NewEvent(window)
	}
}

//Function to draw a layout
func DrawLayout(e IsLayout) {
	//Iterate over the children
	for i := 0; i < e.NumChildren(); i++ {
		//Draw the child
		e.GetChild(i).Draw()
		//Draw the child onto the layout's canvas
		DrawCanvasOntoParent(
			e.GetChild(i).GetCanvas(), e.GetCanvas())
	}
}