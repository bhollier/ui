package element

import (
	"github.com/bhollier/ui/pkg/ui/util"
	"github.com/faiface/pixel"
	"github.com/srwiley/oksvg"
	"image/color"
	"log"
	"math"
	"path/filepath"
	"strings"
)

// Interface for an image
type Image interface {
	// Function to get the image's
	// XML field
	GetField() string
	// Function to get the image's
	// scale options
	GetScale() util.ScaleOption
	// Function to get the image's
	// colour
	GetColor() color.RGBA

	// Function to determine whether
	// the image looks like an SVG
	IsSVG() bool
	// Function to get the image's
	// SVG. Always returns nil if
	// IsSVG returns false
	GetSVG() *oksvg.SvgIcon
	// Function to set the image's
	// SVG
	SetSVG(*oksvg.SvgIcon)

	// Function to get the image's
	// sprite
	GetSprite() *pixel.Sprite
	// Function to set the image's
	// sprite
	SetSprite(*pixel.Sprite)
}

// Implementation of the
// image interface
type ImageImpl struct {
	// The image field from
	// xml
	Field string `uixml:"http://github.com/bhollier/ui/api/schema source,optional"`
	// The image's scale
	// option
	Scale util.ScaleOption `uixml:"http://github.com/bhollier/ui/api/schema scale,optional"`
	// The image's colour
	Color color.RGBA `uixml:"http://github.com/bhollier/ui/api/schema color,optional"`
	// The image's svg (if
	// applicable)
	svg *oksvg.SvgIcon
	// The image's sprite
	sprite *pixel.Sprite
}

// Function to get the image's
// XML field
func (i *ImageImpl) GetField() string { return i.Field }

// Function to get the image's
// image scale options
func (i *ImageImpl) GetScale() util.ScaleOption { return i.Scale }

// Function to get the image's
// colour
func (i *ImageImpl) GetColor() color.RGBA { return i.Color }

// Function to determine whether
// the image looks like an SVG
func (i *ImageImpl) IsSVG() bool {
	// If the image field isn't empty
	return i.GetField() != "" &&
		// And the extension is ".svg"
		strings.ToLower(filepath.Ext(i.GetField())) == ".svg"
}

// Function to get the image's
// SVG
func (i *ImageImpl) GetSVG() *oksvg.SvgIcon { return i.svg }

// Function to set the image's
// SVG
func (i *ImageImpl) SetSVG(s *oksvg.SvgIcon) { i.svg = s }

// Function to get the image's
// sprite
func (i *ImageImpl) GetSprite() *pixel.Sprite { return i.sprite }

// Function to set the image's
// sprite
func (i *ImageImpl) SetSprite(s *pixel.Sprite) { i.sprite = s }

// Function to reset the image
func (i *ImageImpl) Reset() {
	// If the image is an SVG
	if i.IsSVG() {
		// Reset the sprite
		i.SetSprite(nil)
	}
}

// Function to determine whether the
// image has been initialised, by whether
// its image sprite has been set (assuming
// it's meant to be set). This function
// doesn't call element.IsInitialised
func (i *ImageImpl) IsInitialised() bool {
	// If the element doesn't have an image
	return i.GetField() == "" ||
		// Or the image has been initialised
		i.GetSprite() != nil
}

// Function to initialise an element's
// image. Doesn't call element.Init
func InitImage(e Element, i Image) error {
	// todo set width should be a max

	// If the image looks like an SVG
	if i.IsSVG() {
		// Defer to the svg init function
		return initIcon(e, i)
	}

	// If the image hasn't been made yet
	if i.GetSprite() == nil {
		// Load the image
		picture, err := util.CreatePictureFromField(e.GetFS(), i.GetField())
		if err != nil {
			return err
		}
		if picture != nil {
			i.SetSprite(pixel.NewSprite(picture, picture.Bounds()))
		}
	}

	// If the sprite has been created
	if i.GetSprite() != nil {
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
						scale := *e.GetActualHeight() / i.GetSprite().Frame().Size().Y
						// Set the width as the image's
						// width multiplied by the scale factor
						newWidth := i.GetSprite().Frame().Size().X * scale
						e.SetActualWidth(&newWidth)
					}
				} else {
					// If it isn't knowable, just set the
					// width as the width of the image
					newWidth := i.GetSprite().Frame().Size().X
					e.SetActualWidth(&newWidth)
				}
			default:
				// Set the actual width as the size of the image
				newWidth := i.GetSprite().Frame().Size().X
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
						scale := *e.GetActualWidth() / i.GetSprite().Frame().Size().X
						// Set the height as the image's
						// height multiplied by the scale factor
						newHeight := i.GetSprite().Frame().Size().Y * scale
						e.SetActualHeight(&newHeight)
					}
				} else {
					// If it isn't knowable, just set the
					// height as the height of the image
					newHeight := i.GetSprite().Frame().Size().Y
					e.SetActualWidth(&newHeight)
				}
			default:
				// Set the actual height as the size of the image
				newHeight := i.GetSprite().Frame().Size().Y
				e.SetActualWidth(&newHeight)
			}
		}
	}

	return nil
}

// Function to draw an element's
// image
func DrawImage(e Element, i Image) {
	// If the image is an SVG
	if i.IsSVG() {
		// Don't perform any scaling
		mat := pixel.IM
		// Move it to the center of the canvas
		mat = mat.Moved(e.GetCanvas().Bounds().Center())
		// todo use gravity
		// Draw the sprite
		i.GetSprite().Draw(e.GetCanvas(), mat)

	} else {
		// If the sprite exists
		if i.GetSprite() != nil {
			// If the scale is zero
			scale := i.GetScale()
			if scale == util.ZeroScaleOption {
				if i.GetField()[0] == '#' {
					scale = util.Stretch
				} else {
					// Set the scale as the default
					scale = util.DefaultScaleOption
				}
			}

			// If the sprite should repeat
			if scale == util.Tiled {
				// Get the size of the sprite
				spriteSize := i.GetSprite().Frame().Size()
				// Iterate over the y coords of each tile
				for y := spriteSize.Y / 2; y < e.GetCanvas().Bounds().Size().Y; y += spriteSize.Y {
					// Iterate over the x coords of each tile
					for x := spriteSize.X / 2; x < e.GetCanvas().Bounds().Size().X; x += spriteSize.X {
						mat := pixel.IM
						// Move the tile to the position
						mat = mat.Moved(pixel.V(x, y))
						// Draw it
						i.GetSprite().Draw(e.GetCanvas(), mat)
					}
				}
			} else {
				mat := pixel.IM
				// Switch over the scale options
				switch scale {
				case util.NoScale:
					// Nothing needs to be done
				case util.ScaleToFill:
					mat = mat.Scaled(pixel.ZV, math.Max(
						e.GetCanvas().Bounds().Size().X/i.GetSprite().Frame().Size().X,
						e.GetCanvas().Bounds().Size().Y/i.GetSprite().Frame().Size().Y))
				case util.ScaleToFit:
					mat = mat.Scaled(pixel.ZV, math.Min(
						e.GetCanvas().Bounds().Size().X/i.GetSprite().Frame().Size().X,
						e.GetCanvas().Bounds().Size().Y/i.GetSprite().Frame().Size().Y))
				case util.Stretch:
					mat = mat.ScaledXY(pixel.ZV, pixel.V(
						e.GetCanvas().Bounds().Size().X/i.GetSprite().Frame().Size().X,
						e.GetCanvas().Bounds().Size().Y/i.GetSprite().Frame().Size().Y))
				default:
					log.Printf("unknown scale option '%s'", scale)
				}
				// Move it to the center of the canvas
				mat = mat.Moved(e.GetCanvas().Bounds().Center())
				// todo use gravity
				// Draw the sprite
				i.GetSprite().Draw(e.GetCanvas(), mat)
			}
		}
	}
}
