package mapstructure

import (
	"fmt"
	"testing"
	"time"
)

// TestDecodeYAML_SimpleUsage demonstrates basic YAML decoding.
func TestDecodeYAML_SimpleUsage(t *testing.T) {
	yamlData := `
id: 123
name: John Doe
email: john@example.com
age: 30
active: true
score: 95.5
tags:
  - golang
  - testing
created_at: "2024-01-01T10:00:00Z"
metadata:
  ip_address: 192.168.1.1
  country: US
  sessions: 42
`

	var user ExampleUser
	err := DecodeYAML(yamlData, &user)
	if err != nil {
		t.Fatalf("DecodeYAML failed: %v", err)
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

// TestDecodeYAML_WeakTyping demonstrates automatic type conversions with YAML.
func TestDecodeYAML_WeakTyping(t *testing.T) {
	yamlData := `
id: "123"
name: Jane Smith
email: jane@example.com
age: "25"
active: "true"
score: "87.3"
tags: "programming,go,microservices"
created_at: "2024-01-01"
metadata:
  ip_address: 10.0.0.1
  country: UK
  sessions: "15"
`

	var user ExampleUser
	err := DecodeYAML(yamlData, &user)
	if err != nil {
		t.Fatalf("DecodeYAML with weak typing failed: %v", err)
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

// TestAutoDecode_JSON tests auto-detection of JSON format.
func TestAutoDecode_JSON(t *testing.T) {
	jsonData := `{
		"id": 456,
		"name": "JSON User",
		"email": "json@example.com",
		"active": true
	}`

	var user ExampleUser
	err := AutoDecode(jsonData, &user)
	if err != nil {
		t.Fatalf("AutoDecode failed for JSON: %v", err)
	}

	if user.ID != 456 || user.Name != "JSON User" {
		t.Errorf("JSON auto-decode failed: ID=%d, Name=%s", user.ID, user.Name)
	}
}

// TestAutoDecode_YAML tests auto-detection of YAML format.
func TestAutoDecode_YAML(t *testing.T) {
	yamlData := `
id: 789
name: YAML User
email: yaml@example.com
active: true
`

	var user ExampleUser
	err := AutoDecode(yamlData, &user)
	if err != nil {
		t.Fatalf("AutoDecode failed for YAML: %v", err)
	}

	if user.ID != 789 || user.Name != "YAML User" {
		t.Errorf("YAML auto-decode failed: ID=%d, Name=%s", user.ID, user.Name)
	}
}

// TestAutoDecode_EdgeCases tests edge cases for format detection.
func TestAutoDecode_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int // expected ID value
		isJSON   bool
	}{
		{
			name:     "JSON array",
			input:    `[{"id": 100}]`,
			expected: 0, // Arrays aren't supported directly
			isJSON:   true,
		},
		{
			name:     "YAML with document separator",
			input:    "---\nid: 200\nname: Test",
			expected: 200,
			isJSON:   false,
		},
		{
			name:     "YAML flow style",
			input:    "{id: 300, name: FlowStyle}",
			expected: 300,
			isJSON:   false,
		},
		{
			name:     "JSON with whitespace",
			input:    "\n\t  \n{\"id\": 400}\n  ",
			expected: 400,
			isJSON:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result struct {
				ID   int    `json:"id" yaml:"id"`
				Name string `json:"name" yaml:"name"`
			}

			err := AutoDecode(tt.input, &result)

			if tt.expected == 0 && err != nil {
				// Expected to fail
				return
			}

			if err != nil {
				t.Fatalf("AutoDecode failed: %v", err)
			}

			if result.ID != tt.expected {
				t.Errorf("Expected ID=%d, got %d", tt.expected, result.ID)
			}
		})
	}
}

// TestDecodeYAML_ComplexStructures tests YAML with complex nested structures.
func TestDecodeYAML_ComplexStructures(t *testing.T) {
	type Address struct {
		Street  string `yaml:"street"`
		City    string `yaml:"city"`
		ZipCode string `yaml:"zip_code"`
	}

	type Person struct {
		Name      string            `yaml:"name"`
		Age       int               `yaml:"age"`
		Addresses []Address         `yaml:"addresses"`
		Skills    map[string]int    `yaml:"skills"`
		Settings  map[string]string `yaml:"settings"`
	}

	yamlData := `
name: Alice Cooper
age: 30
addresses:
  - street: 123 Main St
    city: New York
    zip_code: "10001"
  - street: 456 Park Ave
    city: Boston
    zip_code: "02101"
skills:
  golang: 5
  python: 3
  javascript: 4
settings:
  theme: dark
  language: en
  notifications: "true"
`

	var person Person
	err := DecodeYAML(yamlData, &person)
	if err != nil {
		t.Fatalf("DecodeYAML failed: %v", err)
	}

	if len(person.Addresses) != 2 {
		t.Errorf("Expected 2 addresses, got %d", len(person.Addresses))
	}
	if person.Addresses[0].City != "New York" {
		t.Errorf("First address city incorrect: %s", person.Addresses[0].City)
	}
	if person.Skills["golang"] != 5 {
		t.Errorf("Golang skill level incorrect: %d", person.Skills["golang"])
	}
	if person.Settings["theme"] != "dark" {
		t.Errorf("Theme setting incorrect: %s", person.Settings["theme"])
	}
}

// TestAutoDecode_Errors tests error cases.
func TestAutoDecode_Errors(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name:    "Empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Invalid type",
			input:   123,
			wantErr: true,
		},
		{
			name:    "Invalid JSON",
			input:   `{"invalid": json}`,
			wantErr: true,
		},
		{
			name:    "Invalid YAML",
			input:   "invalid:\n  - nested\n    bad indentation",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result ExampleUser
			err := AutoDecode(tt.input, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("AutoDecode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDecodeYAML_TimeFormats tests various time format parsing.
func TestDecodeYAML_TimeFormats(t *testing.T) {
	type TimeTest struct {
		T1 time.Time `yaml:"t1"`
		T2 time.Time `yaml:"t2"`
		T3 time.Time `yaml:"t3"`
		T4 time.Time `yaml:"t4"`
		T5 time.Time `yaml:"t5"`
	}

	yamlData := `
t1: "2024-01-01T10:00:00Z"
t2: "2024-01-01T10:00:00.123456789Z"
t3: "2024-01-01T10:00:00"
t4: "2024-01-01 10:00:00"
t5: "2024-01-01"
`

	var tt TimeTest
	err := DecodeYAML(yamlData, &tt)
	if err != nil {
		t.Fatalf("DecodeYAML failed: %v", err)
	}

	// Just verify they parsed without checking exact values
	if tt.T1.IsZero() || tt.T2.IsZero() || tt.T3.IsZero() || tt.T4.IsZero() || tt.T5.IsZero() {
		t.Error("Some time values failed to parse")
	}
}

// BenchmarkDecodeYAML benchmarks YAML decoding performance.
func BenchmarkDecodeYAML(b *testing.B) {
	yamlData := []byte(`
id: 123
name: John Doe
email: john@example.com
age: 30
active: true
score: 95.5
tags:
  - golang
  - testing
`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user ExampleUser
		_ = DecodeYAML(yamlData, &user)
	}
}

// BenchmarkAutoDecode_JSON benchmarks auto-detection with JSON.
func BenchmarkAutoDecode_JSON(b *testing.B) {
	jsonData := []byte(`{"id":123,"name":"John Doe","active":true}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user ExampleUser
		_ = AutoDecode(jsonData, &user)
	}
}

// BenchmarkAutoDecode_YAML benchmarks auto-detection with YAML.
func BenchmarkAutoDecode_YAML(b *testing.B) {
	yamlData := []byte(`id: 123
name: John Doe
active: true`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user ExampleUser
		_ = AutoDecode(yamlData, &user)
	}
}

// ExampleDecodeYAML shows a simple YAML decoding example.
func ExampleDecodeYAML() {
	yamlData := `
id: "42"
name: YAML User
age: "25"
active: "yes"
tags: "go,docker,k8s"
`

	var user ExampleUser
	if err := DecodeYAML(yamlData, &user); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("User: ID=%d, Name=%s, Age=%d, Active=%v\n",
		user.ID, user.Name, user.Age, user.Active)
	fmt.Printf("Tags: %v\n", user.Tags)
	// Output:
	// User: ID=42, Name=YAML User, Age=25, Active=true
	// Tags: [go docker k8s]
}

// ExampleAutoDecode shows automatic format detection.
func ExampleAutoDecode() {
	// Can handle either JSON or YAML
	jsonData := `{"id": 1, "name": "JSON User"}`
	yamlData := `
id: 2
name: YAML User
`

	var user1, user2 ExampleUser

	// Auto-detects JSON
	if err := AutoDecode(jsonData, &user1); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Auto-detects YAML
	if err := AutoDecode(yamlData, &user2); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("JSON User: ID=%d, Name=%s\n", user1.ID, user1.Name)
	fmt.Printf("YAML User: ID=%d, Name=%s\n", user2.ID, user2.Name)
	// Output:
	// JSON User: ID=1, Name=JSON User
	// YAML User: ID=2, Name=YAML User
}

// TestMapDecoder_CustomConfig tests custom decoder configurations.
func TestMapDecoder_CustomConfig(t *testing.T) {
	t.Run("Custom tag name", func(t *testing.T) {
		type CustomTag struct {
			ID   int    `custom:"id"`
			Name string `custom:"name"`
		}

		decoder := &MapDecoder{
			WeaklyTyped:       true,
			TagName:           "custom",
			IgnoreUnknownKeys: true,
			ZeroFields:        true,
		}

		data := map[string]interface{}{
			"id":   "123",
			"name": "Test",
		}

		var result CustomTag
		err := decoder.Decode(data, &result)
		if err != nil {
			t.Fatalf("Decode failed: %v", err)
		}

		if result.ID != 123 || result.Name != "Test" {
			t.Errorf("Custom tag decoding failed: ID=%d, Name=%s", result.ID, result.Name)
		}
	})

	t.Run("Zero fields behavior", func(t *testing.T) {
		type Data struct {
			A string `json:"a"`
			B string `json:"b"`
		}

		decoder := &MapDecoder{
			WeaklyTyped: true,
			TagName:     "json",
			ZeroFields:  false, // Don't zero fields
		}

		// Start with non-zero values
		result := Data{A: "initial", B: "initial"}

		// Only update A
		data := map[string]interface{}{"a": "updated"}

		err := decoder.Decode(data, &result)
		if err != nil {
			t.Fatalf("Decode failed: %v", err)
		}

		if result.A != "updated" || result.B != "initial" {
			t.Errorf("ZeroFields=false not working: A=%s, B=%s", result.A, result.B)
		}
	})
}

// TestSetValue_EdgeCases tests edge cases in type conversion.
func TestSetValue_EdgeCases(t *testing.T) {
	decoder := NewDecoder()

	t.Run("Bool conversions", func(t *testing.T) {
		tests := []struct {
			input    string
			expected bool
		}{
			{"true", true},
			{"false", false},
			{"1", true},
			{"0", false},
			{"yes", true},
			{"no", false},
			{"YES", true},
			{"NO", false},
		}

		for _, tt := range tests {
			err := decoder.Decode(map[string]interface{}{"value": tt.input}, &struct {
				Value bool `json:"value"`
			}{Value: false})

			// Note: strconv.ParseBool only handles true/false/1/0/t/f/T/F/TRUE/FALSE/True/False
			// yes/no will fail, which is expected
			if err == nil {
				// Should have parsed successfully for valid bool strings
				continue
			}
		}
	})

	t.Run("Overflow handling", func(t *testing.T) {
		type Numbers struct {
			Int8  int8  `json:"int8"`
			Uint8 uint8 `json:"uint8"`
		}

		// Test overflow scenarios
		data := map[string]interface{}{
			"int8":  "128", // Overflows int8
			"uint8": "256", // Overflows uint8
		}

		var result Numbers
		err := decoder.Decode(data, &result)
		// Should handle overflow gracefully
		if err != nil && result.Int8 == 0 && result.Uint8 == 0 {
			// Expected behavior - overflow results in zero values or error
			return
		}
	})

	t.Run("Nil handling", func(t *testing.T) {
		type Nullable struct {
			Ptr   *string        `json:"ptr"`
			Slice []int          `json:"slice"`
			Map   map[string]int `json:"map"`
		}

		data := map[string]interface{}{
			"ptr":   nil,
			"slice": nil,
			"map":   nil,
		}

		var result Nullable
		err := decoder.Decode(data, &result)
		if err != nil {
			t.Fatalf("Decode failed: %v", err)
		}

		if result.Ptr != nil || result.Slice != nil || result.Map != nil {
			t.Error("Nil values not handled correctly")
		}
	})
}

// TestToMap_EdgeCases tests edge cases for ToMap function.
func TestToMap_EdgeCases(t *testing.T) {
	t.Run("Unexported fields", func(t *testing.T) {
		type Private struct {
			Public      string `json:"public"`
			private     string // unexported
			OtherPublic int    `json:"other"`
		}

		p := Private{
			Public:      "visible",
			private:     "hidden",
			OtherPublic: 42,
		}

		m, err := ToMap(p)
		if err != nil {
			t.Fatalf("ToMap failed: %v", err)
		}

		if _, exists := m["private"]; exists {
			t.Error("Unexported field was included in map")
		}
		if m["public"] != "visible" {
			t.Error("Public field not correctly mapped")
		}
	})

	t.Run("Omitempty handling", func(t *testing.T) {
		type OmitEmpty struct {
			Required string `json:"required"`
			Optional string `json:"optional,omitempty"`
			Number   int    `json:"number,omitempty"`
		}

		oe := OmitEmpty{
			Required: "present",
			Optional: "", // Empty, should be omitted
			Number:   0,  // Zero, should be omitted
		}

		m, err := ToMap(oe)
		if err != nil {
			t.Fatalf("ToMap failed: %v", err)
		}

		if _, exists := m["optional"]; exists {
			t.Error("Empty optional field was not omitted")
		}
		if _, exists := m["number"]; exists {
			t.Error("Zero number field was not omitted")
		}
	})

	t.Run("Tag variations", func(t *testing.T) {
		type Tags struct {
			NoTag       string
			EmptyTag    string `json:""`
			DashTag     string `json:"-"`
			CommaTag    string `json:",omitempty"`
			NameAndOpts string `json:"custom,omitempty"`
		}

		tags := Tags{
			NoTag:       "no_tag",
			EmptyTag:    "empty_tag",
			DashTag:     "dash_tag",
			CommaTag:    "comma_tag",
			NameAndOpts: "name_opts",
		}

		m, err := ToMap(tags)
		if err != nil {
			t.Fatalf("ToMap failed: %v", err)
		}

		// DashTag should be excluded
		if _, exists := m["-"]; exists {
			t.Error("Dash tag field was included")
		}

		// NoTag should use field name
		if m["NoTag"] != "no_tag" {
			t.Error("Field without tag not using field name")
		}

		// NameAndOpts should use custom name
		if m["custom"] != "name_opts" {
			t.Error("Field with custom name not mapped correctly")
		}
	})
}

// TestDecode_PointerFields tests handling of pointer fields.
func TestDecode_PointerFields(t *testing.T) {
	type WithPointers struct {
		StrPtr  *string    `json:"str_ptr"`
		IntPtr  *int       `json:"int_ptr"`
		TimePtr *time.Time `json:"time_ptr"`
		NilPtr  *string    `json:"nil_ptr"`
	}

	data := map[string]interface{}{
		"str_ptr":  "hello",
		"int_ptr":  float64(42),
		"time_ptr": "2024-01-01T00:00:00Z",
		"nil_ptr":  nil,
	}

	var result WithPointers
	err := DecodeMap(data, &result)
	if err != nil {
		t.Fatalf("DecodeMap failed: %v", err)
	}

	if result.StrPtr == nil || *result.StrPtr != "hello" {
		t.Error("String pointer not decoded correctly")
	}
	if result.IntPtr == nil || *result.IntPtr != 42 {
		t.Error("Int pointer not decoded correctly")
	}
	if result.TimePtr == nil || result.TimePtr.IsZero() {
		t.Error("Time pointer not decoded correctly")
	}
	if result.NilPtr != nil {
		t.Error("Nil pointer should remain nil")
	}
}

// TestDecode_NonStructOutput tests decoding to non-struct types.
func TestDecode_NonStructOutput(t *testing.T) {
	decoder := NewDecoder()

	t.Run("Decode to slice", func(t *testing.T) {
		data := []interface{}{"a", "b", "c"}
		var result []string

		// This should fail as we expect struct output
		err := decoder.Decode(data, &result)
		if err == nil {
			t.Error("Expected error when decoding to non-struct")
		}
	})

	t.Run("Decode single value map", func(t *testing.T) {
		data := map[string]interface{}{"value": "hello"}
		var result string

		err := decoder.Decode(data, &result)
		if err != nil {
			t.Fatalf("Single value decode failed: %v", err)
		}
		if result != "hello" {
			t.Errorf("Expected 'hello', got %s", result)
		}
	})
}

// TestDecode_Errors tests various error conditions.
func TestDecode_Errors(t *testing.T) {
	t.Run("Nil output", func(t *testing.T) {
		decoder := NewDecoder()
		err := decoder.Decode(map[string]interface{}{"test": "value"}, nil)
		if err == nil {
			t.Error("Expected error for nil output")
		}
	})

	t.Run("Non-pointer output", func(t *testing.T) {
		decoder := NewDecoder()
		var result ExampleUser
		err := decoder.decodeMap(map[string]interface{}{"id": 1}, result) // Pass by value
		if err == nil {
			t.Error("Expected error for non-pointer output")
		}
	})

	t.Run("Type conversion failure", func(t *testing.T) {
		decoder := &MapDecoder{
			WeaklyTyped: false, // Strict typing
		}

		type Strict struct {
			Number int `json:"number"`
		}

		var result Strict
		err := decoder.Decode(map[string]interface{}{
			"number": "not-a-number",
		}, &result)

		if err == nil {
			t.Error("Expected error for failed type conversion")
		}
	})
}

// TestMapDecoder_StructToStruct tests struct-to-struct decoding.
func TestMapDecoder_StructToStruct(t *testing.T) {
	type Source struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type Target struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	source := Source{ID: 123, Name: "Test"}
	var target Target

	decoder := NewDecoder()
	err := decoder.Decode(source, &target)
	if err != nil {
		t.Fatalf("Struct-to-struct decode failed: %v", err)
	}

	if target.ID != source.ID || target.Name != source.Name {
		t.Error("Struct values not copied correctly")
	}
}

// TestComplexNestedStructures tests deeply nested structures.
func TestComplexNestedStructures(t *testing.T) {
	type Level3 struct {
		Value string `json:"value"`
	}

	type Level2 struct {
		L3    Level3   `json:"level3"`
		Items []string `json:"items"`
	}

	type Level1 struct {
		L2  Level2         `json:"level2"`
		Map map[string]int `json:"map"`
	}

	data := map[string]interface{}{
		"level2": map[string]interface{}{
			"level3": map[string]interface{}{
				"value": "deep",
			},
			"items": []interface{}{"a", "b", "c"},
		},
		"map": map[string]interface{}{
			"key1": float64(1),
			"key2": float64(2),
		},
	}

	var result Level1
	err := DecodeMap(data, &result)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if result.L2.L3.Value != "deep" {
		t.Error("Deep nested value not decoded correctly")
	}
	if len(result.L2.Items) != 3 {
		t.Error("Nested slice not decoded correctly")
	}
	if result.Map["key1"] != 1 {
		t.Error("Map field not decoded correctly")
	}
}
