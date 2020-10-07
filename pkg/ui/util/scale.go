package util

import (
	"errors"
	"strings"
)

// Type for the scale option of
// something (usually an image)
type ScaleOption string

// Const for not scaling an image
const NoScale = ScaleOption("none")

// Const for image scaling to fill
const ScaleToFill = ScaleOption("fill")

// Const for image scaling to fit
const ScaleToFit = ScaleOption("fit")

// Const for stretching an image
const Stretch = ScaleOption("stretch")

// Const for tiling an image
const Tiled = ScaleOption("tiled")

// A "zero" scale option, as in one where
// it is a zero value
const ZeroScaleOption = ScaleOption("")

// Const for the default scale option
const DefaultScaleOption = ScaleToFit

// Function to parse a string into a scale option type.
// If value is not a valid scale option, the function
// returns an error
func ParseScaleOption(value string) (ScaleOption, error) {
	ret := ScaleOption(strings.ToLower(value))
	if ret != NoScale && ret != ScaleToFill &&
		ret != ScaleToFit && ret != Stretch && ret != Tiled {
		return "", errors.New("invalid scale option '" + value + "'")
	}
	return ret, nil
}
