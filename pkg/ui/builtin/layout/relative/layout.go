package relative

import (
	"encoding/xml"
	"errors"
	"github.com/bhollier/ui/pkg/ui/element"
	"github.com/bhollier/ui/pkg/ui/util"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"net/http"
)

// Layout type for displaying elements
// relative to each other
type Layout struct {
	// A relative layout is an
	// element
	element.Impl

	// The layout's child elements
	// (in order)
	children []relativeElement
}

// Function to create a new relative layout
func NewLayout(fs http.FileSystem, name xml.Name, parent element.Layout) element.Element {
	return &Layout{Impl: element.NewElement(fs, name, parent)}
}

// The XML name of the element
var LayoutTypeName = xml.Name{Space: "http://github.com/bhollier/ui/api/schema", Local: "RelativeLayout"}

// Function to get one of a layout's
// child elements
func (e *Layout) GetChild(n int) element.Element { return e.children[n].Element }

// Function to get the number element
// elements a layout has
func (e *Layout) NumChildren() int { return len(e.children) }

// Function to get one of a layout's child
// elements by its ID. Returns nil if no
// child could be found
func (e *Layout) GetChildByID(id string) element.Element {
	for _, child := range e.children {
		if child.GetID() != nil && *child.GetID() == id {
			return child.Element
		}
	}
	return nil
}

// Function to unmarshal an XML element into
// an element. This function is usually only
// called by xml.Unmarshal
func (e *Layout) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	// Unmarshal the element part of the layout
	err = e.Impl.UnmarshalXML(d, start)
	if err != nil {
		return err
	}

	// Set the element's attributes
	err = element.SetAttrs(e, start.Attr)
	if err != nil {
		return err
	}

	// If the layout is meant to match
	// the content, throw an error
	// todo maybe allow the relativelayout to match content
	if e.GetRelWidth().MatchContent {
		return errors.New("invalid width attribute value 'match_content' on XML element '" +
			element.FullName(e, ".", false) + "'")
	} else if e.GetRelHeight().MatchContent {
		return errors.New("invalid height attribute value 'match_content' on XML element '" +
			element.FullName(e, ".", false) + "'")
	}

	// Create the array of children
	e.children = make([]relativeElement, 0)
	// Loop over the child xml elements
Loop:
	for {
		// Get the next token
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch tt := t.(type) {
		// If this is the start of an element
		case xml.StartElement:
			// Create a relative element
			elem := newRelativeElement(e.GetFS(), e)
			// Decode the XML element into it
			err = d.DecodeElement(&elem, &tt)
			if err != nil {
				return err
			}
			// Add it to the children array
			e.children = append(e.children, elem)

			// If this is the end of the element
		case xml.EndElement:
			if tt == start.End() {
				break Loop
			}
		}
	}

	// Iterate over the children
	for _, child := range e.children {
		// If the top of element exists but the ID leads nowhere
		if child.TopOf != zeroRelativePosition && child.TopOf.ElementID != "" &&
			e.GetChildByID(child.TopOf.ElementID) == nil {
			return element.NewNoElemError(child.Element, child.TopOf.ElementID, "top-of")
		}
		// If the bottom of element exists but the ID leads nowhere
		if child.BottomOf != zeroRelativePosition && child.BottomOf.ElementID != "" &&
			e.GetChildByID(child.BottomOf.ElementID) == nil {
			return element.NewNoElemError(child.Element, child.BottomOf.ElementID, "bottom-of")
		}
		// If the left of element exists but the ID leads nowhere
		if child.LeftOf != zeroRelativePosition && child.LeftOf.ElementID != "" &&
			e.GetChildByID(child.LeftOf.ElementID) == nil {
			return element.NewNoElemError(child.Element, child.LeftOf.ElementID, "left-of")
		}
		// If the right of element exists but the ID leads nowhere
		if child.RightOf != zeroRelativePosition && child.RightOf.ElementID != "" &&
			e.GetChildByID(child.RightOf.ElementID) == nil {
			return element.NewNoElemError(child.Element, child.RightOf.ElementID, "right-of")
		}

		// todo  make sure there are no circular
		// todo  dependencies
	}

	return nil
}

// Function to reset the child
// element's positions
func (e *Layout) ResetPosition() {
	e.Impl.ResetPosition()
	for _, child := range e.children {
		child.ResetPosition()
	}
}

// Function to reset the child
// elements
func (e *Layout) Reset() {
	e.Impl.Reset()
	for _, child := range e.children {
		child.Reset()
	}
}

// Function to determine whether
// the element is initialised
func (e *Layout) IsInitialised() bool {
	return e.Impl.IsInitialised() &&
		element.ChildrenAreInitialised(e)
}

// Function to initialise the element
func (e *Layout) Init(window *pixelgl.Window, bounds *pixel.Rect) error {
	// Initialise the element part of the layout
	err := e.Impl.Init(window, bounds)
	if err != nil {
		return err
	}

	// Iterate over the elements
	for _, child := range e.children {
		// Create the child's bounds
		var childBounds *pixel.Rect
		// If the bounds are known, then the
		// child bounds' default is the parent
		if bounds != nil {
			childBounds = new(pixel.Rect)
			*childBounds = *bounds
		}

		xSet := false
		ySet := false

		// Create a function to modify the child bounds
		// for each of the child's relative position
		type location string
		const topOf = location("top-of")
		const bottomOf = location("bottom-of")
		const leftOf = location("left-of")
		const rightOf = location("right-of")
		modBounds := func(child *relativeElement, loc location, pos relativePosition) error {
			// If the child has a relative position attribute
			if pos != zeroRelativePosition {
				// If the child is aligned to the parent
				if pos.Parent {
					// Switch over the locations
					switch loc {
					case topOf:
						if !child.GetRelHeight().MatchBounds {
							if child.GetActualHeight() != nil {
								childBounds.Min.Y = childBounds.Max.Y -
									*child.GetActualHeight()
							} else {
								childBounds = nil
								return nil
							}
						}
						ySet = true
					case bottomOf:
						if !child.GetRelHeight().MatchBounds {
							if child.GetActualHeight() != nil {
								childBounds.Max.Y = childBounds.Min.Y +
									*child.GetActualHeight()
							} else {
								childBounds = nil
								return nil
							}
						}
						ySet = true
					case leftOf:
						if !child.GetRelWidth().MatchBounds {
							if child.GetActualWidth() != nil {
								childBounds.Max.X = childBounds.Min.X +
									*child.GetActualWidth()
							} else {
								childBounds = nil
								return nil
							}
						}
						xSet = true
					case rightOf:
						if !child.GetRelWidth().MatchBounds {
							if child.GetActualWidth() != nil {
								childBounds.Min.X = childBounds.Max.X -
									*child.GetActualWidth()
							} else {
								childBounds = nil
								return nil
							}
						}
						xSet = true
					}

					// If it's relative to a position
				} else if pos.Pos != util.ZeroRelativeQuantity {
					// Switch over the locations
					switch loc {
					case topOf:
						// If the position is a percentage
						if pos.Pos.Unit == util.Percent {
							childBounds.Max.Y = bounds.Max.Y - bounds.Size().Y*
								(float64(pos.Pos.Quantity)/100)
							// If the position is just in pixels
						} else {
							childBounds.Max.Y = bounds.Max.Y - float64(pos.Pos.Quantity)
						}
						if !child.GetRelHeight().MatchBounds {
							if child.GetActualHeight() != nil {
								childBounds.Min.Y = childBounds.Max.Y +
									*child.GetActualHeight()
							} else {
								childBounds = nil
								return nil
							}
						}
						ySet = true
					case bottomOf:
						// If the position is a percentage
						if pos.Pos.Unit == util.Percent {
							childBounds.Min.Y = bounds.Min.Y + bounds.Size().Y*
								(float64(pos.Pos.Quantity)/100)
							// If the position is just in pixels
						} else {
							childBounds.Min.Y = bounds.Min.Y + float64(pos.Pos.Quantity)
						}
						if !child.GetRelHeight().MatchBounds {
							if child.GetActualHeight() != nil {
								childBounds.Max.Y = childBounds.Min.Y +
									*child.GetActualHeight()
							} else {
								childBounds = nil
								return nil
							}
						}
						ySet = true
					case leftOf:
						// If the position is a percentage
						if pos.Pos.Unit == util.Percent {
							childBounds.Max.X = bounds.Min.X + bounds.Size().X*
								(float64(pos.Pos.Quantity)/100)
							// If the position is just in pixels
						} else {
							childBounds.Max.X = bounds.Min.X + float64(pos.Pos.Quantity)
						}
						if !child.GetRelWidth().MatchBounds {
							if child.GetActualWidth() != nil {
								childBounds.Min.X = childBounds.Max.X -
									*child.GetActualWidth()
							} else {
								childBounds = nil
								return nil
							}
						}
						xSet = true
					case rightOf:
						// If the position is a percentage
						if pos.Pos.Unit == util.Percent {
							childBounds.Min.X = bounds.Min.X + bounds.Size().X*
								(float64(pos.Pos.Quantity)/100)
							// If the position is just in pixels
						} else {
							childBounds.Min.X = bounds.Min.X + float64(pos.Pos.Quantity)
						}
						if !child.GetRelWidth().MatchBounds {
							if child.GetActualWidth() != nil {
								childBounds.Max.X = childBounds.Min.X +
									*child.GetActualWidth()
							} else {
								childBounds = nil
								return nil
							}
						}
						xSet = true
					}

					// If it's relative to an element
				} else {
					// Get the element
					relativeElem := e.GetChildByID(pos.ElementID)
					// This shouldn't happen, but check it anyways
					if relativeElem == nil {
						return element.NewNoElemError(child.Element, child.TopOf.ElementID, string(loc))
					}
					if relativeElem.GetMin() != nil &&
						relativeElem.GetMax() != nil {
						switch loc {
						case topOf:
							// Set the bounds to be the same as the element's
							childBounds.Max.Y = relativeElem.GetMax().Y
							if !child.GetRelHeight().MatchBounds {
								if child.GetActualHeight() != nil {
									childBounds.Min.Y = childBounds.Max.Y +
										*child.GetActualHeight()
								} else {
									childBounds = nil
									return nil
								}
							}
							ySet = true
							if !xSet {
								// Set the X bounds as well
								childBounds.Min.X = relativeElem.GetMin().X
								childBounds.Max.X = relativeElem.GetMax().X
							}
						case bottomOf:
							// Set the bounds to be the same as the element's
							childBounds.Min.Y = relativeElem.GetMin().Y
							if !child.GetRelHeight().MatchBounds {
								if child.GetActualHeight() != nil {
									childBounds.Max.Y = childBounds.Min.Y -
										*child.GetActualHeight()
								} else {
									childBounds = nil
									return nil
								}
							}
							ySet = true
							if !xSet {
								// Set the X bounds as well
								childBounds.Min.X = relativeElem.GetMin().X
								childBounds.Max.X = relativeElem.GetMax().X
							}
						case leftOf:
							// Set the bounds to be the same as the element's
							childBounds.Max.X = relativeElem.GetMin().X
							if !child.GetRelWidth().MatchBounds {
								if child.GetActualWidth() != nil {
									childBounds.Min.X = childBounds.Max.X -
										*child.GetActualWidth()
								} else {
									childBounds = nil
									return nil
								}
							}
							xSet = true
							if !ySet {
								// Set the Y bounds
								childBounds.Min.Y = relativeElem.GetMin().Y
								childBounds.Max.Y = relativeElem.GetMax().Y
							}
						case rightOf:
							// Set the bounds to be the same as the element's
							childBounds.Min.X = relativeElem.GetMax().X
							if !child.GetRelWidth().MatchBounds {
								if child.GetActualWidth() != nil {
									childBounds.Max.X = childBounds.Min.X +
										*child.GetActualWidth()
								}
							}
							xSet = true
							if !ySet {
								// Set the Y bounds
								childBounds.Min.Y = relativeElem.GetMin().Y
								childBounds.Max.Y = relativeElem.GetMax().Y
							}
						}
					} else {
						childBounds = nil
						return nil
					}
				}
			}
			return nil
		}

		if childBounds != nil {
			err = modBounds(&child, topOf, child.TopOf)
			if err != nil {
				return err
			}
		}
		if childBounds != nil {
			err = modBounds(&child, bottomOf, child.BottomOf)
			if err != nil {
				return err
			}
		}
		if childBounds != nil {
			err = modBounds(&child, leftOf, child.LeftOf)
			if err != nil {
				return err
			}
		}
		if childBounds != nil {
			err = modBounds(&child, rightOf, child.RightOf)
			if err != nil {
				return err
			}
		}

		// Initialise the child
		err := child.Init(window, childBounds)
		if err != nil {
			return err
		}
	}

	return nil
}

// Function that is called when there
// is a new event. This function only
// calls NewEvent on the child elements
func (e *Layout) NewEvent(window *pixelgl.Window) {
	e.Impl.NewEvent(window)
	for _, child := range e.children {
		child.NewEvent(window)
	}
}

// Function to draw the element
func (e *Layout) Draw() {
	// Draw the element
	e.Impl.Draw()
	// Draw the layout
	element.DrawLayout(e)
}
