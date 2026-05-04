// Package files - provides functions for working with files.
package files

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// YAMLFromFile reads a yaml file and unmarshals it into a struct.
//
// Arguments:
//   - path: the path to the file to read
//
// Returns:
//   - The unmarshalled struct
//   - An error if there is an error unmarshalling the file
func YAMLFromFile[T any](path string) (*T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed reading file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed reading file: %w", err)
	}

	var config T

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// JSONFromFile reads a json file and unmarshals it into a struct.
//
// Arguments:
//   - path: the path to the file to read
//
// Returns:
//   - The unmarshalled struct
//   - An error if there is an error unmarshalling the file
func JSONFromFile[T any](path string) (*T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed reading file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed reading file: %w", err)
	}

	var config T

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
