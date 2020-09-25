package builtin

import (
	"encoding/xml"
	"errors"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/orfby/ui/pkg/ui/element"
	"log"
)

//Type for an element that scrolls
type Scroll struct {
	//The scroll is an element
	element.Impl
	//It is also (technically) a layout
	element.LayoutImpl

	//The scroll rate
	ScrollRate uint `uixml:"http://github.com/orfby/ui/api/schema scroll-speed,optional"`

	//The parent bounds
	parentBounds *pixel.Rect
	//The child's bounds
	childBounds *pixel.Rect
}

//Function to create a new import element
func NewScroll(name xml.Name, parent element.Layout) element.Element {
	return &Scroll{Impl: element.NewElement(name, parent)}
}

//The XML name of the import element
var ScrollTypeName = xml.Name{Space: "http://github.com/orfby/ui/api/schema", Local: "Scroll"}

//Function to unmarshal an XML element into
//an element. This function is usually only
//called by xml.Unmarshal
func (e *Scroll) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	//Unmarshal the element part of the layout
	err = e.Impl.UnmarshalXML(d, start)
	if err != nil {
		return err
	}

	//Set the element's attributes
	err = element.SetAttrs(e, start.Attr)
	if err != nil {
		return err
	}

	//If the scroll rate wasn't given
	if e.ScrollRate == 0 {
		//Set it to the default (10)
		e.ScrollRate = 10
	}

	//Unmarshal the layout's children
	e.LayoutImpl.Children, err = element.ChildrenUnmarshalXML(e, d, start)
	if err != nil {
		return err
	}

	//If there are no children
	if len(e.LayoutImpl.Children) == 0 {
		return errors.New("no children on XML element '" +
			element.FullName(e, ".", true) + "'")

		//If there are multiple
	} else if len(e.LayoutImpl.Children) > 1 {
		return errors.New("multiple children on XML element '" +
			element.FullName(e, ".", true) + "'")
	}

	return nil
}

//Function to reset the element
func (e *Scroll) Reset() {
	e.Impl.Reset()
	e.LayoutImpl.Reset()
	e.parentBounds = nil
	e.childBounds = nil
}

//Function to determine whether
//the element is initialised
func (e *Scroll) IsInitialised() bool {
	return e.Impl.IsInitialised() &&
		element.ChildrenAreInitialised(e)
}

//Function to initialise an element's
//position, width and height. Because
//it doesn't know the element's actual
//size, it won't set the width or height
//if the relative width or height is
//"match_content"
func (e *Scroll) Init(window *pixelgl.Window, bounds *pixel.Rect) (err error) {
	//Initialise the element part of the import
	err = e.Impl.Init(window, bounds)
	if err != nil {
		return err
	}

	//If the parent bounds aren't known
	if e.parentBounds == nil && bounds != nil {
		e.parentBounds = new(pixel.Rect)
		e.childBounds = new(pixel.Rect)
		*e.parentBounds = *bounds
		*e.childBounds = *bounds
	}

	//Initialise the child
	err = e.GetChild(0).Init(window, e.childBounds)
	if err != nil {
		return err
	}

	//If the width is meant to match the content size
	if e.GetRelWidth().MatchContent {
		//Set the width as the child's
		e.SetActualWidth(e.GetChild(0).GetActualWidth())
	}

	//If the height is meant to match the content size
	if e.GetRelHeight().MatchContent {
		//Set the height as the child's
		e.SetActualHeight(e.GetChild(0).GetActualHeight())
	}

	return nil
}

//Function that is called when there
//is a new event
func (e *Scroll) NewEvent(window *pixelgl.Window) {
	e.Impl.NewEvent(window)
	e.LayoutImpl.NewEvent(window)

	//Get the mouse scroll
	scroll := window.MouseScroll()
	//If the scroll wheel moved
	if scroll != pixel.V(0, 0) {
		//Move the child's bounds
		e.childBounds.Min.X += scroll.X * float64(e.ScrollRate)
		e.childBounds.Max.Y -= scroll.Y * float64(e.ScrollRate)

		//If the X child bounds are going too far
		if e.childBounds.Min.X > e.parentBounds.Min.X {
			e.childBounds.Min.X = e.parentBounds.Min.X
		} else if e.childBounds.Min.X < *e.Children[0].GetActualWidth() {
			e.childBounds.Min.X = *e.Children[0].GetActualWidth()
		}
		//If the Y child bounds are going too far
		if e.childBounds.Max.Y < e.parentBounds.Max.Y {
			e.childBounds.Max.Y = e.parentBounds.Max.Y
		} else if e.childBounds.Max.Y > *e.Children[0].GetActualHeight() {
			e.childBounds.Max.Y = *e.Children[0].GetActualHeight()
		}

		//Reset the child element
		e.Children[0].Reset()
		//While the child is uninitialised
		for !e.Children[0].IsInitialised() {
			//Initialise it
			err := e.Children[0].Init(window, e.childBounds)
			if err != nil {
				log.Printf("Error while initialsing XML element '"+
					element.FullName(e.Children[0], ".", true)+
					"': %+v", err)
			}
		}

		//Redraw the UI
		element.DrawUI(e, window)
	}
}

//Function to draw the element
func (e *Scroll) Draw() {
	//Draw the element
	e.Impl.Draw()
	//Draw the layout
	element.DrawLayout(e)
}
