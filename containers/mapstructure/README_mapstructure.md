# MapStructure - High-Performance Flexible Map Decoding

A Go library that combines the best of native Go unmarshaling performance with the flexibility of weak type conversions. Perfect for handling configuration files, API responses, and other scenarios where data types might vary. Now with YAML support and automatic format detection!

## 1. Features
- **Performance First**: Uses native `json.Unmarshal` and `yaml.Unmarshal` when types match perfectly
- **YAML Support**: Full YAML decoding with the same flexibility as JSON
- **Auto-Detection**: Automatically detects JSON vs YAML format
- **Weak Type Conversions**: Automatically converts strings to numbers, booleans, etc.
- **Simple API**: One-line decoding with sensible defaults
- **Flexible Configuration**: Customize behavior as needed
- **Nested Structure Support**: Handles complex nested data
- **Time Parsing**: Multiple date format support out of the box

<!-- TOC -->
- [MapStructure - High-Performance Flexible Map Decoding](#mapstructure---high-performance-flexible-map-decoding)
  - [1. Features](#1-features)
  - [2. Quick Start](#2-quick-start)
  - [3. When to Use This vs Native Go](#3-when-to-use-this-vs-native-go)
    - [3.1. Use Native Go When:](#31-use-native-go-when)
    - [3.2. Use This Library When:](#32-use-this-library-when)
  - [4. Performance](#4-performance)
  - [5. Examples](#5-examples)
    - [5.1. Configuration File Handling](#51-configuration-file-handling)
    - [5.2. YAML Support](#52-yaml-support)
    - [5.3. Automatic Format Detection](#53-automatic-format-detection)
    - [5.4. Database/API Response Handling](#54-databaseapi-response-handling)
    - [5.5. Advanced Usage](#55-advanced-usage)
  - [6. Type Conversion Rules](#6-type-conversion-rules)
  - [7. Comparison with mitchellh/mapstructure](#7-comparison-with-mitchellhmapstructure)
  - [8. License](#8-license)
<!-- /TOC -->

## 2. Quick Start

```go
import "github.com/mateothegreat/go-util/containers"

// Decode JSON with mixed types (common in config files)
jsonData := `{"id": "123", "active": "true", "tags": "go,docker,k8s"}`
var config MyConfig
err := containers.DecodeJSON(jsonData, &config)

// Decode YAML with automatic conversions
yamlData := `
id: "123"
active: "true"
tags: "go,docker,k8s"
`
err := containers.DecodeYAML(yamlData, &config)

// Auto-detect format (JSON or YAML)
err := containers.AutoDecode(configData, &config)

// Decode from map[string]interface{} (common with databases/APIs)
data := map[string]interface{}{"id": float64(123), "name": "John"}
err := containers.DecodeMap(data, &user)

// Convert struct to map for manipulation
m, err := containers.ToMap(myStruct)
```

## 3. When to Use This vs Native Go

### 3.1. Use Native Go When:
- Working with well-defined APIs where types are guaranteed
- Performance is critical and types always match
- Simple JSON/struct unmarshaling without type variations

### 3.2. Use This Library When:
- Handling configuration files (YAML, JSON, TOML) with string values
- Processing form data or query parameters 
- Working with APIs that return inconsistent types
- Dealing with databases that return generic maps
- Need automatic type conversions (string "123" → int 123)
- Need to auto-detect whether data is JSON or YAML

## 4. Performance

The library optimizes for the common case:

1. **Fast Path**: When decoding JSON/YAML with matching types, it uses native unmarshalers directly
2. **Flexible Path**: Only uses reflection-based conversion when types don't match

Benchmarks show:
- Perfect type match: Same performance as native JSON/YAML
- Type conversions: 2-3x slower than native (but native would fail)

## 5. Examples

### 5.1. Configuration File Handling

```go
type Config struct {
    Port     int      `json:"port" yaml:"port"`
    Debug    bool     `json:"debug" yaml:"debug"`
    Timeout  int      `json:"timeout" yaml:"timeout"`
    Hosts    []string `json:"hosts" yaml:"hosts"`
}

// Config file with all string values
configJSON := `{
    "port": "8080",
    "debug": "true", 
    "timeout": "30",
    "hosts": "api1.example.com,api2.example.com"
}`

var cfg Config
// This just works - all conversions handled automatically
err := containers.DecodeJSON(configJSON, &cfg)
// cfg.Port = 8080 (int)
// cfg.Debug = true (bool)
// cfg.Hosts = ["api1.example.com", "api2.example.com"] ([]string)
```

### 5.2. YAML Support

```go
// YAML configuration with mixed types
yamlConfig := `
server:
  port: "8080"
  host: localhost
  ssl: "true"
  
database:
  connections: "10"
  timeout: "30s"
  hosts: "db1.local,db2.local"
`

type ServerConfig struct {
    Server struct {
        Port int      `yaml:"port"`
        Host string   `yaml:"host"`
        SSL  bool     `yaml:"ssl"`
    } `yaml:"server"`
    Database struct {
        Connections int      `yaml:"connections"`
        Timeout     string   `yaml:"timeout"`
        Hosts       []string `yaml:"hosts"`
    } `yaml:"database"`
}

var cfg ServerConfig
err := containers.DecodeYAML(yamlConfig, &cfg)
// All string values are automatically converted to appropriate types
```

### 5.3. Automatic Format Detection

```go
// The AutoDecode function detects format automatically
func LoadConfig(data []byte) (*Config, error) {
    var cfg Config
    
    // Works with both JSON and YAML
    err := containers.AutoDecode(data, &cfg)
    if err != nil {
        return nil, err
    }
    
    return &cfg, nil
}

// Usage:
jsonData := `{"name": "app", "version": "1.0"}`
yamlData := `
name: app
version: "1.0"
`

// Both work with the same function
cfg1, _ := LoadConfig([]byte(jsonData))  // Detects JSON
cfg2, _ := LoadConfig([]byte(yamlData))  // Detects YAML
```

### 5.4. Database/API Response Handling

```go
// JSON numbers are float64 by default
apiResponse := map[string]interface{}{
    "user_id": float64(12345),
    "age": float64(25),
    "score": float64(98.5),
}

type User struct {
    UserID int     `json:"user_id"`
    Age    int     `json:"age"`
    Score  float64 `json:"score"`
}

var user User
err := containers.DecodeMap(apiResponse, &user)
// float64 values automatically converted to int where needed
```

### 5.5. Advanced Usage

```go
// Create a custom decoder
decoder := &containers.MapDecoder{
    WeaklyTyped:       true,  // Enable type conversions
    TagName:           "json", // Use json tags (or "yaml" for YAML)
    IgnoreUnknownKeys: true,  // Don't error on extra fields
    ZeroFields:        true,  // Clear struct before decoding
}

// Use it for decoding
err := decoder.Decode(data, &output)

// Convert struct to map with custom decoder
m, err := containers.ToMap(myStruct)
```

## 6. Type Conversion Rules

When `WeaklyTyped` is enabled (default), the following conversions are supported:

| From → To | int | uint | float | bool | string | []string | time.Time |
| --------- | --- | ---- | ----- | ---- | ------ | -------- | --------- |
| string    | ✓   | ✓    | ✓     | ✓    | ✓      | ✓*       | ✓**       |
| float64   | ✓   | ✓    | ✓     | ✗    | ✓      | ✗        | ✗         |
| int       | ✓   | ✓    | ✓     | ✗    | ✓      | ✗        | ✗         |
| bool      | ✗   | ✗    | ✗     | ✓    | ✓      | ✗        | ✗         |

\* Comma-separated string → []string (e.g., "a,b,c" → ["a","b","c"])  
\** Multiple time formats supported (RFC3339, RFC3339Nano, "2006-01-02", etc.)

## 7. Comparison with mitchellh/mapstructure

This library provides similar functionality but with key differences:

| Feature          | This Library                  | mitchellh/mapstructure     |
| ---------------- | ----------------------------- | -------------------------- |
| JSON Performance | Native speed when types match | Always uses reflection     |
| YAML Support     | Built-in with auto-detection  | Requires separate handling |
| API Simplicity   | One-line functions            | Requires decoder setup     |
| Default Behavior | Weak typing enabled           | Strict by default          |
| Binary Size      | Minimal dependencies          | Larger footprint           |

## 8. License

This library inherits the Apache 2.0 license from the CloudWeGo project dependencies. 