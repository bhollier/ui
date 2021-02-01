package element

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/bhollier/ui/pkg/ui/util"
	"image/color"
	"math/bits"
	"reflect"
	"strconv"
	"strings"
)

// Type for an attribute parser
type AttrParser func(attr string) (reflect.Value, error)

// Type for an attribute parser map
type attributeTypesMap map[reflect.Type]AttrParser

// Map of attribute types, with the
// key being the attribute's (reflect)
// type
var attributeTypes attributeTypesMap

// Function to initialise the attributes types map
func init() {
	// Create the attributes types map,
	// with all the primitive types (except
	// uintptr, complex32 and complex64) and
	// the util.* types
	attributeTypes = attributeTypesMap{
		// "Parsing" a string (function just spits it back)
		reflect.TypeOf((*string)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			return reflect.ValueOf(attr), nil
		},
		// Parsing a boolean type
		reflect.TypeOf((*bool)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseBool(attr)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		},
		// Parsing an int type
		reflect.TypeOf((*int)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.Atoi(attr)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		},
		// Parsing an int8 type
		reflect.TypeOf((*int8)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseInt(attr, 10, 8)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(int8(val)), nil
		},
		// Parsing an int16 type
		reflect.TypeOf((*int16)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseInt(attr, 10, 16)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(int16(val)), nil
		},
		// Parsing an int32 type
		reflect.TypeOf((*int32)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseInt(attr, 10, 32)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(int32(val)), nil
		},
		// Parsing an int64 type
		reflect.TypeOf((*int32)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseInt(attr, 10, 64)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		},
		// Parsing a uint type
		reflect.TypeOf((*uint)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseUint(attr, 10, bits.UintSize)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(uint(val)), nil
		},
		// Parsing a uint8 type
		reflect.TypeOf((*uint8)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseUint(attr, 10, 8)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(uint8(val)), nil
		},
		// Parsing a uint16 type
		reflect.TypeOf((*uint16)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseUint(attr, 10, 16)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(uint16(val)), nil
		},
		// Parsing a uint32 type
		reflect.TypeOf((*uint32)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseUint(attr, 10, 32)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(uint32(val)), nil
		},
		// Parsing a uint64 type
		reflect.TypeOf((*uint32)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseUint(attr, 10, 64)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		},
		// Parsing a byte type
		reflect.TypeOf((*byte)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseUint(attr, 10, 8)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(byte(val)), nil
		},
		// Parsing a rune type
		reflect.TypeOf((*rune)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseInt(attr, 10, 32)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(rune(val)), nil
		},
		// Parsing a float32 type
		reflect.TypeOf((*float32)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseFloat(attr, 32)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(float32(val)), nil
		},
		// Parsing a float64 type
		reflect.TypeOf((*float64)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := strconv.ParseFloat(attr, 64)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		},
		// Parsing a "color.RGBA" type
		reflect.TypeOf((*color.RGBA)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := util.ParseColor(attr)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		},
		// Parsing a "util.Unit" type (as any
		// unit type, there's no way to know
		// whether the unit has to be absolute)
		reflect.TypeOf((*util.Unit)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := util.ParseUnit(attr)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		},
		// Parsing a "util.RelativeSize" type
		reflect.TypeOf((*util.RelativeSize)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := util.ParseRelativeSize(attr)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		},
		// Parsing a "util.AbsoluteQuantity" type
		reflect.TypeOf((*util.AbsoluteQuantity)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := util.ParseAbsoluteQuantity(attr)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		},
		// Parsing a "util.Gravity" type
		reflect.TypeOf((*util.Gravity)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := util.ParseGravity(attr, util.DefaultGravity)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		},
		// Parsing a "util.Ratio" type
		reflect.TypeOf((*util.Ratio)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := util.ParseRatio(attr)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		},
		// Parsing a "util.Orientation" type
		reflect.TypeOf((*util.Orientation)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := util.ParseOrientation(attr, util.DefaultOrientation)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		},
		// Parsing a "util.ScaleOption" type
		reflect.TypeOf((*util.ScaleOption)(nil)).Elem(): func(attr string) (reflect.Value, error) {
			val, err := util.ParseScaleOption(attr)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(val), nil
		},
	}
}

// Function to register an
// attribute type
func RegisterAttrType(t reflect.Type, p AttrParser) {
	// Register the attribute type
	attributeTypes[t] = p
}

// Function to parse an attribute
// string into the given type
func ParseAttr(t reflect.Type, v string) (reflect.Value, error) {
	// Try to get the attribute factory
	parser, ok := attributeTypes[t]
	// If it was found, call the parser
	if ok {
		return parser(v)
	} else {
		return reflect.Value{}, errors.New(
			"unknown attribute type '" + t.Name() + "'")
	}
}

// Function to parse the given xml
// attributes and set the fields of
// the given element using uixml tags.
// This function searches for tags
// recursively. It does not support
// arrays or maps
func SetAttrs(e Element, attrs []xml.Attr) error {
	// Get the element's type info
	t := reflect.TypeOf(e).Elem()
	v := reflect.ValueOf(e).Elem()

	// Firstly, look for a namespace attribute
	for _, attr := range attrs {
		if attr.Name.Space == "xmlns" {
			// Add the namespace to the element
			e.AddNamespace(attr.Value)
		} else if attr.Name.Local == "xmlns" {
			e.AddNamespace(attr.Name.Space)
		}
	}

	// Create a map of fields
	type Field struct {
		Name     string
		Value    reflect.Value
		Optional bool
		Set      bool
	}
	fields := make(map[xml.Name]*Field, 0)

	var findFieldsWithTag func(t reflect.Type, v reflect.Value) error
	findFieldsWithTag = func(t reflect.Type, v reflect.Value) error {
		// Iterate over the element's fields
		for i := 0; i < t.NumField(); i++ {
			// Try to get the field's uixml tag
			tag, ok := t.Field(i).Tag.Lookup("uixml")
			// If it's found
			if ok {
				// If the tag wants the field to be hidden,
				// continue
				if tag == "hidden" {
					continue
				}

				// Create a field with the default values
				field := Field{
					Name:     t.Field(i).Name,
					Value:    v.Field(i),
					Optional: false,
					Set:      false,
				}

				// Split the tag by commas
				commaSepList := strings.Split(tag, ",")

				// Parse the name of the attribute
				spaceIndex := strings.Index(commaSepList[0], " ")
				var attrName xml.Name
				if spaceIndex != -1 {
					// Set the namespace
					attrName.Space = commaSepList[0][:spaceIndex]
					// Set the local
					attrName.Local = commaSepList[0][spaceIndex+1:]

					// Otherwise just set local as the whole name
				} else {
					attrName.Local = commaSepList[0]
				}

				// Iterate over all the tokens but the first
				// (which is the field's name)
				for j := 1; j < len(commaSepList); j++ {
					switch commaSepList[j] {
					// If the tag specifies the field is optional
					case "optional":
						field.Optional = true

						// Otherwise the token is unknown and so return an error
					default:
						return errors.New("unknown token '" + commaSepList[j] +
							"' in uixml tag on XML element '" +
							FullName(e, ".", false) + "'")
					}
				}
				// Add the field to the map
				fields[attrName] = &field

				// If it doesn't have a tag but
				// has subfields that might have a tag
			} else if v.Field(i).Kind() == reflect.Struct {
				// Find fields in the struct
				err := findFieldsWithTag(t.Field(i).Type, v.Field(i))
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Find fields
	err := findFieldsWithTag(t, v)
	if err != nil {
		return err
	}

	// Iterate over the attributes
	for _, attr := range attrs {
		// If the attribute is for the namespace,
		// we've already dealt with it so skip it
		if attr.Name.Space == "xmlns" ||
			attr.Name.Local == "xmlns" {
			continue
		}

		// Look for the attribute in the fields map
		field, ok := fields[attr.Name]
		// If it doesn't exist
		if !ok {
			// If the attribute doesn't
			// have a namespace
			if attr.Name.Space == "" {
				// Look for the attribute in the
				// element's namespaces
				for _, ns := range e.GetNamespaces() {
					field, ok = fields[xml.Name{
						Space: ns, Local: attr.Name.Local}]
					if ok {
						break
					}
				}
			}

			// If it still wasn't found
			if !ok {
				return errors.New("unknown attribute '" + XMLNameToString(attr.Name) +
					"' on XML element '" + FullName(e, ".", false) + "'")
			}
		}

		// Try to parse the attribute
		val, err := ParseAttr(field.Value.Type(), attr.Value)
		if err != nil {
			return errors.New(fmt.Sprintf("error parsing attribute '"+
				XMLNameToString(attr.Name)+"' on XML element '"+
				FullName(e, ".", false)+"': %+v", err))
		}
		// Otherwise set the value
		field.Value.Set(val)
		// Set the field as set
		field.Set = true
	}

	// Iterate over the fields
	for attrName, field := range fields {
		// If the field isn't optional but wasn't set
		if !field.Optional && !field.Set {
			return errors.New("no '" + XMLNameToString(attrName) + "' attribute on " +
				"XML element '" + FullName(e, ".", false) + "'")
		}
	}

	return nil
}
