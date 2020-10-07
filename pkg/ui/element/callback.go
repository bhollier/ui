package element

import (
	"errors"
)

// Type for a callback
type Callback func(Element) error

// Type for a callback map
type callbackMap map[string]Callback

// Map of callbacks, with the
// key being the callback's name
var callbacks callbackMap

// Function to initialise the callback map
func init() {
	// Create the element types map
	callbacks = make(callbackMap, 0)
}

// Function to register a new callback
func RegisterCallback(name string, c Callback) {
	// Add the callback
	callbacks[name] = c
}

// Function to call a callback
func Call(name string, e Element) error {
	// Try to get the callback
	callback, ok := callbacks[name]
	// If it was found, call it
	if ok {
		return callback(e)
	} else {
		return errors.New("unknown callback '" + name + "'")
	}
}
