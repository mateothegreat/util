package ui

import (
	"context"
	"io"
)

// Renderer handles the output of components to a destination.
//
// Arguments:
// - None
//
// Returns:
// - An interface for rendering components to various outputs.
type Renderer interface {
	// Render outputs the component to the destination.
	//
	// Arguments:
	// - ctx: Context for cancellation and timeout control.
	// - component: The component to render.
	//
	// Returns:
	// - Error if rendering fails.
	Render(ctx context.Context, component Component) error

	// Clear removes all rendered content.
	//
	// Arguments:
	// - None
	//
	// Returns:
	// - Error if clearing fails.
	Clear() error
}

// TerminalRenderer renders components to a terminal with ANSI escape codes.
//
// Arguments:
// - None
//
// Returns:
// - A renderer that outputs to terminal with real-time updates.
type TerminalRenderer struct {
	output     io.Writer
	lastHeight int
}

// NewTerminalRenderer creates a new terminal renderer.
//
// Arguments:
// - w: The writer to output to (typically os.Stdout).
//
// Returns:
// - A new terminal renderer instance.
func NewTerminalRenderer(w io.Writer) *TerminalRenderer {
	return &TerminalRenderer{
		output: w,
	}
}

// Render outputs the component to the terminal.
//
// Arguments:
// - ctx: Context for cancellation and timeout control.
// - component: The component to render.
//
// Returns:
// - Error if rendering fails.
func (r *TerminalRenderer) Render(ctx context.Context, component Component) error {
	// Clear previous content
	if err := r.Clear(); err != nil {
		return err
	}

	// Render new content
	content := component.Render(ctx)
	lines := countLines(content)
	r.lastHeight = lines

	_, err := io.WriteString(r.output, content)
	return err
}

// Clear removes all previously rendered content using ANSI escape codes.
//
// Arguments:
// - None
//
// Returns:
// - Error if clearing fails.
func (r *TerminalRenderer) Clear() error {
	if r.lastHeight > 0 {
		// Move cursor up and clear lines
		for i := 0; i < r.lastHeight; i++ {
			if i > 0 {
				_, err := io.WriteString(r.output, "\033[1A") // Move up one line
				if err != nil {
					return err
				}
			}
			_, err := io.WriteString(r.output, "\033[2K\r") // Clear line and return to start
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// countLines counts the number of lines in a string.
//
// Arguments:
// - s: The string to count lines in.
//
// Returns:
// - Number of lines in the string.
func countLines(s string) int {
	if s == "" {
		return 0
	}
	lines := 1
	for _, ch := range s {
		if ch == '\n' {
			lines++
		}
	}
	return lines
}
