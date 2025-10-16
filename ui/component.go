package ui

import "context"

// Component represents a UI element that can be rendered.
//
// Arguments:
// - None
//
// Returns:
// - A UI element that can be rendered to terminal output.
type Component interface {
	// Render returns the string representation of the component.
	//
	// Arguments:
	// - ctx: Context for cancellation and timeout control.
	//
	// Returns:
	// - String representation of the component to be displayed.
	Render(ctx context.Context) string
}

// ComponentFunc is a function adapter for the Component interface.
//
// Arguments:
// - ctx: Context for cancellation and timeout control.
//
// Returns:
// - String representation of the component to be displayed.
type ComponentFunc func(ctx context.Context) string

// Render implements the Component interface for ComponentFunc.
//
// Arguments:
// - ctx: Context for cancellation and timeout control.
//
// Returns:
// - String representation of the component to be displayed.
func (f ComponentFunc) Render(ctx context.Context) string {
	return f(ctx)
}
