package ui

import (
	"context"
	"testing"
)

// TestComponentFunc tests the ComponentFunc adapter.
func TestComponentFunc(t *testing.T) {
	ctx := context.Background()

	// Test basic function
	component := ComponentFunc(func(ctx context.Context) string {
		return "Hello, World!"
	})

	result := component.Render(ctx)
	if result != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", result)
	}

	// Test with context check
	cancelCtx, cancel := context.WithCancel(ctx)
	cancel()

	component = ComponentFunc(func(ctx context.Context) string {
		select {
		case <-ctx.Done():
			return "cancelled"
		default:
			return "running"
		}
	})

	result = component.Render(cancelCtx)
	if result != "cancelled" {
		t.Errorf("Expected 'cancelled', got '%s'", result)
	}
}
