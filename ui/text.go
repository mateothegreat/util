package ui

import (
	"context"
	"sync"
)

// Text creates a static text component.
//
// Arguments:
// - content: The text content to display.
//
// Returns:
// - A component that renders the given text.
func Text(content string) Component {
	return ComponentFunc(func(ctx context.Context) string {
		return content
	})
}

// TextComponent is a dynamic text component that can be updated.
//
// Arguments:
// - None
//
// Returns:
// - A component that can display and update text in real-time.
type TextComponent struct {
	mu      sync.RWMutex
	content string
}

// NewText creates a new dynamic text component.
//
// Arguments:
// - initialContent: The initial text to display.
//
// Returns:
// - A new TextComponent instance.
func NewText(initialContent string) *TextComponent {
	return &TextComponent{
		content: initialContent,
	}
}

// Render returns the current text content.
//
// Arguments:
// - ctx: Context for cancellation and timeout control.
//
// Returns:
// - The current text content of the component.
func (t *TextComponent) Render(ctx context.Context) string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.content
}

// SetText updates the text content.
//
// Arguments:
// - content: The new text content to display.
//
// Returns:
// - None
func (t *TextComponent) SetText(content string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.content = content
}

// GetText returns the current text content.
//
// Arguments:
// - None
//
// Returns:
// - The current text content.
func (t *TextComponent) GetText() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.content
}
