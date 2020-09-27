package builtin

import (
	_ "github.com/orfby/ui/pkg/ui/builtin/button"
	_ "github.com/orfby/ui/pkg/ui/builtin/layout"
	"github.com/orfby/ui/pkg/ui/element"
	"reflect"
)

//Function to register the built-in types
func init() {
	//Register the built-in element types
	element.Register(FixedRatioTypeName,
		reflect.TypeOf((*FixedRatio)(nil)).Elem(), NewFixedRatio)
	element.Register(ImageTypeName,
		reflect.TypeOf((*Image)(nil)).Elem(), NewImage)
	element.Register(ImportTypeName,
		reflect.TypeOf((*Import)(nil)).Elem(), NewImport)
	element.Register(ScrollTypeName,
		reflect.TypeOf((*Scroll)(nil)).Elem(), NewScroll)
}
