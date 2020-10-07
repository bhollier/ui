package util

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// Type for the absolute size
// of something
type AbsoluteQuantity struct {
	// The quantity itself
	Quantity int32

	// The unit. Only absolute units
	// are allowed (For example,
	// Pixel is allowed, Percent
	// isn't)
	Unit Unit
}

// A "zero" absolute quantity, as in one where
// all the fields are zero values
var ZeroAbsoluteQuantity = AbsoluteQuantity{}

// The default absolute quantity
var DefaultAbsoluteQuantity = AbsoluteQuantity{Quantity: 0, Unit: Pixels}

// Function to parse a string and
// convert it to an absolute unit
func ParseAbsoluteUnit(str string) (unit Unit, err error) {
	// Convert the string to lowercase
	str = strings.ToLower(str)
	switch str {
	case string(Pixels):
		return Pixels, nil
	default:
		return "", errors.New("invalid unit '" + str + "'")
	}
}

// Function to parse a string and convert
// it to an AbsoluteQuantity
func ParseAbsoluteQuantity(str string) (size AbsoluteQuantity, err error) {
	// Remove all the whitespace from the string
	strMod := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
	// Convert the string to lowercase
	strMod = strings.ToLower(strMod)

	// Find the last digit in the string
	lastDigit := strings.LastIndexFunc(strMod, func(r rune) bool {
		return unicode.IsDigit(r)
	})
	if lastDigit == -1 {
		return AbsoluteQuantity{}, errors.New("invalid quantity '" + str + "'")
	}

	// Convert the start of the string to an int
	quantity, err := strconv.ParseInt(strMod[:lastDigit+1], 10, 32)
	if err != nil {
		return AbsoluteQuantity{}, errors.New("invalid quantity '" + str + "'")
	}
	size.Quantity = int32(quantity)

	// A size cannot be negative
	if size.Quantity < 0 {
		return AbsoluteQuantity{}, errors.New("invalid quantity '" + str + "'")
	}

	// Convert the unit
	size.Unit, err = ParseAbsoluteUnit(strMod[lastDigit+1:])
	if err != nil {
		return AbsoluteQuantity{}, errors.New("invalid quantity '" + str + "'")
	}
	return size, nil
}
