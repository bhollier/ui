package relative

import (
	"encoding/xml"
	"errors"
	"github.com/orfby/ui/pkg/ui/element"
)

//Wrapper type that stores an
//element.Element and its relative fields
type relativeElement struct {
	//The parent element
	parent element.Layout

	//The attribute for specifying which
	//element this element goes on top of
	TopOf relativePosition `uixml:"http://github.com/orfby/ui/api/schema top-of,optional"`
	//The attribute for specifying which
	//element this element goes on the
	//bottom of
	BottomOf relativePosition `uixml:"http://github.com/orfby/ui/api/schema bottom-of,optional"`
	//The attribute for specifying which
	//element this element goes to the
	//left of
	LeftOf relativePosition `uixml:"http://github.com/orfby/ui/api/schema left-of,optional"`
	//The attribute for specifying which
	//element this element goes to the
	//right of
	RightOf relativePosition `uixml:"http://github.com/orfby/ui/api/schema right-of,optional"`

	//The element itself (the "hidden" tag
	//means element.SetAttrs won't touch it)
	element.Element `uixml:"hidden"`
}

//Function to unmarshal an XML element into
//a relative element. This function is
//only called by xml.Unmarshal
func (e *relativeElement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	//Create an element of the type
	e.Element = element.New(start.Name, e.parent)
	//If the element wasn't created
	if e.Element == nil {
		return errors.New("unknown element type '" +
			element.XMLNameToString(start.Name) + "'")
	}

	//Create an array of the relative
	//attributes
	relativetAttrs := make([]xml.Attr, 0)

	//The names of the relative attributes
	topOfName := xml.Name{Space: "http://github.com/orfby/ui/api/schema", Local: "top-of"}
	bottomOfName := xml.Name{Space: "http://github.com/orfby/ui/api/schema", Local: "bottom-of"}
	leftOfName := xml.Name{Space: "http://github.com/orfby/ui/api/schema", Local: "left-of"}
	rightOfName := xml.Name{Space: "http://github.com/orfby/ui/api/schema", Local: "right-of"}
	//Iterate over the attributes
	for _, attr := range start.Attr {
		//If the attribute isn't a relative attribute
		if element.XMLNameMatch(attr.Name, topOfName) ||
			element.XMLNameMatch(attr.Name, bottomOfName) ||
			element.XMLNameMatch(attr.Name, leftOfName) ||
			element.XMLNameMatch(attr.Name, rightOfName) {
			//Add it to the attributes array
			relativetAttrs = append(relativetAttrs, attr)
		}
	}

	//Set the relative attributes
	err = element.SetAttrs(e, relativetAttrs)
	if err != nil {
		return err
	}

	//If both top and bottom are set
	if e.TopOf != zeroRelativePosition && e.BottomOf != zeroRelativePosition {
		return errors.New("both 'top-of' and 'bottom-of' attributes set on XML element '" +
			element.FullName(e, ".", false) + "'")
	}
	//If both left and right are set
	if e.LeftOf != zeroRelativePosition && e.RightOf != zeroRelativePosition {
		return errors.New("both 'left-of' and 'right-of' attributes set on XML element '" +
			element.FullName(e, ".", false) + "'")
	}

	//If none of the attributes are set
	if e.TopOf == zeroRelativePosition && e.BottomOf == zeroRelativePosition &&
		e.LeftOf == zeroRelativePosition && e.RightOf == zeroRelativePosition {
		return errors.New("XML element '" + element.FullName(e, ".", false) +
			"' has no position attribute, must have at least 'top-of', 'bottom-of', 'left-of' or 'right-of'")
	}

	//Create an array of the element's
	//attributes
	elementAttrs := make([]xml.Attr, 0)

	//Iterate over the attributes
	for _, attr := range start.Attr {
		//If the attribute isn't a relative attribute
		if !element.XMLNameMatch(attr.Name, topOfName) &&
			!element.XMLNameMatch(attr.Name, bottomOfName) &&
			!element.XMLNameMatch(attr.Name, leftOfName) &&
			!element.XMLNameMatch(attr.Name, rightOfName) {
			//Add it to the attributes array
			elementAttrs = append(elementAttrs, attr)
		}
	}

	//Replace the attributes
	start.Attr = elementAttrs

	//Unmarshal the element itself
	return e.Element.UnmarshalXML(d, start)
}
