package relative

import (
	"encoding/xml"
	"errors"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/orfby/ui/pkg/ui/element"
	"github.com/orfby/ui/pkg/ui/util"
)

//Layout type for displaying elements
//relative to each other
type Layout struct {
	//A relative layout is an
	//element
	element.Impl

	//The layout's child elements
	//(in order)
	children []relativeElement
}

//Function to create a new relative layout
func NewLayout(name xml.Name, parent element.Layout) element.Element {
	return &Layout{Impl: element.NewElement(name, parent)}
}

//The XML name of the element
var LayoutTypeName = xml.Name{Space: "http://github.com/orfby/ui/api/schema", Local: "RelativeLayout"}

//Function to get one of a layout's
//child elements
func (e *Layout) GetChild(n int) element.Element { return &e.children[n] }

//Function to get the number element
//elements a layout has
func (e *Layout) NumChildren() int { return len(e.children) }

//Function to get one of a layout's child
//elements by its ID. Returns nil if no
//child could be found
func (e *Layout) GetChildByID(id string) element.Element {
	for _, child := range e.children {
		if child.GetID() != nil && *child.GetID() == id {
			return &child
		}
	}
	return nil
}

//Function to unmarshal an XML element into
//an element. This function is usually only
//called by xml.Unmarshal
func (e *Layout) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
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

	//If the layout is meant to match
	//the content, throw an error
	//todo maybe allow the relativelayout to match content
	if e.GetRelWidth().MatchContent {
		return errors.New("invalid width attribute value 'match_content' on XML element '" +
			element.FullName(e, ".", false) + "'")
	} else if e.GetRelHeight().MatchContent {
		return errors.New("invalid height attribute value 'match_content' on XML element '" +
			element.FullName(e, ".", false) + "'")
	}

	//Create the array of children
	e.children = make([]relativeElement, 0)
	//Loop over the child xml elements
Loop:
	for {
		//Get the next token
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch tt := t.(type) {
		//If this is the start of an element
		case xml.StartElement:
			//Create a relative element
			elem := relativeElement{parent: e}
			//Decode the XML element into it
			err = d.DecodeElement(&elem, &tt)
			if err != nil {
				return err
			}
			//Add it to the children array
			e.children = append(e.children, elem)

			//If this is the end of the element
		case xml.EndElement:
			if tt == start.End() {
				break Loop
			}
		}
	}

	//Iterate over the children
	for _, child := range e.children {
		//If the top of element exists but the ID leads nowhere
		if child.TopOf != zeroRelativePosition && child.TopOf.ElementID != "" &&
			e.GetChildByID(child.TopOf.ElementID) == nil {
			return element.NewNoElemError(child.Element, child.TopOf.ElementID, "top-of")
		}
		//If the bottom of element exists but the ID leads nowhere
		if child.BottomOf != zeroRelativePosition && child.BottomOf.ElementID != "" &&
			e.GetChildByID(child.BottomOf.ElementID) == nil {
			return element.NewNoElemError(child.Element, child.BottomOf.ElementID, "bottom-of")
		}
		//If the left of element exists but the ID leads nowhere
		if child.LeftOf != zeroRelativePosition && child.LeftOf.ElementID != "" &&
			e.GetChildByID(child.LeftOf.ElementID) == nil {
			return element.NewNoElemError(child.Element, child.LeftOf.ElementID, "left-of")
		}
		//If the right of element exists but the ID leads nowhere
		if child.RightOf != zeroRelativePosition && child.RightOf.ElementID != "" &&
			e.GetChildByID(child.RightOf.ElementID) == nil {
			return element.NewNoElemError(child.Element, child.RightOf.ElementID, "right-of")
		}

		//todo  make sure there are no circular
		//todo  dependencies
	}

	return nil
}

//Function to reset the element
func (e *Layout) Reset() {
	e.Impl.Reset()
	for _, child := range e.children {
		child.Reset()
	}
}

//Function to determine whether
//the element is initialised
func (e *Layout) IsInitialised() bool {
	return e.Impl.IsInitialised() &&
		element.ChildrenAreInitialised(e)
}

//Function to initialise the element
func (e *Layout) Init(window *pixelgl.Window, bounds *pixel.Rect) error {
	//Initialise the element part of the layout
	err := e.Impl.Init(window, bounds)
	if err != nil {
		return err
	}

	//Iterate over the elements
	for _, child := range e.children {
		//Create the child's bounds
		childBounds := &pixel.Rect{Min: pixel.V(-1, -1), Max: pixel.V(-1, -1)}
		ySet := false
		xSet := false

		if bounds == nil {
			childBounds = nil
		} else {
			//If the child's width and height are known
			if child.GetActualWidth() != nil && child.GetActualHeight() != nil {
				//If the child has a top of attribute
				if child.TopOf != zeroRelativePosition {
					//If the child is at the top of the parent
					if child.TopOf.Parent {
						//Set the max Y of the child bounds
						childBounds.Max.Y = bounds.Max.Y
						childBounds.Min.Y = childBounds.Max.Y -
							*child.GetActualHeight()

						//If it's relative to a position
					} else if child.TopOf.Pos != util.ZeroRelativeQuantity {
						//If the position is a percentage
						if child.TopOf.Pos.Unit == util.Percent {
							childBounds.Max.Y = bounds.Max.Y - bounds.Size().Y*
								(float64(child.TopOf.Pos.Quantity)/100)
							childBounds.Min.Y = childBounds.Max.Y +
								*child.GetActualHeight()

							//If the position is just in pixels
						} else {
							childBounds.Max.Y = bounds.Max.Y - float64(child.TopOf.Pos.Quantity)
							childBounds.Min.Y = childBounds.Max.Y +
								*child.GetActualHeight()
						}

						//If it's relative to an element
					} else {
						//Get the element
						relativeElem := e.GetChildByID(child.TopOf.ElementID)
						//This shouldn't happen, but check it anyways
						if relativeElem == nil {
							return element.NewNoElemError(child.Element, child.TopOf.ElementID, "top-of")
						}
						if relativeElem.GetMin() != nil &&
							relativeElem.GetMax() != nil {
							//Set the bounds to be the same as the element's
							childBounds.Max.Y = relativeElem.GetMax().Y
							childBounds.Min.Y = childBounds.Max.Y +
								*child.GetActualHeight()

							//Set the X bounds as well
							childBounds.Min.X = relativeElem.GetMin().X
							childBounds.Max.X = relativeElem.GetMax().X
							xSet = true
						} else {
							childBounds = nil
						}
					}
					ySet = true

					//If the child has a bottom of of attribute
				} else if child.BottomOf != zeroRelativePosition {
					//If the child is at the bottom of the parent
					if child.BottomOf.Parent {
						//Set the min Y of the child bounds
						childBounds.Min.Y = bounds.Min.Y
						childBounds.Max.Y = childBounds.Min.Y +
							*child.GetActualHeight()

						//If it's relative to a position
					} else if child.BottomOf.Pos != util.ZeroRelativeQuantity {
						//If the position is a percentage
						if child.BottomOf.Pos.Unit == util.Percent {
							childBounds.Min.Y = bounds.Max.Y - bounds.Size().Y*
								(float64(child.BottomOf.Pos.Quantity)/100)
							childBounds.Max.Y = childBounds.Min.Y -
								*child.GetActualHeight()

							//If the position is just in pixels
						} else {
							childBounds.Max.Y = bounds.Max.Y - float64(child.BottomOf.Pos.Quantity)
							childBounds.Min.Y = childBounds.Max.Y -
								*child.GetActualHeight()
						}

						//If it's relative to an element
					} else {
						//Get the element
						relativeElem := e.GetChildByID(child.BottomOf.ElementID)
						//This shouldn't happen, but check it anyways
						if relativeElem == nil {
							return element.NewNoElemError(child.Element, child.BottomOf.ElementID, "bottom-of")
						}
						if relativeElem.GetMin() != nil &&
							relativeElem.GetMax() != nil {
							//Set the bounds to be the same as the element's
							childBounds.Min.Y = relativeElem.GetMin().Y
							childBounds.Max.Y = childBounds.Min.Y -
								*child.GetActualHeight()

							//Set the X bounds as well
							childBounds.Min.X = relativeElem.GetMin().X
							childBounds.Max.X = relativeElem.GetMax().X
							xSet = true
						} else {
							childBounds = nil
						}
					}
					ySet = true
				}

				//If the child bounds can still be calculated
				if childBounds != nil {
					//If the child has a left of attribute
					if child.LeftOf != zeroRelativePosition {
						//If the child is at the left of the parent
						if child.LeftOf.Parent {
							//Set the max X of the child bounds
							childBounds.Min.X = bounds.Min.X
							childBounds.Max.X = childBounds.Min.X +
								*child.GetActualWidth()

							//If it's relative to a position
						} else if child.LeftOf.Pos != util.ZeroRelativeQuantity {
							//If the position is a percentage
							if child.LeftOf.Pos.Unit == util.Percent {
								childBounds.Max.X = bounds.Min.X + bounds.Size().X*
									(float64(child.LeftOf.Pos.Quantity)/100)
								childBounds.Min.X = childBounds.Max.X -
									*child.GetActualWidth()

								//If the position is just in pixels
							} else {
								childBounds.Max.X = bounds.Min.X + float64(child.LeftOf.Pos.Quantity)
								childBounds.Min.X = childBounds.Max.X -
									*child.GetActualWidth()
							}

							//If it's relative to an element
						} else {
							//Get the element
							relativeElem := e.GetChildByID(child.LeftOf.ElementID)
							//This shouldn't happen, but check it anyways
							if relativeElem == nil {
								return element.NewNoElemError(child.Element, child.LeftOf.ElementID, "left-of")
							}
							if relativeElem.GetMin() != nil &&
								relativeElem.GetMax() != nil {
								//Set the bounds to be the same as the element's
								childBounds.Max.X = relativeElem.GetMin().X
								childBounds.Min.X = childBounds.Max.X -
									*child.GetActualWidth()

								//If the Y bounds aren't set
								if !ySet {
									//Set the Y bounds
									childBounds.Min.Y = relativeElem.GetMin().Y
									childBounds.Max.Y = relativeElem.GetMax().Y
									ySet = true
								}

							} else {
								childBounds = nil
							}
						}
						xSet = true

						//If the child has a right of of attribute
					} else if child.RightOf != zeroRelativePosition {
						//If the child is at the right of the parent
						if child.RightOf.Parent {
							//Set the min X of the child bounds
							childBounds.Max.X = bounds.Max.X
							childBounds.Min.X = childBounds.Max.X -
								*child.GetActualWidth()

							//If it's relative to a position
						} else if child.RightOf.Pos != util.ZeroRelativeQuantity {
							//If the position is a percentage
							if child.RightOf.Pos.Unit == util.Percent {
								childBounds.Min.X = bounds.Min.X + bounds.Size().X*
									(float64(child.RightOf.Pos.Quantity)/100)
								childBounds.Max.X = childBounds.Min.X +
									*child.GetActualWidth()

								//If the position is just in pixels
							} else {
								childBounds.Min.X = bounds.Min.X + float64(child.RightOf.Pos.Quantity)
								childBounds.Max.X = childBounds.Max.X +
									*child.GetActualWidth()
							}

							//If it's relative to an element
						} else {
							//Get the element
							relativeElem := e.GetChildByID(child.RightOf.ElementID)
							//This shouldn't happen, but check it anyways
							if relativeElem == nil {
								return element.NewNoElemError(child.Element, child.BottomOf.ElementID, "right-of")
							}
							if relativeElem.GetMin() != nil &&
								relativeElem.GetMax() != nil {
								//Set the bounds to be the same as the element's
								childBounds.Min.X = relativeElem.GetMax().X
								childBounds.Max.X = childBounds.Min.X +
									*child.GetActualWidth()

								//If the Y bounds aren't set
								if !ySet {
									//Set the Y bounds
									childBounds.Min.Y = relativeElem.GetMin().Y
									childBounds.Max.Y = relativeElem.GetMax().Y
									ySet = true
								}
							} else {
								childBounds = nil
							}
						}
						xSet = true
					}
				}
			} else {
				childBounds = nil
			}

			//If the child bounds can still be calculated
			if childBounds != nil {
				//If the Y wasn't set
				if !ySet {
					//Set the Y child bounds as
					//the whole height of the parent
					childBounds.Min.Y = bounds.Min.Y
					childBounds.Max.Y = bounds.Max.Y
				}
				//If the X wasn't set
				if !xSet {
					//set the X child bounds as
					//the whole width of the parent
					childBounds.Min.X = bounds.Min.X
					childBounds.Max.X = bounds.Max.X
				}
			}
		}

		//Initialise the child
		err := child.Init(window, childBounds)
		if err != nil {
			return err
		}
	}

	return nil
}

//Function that is called when there
//is a new event. This function only
//calls NewEvent on the child elements
func (e *Layout) NewEvent(window *pixelgl.Window) {
	e.Impl.NewEvent(window)
	for _, child := range e.children {
		child.NewEvent(window)
	}
}

//Function to draw the element
func (e *Layout) Draw() {
	//Draw the element
	e.Impl.Draw()
	//Draw the layout
	element.DrawLayout(e)
}
