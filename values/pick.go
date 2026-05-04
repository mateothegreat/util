package values

import "reflect"

// Pick returns the value of the key in the map if it exists. Otherwise, it
// returns the default value.
func Pick[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultValue
}

// PickHasValue returns the first non-zero value from the provided values.
//
// Arguments:
// - values: A variadic list of values to check.
//
// Returns:
// - The first non-zero value found, or zero value if none found.
func PickHasValue[T any](values ...T) T {
	var zero T
	for _, v := range values {
		if !IsZero(v) {
			return v
		}
	}
	return zero
}

// IsZero checks if a value is zero.
//
// Arguments:
// - v: The value to check.
//
// Returns:
// - True if the value is zero, false otherwise.
func IsZero[T any](v T) bool {
	val := any(v)

	if val == nil {
		return true
	}

	// Fast path for common types
	switch ptr := val.(type) {
	case string:
		return ptr == ""
	case int, int8, int16, int32, int64:
		return ptr == 0
	case uint, uint8, uint16, uint32, uint64:
		return ptr == 0
	case float32, float64:
		return ptr == 0
	case bool:
		return ptr == false
	case *string:
		return ptr == nil || *ptr == ""
	// ... other common types

	default:
		// Correct fallback for ALL other types
		return reflect.ValueOf(v).IsZero()
	}
}
