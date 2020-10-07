package element

// Type for an error when an
// attribute's element ID doesn't
// match an actual ID
type NoElemError struct {
	Element      Element
	ReferencedID string
	AttrName     string
}

// Function to create a NoElemError.
// elem is the element that references
// an ID, refElem is the ID of the element
// being referenced (that doesn't exist),
// and attrName is the name of the attribute
// asking for the ID
func NewNoElemError(elem Element, refElem string, attrName string) NoElemError {
	return NoElemError{elem, refElem, attrName}
}

// Function to return the error string
func (e NoElemError) Error() string {
	return "no element found with ID '" + e.ReferencedID +
		"' (referenced by '" + e.AttrName + "' attribute on XML element '" +
		FullName(e.Element, ".", false) + "')"
}

// todo add more
