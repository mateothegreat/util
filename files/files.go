// Package files - provides functions for working with files.
package files

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// GetFileSize returns the size of the file at the given file path.
//
// Arguments:
//   - path: the path to the file
//
// Returns:
//   - The size of the file
//   - -1 if there is an error getting the file size
func GetFileSize(path string) int64 {
	file, err := os.Open(path)
	if err != nil {
		return -1
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return -1
	}

	return stat.Size()
}

// FileExists checks if the file exists at the given file path.
//
// Arguments:
//   - filePath: the path to the file
//
// Returns:
//   - True if the file exists, otherwise false
func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}

// WaitForFileExists waits for the file to exist at the given file path.
//
// Arguments:
//   - filePath: the path to the file
//   - timeout: the timeout to wait for the file to exist
//
// Returns:
//   - True if the file exists within the specified timeout, otherwise false
func WaitForFileExists(filePath string, timeout time.Duration) bool {
	// Create a timer for the timeout
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	// Create a ticker for periodically checking the file existence
	checkInterval := 100 * time.Millisecond
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-timer.C:
			// Timeout occurred
			return false
		case <-ticker.C:
			if _, err := os.Stat(filePath); err == nil {
				// File exists
				return true
			} else if !os.IsNotExist(err) {
				// An error other than "not exist", stop waiting
				return false
			}
			// If file does not exist, continue checking
		}
	}
}

// WaitForNoFileHandlers waits for all file handlers to be closed for the given file path.
//
// Arguments:
//   - filePath: the path to the file
//   - timeout: the timeout to wait for the file handlers to be closed
//   - local: whether to use the local lsof command
//
// Returns:
//   - True if all file handlers are closed within the specified timeout, otherwise false
func WaitForNoFileHandlers(filePath string, timeout time.Duration, local bool) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		var cmd *exec.Cmd
		if local {
			cmd = exec.Command("lsof", filePath)
		} else {
			cmd = exec.Command("sh", "-c", fmt.Sprintf("lsof | grep %s", filePath))
		}

		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			// lsof returns an error if no file handlers are found.
			return true
		}

		if out.Len() == 0 {
			// Also check if the output is empty, indicating no open file handlers.
			return true
		}

		time.Sleep(100 * time.Millisecond) // Wait before trying again.
	}

	return false // Timeout reached.
}

// MoveFile moves a file.
//
// Arguments:
//   - src: the source file
//   - dst: the destination file
//
// Returns:
//   - An error if there is an error moving the file
func MoveFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	err = sourceFile.Close()
	if err != nil {
		return err
	}

	err = os.Remove(src)
	if err != nil {
		return err
	}

	return nil
}

// WalkFile walks up the directory tree from the current directory to find the given file.
// It returns the full path to the file if it is found, otherwise it returns an empty string.
//
// Arguments:
//   - filename: the name of the file to search for
//   - levels: the number of levels to walk up the directory tree
//
// Returns:
//   - the full path to the file if it is found, otherwise an empty string
func WalkFile(filename string, levels int) string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for i := 0; i < levels; i++ {
		path := filepath.Join(dir, filename)
		if _, err := os.Stat(path); err == nil {
			return path
		}
		dir = filepath.Dir(dir)
	}

	return ""
}

// FirstFileExists returns the first file which exists or an error if we can't detect if a file that exists.
//
// Arguments:
//   - paths: the paths to check
//
// Returns:
//   - The first file which exists, or an error if we can't detect if a file that exists.
//   - An error if there is an error checking the paths
func FirstFileExists(paths ...string) (string, error) {
	for _, path := range paths {
		exists := FileExists(path)
		if exists {
			return path, nil
		}
	}
	return "", nil
}

// FileIsEmpty checks if a file is empty.
//
// Arguments:
//   - path: the path to check
//
// Returns:
//   - True if the file is empty, false if the file is not empty
//   - An error if there is an error checking the file
func FileIsEmpty(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return true, fmt.Errorf("getting details of file '%s': %w", path, err)
	}
	return fi.Size() == 0, nil
}

// IsEmpty checks if a file is empty.
//
// Arguments:
//   - name: the name of the file to check
//
// Returns:
//   - True if the file is empty, false if the file is not empty
//   - An error if there is an error checking the file
func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// RenameFile renames a file.
//
// Arguments:
//   - src: the source file
//   - dst: the destination file
//
// Returns:
//   - An error if there is an error renaming the file
func RenameFile(src string, dst string) (err error) {
	if src == dst {
		return nil
	}
	err = CopyFile(src, dst)
	if err != nil {
		return fmt.Errorf("failed to copy source file %s to %s: %s", src, dst, err)
	}
	err = os.RemoveAll(src)
	if err != nil {
		return fmt.Errorf("failed to cleanup source file %s: %s", src, err)
	}
	return nil
}

// CopyFileOrDir copies the source file or directory to the given destination.
//
// Arguments:
//   - src: the source file or directory
//   - dst: the destination file or directory
//   - force: whether to force the copy
//
// Returns:
//   - An error if there is an error copying the file or directory
func CopyFileOrDir(src string, dst string, force bool) (err error) {
	fi, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("getting details of file '%s': %w", src, err)
	}
	if fi.IsDir() {
		return CopyDir(src, dst, force)
	}
	return CopyFile(src, dst)
}

// CopyUnlessSymLink copies a file unless it is a symlink.
//
// Arguments:
//   - entry: the entry to copy
//   - srcPath: the source path
//   - dstPath: the destination path
//
// Returns:
//   - An error if there is an error copying the file
func CopyUnlessSymLink(entry os.DirEntry, srcPath, dstPath string) (err error) {
	// Skip symlinks.
	var info os.FileInfo
	info, err = entry.Info()
	if err != nil {
		return
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return
	}

	err = CopyFile(srcPath, dstPath)
	return
}

// CopyFile copies a file.
//
// Arguments:
//   - src: the source file
//   - dst: the destination file
//
// Returns:
//   - An error if there is an error copying the file
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// LoadBytes loads a file.
//
// Arguments:
//   - dir: the directory to load the file from
//   - name: the name of the file to load
//
// Returns:
//   - The bytes of the file
//   - An error if there is an error loading the file
func LoadBytes(dir, name string) ([]byte, error) {
	path := filepath.Join(dir, name) // relative path
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error loading file %s in directory %s, %v", name, dir, err)
	}
	return bytes, nil
}

// DeleteFile deletes a file. This should NOT be used to delete any sensitive information
// because it can easily be recovered. Use DestroyFile to delete sensitive information.
//
// Arguments:
//   - fileName: the name of the file to delete
//
// Returns:
//   - An error if there is an error deleting the file
func DeleteFile(fileName string) (err error) {
	if fileName != "" {
		exists := FileExists(fileName)

		if exists {
			err = os.Remove(fileName)
			if err != nil {
				return fmt.Errorf("Could not remove file due to %s: %w", fileName, err)
			}
		}
	} else {
		return fmt.Errorf("Filename is not valid")
	}
	return nil
}

// DestroyFile will securely delete a file by first overwriting it with random bytes, then deleting it.
//
// Arguments:
//   - filename: the name of the file to delete
//
// Returns:
//   - An error if there is an error deleting the file
func DestroyFile(filename string) error {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("Could not Destroy %s: %w", filename, err)
	}
	size := fileInfo.Size()
	// Overwrite the file with random data. Doing this multiple times is probably more secure
	randomBytes := make([]byte, size)
	// Avoid false positive G404 of gosec module - https://github.com/securego/gosec/issues/291
	/* #nosec */
	_, _ = rand.Read(randomBytes)
	err = os.WriteFile(filename, randomBytes, DefaultDirWritePermissions)
	if err != nil {
		return fmt.Errorf("Unable to overwrite %s with random data: %w", filename, err)
	}
	return DeleteFile(filename)
}

// FilterFileExists filters out files which do not exist.
//
// Arguments:
//   - paths: the paths to filter
//
// Returns:
//   - The paths which exist
func FilterFileExists(paths []string) []string {
	var answer []string
	for _, path := range paths {
		exists := FileExists(path)
		if exists {
			answer = append(answer, path)
		}
	}
	return answer
}

// IgnoreFile returns true if the path matches any of the ignores. The match is the same as filepath.Match.
//
// Arguments:
//   - path: the path to check
//   - ignores: the ignores to check against
//
// Returns:
//   - True if the path matches any of the ignores, false otherwise
//   - An error if there is an error checking the path
func IgnoreFile(path string, ignores []string) (bool, error) {
	for _, ignore := range ignores {
		if matched, err := filepath.Match(ignore, path); err != nil {
			return false, fmt.Errorf("error when matching ignore %s against path %s: %w", ignore, path, err)
		} else if matched {
			return true, nil
		}
	}
	return false, nil
}

// GlobAllFiles performs a glob on the pattern and then processes all the files found.
// if a folder matches the glob its treated as another glob to recurse into the directory.
//
// Arguments:
//   - basedir: the base directory to start the glob from
//   - pattern: the pattern to glob
//   - fn: the function to call for each file
//
// Returns:
//   - An error if there is an error globbing the files
func GlobAllFiles(basedir string, pattern string, fn func(string) error) error {
	names, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to evaluate glob pattern '%s': %w", pattern, err)
	}
	for _, name := range names {
		fullPath := name
		if basedir != "" {
			fullPath = filepath.Join(basedir, name)
		}
		fi, err := os.Stat(fullPath)
		if err != nil {
			return fmt.Errorf("getting details of file '%s': %w", fullPath, err)
		}
		if fi.IsDir() {
			err = GlobAllFiles("", filepath.Join(fullPath, "*"), fn)
			if err != nil {
				return err
			}
		} else {
			err = fn(fullPath)
			if err != nil {
				return fmt.Errorf("failed processing file '%s': %w", fullPath, err)
			}
		}
	}
	return nil
}
