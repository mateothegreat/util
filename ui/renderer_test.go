package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestCountLines tests the line counting function.
func TestCountLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"empty", "", 0},
		{"single line", "hello", 1},
		{"two lines", "hello\nworld", 2},
		{"three lines", "one\ntwo\nthree", 3},
		{"trailing newline", "hello\n", 2},
		{"multiple newlines", "a\n\n\nb", 4},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := countLines(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %d lines, got %d for input '%s'", tc.expected, result, tc.input)
			}
		})
	}
}

// TestTerminalRenderer tests the terminal renderer.
func TestTerminalRenderer(t *testing.T) {
	ctx := context.Background()

	t.Run("basic render", func(t *testing.T) {
		var buf bytes.Buffer
		renderer := NewTerminalRenderer(&buf)

		comp := Text("Hello, Terminal!")
		err := renderer.Render(ctx, comp)
		if err != nil {
			t.Fatalf("Render failed: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "Hello, Terminal!") {
			t.Errorf("Expected output to contain 'Hello, Terminal!', got '%s'", output)
		}
	})

	t.Run("clear and re-render", func(t *testing.T) {
		var buf bytes.Buffer
		renderer := NewTerminalRenderer(&buf)

		comp := Text("Hello, Terminal!")
		err := renderer.Render(ctx, comp)
		if err != nil {
			t.Fatalf("Render failed: %v", err)
		}
	})
}
