// Package files - provides functions for working with symlinks.
package files

import (
	"fmt"
	"os"
)

// RecreateSymLink removes any old symlink/binary if they exist and creates a new symlink from the given source.
//
// Arguments:
//   - src: the source file
//   - target: the target file
//
// Returns:
//   - An error if there is an error creating the symlink
func RecreateSymLink(src, target string) error {
	exists := FileExists(target)
	if exists {
		err := os.Remove(target)
		if err != nil {
			return fmt.Errorf("failed to remove %s: %w", target, err)
		}
	}
	err := os.Symlink(src, target)
	if err != nil {
		return fmt.Errorf("failed to create symlink from %s to %s: %w", src, target, err)
	}
	return nil
}
