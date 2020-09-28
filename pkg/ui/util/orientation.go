package util

import (
	"errors"
	"strings"
)

//Type for the orientation of something
type Orientation string

//Const of a horizontal orientation
const HorizontalOrientation = Orientation("horizontal")

//Const of a vertical orientation
const VerticalOrientation = Orientation("vertical")

//Const of the default orientation
const DefaultOrientation = VerticalOrientation

//Function to parse a string into an orientation type.
//If value is not a valid orientation, the function
//returns an error. If the value is an empty string,
//the default orientation 'def' is used instead
func ParseOrientation(value string, def Orientation) (Orientation, error) {
	//Convert the string to lowercase
	value = strings.ToLower(value)
	if value == "" {
		return def, nil
	} else if Orientation(value) != HorizontalOrientation &&
		Orientation(value) != VerticalOrientation {
		return "", errors.New("invalid orientation '" + value + "'")
	}
	return Orientation(value), nil
}
