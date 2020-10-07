package element

import (
	"encoding/xml"
	"errors"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/orfby/ui/pkg/ui/util"
	"github.com/xlab/treeprint"
	"net/http"
)

// Interface type for something
// that is an element
type Element interface {
	// Function to get the element's parent
	GetParent() Layout

	// Function to get the filesystem to use
	// for opening files
	GetFS() http.FileSystem

	// Function to get the element's XML name
	GetName() xml.Name
	// Function to get the element's XML
	// namespaces. This should include
	// GetName().Space
	GetNamespaces() []string
	// Function to add a namespace to the
	// element
	AddNamespace(string)

	// Function to get the element's
	// ID (or nil, if it doesn't have
	// one)
	GetID() *string

	// Function to get the element's
	// relative width
	GetRelWidth() util.RelativeSize
	// Function to get the element's
	// relative height
	GetRelHeight() util.RelativeSize

	// Function to get the element's
	// actual width. If nil, it isn't
	// known yet
	GetActualWidth() *float64
	// Function to set the element's
	// actual width
	SetActualWidth(*float64)

	// Function to get the element's
	// actual height. If nil, it isn't
	// known yet
	GetActualHeight() *float64
	// Function to set the element's
	// actual height
	SetActualHeight(*float64)

	// Function to get the element's
	// minimum point. If nil, it isn't
	// known yet
	GetMin() *pixel.Vec
	// Function to set the element's
	// minimum point
	SetMin(*pixel.Vec)

	// Function to get the element's
	// maximum point. If nil, it isn't
	// known yet
	GetMax() *pixel.Vec
	// Function to set the element's
	// maximum point
	SetMax(*pixel.Vec)

	// Function to get the element's
	// bounds. If nil, either the min
	// or max isn't known
	GetBounds() *pixel.Rect

	// Function to get the element's
	// padding
	GetPadding() util.AbsoluteQuantity

	// Function to get the element's
	// gravity
	GetGravity() util.Gravity

	// todo tint (or maybe foreground?)

	// Function to get the element's
	// background
	GetBkg() Image

	// Function to get the element's
	// canvas
	GetCanvas() *pixelgl.Canvas

	// Function to unmarshal an XML element
	// into a UI element. This function is
	// usually only called by xml.Unmarshal
	UnmarshalXML(d *xml.Decoder, start xml.StartElement) error

	// Function to reset the element's
	// position
	ResetPosition()

	// Function to reset the element
	// completely
	Reset()

	// Function to determine whether
	// the element is initialised
	IsInitialised() bool

	// Function to initialise the element
	// (load textures, create sprites, set
	// sprite locations, etc.)
	Init(window *pixelgl.Window, bounds *pixel.Rect) error

	// Function that is called when there
	// is a new event
	NewEvent(*pixelgl.Window)

	// Function to draw the element
	// to its canvas
	Draw()
}

// Type for the implementation of a UI
// element. Most elements should include
// this at some point (but only once)
type Impl struct {
	// The element's parent
	parent Layout

	// The filesystem to use
	fs http.FileSystem

	// The element's XML name
	name xml.Name
	// The element's XML namespaces
	namespaces []string

	// The element's ID
	ID string `uixml:"http://github.com/orfby/ui/api/schema id,optional"`

	// The element's relative width
	RelativeWidth util.RelativeSize `uixml:"http://github.com/orfby/ui/api/schema width"`
	// The element's relative height
	RelativeHeight util.RelativeSize `uixml:"http://github.com/orfby/ui/api/schema height"`

	// The element's width
	// (or nil, if unknown)
	width *float64
	// The element's height
	// (or nil, if unknown)
	height *float64

	// The element's minimum point
	// (or nil, if unknown)
	min *pixel.Vec
	// The element's maximum point
	// (or nil, if unknown)
	max *pixel.Vec

	// The element's padding
	Padding util.AbsoluteQuantity `uixml:"http://github.com/orfby/ui/api/schema padding,optional"`

	// The element's background
	Bkg Background

	// The element's canvas
	Canvas *pixelgl.Canvas

	// The element's gravity
	Gravity util.Gravity `uixml:"http://github.com/orfby/ui/api/schema gravity,optional"`
}

// Function to create an element
func NewElement(fs http.FileSystem, name xml.Name, parent Layout) Impl {
	e := Impl{
		name:    name,
		parent:  parent,
		fs:      fs,
		Gravity: util.DefaultGravity,
		Padding: util.DefaultAbsoluteQuantity,
	}

	// If the parent was actually given,
	// set the namespace as the parents.
	// Otherwise make an empty array
	if parent != nil {
		e.namespaces = parent.GetNamespaces()
	} else {
		e.namespaces = make([]string, 0)
	}

	// If the name has a namespace, then add it
	if name.Space != "" {
		e.AddNamespace(name.Space)
	}

	return e
}

// Function to get the element's parent
func (e *Impl) GetParent() Layout { return e.parent }

// Function to get the filesystem to use
// for opening files
func (e *Impl) GetFS() http.FileSystem { return e.fs }

// Function to get the element's name
func (e *Impl) GetName() xml.Name { return e.name }

// Function to get the element's namespaces
func (e *Impl) GetNamespaces() []string { return e.namespaces }

// Function to add a namespace to the
// element
func (e *Impl) AddNamespace(namespace string) {
	for _, ns := range e.namespaces {
		if ns == namespace {
			return
		}
	}
	e.namespaces = append(e.namespaces, namespace)
}

// Function to get the element's
// ID (or nil, if it doesn't have
// one)
func (e *Impl) GetID() *string {
	if e.ID == "" {
		return nil
	} else {
		return &e.ID
	}
}

// Function to get the element's
// relative width
func (e *Impl) GetRelWidth() util.RelativeSize { return e.RelativeWidth }

// Function to get the element's
// relative height
func (e *Impl) GetRelHeight() util.RelativeSize { return e.RelativeHeight }

// Function to get the element's
// actual width. If nil, it isn't
// known yet
func (e *Impl) GetActualWidth() *float64 { return e.width }

// Function to set the element's
// actual width
func (e *Impl) SetActualWidth(width *float64) { e.width = width }

// Function to get the element's
// actual height. If nil, it isn't
// known yet
func (e *Impl) GetActualHeight() *float64 { return e.height }

// Function to set the element's
// actual height
func (e *Impl) SetActualHeight(height *float64) { e.height = height }

// Function to get the element's
// minimum point. If nil, it isn't
// known yet
func (e *Impl) GetMin() *pixel.Vec { return e.min }

// Function to set the element's
// minimum point
func (e *Impl) SetMin(min *pixel.Vec) { e.min = min }

// Function to get the element's
// maximum point. If nil, it isn't
// known yet
func (e *Impl) GetMax() *pixel.Vec { return e.max }

// Function to set the element's
// maximum point
func (e *Impl) SetMax(max *pixel.Vec) { e.max = max }

// Function to get the element's
// bounds. If nil, either the min
// or max isn't known
func (e *Impl) GetBounds() *pixel.Rect {
	if e.min != nil && e.max != nil {
		return &pixel.Rect{
			Min: *e.min,
			Max: *e.max,
		}
	} else {
		return nil
	}
}

// Function to get the element's
// padding
func (e *Impl) GetPadding() util.AbsoluteQuantity { return e.Padding }

// Function to get the element's
// gravity
func (e *Impl) GetGravity() util.Gravity { return e.Gravity }

// Function to get the element's
// background
func (e *Impl) GetBkg() Image { return &e.Bkg }

// Function to get the element's
// canvas
func (e *Impl) GetCanvas() *pixelgl.Canvas { return e.Canvas }

// Function to unmarshal an XML element into
// an element. SetAttrs should've been called
// before this function
func (e *Impl) UnmarshalXML(d *xml.Decoder, s xml.StartElement) (err error) {
	// Unmarshal the background
	return e.Bkg.UnmarshalXML(d, s)
}

// Function to reset the element's
// position
func (e *Impl) ResetPosition() {
	// Set the min and max to nil
	e.min = nil
	e.max = nil
}

// Function to reset the element.
// This function doesn't reset
// the width and height of the
// element
func (e *Impl) Reset() {
	// Reset the position
	e.ResetPosition()
	// Reset the background
	e.Bkg.Reset()
	// Reset the width and height
	e.width = nil
	e.height = nil
}

// Function to determine whether the
// element has been initialised by
// whether its width, height, position
// and canvas are set (not nil) and if
// it's background is initialised
func (e *Impl) IsInitialised() bool {
	return e.GetMin() != nil &&
		e.GetMax() != nil &&
		e.GetCanvas() != nil &&
		e.GetCanvas().Bounds() == *e.GetBounds() &&
		e.Bkg.IsInitialised()
}

// Function to calculate an element's
// width
func CalculateWidth(parent Element, window *pixelgl.Window,
	bounds *pixel.Rect, relWidth util.RelativeSize) (width *float64) {
	// If the width is just in pixels
	if relWidth.Unit == util.Pixels {
		// Set the actual width as the number of pixels
		newWidth := float64(relWidth.Quantity)
		width = &newWidth

		// If the width depends on the bounds
	} else if relWidth.MatchBounds {
		if bounds != nil {
			newWidth := bounds.Size().X
			return &newWidth
		} else {
			// todo check that the parent doesn't match_content
			return
		}

		// If the width depends on the parent
	} else if relWidth.MatchParent ||
		relWidth.Unit == util.Percent {
		// Go up the hierarchy until a parent
		// is found that doesn't depend on it's
		// child's width
		next := parent
		for next != nil && next.GetRelWidth().MatchContent {
			next = next.GetParent()
		}

		// If the parent was found
		if next != nil {
			if relWidth.MatchParent {
				// Set the size as the parent's
				width = next.GetActualWidth()

			} else if next.GetActualWidth() != nil &&
				relWidth.Unit == util.Percent {
				newWidth := *next.GetActualWidth() * (float64(relWidth.Quantity) / 100)
				width = &newWidth
			}
		} else {
			newWidth := window.Bounds().Max.X
			if relWidth.Unit == util.Percent {
				newWidth *= float64(relWidth.Quantity) / 100
			}
			width = &newWidth
		}
	}
	return
}

// Function to calculate an element's
// height
func CalculateHeight(parent Element, window *pixelgl.Window,
	bounds *pixel.Rect, relHeight util.RelativeSize) (height *float64) {
	// If the height is just in pixels
	if relHeight.Unit == util.Pixels {
		// Set the actual height as the number of pixels
		newHeight := float64(relHeight.Quantity)
		height = &newHeight

		// If the height depends on the bounds
	} else if relHeight.MatchBounds {
		if bounds != nil {
			newHeight := bounds.Size().Y
			return &newHeight
		} else {
			return
		}

		// If the height depends on the parent
	} else if relHeight.MatchParent ||
		relHeight.Unit == util.Percent {
		// Go up the hierarchy until a parent
		// is found that doesn't depend on it's
		// child's height
		next := parent
		for next != nil && next.GetRelHeight().MatchContent {
			next = next.GetParent()
		}

		// If the parent was found
		if next != nil {
			if relHeight.MatchParent {
				// Set the size as the parent's
				height = next.GetActualHeight()

			} else if next.GetActualHeight() != nil &&
				relHeight.Unit == util.Percent {
				newHeight := *next.GetActualHeight() * (float64(relHeight.Quantity) / 100)
				height = &newHeight
			}
		} else {
			newHeight := window.Bounds().Max.Y
			if relHeight.Unit == util.Percent {
				newHeight *= float64(relHeight.Quantity) / 100
			}
			height = &newHeight
		}
	}
	return
}

// Function to calculate an element's
// minimum point
func CalculateMin(e Element, bounds *pixel.Rect, size pixel.Vec) (min *pixel.Vec) {
	// If the bounds and the element's
	// width and height are known
	if bounds != nil {
		min = new(pixel.Vec)
		// Set the X position
		if e.GetGravity().HorizGravity == util.GravNeg {
			min.X = bounds.Min.X
		} else if e.GetGravity().HorizGravity == util.GravPos {
			min.X = bounds.Max.X - size.X
		} else {
			min.X = bounds.Center().X - (size.X / 2)
		}

		// Set the Y position
		if e.GetGravity().VertGravity == util.GravNeg {
			min.Y = bounds.Max.Y - size.Y
		} else if e.GetGravity().VertGravity == util.GravPos {
			min.Y = bounds.Min.Y
		} else {
			min.Y = bounds.Center().Y - (size.Y / 2)
		}
	}
	return
}

// Function to initialise an element's
// position, width and height. Because
// it doesn't know the element's actual
// size, it won't set the width or height
// if the relative width or height is
// "match_content"
func (e *Impl) Init(window *pixelgl.Window, bounds *pixel.Rect) error {
	// If the width isn't known, try to calculate it
	if e.width == nil {
		e.width = CalculateWidth(e.GetParent(),
			window, bounds, e.GetRelWidth())
	}
	// If the height isn't known, try to calculate it
	if e.height == nil {
		e.height = CalculateHeight(e.GetParent(),
			window, bounds, e.GetRelHeight())
	}

	// If the bounds aren't known and
	// the size is known
	if e.GetBounds() == nil &&
		e.width != nil &&
		e.height != nil {
		// Create a vector for the size
		size := pixel.Vec{X: *e.width, Y: *e.height}
		// Calculate the min point
		min := CalculateMin(e, bounds, size)
		// If it was successful
		if min != nil {
			// Set the minimum point
			e.min = min
			// Set the maximum point
			max := min.Add(size)
			e.max = &max
		}
	}

	// If the bounds are known
	if e.GetBounds() != nil {
		// If the canvas hasn't been made
		if e.Canvas == nil {
			// Create a canvas
			e.Canvas = pixelgl.NewCanvas(*e.GetBounds())
		} else {
			e.Canvas.SetBounds(*e.GetBounds())
		}
	}

	// Initialise the background
	err := InitBkg(e, &e.Bkg)
	if err != nil {
		return err
	}

	return nil
}

// Function that is called when there
// is a new event. This function does
// nothing
func (e *Impl) NewEvent(*pixelgl.Window) {}

// Function to draw the element.
// This function should be called
// first
func (e *Impl) Draw() {
	// Draw the background
	DrawBkg(e, &e.Bkg)
}

// Function to draw a canvas (usually an
// element's) onto a parent canvas. This
// function is generally used by layouts,
// or just anything that has child
// element(s)
func DrawCanvasOntoParent(child *pixelgl.Canvas, parent *pixelgl.Canvas) {
	mat := pixel.IM
	// Move it to where the canvas wants to be
	mat = mat.Moved(child.Bounds().Center())
	// Draw the child canvas onto the parent
	child.Draw(parent, mat)
}

// Function to reset and initialise the
// entire UI element tree, by traversing
// up the given element's parents
func InitUI(e Element, window *pixelgl.Window, bounds *pixel.Rect) error {
	// While the element has a parent, go up the tree
	for e.GetParent() != nil {
		e = e.GetParent()
	}

	// Reset the element (and its children)
	e.Reset()

	// Keep initialising the element until it's all done
	// or until the element has been initialised for the 1000th time
	for i := 0; i < 1000 && !e.IsInitialised(); i++ {
		// Initialise the element (and therefore all its children)
		err := e.Init(window, bounds)
		if err != nil {
			return err
		}
	}

	// If the element still isn't initialised (because the
	// loop limit was reached) return an error
	if !e.IsInitialised() {
		// Make a tree of uninitialised elements
		tree := treeprint.New()

		// Recursive function to find uninitialised elements
		var getUninitialisedElements func(treeprint.Tree, Element)
		getUninitialisedElements = func(branch treeprint.Tree, e Element) {
			// If the element isn't initialised
			if !e.IsInitialised() {
				elemName := Name(e, true)
				// Try to convert to a layout
				layout, ok := e.(Layout)
				// If it is a layout, iterate over the children
				if ok {
					newBranch := branch.AddBranch(elemName)
					for i := 0; i < layout.NumChildren(); i++ {
						// Get the child element's uninitialised elements
						getUninitialisedElements(newBranch, layout.GetChild(i))
					}
				} else {
					branch.AddNode(elemName)
				}
			}
		}

		// Call the function on root
		getUninitialisedElements(tree, e)

		// Return an error with the uninitialised elements
		return errors.New("infinite loop in element init detected. " +
			"The following element(s) are still uninitialised: \n" + tree.String())
	}
	return nil
}

// Function to draw the entire UI element
// tree, by traversing up the given
// element's parents
func DrawUI(e Element, window *pixelgl.Window) {
	// While the element has a parent, go up the tree
	for e.GetParent() != nil {
		e = e.GetParent()
	}
	// Draw the element
	e.Draw()
	// Draw the element onto the window
	DrawCanvasOntoParent(e.GetCanvas(), window.Canvas())
	// Swap the window's buffers
	window.SwapBuffers()
}
