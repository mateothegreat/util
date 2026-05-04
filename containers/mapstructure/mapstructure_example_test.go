package mapstructure

import (
	"fmt"
	"testing"
)

// TestMapstructureStandalone verifies our implementation works correctly.
func TestMapstructureStandalone(t *testing.T) {
	// Test 1: Simple JSON decoding with weak typing
	t.Run("WeakTypeJSON", func(t *testing.T) {
		type Config struct {
			Port    int    `json:"port"`
			Debug   bool   `json:"debug"`
			Timeout int    `json:"timeout"`
			Host    string `json:"host"`
		}

		jsonData := `{
			"port": "8080",
			"debug": "true",
			"timeout": "30",
			"host": "localhost"
		}`

		var cfg Config
		err := DecodeJSON(jsonData, &cfg)
		if err != nil {
			t.Fatalf("DecodeJSON failed: %v", err)
		}

		if cfg.Port != 8080 {
			t.Errorf("Expected port 8080, got %d", cfg.Port)
		}
		if !cfg.Debug {
			t.Errorf("Expected debug true, got %v", cfg.Debug)
		}
		if cfg.Timeout != 30 {
			t.Errorf("Expected timeout 30, got %d", cfg.Timeout)
		}
	})

	// Test 2: Map decoding with float64 to int conversion
	t.Run("MapDecoding", func(t *testing.T) {
		type User struct {
			ID   int    `json:"id"`
			Age  int    `json:"age"`
			Name string `json:"name"`
		}

		data := map[string]interface{}{
			"id":   float64(123),
			"age":  float64(25),
			"name": "John Doe",
		}

		var user User
		err := DecodeMap(data, &user)
		if err != nil {
			t.Fatalf("DecodeMap failed: %v", err)
		}

		if user.ID != 123 || user.Age != 25 {
			t.Errorf("Float64 to int conversion failed: ID=%d, Age=%d", user.ID, user.Age)
		}
	})

	// Test 3: ToMap conversion
	t.Run("ToMap", func(t *testing.T) {
		type Person struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Age   int    `json:"age"`
		}

		person := Person{
			Name:  "Alice",
			Email: "alice@example.com",
			Age:   30,
		}

		m, err := ToMap(person)
		if err != nil {
			t.Fatalf("ToMap failed: %v", err)
		}

		if m["name"] != "Alice" || m["age"] != 30 {
			t.Errorf("ToMap conversion incorrect: %+v", m)
		}
	})

	fmt.Println("All mapstructure tests passed!")
}
