// Package maps provides utility functions for working with maps.
package maps

import (
	"strings"
)

// Get returns the value in the map for the provided selector.
//
// The value is returned as a pointer to the type T.
// If the selector is not found, the returned pointer is nil and the boolean is false.
//
// Example:
//
// Arguments:
// - m: the map to query
// - selector: the selector to use
//
// Returns:
// - the value in the map for the provided selector
// - a boolean indicating whether the selector was found
func Get[T any](m map[string]interface{}, selector string) (*T, bool) {
	pathElements := strings.Split(selector, ".")
	valueMap := m
	value := new(T)
	for i, element := range pathElements {
		tmp, exists := valueMap[element]
		if !exists {
			return nil, false
		}

		if i == len(pathElements)-1 {
			*value = tmp.(T)
		} else {
			switch v := tmp.(type) {
			case map[string]interface{}:
				valueMap = v
			default:
				return nil, false
			}
		}
	}

	return value, true
}
