package mapstructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// MapDecoder provides flexible decoding of maps and JSON data into Go structures.
// It combines the performance of native Go unmarshaling with the flexibility of weak type conversions.
//
// Arguments:
// - None
//
// Returns:
// - A configured MapDecoder instance ready for decoding operations.
type MapDecoder struct {
	// WeaklyTyped allows string representations to be converted to their appropriate types.
	WeaklyTyped bool

	// TagName specifies the struct tag name to use for field mapping.
	TagName string

	// IgnoreUnknownKeys ignores keys that don't match any struct fields.
	IgnoreUnknownKeys bool

	// ZeroFields zeroes fields before decoding to ensure clean state.
	ZeroFields bool
}

// NewDecoder creates a new MapDecoder with sensible defaults for maximum ease of use.
//
// Arguments:
// - None
//
// Returns:
// - *MapDecoder: A decoder configured with weak typing enabled and using "json" tags by default.
func NewDecoder() *MapDecoder {
	return &MapDecoder{
		WeaklyTyped:       true,
		TagName:           "json",
		IgnoreUnknownKeys: true,
		ZeroFields:        true,
	}
}

// Decode decodes input data into the output structure with automatic type conversions.
// It first attempts native JSON unmarshaling for performance, then falls back to flexible decoding.
//
// Arguments:
// - input: The input data ([]byte, string, map[string]interface{}, or struct)
// - output: Pointer to the structure to decode into
//
// Returns:
// - error: nil if successful, error describing what went wrong otherwise
func (d *MapDecoder) Decode(input interface{}, output interface{}) error {
	// Fast path: try native JSON unmarshaling first
	switch v := input.(type) {
	case []byte:
		// Check if it might be YAML
		trimmed := bytes.TrimSpace(v)
		if len(trimmed) > 0 && trimmed[0] != '{' && trimmed[0] != '[' {
			// Might be YAML, try YAML parsing
			var m map[string]interface{}
			if err := yaml.Unmarshal(v, &m); err == nil {
				return d.decodeMap(m, output)
			}
		}

		// Try JSON
		if err := json.Unmarshal(v, output); err == nil {
			return nil
		}
		// Fall back to flexible JSON decoding
		var m map[string]interface{}
		if err := json.Unmarshal(v, &m); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
		return d.decodeMap(m, output)

	case string:
		return d.Decode([]byte(v), output)

	case map[string]interface{}:
		// Direct map decoding with type conversions
		return d.decodeMap(v, output)

	default:
		// Handle struct-to-struct copying or other types
		return d.decodeValue(reflect.ValueOf(input), reflect.ValueOf(output))
	}
}

// DecodeJSON is a convenience method that decodes JSON data with weak typing support.
//
// Arguments:
// - jsonData: JSON data as []byte or string
// - output: Pointer to structure to decode into
//
// Returns:
// - error: nil if successful, error otherwise
func DecodeJSON(jsonData interface{}, output interface{}) error {
	return NewDecoder().Decode(jsonData, output)
}

// DecodeYAML is a convenience method that decodes YAML data with weak typing support.
//
// Arguments:
// - yamlData: YAML data as []byte or string
// - output: Pointer to structure to decode into
//
// Returns:
// - error: nil if successful, error otherwise
func DecodeYAML(yamlData interface{}, output interface{}) error {
	decoder := NewDecoder()
	decoder.TagName = "yaml" // Use yaml tags by default for YAML

	var data []byte
	switch v := yamlData.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("yamlData must be []byte or string, got %T", yamlData)
	}

	// Try native YAML unmarshaling first
	if err := yaml.Unmarshal(data, output); err == nil {
		return nil
	}

	// Fall back to flexible decoding
	var m map[string]interface{}
	if err := yaml.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	return decoder.decodeMap(m, output)
}

// AutoDecode automatically detects whether the input data is JSON or YAML and decodes accordingly.
// It uses weak typing by default to handle string-to-type conversions common in config files.
//
// Arguments:
// - data: Input data as []byte or string (JSON or YAML format)
// - output: Pointer to structure to decode into
//
// Returns:
// - error: nil if successful, error describing what went wrong otherwise
func AutoDecode(data interface{}, output interface{}) error {
	var bytes []byte
	switch v := data.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("data must be []byte or string, got %T", data)
	}

	// Trim whitespace
	bytes = []byte(strings.TrimSpace(string(bytes)))
	if len(bytes) == 0 {
		return fmt.Errorf("empty input data")
	}

	// Detect format based on first non-whitespace character
	firstChar := bytes[0]

	// JSON typically starts with { or [
	if firstChar == '{' || firstChar == '[' {
		// Try JSON first
		err := DecodeJSON(bytes, output)
		if err == nil {
			return nil
		}

		// If JSON failed, could be YAML flow style, try YAML
		yamlErr := DecodeYAML(bytes, output)
		if yamlErr == nil {
			return nil
		}

		// Return JSON error as it's more likely to be JSON if it starts with { or [
		return err
	}

	// Try YAML for everything else (YAML is more permissive)
	// This includes YAML documents that start with ---, or key: value pairs
	return DecodeYAML(bytes, output)
}

// DecodeMap is a convenience method that decodes a map with weak typing support.
//
// Arguments:
// - m: map[string]interface{} to decode
// - output: Pointer to structure to decode into
//
// Returns:
// - error: nil if successful, error otherwise
func DecodeMap(m map[string]interface{}, output interface{}) error {
	return NewDecoder().Decode(m, output)
}

// ToMap converts a struct to map[string]interface{} for easy manipulation.
//
// Arguments:
// - input: The struct to convert
//
// Returns:
// - map[string]interface{}: The resulting map
// - error: nil if successful, error otherwise
func ToMap(input interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	decoder := NewDecoder()

	rv := reflect.ValueOf(input)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct, got %s", rv.Kind())
	}

	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get field name from tag or use field name
		name := field.Name
		if tag := field.Tag.Get(decoder.TagName); tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" && parts[0] != "-" {
				name = parts[0]
			} else if parts[0] == "-" {
				continue
			}
		}

		// Handle zero values based on omitempty
		if value.IsZero() && strings.Contains(field.Tag.Get(decoder.TagName), "omitempty") {
			continue
		}

		result[name] = value.Interface()
	}

	return result, nil
}

// decodeMap performs the actual map to struct decoding with type conversions.
func (d *MapDecoder) decodeMap(m map[string]interface{}, output interface{}) error {
	rv := reflect.ValueOf(output)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("output must be a non-nil pointer")
	}

	rv = rv.Elem()
	rt := rv.Type()

	if rv.Kind() != reflect.Struct {
		// If output is not a struct, try direct assignment
		if len(m) == 1 {
			for _, v := range m {
				return d.setValue(rv, v)
			}
		}
		return fmt.Errorf("output must be a struct for map decoding, got %s", rv.Kind())
	}

	// Create a map of field names to field indices for fast lookup
	fieldMap := make(map[string]int)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}

		// Primary name from tag
		if tag := field.Tag.Get(d.TagName); tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" && parts[0] != "-" {
				fieldMap[parts[0]] = i
			}
		}

		// Also map the field name itself (case-insensitive fallback)
		fieldMap[strings.ToLower(field.Name)] = i
	}

	// Clear fields if requested
	if d.ZeroFields {
		rv.Set(reflect.Zero(rt))
	}

	// Track unknown keys
	var unknownKeys []string

	// Process each key-value pair
	for key, value := range m {
		fieldIndex, found := fieldMap[key]
		if !found {
			// Try case-insensitive match
			fieldIndex, found = fieldMap[strings.ToLower(key)]
		}

		if !found {
			unknownKeys = append(unknownKeys, key)
			continue
		}

		fieldValue := rv.Field(fieldIndex)
		if !fieldValue.CanSet() {
			continue
		}

		if err := d.setValue(fieldValue, value); err != nil {
			return fmt.Errorf("error setting field %s: %w", rt.Field(fieldIndex).Name, err)
		}
	}

	// Check for unknown keys if not ignoring them
	if !d.IgnoreUnknownKeys && len(unknownKeys) > 0 {
		return fmt.Errorf("unknown field: %s", unknownKeys[0])
	}

	return nil
}

// setValue sets a reflect.Value with automatic type conversion.
func (d *MapDecoder) setValue(target reflect.Value, source interface{}) error {
	if source == nil {
		target.Set(reflect.Zero(target.Type()))
		return nil
	}

	sourceValue := reflect.ValueOf(source)
	targetType := target.Type()

	// Direct assignment if types match
	if sourceValue.Type().AssignableTo(targetType) {
		target.Set(sourceValue)
		return nil
	}

	// Handle conversions
	if d.WeaklyTyped {
		// String to various types
		if sourceValue.Kind() == reflect.String {
			str := sourceValue.String()

			switch targetType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if i, err := strconv.ParseInt(str, 10, targetType.Bits()); err == nil {
					target.SetInt(i)
					return nil
				}

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if i, err := strconv.ParseUint(str, 10, targetType.Bits()); err == nil {
					target.SetUint(i)
					return nil
				}

			case reflect.Float32, reflect.Float64:
				if f, err := strconv.ParseFloat(str, targetType.Bits()); err == nil {
					target.SetFloat(f)
					return nil
				}

			case reflect.Bool:
				// Handle more bool variations
				switch strings.ToLower(str) {
				case "true", "1", "t", "yes", "y", "on":
					target.SetBool(true)
					return nil
				case "false", "0", "f", "no", "n", "off":
					target.SetBool(false)
					return nil
				default:
					if b, err := strconv.ParseBool(str); err == nil {
						target.SetBool(b)
						return nil
					}
				}

			case reflect.Slice:
				// Handle comma-separated strings to string slice
				if targetType.Elem().Kind() == reflect.String {
					parts := strings.Split(str, ",")
					slice := reflect.MakeSlice(targetType, len(parts), len(parts))
					for i, part := range parts {
						slice.Index(i).SetString(strings.TrimSpace(part))
					}
					target.Set(slice)
					return nil
				}
			}

			// Handle time.Time
			if targetType == reflect.TypeOf(time.Time{}) {
				formats := []string{
					time.RFC3339,
					time.RFC3339Nano,
					"2006-01-02T15:04:05",
					"2006-01-02 15:04:05",
					"2006-01-02",
				}
				for _, format := range formats {
					if t, err := time.Parse(format, str); err == nil {
						target.Set(reflect.ValueOf(t))
						return nil
					}
				}
			}
		}

		// Number to various types (handle JSON's float64 default)
		if sourceValue.Kind() == reflect.Float64 {
			f := sourceValue.Float()

			switch targetType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				target.SetInt(int64(f))
				return nil

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				target.SetUint(uint64(f))
				return nil

			case reflect.Float32:
				target.SetFloat(f)
				return nil

			case reflect.String:
				target.SetString(strconv.FormatFloat(f, 'f', -1, 64))
				return nil
			}
		}
	}

	// Handle slices and arrays
	if targetType.Kind() == reflect.Slice && sourceValue.Kind() == reflect.Slice {
		slice := reflect.MakeSlice(targetType, sourceValue.Len(), sourceValue.Len())
		for i := 0; i < sourceValue.Len(); i++ {
			if err := d.setValue(slice.Index(i), sourceValue.Index(i).Interface()); err != nil {
				return err
			}
		}
		target.Set(slice)
		return nil
	}

	// Handle maps
	if targetType.Kind() == reflect.Map && sourceValue.Kind() == reflect.Map {
		targetKeyType := targetType.Key()
		targetValueType := targetType.Elem()
		newMap := reflect.MakeMap(targetType)

		iter := sourceValue.MapRange()
		for iter.Next() {
			// Convert key
			newKey := reflect.New(targetKeyType).Elem()
			if err := d.setValue(newKey, iter.Key().Interface()); err != nil {
				return fmt.Errorf("error converting map key: %w", err)
			}

			// Convert value
			newValue := reflect.New(targetValueType).Elem()
			if err := d.setValue(newValue, iter.Value().Interface()); err != nil {
				return fmt.Errorf("error converting map value: %w", err)
			}

			newMap.SetMapIndex(newKey, newValue)
		}
		target.Set(newMap)
		return nil
	}

	// Handle nested maps to structs
	if targetType.Kind() == reflect.Struct {
		if m, ok := source.(map[string]interface{}); ok {
			return d.decodeMap(m, target.Addr().Interface())
		}
	}

	// Handle pointers
	if targetType.Kind() == reflect.Ptr {
		if sourceValue.Kind() == reflect.Ptr {
			if sourceValue.IsNil() {
				target.Set(reflect.Zero(targetType))
				return nil
			}
			sourceValue = sourceValue.Elem()
		}

		ptr := reflect.New(targetType.Elem())
		if err := d.setValue(ptr.Elem(), sourceValue.Interface()); err != nil {
			return err
		}
		target.Set(ptr)
		return nil
	}

	// Last resort: try direct conversion
	if sourceValue.Type().ConvertibleTo(targetType) {
		target.Set(sourceValue.Convert(targetType))
		return nil
	}

	return fmt.Errorf("cannot convert %v (type %s) to %s", source, sourceValue.Type(), targetType)
}

// decodeValue handles struct-to-struct and other value decoding.
func (d *MapDecoder) decodeValue(source, target reflect.Value) error {
	// Dereference pointers
	if source.Kind() == reflect.Ptr {
		if source.IsNil() {
			target.Set(reflect.Zero(target.Type()))
			return nil
		}
		source = source.Elem()
	}

	if target.Kind() == reflect.Ptr {
		if target.IsNil() {
			target.Set(reflect.New(target.Type().Elem()))
		}
		target = target.Elem()
	}

	// Convert to map first, then decode
	if source.Kind() == reflect.Struct {
		m, err := ToMap(source.Interface())
		if err != nil {
			return err
		}
		return d.decodeMap(m, target.Addr().Interface())
	}

	return d.setValue(target, source.Interface())
}
