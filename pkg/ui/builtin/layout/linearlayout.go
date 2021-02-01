package layout

import (
	"encoding/xml"
	"github.com/bhollier/ui/pkg/ui/element"
	"github.com/bhollier/ui/pkg/ui/util"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"net/http"
)

// Layout type for displaying elements as a
// list (either vertically or horizontally)
type LinearLayout struct {
	// A linear layout is an element
	element.Impl
	// It is also a layout
	element.LayoutImpl

	// The element's orientation
	Orientation util.Orientation `uixml:"http://github.com/bhollier/ui/api/schema orientation,optional"`
}

// Function to create a new linear layout
func NewLinearLayout(fs http.FileSystem, name xml.Name, parent element.Layout) element.Element {
	return &LinearLayout{
		Impl:        element.NewElement(fs, name, parent),
		Orientation: util.DefaultOrientation,
	}
}

// The XML name of the element
var LinearLayoutTypeName = xml.Name{Space: "http://github.com/bhollier/ui/api/schema", Local: "LinearLayout"}

// Function to unmarshal an XML element into
// an element. This function is usually only
// called by xml.Unmarshal
func (e *LinearLayout) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
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

	// Unmarshal the layout's children
	e.LayoutImpl.Children, err = element.ChildrenUnmarshalXML(e.GetFS(), e, d, start)
	if err != nil {
		return err
	}

	return nil
}

// Function to reset the element's
// position
func (e *LinearLayout) ResetPosition() {
	e.Impl.ResetPosition()
	e.LayoutImpl.ResetPosition()
}

// Function to reset the element
func (e *LinearLayout) Reset() {
	e.Impl.Reset()
	e.LayoutImpl.Reset()
}

// Function to determine whether
// the element is initialised
func (e *LinearLayout) IsInitialised() bool {
	return e.Impl.IsInitialised() &&
		element.ChildrenAreInitialised(e)
}

// Function to initialise the element
func (e *LinearLayout) Init(window *pixelgl.Window, bounds *pixel.Rect) error {
	// If the layout's width isn't known and
	// the width is meant to match the content size
	if e.GetActualWidth() == nil && e.GetRelWidth().MatchContent {
		// If the orientation is horizontal
		if e.Orientation == util.HorizontalOrientation {
			var width float64
			if e.GetPadding().Unit == util.Pixels {
				width = float64(e.GetPadding().Quantity * 2)
			}
			allChildrenInit := true
			for i := 0; i < e.NumChildren(); i++ {
				if e.GetChild(i).GetActualWidth() == nil {
					allChildrenInit = false
					break
				}

				// Add the child's width to the width
				width += *e.GetChild(i).GetActualWidth()
			}
			// If all the children were considered,
			// set the actual width
			if allChildrenInit {
				e.SetActualWidth(&width)
			}

		} else {
			var maxWidth float64
			if e.GetPadding().Unit == util.Pixels {
				maxWidth = float64(e.GetPadding().Quantity * 2)
			}
			allChildrenInit := true
			for i := 0; i < e.NumChildren(); i++ {
				if e.GetChild(i).GetActualWidth() == nil {
					allChildrenInit = false
					break
				}

				// Get the max width
				maxWidth = math.Max(maxWidth, *e.GetChild(i).GetActualWidth())
			}
			// If all the children were considered,
			// set the actual width
			if allChildrenInit {
				e.SetActualWidth(&maxWidth)
			}
		}
	}

	// If the layout's height isn't known and
	// the height is meant to match the content size
	if e.GetActualHeight() == nil && e.GetRelHeight().MatchContent {
		// If the orientation is horizontal
		if e.Orientation == util.VerticalOrientation {
			var height float64
			if e.GetPadding().Unit == util.Pixels {
				height = float64(e.GetPadding().Quantity * 2)
			}
			allChildrenInit := true
			for i := 0; i < e.NumChildren(); i++ {
				if e.GetChild(i).GetActualHeight() == nil {
					allChildrenInit = false
					break
				}

				// Add the child's height to the height
				height += *e.GetChild(i).GetActualHeight()
			}
			// If all the children were considered,
			// set the actual height
			if allChildrenInit {
				e.SetActualHeight(&height)
			}

		} else {
			var maxHeight float64
			if e.GetPadding().Unit == util.Pixels {
				maxHeight = float64(e.GetPadding().Quantity * 2)
			}
			allChildrenInit := true
			for i := 0; i < e.NumChildren(); i++ {
				if e.GetChild(i).GetActualHeight() == nil {
					allChildrenInit = false
					break
				}

				// Get the max height
				maxHeight = math.Max(maxHeight, *e.GetChild(i).GetActualHeight())
			}
			// If all the children were considered,
			// set the actual height
			if allChildrenInit {
				e.SetActualHeight(&maxHeight)
			}
		}
	}

	// Initialise the element part of the layout
	err := e.Impl.Init(window, bounds)
	if err != nil {
		return err
	}

	// The child's position
	var childPos *pixel.Vec

	// If the minimum is known
	if bounds != nil &&
		e.GetMin() != nil &&
		e.GetMax() != nil {
		// Get the padding
		var padding float64
		if e.GetPadding().Unit == util.Pixels {
			padding = float64(e.GetPadding().Quantity)
		}
		// Set the position
		childPos = &pixel.Vec{
			X: e.GetMin().X + padding,
			Y: e.GetMax().Y + padding,
		}
	} else {
		childPos = nil
	}

	// Initialise the children
	var child element.Element
	for i := 0; i < e.NumChildren(); i++ {
		child = e.GetChild(i)

		// If the child hasn't been initialised yet
		if !child.IsInitialised() {
			// The child's bounds (default is nil)
			var childBounds *pixel.Rect
			// If the child's position can be calculated
			// (if the previous position is known and the
			// child's width/height is known)
			if childPos != nil &&
				child.GetActualWidth() != nil &&
				child.GetActualHeight() != nil {
				// Calculate the child's size
				childSize := pixel.V(*child.GetActualWidth(),
					*child.GetActualHeight())
				// Minus the Y by the size of the child
				childPos.Y -= childSize.Y
				// Set the child bounds
				childBounds = &pixel.Rect{
					Min: *childPos,
					Max: childPos.Add(childSize),
				}

				// Increase childPos (for the next child)
				if e.Orientation == util.HorizontalOrientation {
					childPos.X += childSize.X
				}
			} else {
				childPos = nil
			}

			// Initialise the child
			err := child.Init(window, childBounds)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Function that is called when there
// is a new event
func (e *LinearLayout) NewEvent(window *pixelgl.Window) {
	e.Impl.NewEvent(window)
	e.LayoutImpl.NewEvent(window)
}

// Function to draw the element
func (e *LinearLayout) Draw() {
	// Draw the element
	e.Impl.Draw()
	// Draw the layout
	element.DrawLayout(e)
}
