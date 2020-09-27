package util

import (
	"errors"
	"github.com/faiface/pixel"
	"strconv"
	"strings"
)

//Type for a ratio (such as for
//an aspect ratio)
type Ratio struct {
	//The left part of the ratio
	Left int
	//The right part of the ratio
	Right int
}

//Function to parse a string into a ratio type.
//If value is not a valid ratio, the function
//returns an error
func ParseRatio(value string) (r Ratio, err error) {
	//Find the colon
	colon := strings.IndexRune(value, ':')
	if colon == -1 {
		return Ratio{}, errors.New("invalid ratio '" + value + "'")
	}
	//Convert the left side
	r.Left, err = strconv.Atoi(value[:colon])
	if err != nil {
		return Ratio{}, errors.New("invalid ratio '" + value + "'")
	}
	//Convert the right side
	r.Right, err = strconv.Atoi(value[colon+1:])
	if err != nil {
		return Ratio{}, errors.New("invalid ratio '" + value + "'")
	}
	return r, nil
}

//Function to restrict the given
//dimensions so they fit the ratio
//todo this works but there's gotta be a better way
func (r *Ratio) RestrictDimensions(dimensions pixel.Vec) pixel.Vec {
	//Calculate the width option
	option1 := pixel.V(dimensions.X, dimensions.X/(float64(r.Left)/float64(r.Right)))
	//Calculate the height option
	option2 := pixel.V(dimensions.Y/(float64(r.Right)/float64(r.Left)), dimensions.Y)

	//If the first option is too large
	if option1.X > dimensions.X || option1.Y > dimensions.Y {
		//Use the second
		return option2

		//Otherwise use the first
	} else {
		return option1
	}
}
