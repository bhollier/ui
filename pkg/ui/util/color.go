package util

import (
	"encoding/hex"
	"errors"
	"image/color"
)

//Function to parse a colour string to a color.RGBA type
func ParseColor(str string) (color.RGBA, error) {
	//If the first character is a hash
	if str[0] == '#' {
		//Get rid of it
		str = str[1:]
	}
	//If the string is too long
	if len(str) > 8 {
		return color.RGBA{}, errors.New("invalid colour format")

		//If the string is way too short
	} else if len(str) < 6 {
		return color.RGBA{}, errors.New("invalid colour format")

		//If the string is missing the alpha
	} else {
		for len(str) < 8 {
			str = str + "F"
		}
	}

	//Decode the string with hex
	fields, err := hex.DecodeString(str)
	if err != nil {
		return color.RGBA{}, errors.New("invalid colour format")
	}
	//Return a colour type from the hex fields
	return color.RGBA{R:fields[0], G:fields[1], B:fields[2], A:fields[3]}, nil
}