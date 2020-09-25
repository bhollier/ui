package util

import (
	"github.com/faiface/pixel"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

//Function to load an image and
//convert it to a pixel.Picture
func LoadPicture(path string) (pixel.Picture, error) {
	//Open the file
	file, err := os.Open(path)
	if err != nil {return nil, err}
	defer file.Close()
	//Decode it as an image
	img, _, err := image.Decode(file)
	if err != nil {return nil, err}

	//Convert it to a pixel.Picture and return
	return pixel.PictureDataFromImage(img), nil
}