package relative

import (
	"github.com/bhollier/ui/pkg/ui/element"
	"reflect"
)

// Function to register the relative types
func init() {
	// Register the relative attribute types
	element.RegisterAttrType(
		reflect.TypeOf((*relativePosition)(nil)).Elem(), func(attr string) (reflect.Value, error) {
			val, err := parseRelativePosition(attr)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		})

	// Register the relative element types
	element.Register(LayoutTypeName,
		reflect.TypeOf((*Layout)(nil)).Elem(), NewLayout)
}
