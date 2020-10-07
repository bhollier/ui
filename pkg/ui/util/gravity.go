package util

import (
	"errors"
	"strings"
)

type GravSide int8

const (
	// Const for the middle of a gravity axis
	GravCenter = GravSide(0)
	// Const for the negative end of
	// a gravity axis (left, top, etc.)
	GravNeg = GravSide(-1)
	// Const for the positive end of
	// a gravity axis (right, bottom, etc.)
	GravPos = GravSide(1)
)

// Type for the gravity of something
type Gravity struct {
	// The horizontal gravity. Either
	// GravNeg (left), GravCenter
	// (center), GravPos (right)
	HorizGravity GravSide

	// The vertical gravity. Either
	// GravNeg (top), GravCenter
	// (center), GravPos (bottom)
	VertGravity GravSide
}

// The possible gravity types as a map
var GravityTypes = map[string]Gravity{
	"center":       {GravCenter, GravCenter},
	"top":          {GravCenter, GravNeg},
	"bottom":       {GravCenter, GravPos},
	"left":         {GravNeg, GravCenter},
	"right":        {GravPos, GravCenter},
	"top-left":     {GravNeg, GravNeg},
	"top-right":    {GravPos, GravNeg},
	"bottom-left":  {GravNeg, GravPos},
	"bottom-right": {GravPos, GravPos},
}

// The default gravity
var DefaultGravity = GravityTypes["top-left"]

// Function to convert a gravity to a string.
// Returns an empty string if the gravity type
// is unknown (or invalid)
func (g *Gravity) String() string {
	// Iterate over the gravity types
	for name, gravity := range GravityTypes {
		// If the horizontal and vertical gravity is a match
		if g.HorizGravity == gravity.HorizGravity &&
			g.VertGravity == gravity.VertGravity {
			return name
		}
	}
	// Otherwise return an empty string
	return ""
}

// Function to parse a string into a gravity type.
// If value is not a valid gravity, the function
// returns an error. If the value is an empty string,
// the default gravity 'def' is used instead
func ParseGravity(value string, def Gravity) (Gravity, error) {
	// Convert the string to lowercase
	value = strings.ToLower(value)
	// If the value isn't given, use the default
	if value == "" {
		return def, nil
	}
	// Try to get the gravity type from the map
	gravity, ok := GravityTypes[value]
	if ok {
		return gravity, nil
	} else {
		return Gravity{}, errors.New("invalid gravity '" + value + "'")
	}
}
