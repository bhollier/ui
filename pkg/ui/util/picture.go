package util

import (
	"errors"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
)

//Function to load an SVG
func LoadSVG(fs http.FileSystem, path string,
	color color.RGBA) (*oksvg.SvgIcon, error) {
	//Open the file
	file, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//Read the whole file
	//todo could possibly be better
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	//Convert to a string
	str := string(bytes)
	//Replace any instances of "currentColor" with the given colour
	str = strings.ReplaceAll(str, "currentColor",
		fmt.Sprintf("rgb(%d,%d,%d)",
			color.R, color.G, color.B))
	//todo there are probably other attributes to replace

	//Convert the string to an SVG icon and return
	return oksvg.ReadIconStream(strings.NewReader(str))
}

//Function to create a pixel.Picture
//from a given SVG
func CreatePictureFromSVG(svg *oksvg.SvgIcon, scale ScaleOption,
	w, h float64) pixel.Picture {
	//Convert the SVG's view box to a pixel rect
	viewbox := pixel.R(svg.ViewBox.X, svg.ViewBox.Y,
		svg.ViewBox.W, svg.ViewBox.H)
	//If the view box's size is 0, 0
	if viewbox.Size() == pixel.ZV {
		//Change it to a 16x16 box
		viewbox = pixel.R(
			0, 0, 16, 16)
	}

	//If the sprite should repeat
	if scale == Tiled {
		//Create a blank image
		img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

		//Get the size of the svg
		svgSize := viewbox.Size()
		//Iterate over the y coords of each tile
		for y := 0.0; y < h; y += svgSize.Y {
			//Iterate over the x coords of each tile
			for x := 0.0; x < w; x += svgSize.X {
				//Set the target area of the SVG
				//todo transform
				svg.SetTarget(viewbox.Min.X, viewbox.Min.Y,
					viewbox.Max.X, viewbox.Max.Y)

				//Draw the SVG onto the image
				svg.Draw(rasterx.NewDasher(
					int(svgSize.X), int(svgSize.Y),
					rasterx.NewScannerGV(int(svgSize.X), int(svgSize.Y),
						img, img.Bounds())), 1)
			}
		}

		//Create a pixel picture from the image
		return pixel.PictureDataFromImage(img)

	} else {
		//The width and height of the image
		//var xScale, yScale float64
		var newWidth, newHeight float64
		//Switch over the scale options
		switch scale {
		case NoScale:
			//Just use the width and height of the SVG
			newWidth, newHeight =
				viewbox.Size().X, viewbox.Size().Y
		case ScaleToFill:
			scale := math.Max(
				w/viewbox.Size().X,
				h/viewbox.Size().Y)
			newWidth = viewbox.Size().X * scale
			newHeight = viewbox.Size().Y * scale
		case ScaleToFit:
			scale := math.Min(
				w/viewbox.Size().X,
				h/viewbox.Size().Y)
			newWidth = viewbox.Size().X * scale
			newHeight = viewbox.Size().Y * scale
		case Stretch:
			newWidth, newHeight = w, h
		default:
			log.Printf("unknown scale option '%s'", scale)
		}

		//Create a blank image
		img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

		//Set the target area of the SVG
		svg.SetTarget((w/2)-(newWidth/2),
			(h/2)-(newHeight/2),
			newWidth, newHeight)

		//Draw the SVG onto the image
		svg.Draw(rasterx.NewDasher(int(w), int(h),
			rasterx.NewScannerGV(int(w), int(h),
				img, img.Bounds())), 1)

		//Create a pixel picture from the image
		return pixel.PictureDataFromImage(img)
	}
}

//Function to load an image and
//convert it to a pixel.Picture
func LoadPicture(fs http.FileSystem, path string) (pixel.Picture, error) {
	//Open the file
	file, err := fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//Decode it as an image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	//Convert it to a pixel.Picture and return
	return pixel.PictureDataFromImage(img), nil
}

//Function to create a picture
//from an XML string
func CreatePictureFromField(fs http.FileSystem, field string) (pixel.Picture, error) {
	if field != "" {
		//If the first character is a hash
		if field[0] == '#' {
			//Convert the field to a colour type
			colour, err := ParseColor(field)
			if err != nil {
				return nil, errors.New("invalid colour attribute value '" + field + "'")
			}
			//Create a 1x1 image
			img := image.NewRGBA(image.Rect(0, 0, 2, 2))
			for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
				for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
					img.SetRGBA(x, y, colour)
				}
			}
			//Convert it to a pixel picture
			pic := pixel.PictureDataFromImage(img)
			//Return the picture
			return pic, nil
		} else {
			//Load the picture
			pic, err := LoadPicture(fs, field)
			if err != nil {
				return nil, err
			}
			//Return the picture
			return pic, nil
		}
	}
	return nil, nil
}
