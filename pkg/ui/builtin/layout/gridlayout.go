package layout

import (
	"encoding/xml"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/orfby/ui/pkg/ui/element"
	"github.com/orfby/ui/pkg/ui/util"
	"math"
	"net/http"
)

//Layout type for displaying elements in a
//grid (either vertically or horizontally)
type GridLayout struct {
	//A grid layout is an element
	element.Impl
	//It is also a layout
	element.LayoutImpl

	//The element's orientation
	Orientation util.Orientation `uixml:"http://github.com/orfby/ui/api/schema orientation,optional"`

	//The number of columns on each row
	Columns uint `uixml:"http://github.com/orfby/ui/api/schema columns,optional"`
	//The minimum width of a cell
	CellWidth util.RelativeSize `uixml:"http://github.com/orfby/ui/api/schema cell-width,optional"`
	//The minimum height of a cell
	CellHeight util.RelativeSize `uixml:"http://github.com/orfby/ui/api/schema cell-height,optional"`

	//The children in a grid format
	grid [][]element.Element
}

//Function to create a new grid layout
func NewGridLayout(fs http.FileSystem, name xml.Name, parent element.Layout) element.Element {
	return &GridLayout{
		Impl:        element.NewElement(fs, name, parent),
		Orientation: util.DefaultOrientation,
		Columns:     0,
		CellWidth:   util.ZeroRelativeSize,
		CellHeight:  util.ZeroRelativeSize,
	}
}

//The XML name of the element
var GridLayoutTypeName = xml.Name{Space: "http://github.com/orfby/ui/api/schema", Local: "GridLayout"}

//Function to unmarshal an XML element into
//an element. This function is usually only
//called by xml.Unmarshal
func (e *GridLayout) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
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

	//Unmarshal the layout's children
	e.LayoutImpl.Children, err = element.ChildrenUnmarshalXML(e.GetFS(), e, d, start)
	if err != nil {
		return err
	}

	//If there are children
	if len(e.LayoutImpl.Children) > 0 {
		if e.Columns == 0 {
			if e.Orientation == util.HorizontalOrientation {
				e.Columns = uint(len(e.Children))
			} else {
				e.Columns = 1
			}
		}

		//Create the grid
		e.grid = make([][]element.Element, 1)
		e.grid[0] = make([]element.Element, 0)

		//Iterate over the children
		row := 0
		column := 0
		addedChildren := 0
		for _, child := range e.Children {
			//Add the child to the row
			e.grid[row] = append(e.grid[row], child)
			addedChildren++

			//Go to the next column
			column++
			//If you should go to a new row
			if uint(column) == e.Columns {
				//Reset the column number
				column = 0
				//Go to the next column
				row++
				//If there will actually be a new column
				if addedChildren < e.NumChildren() {
					//Append a new array
					e.grid = append(e.grid, make([]element.Element, 0))
				}
			}
		}
	}

	//If the cell width wasn't given
	if e.CellWidth == util.ZeroRelativeSize {
		//Set the width as a percentage (so
		//it takes up the whole width of the parent)
		e.CellWidth.Quantity = int32(100 / len(e.grid[0]))
		e.CellWidth.Unit = util.Percent
	}

	//If the row height wasn't given
	if e.CellHeight == util.ZeroRelativeSize {
		//Set the height as a percentage (so
		//it takes up the whole height of the parent)
		e.CellHeight.Quantity = int32(100 / len(e.grid))
		e.CellHeight.Unit = util.Percent
	}

	return nil
}

//Function to reset the element
func (e *GridLayout) Reset() {
	e.Impl.Reset()
	e.LayoutImpl.Reset()
}

//Function to determine whether
//the element is initialised
func (e *GridLayout) IsInitialised() bool {
	return e.Impl.IsInitialised() &&
		element.ChildrenAreInitialised(e)
}

//Function to initialise the element
func (e *GridLayout) Init(window *pixelgl.Window, bounds *pixel.Rect) error {
	//Initialise the element part of the layout
	err := e.Impl.Init(window, bounds)
	if err != nil {
		return err
	}

	//The actual width of a cell
	actualCellWidth := new(float64)

	//If the cell width is to match the content
	if e.CellWidth.MatchContent {
		//Just set the minimum to 0
		*actualCellWidth = 0
	} else {
		//Otherwise calculate the minimum from the relative width
		//(with the layout itself as the parent)
		//todo cell width can't be match_bounds
		actualCellWidth = element.CalculateWidth(e, window, nil, e.CellWidth)
	}

	//If the actual cell width is still known
	if actualCellWidth != nil {
		//Iterate over the children
		for _, child := range e.Children {
			//If the child's width isn't known
			if child.GetActualWidth() == nil {
				//Reset the width
				actualCellWidth = nil
				break
			}
			//Get the max width
			*actualCellWidth = math.Max(*actualCellWidth, *child.GetActualWidth())
		}

		//If the width is meant to match the content size
		//and the cell width was calculated
		if e.GetRelWidth().MatchContent && actualCellWidth != nil {
			//Set the actual width as the column
			//width multiplied by the number of columns
			actualWidth := *actualCellWidth * float64(len(e.grid[0]))
			e.SetActualWidth(&actualWidth)
		}
	}

	//The actual height of a cell
	actualCellHeight := new(float64)

	//If the cell height is to match the content
	if e.CellWidth.MatchContent {
		//Just set the minimum to 0
		*actualCellHeight = 0
	} else {
		//Otherwise calculate the minimum from the relative height
		//(with the layout itself as the parent)
		//todo cell width can't be match_bounds
		actualCellHeight = element.CalculateHeight(e, window, nil, e.CellHeight)
	}

	//If the actual cell height is still known
	if actualCellHeight != nil {
		//Iterate over the children
		for _, child := range e.Children {
			//If the child's height isn't known
			if child.GetActualHeight() == nil {
				//Reset the height
				actualCellHeight = nil
				break
			}
			//Get the max height
			*actualCellHeight = math.Max(*actualCellHeight, *child.GetActualHeight())
		}

		//If the height is meant to match the content size
		//and the cell height was calculated
		if e.GetRelHeight().MatchContent && actualCellHeight != nil {
			//Set the actual width as the row
			//height multiplied by the number of rows
			actualHeight := *actualCellHeight * float64(len(e.grid))
			e.SetActualHeight(&actualHeight)
		}
	}

	//Iterate over the rows of the grid
	for y, ys := range e.grid {
		//Iterate over the children of the row
		for x, child := range ys {
			//If the child hasn't been initialised yet
			if !child.IsInitialised() {
				childBounds := (*pixel.Rect)(nil)
				//If the layout's minimum,
				//cell width and height are known
				if e.GetMin() != nil &&
					actualCellWidth != nil &&
					actualCellHeight != nil {
					//Set the child bounds to a valid place
					childBounds = &pixel.Rect{}

					//Get the padding
					var padding float64
					if e.GetPadding().Unit == util.Pixels {
						padding = float64(e.GetPadding().Quantity)
					}

					//Set the child's position
					childBounds.Min.X = e.GetMin().X + padding + (float64(x) * *actualCellWidth)
					childBounds.Min.Y = e.GetMax().Y - padding - (float64(y+1) * *actualCellHeight)
					//Set the child's max position
					//(as the position + the cell size)
					childBounds.Max = childBounds.Min.Add(pixel.V(*actualCellWidth, *actualCellHeight))

					//Otherwise set child bounds to nil
				} else {
					childBounds = nil
				}

				//Initialise the child
				err := child.Init(window, childBounds)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

//Function that is called when there
//is a new event
func (e *GridLayout) NewEvent(window *pixelgl.Window) {
	e.Impl.NewEvent(window)
	e.LayoutImpl.NewEvent(window)
}

//Function to draw the element
func (e *GridLayout) Draw() {
	//Draw the element
	e.Impl.Draw()
	//Draw the layout
	element.DrawLayout(e)
}
