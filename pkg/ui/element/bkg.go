package element

import (
	"encoding/xml"
	"github.com/bhollier/ui/pkg/ui/util"
	"github.com/faiface/pixel"
	"github.com/srwiley/oksvg"
	"image/color"
	"path/filepath"
	"strings"
)

// Implementation of the image
// interface for a background
type Background struct {
	// The background field
	// from xml
	Field string `uixml:"http://github.com/bhollier/ui/api/schema background,optional"`
	// The background's scale option
	Scale util.ScaleOption `uixml:"http://github.com/bhollier/ui/api/schema bkg-scale,optional"`
	// The background's colour
	Color color.RGBA `uixml:"http://github.com/bhollier/ui/api/schema bkg-color,optional"`
	// The background's svg (if
	// applicable)
	svg *oksvg.SvgIcon
	// The element's background
	// sprite
	sprite *pixel.Sprite
}

// Function to get the background's
// XML field
func (b *Background) GetField() string { return b.Field }

// Function to get the background's
// scale option
func (b *Background) GetScale() util.ScaleOption { return b.Scale }

// Function to get the background's
// colour
func (b *Background) GetColor() color.RGBA { return b.Color }

// Function to determine whether
// the background looks like an SVG
func (b *Background) IsSVG() bool {
	// If the image field isn't empty
	return b.GetField() != "" &&
		// And the extension is ".svg"
		strings.ToLower(filepath.Ext(b.GetField())) == ".svg"
}

// Function to get the background's
// SVG
func (b *Background) GetSVG() *oksvg.SvgIcon { return b.svg }

// Function to set the background's
// SVG
func (b *Background) SetSVG(s *oksvg.SvgIcon) { b.svg = s }

// Function to get the background's
// sprite
func (b *Background) GetSprite() *pixel.Sprite { return b.sprite }

// Function to set the background's
// sprite
func (b *Background) SetSprite(s *pixel.Sprite) { b.sprite = s }

// Function to unmarshal an XML element into
// an element. SetAttrs should've been called
// before this function
func (b *Background) UnmarshalXML(*xml.Decoder, xml.StartElement) (err error) {
	// If the scale wasn't given
	if b.Scale == util.ZeroScaleOption {
		// Default for a background is scale to fill
		b.Scale = util.ScaleToFill
	}
	return nil
}

// Function to reset the background
func (b *Background) Reset() {
	// If the image is an SVG
	if b.IsSVG() {
		// Reset the sprite
		b.SetSprite(nil)
	}
}

// Function to determine whether the
// background has been initialised, by
// whether its background sprite has been
// set (assuming it's meant to be set).
// This function doesn't call
// element.IsInitialised
func (b *Background) IsInitialised() bool {
	// If the element doesn't have a background
	return b.GetField() == "" ||
		// Or the background has been initialised
		b.GetSprite() != nil
}

// Function to initialise an element's
// background. Doesn't call element.Init
func InitBkg(e Element, i Image) error {
	// If the image looks like an SVG
	if i.IsSVG() {
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

		// If the svg has been loaded
		if i.GetSVG() != nil {
			// If the image hasn't been created
			// and the width and height are known
			if i.GetSprite() == nil &&
				e.GetActualWidth() != nil &&
				e.GetActualHeight() != nil {
				// Create a picture from the SVG
				pic := util.CreatePictureFromSVG(i.GetSVG(), i.GetScale(),
					*e.GetActualWidth(), *e.GetActualHeight())
				// Create a sprite with the picture and set it
				i.SetSprite(pixel.NewSprite(pic, pic.Bounds()))
			}
		}
		return nil
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
	return nil
}

// Function to draw the background
// of an element. This function
// should be called first
func DrawBkg(e Element, b Image) {
	// Clear the canvas (with a transparent background)
	// (this is why it's important the background is
	// drawn first)
	e.GetCanvas().Clear(color.Transparent)

	// Draw the background
	DrawImage(e, b)
}
