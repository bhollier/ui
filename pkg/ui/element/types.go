package element

import (
	"encoding/xml"
	"net/http"
	"reflect"
)

//Type for an element factory
type Factory func(fs http.FileSystem, name xml.Name, parent Layout) Element

//Type for storing information
//about an element type
type elementType struct {
	//The element's XML name
	Name xml.Name
	//The element's reflection type
	ReflectType reflect.Type
	//A factory function to create
	//an element of this type
	Factory Factory
}

//Type for an element types map
type elementTypesMap map[xml.Name]elementType

//Map of element types, with the
//key being the element's XML name
var elementTypes elementTypesMap

//Function to initialise the element types map
func init() {
	//Create the element types map
	elementTypes = make(elementTypesMap, 0)
}

//Function to convert an xml name to a string
func XMLNameToString(name xml.Name) string {
	if name.Space == "" {
		return name.Local
	} else {
		return name.Space + ":" + name.Local
	}
}

//Function to determine if
//two XML names match. Returns
//true if n1 and n2 match exactly
//or if n1.Space == "" and n1.Local
//and n2.Local match
func XMLNameMatch(n1 xml.Name, n2 xml.Name) bool {
	return n1.Local == n2.Local &&
		(n1.Space == n2.Space || n1.Space == "")
}

//Function to register a UI element type
func Register(name xml.Name, reflectType reflect.Type, factory Factory) {
	//Add the factory to the map
	elementTypes[name] = elementType{
		ReflectType: reflectType,
		Factory:     factory,
	}
}

//Function to get an element's name.
//The function returns an empty string
//if the element isn't registered
func Name(e Element, includeNamespace bool) (name string) {
	if includeNamespace {
		name = XMLNameToString(e.GetName())
	} else {
		name = e.GetName().Local
	}
	if e.GetID() != nil {
		name += "(id=" + *e.GetID() + ")"
	}
	return
}

//Function to get an element's full name,
//which includes its parent's and
//grandparent's etc. element names.
//The function calls any unknown elements
//"unknown" in the path
func FullName(e Element, sep string, includeNamespace bool) (name string) {
	name = Name(e, includeNamespace)
	e = e.GetParent()
	//While the current element isn't nil
	for e != nil {
		//Add it to the name
		name = Name(e, includeNamespace) + sep + name
		//Go to the next element up the tree
		e = e.GetParent()
	}
	return
}

//Function to create an element with the
//given XML name. If none is found and
//name.Space is empty, the function will
//search for an element type with a
//matching name.Local and with a matching
//name.Space in parent.GetNamespaces(),
//then if still not found will search for
//an element type with a matching name.Local
func New(fs http.FileSystem, name xml.Name, parent Layout) Element {
	//Try to get the element type
	elemType, ok := elementTypes[name]
	//If it was found, call the type's factory
	if ok {
		return elemType.Factory(fs, name, parent)
		//If it wasn't found, but the namespace is empty
	} else if name.Space == "" {
		if parent != nil {
			//Iterate over the parent's namespaces
			for _, ns := range parent.GetNamespaces() {
				elemType, ok = elementTypes[xml.Name{
					Local: name.Local, Space: ns}]
				if ok {
					return elemType.Factory(fs, name, parent)
				}
			}
		}

		//Iterate over the element types
		for elemTypeName, elemType := range elementTypes {
			//If the local name is a match
			if elemTypeName.Local == name.Local {
				//Assume that it's a match and call
				//the type's factory
				return elemType.Factory(fs, name, parent)
			}
		}
	}
	return nil
}
