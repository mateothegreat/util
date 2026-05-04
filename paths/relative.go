// Package paths provides utilities for working with file paths.
package paths

import (
	"errors"
	"os"
	"path/filepath"
)

// FindRootByMarker walks up from the current working directory to locate the root of the project,
// identified by a marker file or directory (e.g., "go.mod" or ".git").
//
// Arguments:
// - markerFile: the name of the file or directory that signifies project root (e.g., "go.mod" or ".git")
// Returns:
// - absolute path to the root directory if found, otherwise an error
//
// Example:
//
//	root, err := paths.FindRootByMarker("go.mod")
//	if err != nil {
//		log.Fatal(err) // project root not found, or error getting current working directory
//	}
//	println(root) // prints the absolute path to the project root
func FindRootByMarker(markerFile string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		marker := filepath.Join(dir, markerFile)
		if _, err := os.Stat(marker); err == nil {
			return dir, nil // found the marker
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached the root of filesystem
		}
		dir = parent
	}
	return "", errors.New("project root not found")
}
