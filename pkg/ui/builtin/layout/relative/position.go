package relative

import (
	"github.com/orfby/ui/pkg/ui/util"
	"strings"
)

//Type for a relative position
type relativePosition struct {
	//Whether the position is relative
	//to the parent
	Parent bool

	//The ID of the element that the
	//position is relative to. Empty
	//if Parent and AbsolutePos are
	//non-zero
	ElementID string

	//The position the element is
	//relative to. Empty if
	//Parent or ElementID are
	//non-zero
	Pos util.RelativeQuantity
}

//A "zero" relative position, as in one where
//all the fields are zero values
var zeroRelativePosition = relativePosition{}

//Function to parse a string into a
//relative position (as a
//reflect.Value type)
func parseRelativePosition(attr string) (pos relativePosition, err error) {
	//Convert the attribute to lowercase
	attr = strings.ToLower(attr)
	//If the attribute is "parent"
	if attr == "parent" {return relativePosition{Parent: true}, nil}

	//Try to parse it as a relative quantity
	quantity, err := util.ParseRelativeQuantity(attr)
	//If it failed
	if err != nil {
		//Just set the ID to the value of the attribute
		return relativePosition{ElementID: attr}, nil
	}
	//Return the position as the relative quantity
	return relativePosition{Pos: quantity}, nil
}