package ui

import (
	"context"
	"sync"
	"testing"
)

// TestText tests the static Text component.
func TestText(t *testing.T) {
	ctx := context.Background()

	// Test empty text
	comp := Text("")
	if result := comp.Render(ctx); result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}

	// Test with content
	comp = Text("Hello, Terminal!")
	if result := comp.Render(ctx); result != "Hello, Terminal!" {
		t.Errorf("Expected 'Hello, Terminal!', got '%s'", result)
	}

	// Test with newlines
	comp = Text("Line 1\nLine 2\nLine 3")
	if result := comp.Render(ctx); result != "Line 1\nLine 2\nLine 3" {
		t.Errorf("Expected multiline text, got '%s'", result)
	}
}

// TestTextComponent tests the dynamic TextComponent.
func TestTextComponent(t *testing.T) {
	ctx := context.Background()

	// Test creation and initial content
	comp := NewText("Initial")
	if result := comp.Render(ctx); result != "Initial" {
		t.Errorf("Expected 'Initial', got '%s'", result)
	}

	// Test GetText
	if text := comp.GetText(); text != "Initial" {
		t.Errorf("GetText: Expected 'Initial', got '%s'", text)
	}

	// Test SetText
	comp.SetText("Updated")
	if result := comp.Render(ctx); result != "Updated" {
		t.Errorf("Expected 'Updated', got '%s'", result)
	}

	// Test empty text
	comp.SetText("")
	if result := comp.Render(ctx); result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

// TestTextComponentConcurrency tests thread safety of TextComponent.
func TestTextComponentConcurrency(t *testing.T) {
	ctx := context.Background()
	comp := NewText("Start")

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	// Multiple readers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = comp.Render(ctx)
				_ = comp.GetText()
			}
		}()
	}

	// Multiple writers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				comp.SetText(string(rune('A' + id)))
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for any errors (there shouldn't be any race conditions)
	for err := range errors {
		t.Errorf("Concurrency error: %v", err)
	}
}
