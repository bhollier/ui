package util

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

//Type for a size unit
type Unit string

//Const for a unit of pixels
const Pixels = Unit("px")
//Const for a percentage "unit"
const Percent = Unit("%")

//Function to parse a string and
//convert it to a unit
func ParseUnit(str string) (unit Unit, err error) {
	switch str {
	case string(Pixels):
		return Pixels, nil
	case string(Percent):
		return Percent, nil
	default:
		return "", errors.New("invalid unit '" + str + "'")
	}
}

//Type for a relative quantity
type RelativeQuantity struct {
	//The quantity itself. Zero if
	//MatchParent or MatchContent
	//are non-zero
	Quantity int32

	//The quantity's unit. Zero if
	//MatchParent or MatchContent
	//are non-zero
	Unit Unit
}

//A "zero" relative quantity, as in one where
//all the fields are zero values
var ZeroRelativeQuantity = RelativeQuantity{}
//The default relative quantity
var DefaultRelativeQuantity = RelativeQuantity{Quantity: 0, Unit: Pixels}

//Function to parse a string and convert
//it to a RelativeQuantity
func ParseRelativeQuantity(str string) (size RelativeQuantity, err error) {
	//Remove all the whitespace from the string
	strMod := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {return -1}
		return r
	}, str)
	//Convert the string to lowercase
	strMod = strings.ToLower(strMod)

	//Find the last digit in the string
	lastDigit := strings.LastIndexFunc(strMod, func(r rune) bool {
		return unicode.IsDigit(r)
	})
	if lastDigit == -1 {
		return RelativeQuantity{}, errors.New("invalid quantity '" + str + "'")
	}

	//Convert the start of the string to an int
	quantity, err := strconv.ParseInt(strMod[:lastDigit + 1], 10, 32)
	if err != nil {
		return RelativeQuantity{}, errors.New("invalid quantity '" + str + "'")
	}
	size.Quantity = int32(quantity)

	//A size cannot be negative
	if size.Quantity < 0 {
		return RelativeQuantity{}, errors.New("invalid quantity '" + str + "'")
	}

	//Convert the unit
	size.Unit, err = ParseUnit(strMod[lastDigit + 1:])
	if err != nil {
		return RelativeQuantity{}, errors.New("invalid quantity '" + str + "'")
	}
	return size, nil
}

//Type for the relative size
//(usually width or height) of
//an element, which can match
//its parent, content or have
//a relative quantity
type RelativeSize struct {
	//Whether the quantity is the
	//same as the parent element
	MatchParent bool

	//Whether the quantity is the
	//same as the element's content
	MatchContent bool

	//The quantity itself. Zero if
	//MatchParent or MatchContent
	//are non-zero
	RelativeQuantity
}

//A "zero" relative quantity, as in one where
//all the fields are zero values
var ZeroRelativeSize = RelativeSize{}
//The default relative quantity
var DefaultRelativeSize = RelativeSize{RelativeQuantity: DefaultRelativeQuantity}

//Function to parse a string and convert
//it to a RelativeSize
func ParseRelativeSize(str string) (size RelativeSize, err error) {
	//Remove all the whitespace from the string
	strMod := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {return -1}
		return r
	}, str)
	//Convert the string to lowercase
	strMod = strings.ToLower(strMod)
	//If the string is equal to "match_parent"
	if strMod == "match_parent" {
		size.MatchParent = true
		return size, nil
	}
	//If the string is equal to "match_content"
	if strMod == "match_content" {
		size.MatchContent = true
		return size, nil
	}

	//Otherwise convert it to a relative quantity
	quantity, err := ParseRelativeQuantity(strMod)
	if err != nil {return RelativeSize{}, err}
	return RelativeSize{RelativeQuantity: quantity}, nil
}