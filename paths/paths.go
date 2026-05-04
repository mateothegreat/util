// Package paths - provides functions for working with paths.
package paths

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Expand expands a path such as ~/workspace/foo/bar to /Users/mateo/workspace/foo/bar.
//
// Arguments:
//   - path: the path to expand
//
// Returns:
//   - The expanded path
//   - An error if there is an error expanding the path
func Expand(path string) (string, error) {
	var err error
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path, fmt.Errorf("failed to get user home directory: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	if filepath.IsAbs(path) {
		return path, nil
	}

	abs, err := filepath.Abs(os.ExpandEnv(path))
	if err != nil {
		return path, fmt.Errorf("failed to get absolute path: %w", err)
	}
	return abs, nil
}

// IsSub checks if a path is a subpath of another path.
//
// Arguments:
//   - path: the path to check
//   - basePath: the base path
//
// Returns:
//   - True if the path is a subpath of the base path, false otherwise
func IsSub(path, basePath string) bool {
	rel, err := filepath.Rel(basePath, path)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
