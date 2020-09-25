package button

import (
	"github.com/orfby/ui/pkg/ui/element"
	"reflect"
)

//Function to register the button types
func init() {
	//Register the button types
	element.Register(ImageButtonTypeName,
		reflect.TypeOf((*ImageButton)(nil)).Elem(), NewImageButton)
	element.Register(TextButtonTypeName,
		reflect.TypeOf((*TextButton)(nil)).Elem(), NewTextButton)
}