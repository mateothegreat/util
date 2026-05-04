package values

// Ptr returns a pointer to a shallow copy of v.
//
// This is a helper function to create pointers to values.
// It is useful when you need to pass a pointer to a value to a function.
//
// It is equivalent to the following code:
//
//	func Ptr[T any](v T) *T { return &v }
//
// Arguments:
//
//	v: The value to create a pointer to.
//
// Returns:
//
//	A pointer to a shallow copy of v.
//
// Example:
//
//	ptr := values.Ptr(1)
//	fmt.Println(*ptr) // 1
//
//	ptr = values.Ptr("hello")
//	fmt.Println(*ptr) // "hello"
func Ptr[T any](v T) *T {
	return &v
}
