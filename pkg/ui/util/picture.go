package util

import (
	"errors"
	"github.com/faiface/pixel"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

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
