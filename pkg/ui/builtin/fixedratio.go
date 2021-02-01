package builtin

import (
	"encoding/xml"
	"errors"
	"github.com/bhollier/ui/pkg/ui/element"
	"github.com/bhollier/ui/pkg/ui/util"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"net/http"
)

// Type for an element that enforces a fixed ratio
// on its child element
type FixedRatio struct {
	// The fixed ratio element is an element
	element.Impl
	// It is also (technically) a layout
	element.LayoutImpl

	// The ratio itself
	Ratio util.Ratio `uixml:"http://github.com/bhollier/ui/api/schema ratio"`
}

// Function to create a new fixed ratio element
func NewFixedRatio(fs http.FileSystem, name xml.Name, parent element.Layout) element.Element {
	return &FixedRatio{Impl: element.NewElement(fs, name, parent)}
}

// The XML name of the import element
var FixedRatioTypeName = xml.Name{Space: "http://github.com/bhollier/ui/api/schema", Local: "FixedRatio"}

// Function to unmarshal an XML element into
// an element. This function is usually only
// called by xml.Unmarshal
func (e *FixedRatio) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
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

	// If there are no children
	if len(e.LayoutImpl.Children) == 0 {
		return errors.New("no children on XML element '" +
			element.FullName(e, ".", true) + "'")

		// If there are multiple
	} else if len(e.LayoutImpl.Children) > 1 {
		return errors.New("multiple children on XML element '" +
			element.FullName(e, ".", true) + "'")
	}

	return nil
}

// Function to reset the element's
// position
func (e *FixedRatio) ResetPosition() {
	e.Impl.ResetPosition()
	e.LayoutImpl.ResetPosition()
}

// Function to reset the element
func (e *FixedRatio) Reset() {
	e.Impl.Reset()
	e.LayoutImpl.Reset()
}

// Function to determine whether
// the element is initialised
func (e *FixedRatio) IsInitialised() bool {
	return e.Impl.IsInitialised() &&
		element.ChildrenAreInitialised(e)
}

// Function to initialise an element's
// position, width and height. Because
// it doesn't know the element's actual
// size, it won't set the width or height
// if the relative width or height is
// "match_content"
func (e *FixedRatio) Init(window *pixelgl.Window, bounds *pixel.Rect) (err error) {
	// Initialise the element part of the fixed ratio
	err = e.Impl.Init(window, bounds)
	if err != nil {
		return err
	}

	// If the bounds are known
	if bounds != nil {
		// Calculate the correct dimensions of the child
		dimensions := e.Ratio.RestrictDimensions(bounds.Size())

		// If the bounds are integers
		boundsSize := bounds.Size()
		if boundsSize == boundsSize.Floor() {
			// While the dimensions aren't integers
			for dimensions.Floor() != dimensions {
				if dimensions.X == bounds.Size().X {
					boundsSize.X -= 1
					dimensions = e.Ratio.RestrictDimensions(boundsSize)
				} else {
					boundsSize.Y -= 1
					dimensions = e.Ratio.RestrictDimensions(boundsSize)
				}
			}
		}

		// Calculate the minimum point for the child (based on gravity)
		min := element.CalculateMin(e, bounds, dimensions)

		// If the minimum point was calculated (it should've been)
		if min != nil {
			// Create the bounds of the child
			childBounds := pixel.Rect{
				Min: *min,
				Max: min.Add(dimensions),
			}

			// Initialise the child with the bounds
			err = e.GetChild(0).Init(window, &childBounds)
			if err != nil {
				return err
			}
		} else {
			// Initialise the child with the bounds
			err = e.GetChild(0).Init(window, nil)
			if err != nil {
				return err
			}
		}

		// Otherwise
	} else {
		// Initialise the child with no bounds
		err = e.GetChild(0).Init(window, nil)
		if err != nil {
			return err
		}
	}

	// If the width is meant to match the content size
	if e.GetRelWidth().MatchContent {
		// Set the width as the child's
		e.SetActualWidth(e.GetChild(0).GetActualWidth())
	}

	// If the height is meant to match the content size
	if e.GetRelHeight().MatchContent {
		// Set the height as the child's
		e.SetActualHeight(e.GetChild(0).GetActualHeight())
	}

	return nil
}

// Function that is called when there
// is a new event
func (e *FixedRatio) NewEvent(window *pixelgl.Window) {
	e.Impl.NewEvent(window)
	e.LayoutImpl.NewEvent(window)
}

// Function to draw the element
func (e *FixedRatio) Draw() {
	// Draw the element
	e.Impl.Draw()
	// Draw the layout
	element.DrawLayout(e)
}
