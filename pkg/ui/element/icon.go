package element

import (
	"github.com/faiface/pixel"
	"github.com/orfby/ui/pkg/ui/util"
)

// Function to initialise an image
// as if it were an SVG icon. Should
// be called if ImageIsSVG is true
func initIcon(e Element, i Image) error {
	// If the SVG hasn't been loaded yet
	if i.GetSVG() == nil {
		// Load the svg
		svg, err := util.LoadSVG(
			e.GetFS(), i.GetField(), i.GetColor())
		if err != nil {
			return err
		}
		// Set the SVG
		i.SetSVG(svg)
	}

	// Get the scale
	scale := i.GetScale()
	if scale == util.ZeroScaleOption {
		scale = util.DefaultScaleOption
	}

	// If the svg has been loaded
	if i.GetSVG() != nil {
		// Convert the SVG's view box to a pixel rect
		viewbox := pixel.R(i.GetSVG().ViewBox.X, i.GetSVG().ViewBox.Y,
			i.GetSVG().ViewBox.W, i.GetSVG().ViewBox.H)
		// If the view box's size is 0, 0
		if viewbox.Size() == pixel.ZV {
			// Change it to a 16x16 box
			viewbox = pixel.R(
				0, 0, 16, 16)
		}

		scale := i.GetScale()
		if scale == util.ZeroScaleOption {
			scale = util.DefaultScaleOption
		}

		// If the element's width isn't known
		// and the width should be the content
		if e.GetActualWidth() == nil && e.GetRelWidth().MatchContent {
			switch scale {
			case util.ScaleToFill:
				fallthrough
			case util.ScaleToFit:
				fallthrough
			case util.Stretch:
				// If the height is knowable
				if !e.GetRelHeight().MatchContent {
					// If the height is known
					if e.GetActualHeight() != nil {
						// Calculate the scale factor of the height
						scale := *e.GetActualHeight() / viewbox.Size().Y
						// Set the width as the image's
						// width multiplied by the scale factor
						newWidth := viewbox.Size().X * scale
						e.SetActualWidth(&newWidth)
					}
				} else {
					// If it isn't knowable, just set the
					// width as the width of the image
					newWidth := viewbox.Size().X
					e.SetActualWidth(&newWidth)
				}
			default:
				// Set the actual width as the size of the image
				newWidth := viewbox.Size().X
				e.SetActualWidth(&newWidth)
			}
		}

		// If the element's height isn't known
		// and the height should be the content
		if e.GetActualHeight() == nil && e.GetRelHeight().MatchContent {
			switch scale {
			case util.ScaleToFill:
				fallthrough
			case util.ScaleToFit:
				fallthrough
			case util.Stretch:
				// If the width is knowable
				if !e.GetRelWidth().MatchContent {
					// If the width is known
					if e.GetActualWidth() != nil {
						// Calculate the scale factor of the width
						scale := *e.GetActualWidth() / viewbox.Size().X
						// Set the height as the image's
						// height multiplied by the scale factor
						newHeight := viewbox.Size().Y * scale
						e.SetActualHeight(&newHeight)
					}
				} else {
					// If it isn't knowable, just set the
					// height as the height of the image
					newHeight := viewbox.Size().Y
					e.SetActualWidth(&newHeight)
				}
			default:
				// Set the actual height as the size of the image
				newHeight := viewbox.Size().Y
				e.SetActualWidth(&newHeight)
			}
		}

		// If the image hasn't been created
		// and the width and height are known
		if i.GetSprite() == nil &&
			e.GetActualWidth() != nil &&
			e.GetActualHeight() != nil {
			// Create a picture from the SVG
			pic := util.CreatePictureFromSVG(i.GetSVG(), scale,
				*e.GetActualWidth(), *e.GetActualHeight())
			// Create a sprite with the picture and set it
			i.SetSprite(pixel.NewSprite(pic, pic.Bounds()))
		}
	}

	return nil
}
