// Package files - provides functions for working with files.
package files

import (
	"mime"
	"path/filepath"
	"strings"
)

// SafeName converts the name to one that can safely be used on the filesystem.
//
// Arguments:
//   - name: the name to convert
//
// Returns:
//   - The converted name
func SafeName(name string) string {
	replacer := strings.NewReplacer(".", "_", "/", "_")
	return replacer.Replace(name)
}

// ContentType returns the MIME type for the given file name.
//
// Arguments:
//   - name: the name of the file to get the MIME type for
//
// Returns:
//   - The MIME type for the given file name
func ContentType(name string) string {
	ext := filepath.Ext(name)
	answer := mime.TypeByExtension(ext)
	if answer == "" {
		if ext == ".log" || ext == ".txt" {
			return "text/plain; charset=utf-8"
		}
	}
	return answer
}
