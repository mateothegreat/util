package mapstructure

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// ExampleUser demonstrates a typical user struct with various field types.
type ExampleUser struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Active    bool      `json:"active"`
	Score     float64   `json:"score"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	Metadata  Metadata  `json:"metadata"`
}

// Metadata represents nested structure support.
type Metadata struct {
	IP       string `json:"ip_address"`
	Country  string `json:"country"`
	Sessions int    `json:"sessions"`
}

// TestDecodeJSON_SimpleUsage demonstrates the simplest usage pattern.
func TestDecodeJSON_SimpleUsage(t *testing.T) {
	// JSON with perfect type matching - uses native unmarshaling for speed
	jsonData := `{
		"id": 123,
		"name": "John Doe",
		"email": "john@example.com",
		"age": 30,
		"active": true,
		"score": 95.5,
		"tags": ["golang", "testing"],
		"created_at": "2024-01-01T10:00:00Z",
		"metadata": {
			"ip_address": "192.168.1.1",
			"country": "US",
			"sessions": 42
		}
	}`

	var user ExampleUser
	err := DecodeJSON(jsonData, &user)
	if err != nil {
		t.Fatalf("DecodeJSON failed: %v", err)
	}

	// Verify the results
	if user.ID != 123 || user.Name != "John Doe" {
		t.Errorf("Basic fields not decoded correctly: ID=%d, Name=%s", user.ID, user.Name)
	}
	if len(user.Tags) != 2 || user.Tags[0] != "golang" {
		t.Errorf("Array field not decoded correctly: %v", user.Tags)
	}
	if user.Metadata.Sessions != 42 {
		t.Errorf("Nested structure not decoded correctly: sessions=%d", user.Metadata.Sessions)
	}
}

// TestDecodeJSON_WeakTyping demonstrates automatic type conversions.
func TestDecodeJSON_WeakTyping(t *testing.T) {
	// Real-world scenario: config file or form data with string types
	jsonData := `{
		"id": "123",
		"name": "Jane Smith",
		"email": "jane@example.com", 
		"age": "25",
		"active": "true",
		"score": "87.3",
		"tags": "programming,go,microservices",
		"created_at": "2024-01-01",
		"metadata": {
			"ip_address": "10.0.0.1",
			"country": "UK",
			"sessions": "15"
		}
	}`

	var user ExampleUser
	err := DecodeJSON(jsonData, &user)
	if err != nil {
		t.Fatalf("DecodeJSON with weak typing failed: %v", err)
	}

	// All string values should be converted to their proper types
	if user.ID != 123 {
		t.Errorf("String to int conversion failed: ID=%d", user.ID)
	}
	if user.Age != 25 {
		t.Errorf("String to int conversion failed: Age=%d", user.Age)
	}
	if !user.Active {
		t.Errorf("String to bool conversion failed: Active=%v", user.Active)
	}
	if user.Score != 87.3 {
		t.Errorf("String to float conversion failed: Score=%f", user.Score)
	}
	if len(user.Tags) != 3 || user.Tags[1] != "go" {
		t.Errorf("Comma-separated string to slice conversion failed: %v", user.Tags)
	}
	if user.Metadata.Sessions != 15 {
		t.Errorf("Nested string to int conversion failed: sessions=%d", user.Metadata.Sessions)
	}
}

// TestDecodeMap demonstrates decoding from map[string]interface{}.
func TestDecodeMap(t *testing.T) {
	// Common scenario: data from database or API that returns generic maps
	data := map[string]interface{}{
		"id":     float64(456), // JSON numbers are float64
		"name":   "Bob Wilson",
		"email":  "bob@example.com",
		"age":    float64(35),
		"active": true,
		"score":  float64(92.1),
		"tags":   []interface{}{"backend", "api", "rest"},
		"metadata": map[string]interface{}{
			"ip_address": "172.16.0.1",
			"country":    "CA",
			"sessions":   float64(28),
		},
	}

	var user ExampleUser
	err := DecodeMap(data, &user)
	if err != nil {
		t.Fatalf("DecodeMap failed: %v", err)
	}

	// Verify float64 to int conversions (common JSON issue)
	if user.ID != 456 || user.Age != 35 {
		t.Errorf("Float64 to int conversion failed: ID=%d, Age=%d", user.ID, user.Age)
	}
	if user.Metadata.Sessions != 28 {
		t.Errorf("Nested float64 to int conversion failed: sessions=%d", user.Metadata.Sessions)
	}
}

// TestToMap demonstrates struct to map conversion.
func TestToMap(t *testing.T) {
	user := ExampleUser{
		ID:     789,
		Name:   "Alice Cooper",
		Email:  "alice@example.com",
		Age:    28,
		Active: true,
		Score:  98.7,
		Tags:   []string{"frontend", "react", "typescript"},
		Metadata: Metadata{
			IP:       "192.168.0.1",
			Country:  "AU",
			Sessions: 52,
		},
	}

	m, err := ToMap(user)
	if err != nil {
		t.Fatalf("ToMap failed: %v", err)
	}

	// Verify the map contains all fields
	if m["id"] != 789 || m["name"] != "Alice Cooper" {
		t.Errorf("Basic fields not in map: id=%v, name=%v", m["id"], m["name"])
	}

	// Check nested structure
	if metadata, ok := m["metadata"].(Metadata); ok {
		if metadata.Sessions != 52 {
			t.Errorf("Nested structure not preserved: sessions=%d", metadata.Sessions)
		}
	} else {
		t.Errorf("Metadata not found or wrong type in map")
	}
}

// TestCustomDecoder demonstrates using a custom configured decoder.
func TestCustomDecoder(t *testing.T) {
	// Create decoder with custom settings
	decoder := &MapDecoder{
		WeaklyTyped:       false, // Strict typing
		TagName:           "json",
		IgnoreUnknownKeys: false, // Will error on unknown fields
		ZeroFields:        true,
	}

	// This should fail with strict typing
	jsonData := `{"id": "not-a-number", "name": "Test"}`
	var user ExampleUser
	err := decoder.Decode([]byte(jsonData), &user)
	if err == nil {
		t.Error("Expected error with strict typing, but got none")
	}

	// This should fail with unknown keys
	jsonData = `{"id": 123, "unknown_field": "value"}`
	err = decoder.Decode([]byte(jsonData), &user)
	if err == nil {
		t.Error("Expected error with unknown field, but got none")
	}
}

// BenchmarkNativeJSON benchmarks native JSON unmarshaling.
func BenchmarkNativeJSON(b *testing.B) {
	jsonData := []byte(`{
		"id": 123,
		"name": "John Doe",
		"email": "john@example.com",
		"age": 30,
		"active": true,
		"score": 95.5,
		"tags": ["golang", "testing"]
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user ExampleUser
		_ = json.Unmarshal(jsonData, &user)
	}
}

// BenchmarkMapstructureJSON benchmarks our implementation with perfect types.
func BenchmarkMapstructureJSON(b *testing.B) {
	jsonData := []byte(`{
		"id": 123,
		"name": "John Doe",
		"email": "john@example.com",
		"age": 30,
		"active": true,
		"score": 95.5,
		"tags": ["golang", "testing"]
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user ExampleUser
		_ = DecodeJSON(jsonData, &user)
	}
}

// BenchmarkMapstructureWeakTyping benchmarks with type conversions.
func BenchmarkMapstructureWeakTyping(b *testing.B) {
	jsonData := []byte(`{
		"id": "123",
		"name": "John Doe",
		"email": "john@example.com",
		"age": "30",
		"active": "true",
		"score": "95.5",
		"tags": "golang,testing"
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user ExampleUser
		_ = DecodeJSON(jsonData, &user)
	}
}

// ExampleDecodeJSON shows a simple usage example.
func ExampleDecodeJSON() {
	// Config file with mixed types (common in real applications)
	configJSON := `{
		"id": "42",
		"name": "Example User",
		"age": "25",
		"active": "yes",
		"tags": "go,docker,k8s"
	}`

	var user ExampleUser
	if err := DecodeJSON(configJSON, &user); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("User: ID=%d, Name=%s, Age=%d, Active=%v\n",
		user.ID, user.Name, user.Age, user.Active)
	fmt.Printf("Tags: %v\n", user.Tags)
	// Output:
	// User: ID=42, Name=Example User, Age=25, Active=true
	// Tags: [go docker k8s]
}

// ExampleDecodeMap shows decoding from a map.
func ExampleDecodeMap() {
	// Data from database query or API response
	userData := map[string]interface{}{
		"id":     float64(100), // Numbers often come as float64
		"name":   "Map User",
		"age":    float64(30),
		"active": true,
	}

	var user ExampleUser
	if err := DecodeMap(userData, &user); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("User: ID=%d, Name=%s\n", user.ID, user.Name)
	// Output:
	// User: ID=100, Name=Map User
}

// ExampleToMap shows converting struct to map.
func ExampleToMap() {
	user := ExampleUser{
		ID:     1,
		Name:   "Test User",
		Email:  "test@example.com",
		Active: true,
	}

	m, err := ToMap(user)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Easy to manipulate or serialize
	fmt.Printf("Map keys: ")
	for k := range m {
		fmt.Printf("%s ", k)
	}
	// Output:
	// Map keys: id name email age active score tags created_at metadata
}
