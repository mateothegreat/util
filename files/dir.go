package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultDirWritePermissions default permissions when creating a directory
	DefaultDirWritePermissions = 0766
	// DefaultFileWritePermissions default permissions when creating a file
	DefaultFileWritePermissions = 0644
	MaximumNewDirectoryAttempts = 1000
)

// FileExists checks if path exists and is a file
// func FileExists(path string) (bool, error) {
// 	fileInfo, err := os.Stat(path)
// 	if err == nil {
// 		return !fileInfo.IsDir(), nil
// 	}
// 	if os.IsNotExist(err) {
// 		return false, nil
// 	}
// 	return false, fmt.Errorf("failed to check if file exists %s: %w", path, err)
// }

// DirExists checks if path exists and is a directory.
//
// Arguments:
//   - path: the path to check
//
// Returns:
//   - True if the path exists and is a directory,false if the path does not exist or is not a directory
//   - An error if there is an error checking the path
func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CreateUniqueDirectory creates a new directory but if the combination of dir and name exists
// then append a number until a unique name is found.
//
// Arguments:
//   - dir: the directory to create the unique directory in
//   - name: the name of the directory to create
//   - maximumAttempts: the maximum number of attempts to create a unique directory
//
// Returns:
//   - The full path to the unique directory
//   - An error if there is an error creating the unique directory
func CreateUniqueDirectory(dir string, name string, maximumAttempts int) (string, error) {
	for i := 0; i < maximumAttempts; i++ {
		n := name
		if i > 0 {
			n += strconv.Itoa(i)
		}
		p := filepath.Join(dir, n)
		exists := FileExists(p)
		if !exists {
			err := os.MkdirAll(p, DefaultDirWritePermissions)
			if err != nil {
				return "", fmt.Errorf("Failed to create directory %s due to %s", p, err)
			}
			return p, nil
		}
	}
	return "", fmt.Errorf("Could not create a unique file in %s starting with %s after %d attempts", dir, name, maximumAttempts)
}

// RenameDir renames a directory.
//
// Arguments:
//   - src: the source directory
//   - dst: the destination directory
//   - force: whether to force the rename
//
// Returns:
//   - An error if there is an error renaming the directory
func RenameDir(src string, dst string, force bool) (err error) {
	err = CopyDir(src, dst, force)
	if err != nil {
		return fmt.Errorf("failed to copy source dir %s to %s: %w", src, dst, err)
	}
	err = os.RemoveAll(src)
	if err != nil {
		return fmt.Errorf("failed to cleanup source dir %s: %w", src, err)
	}
	return nil
}

// CopyDir copies a directory.
//
// Arguments:
//   - src: the source directory
//   - dst: the destination directory
//   - force: whether to force the copy
//
// Returns:
//   - An error if there is an error copying the directory
func CopyDir(src string, dst string, force bool) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		if force {
			os.RemoveAll(dst)
		} else {
			return os.ErrExist
		}
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath, force)
			if err != nil {
				return
			}
		} else {
			err = CopyUnlessSymLink(entry, srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

// CopyDirPreserve copies from the src dir to the dst dir if the file does NOT already exist in dst.
//
// Arguments:
//   - src: the source directory
//   - dst: the destination directory
//
// Returns:
//   - An error if there is an error copying the directory
func CopyDirPreserve(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("checking %s exists: %w", src, err)
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("checking %s exists: %w", dst, err)
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return fmt.Errorf("creating %s: %w", dst, err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("reading files in %s: %w", src, err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDirPreserve(srcPath, dstPath)
			if err != nil {
				return fmt.Errorf("recursively copying %s: %w", entry.Name(), err)
			}
		} else {
			// Skip symlinks.
			info, err := entry.Info()
			if err != nil {
				return err
			}
			if info.Mode()&os.ModeSymlink != 0 {
				continue
			}
			if _, err := os.Stat(dstPath); os.IsNotExist(err) {
				err = CopyFile(srcPath, dstPath)
				if err != nil {
					return fmt.Errorf("copying %s to %s: %w", srcPath, dstPath, err)
				}
			} else if err != nil {
				return fmt.Errorf("checking if %s exists: %w", dstPath, err)
			}
		}
	}
	return nil
}

// CopyDirOverwrite copies from the source dir to the destination dir overwriting files along the way.
//
// Arguments:
//   - src: the source directory
//   - dst: the destination directory
//
// Returns:
//   - An error if there is an error copying the directory
func CopyDirOverwrite(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDirOverwrite(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			err = CopyUnlessSymLink(entry, srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}
	return
}

// DeleteDirContents removes all the contents of the given directory.
//
// Arguments:
//   - dir: the directory to delete the contents of
//
// Returns:
//   - An error if there is an error deleting the directory contents
func DeleteDirContents(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		// lets ignore the top level dir
		if dir != file {
			err = os.RemoveAll(file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteDirContentsExcept removes all the contents of the given directory except for the given directory.
//
// Arguments:
//   - dir: the directory to delete the contents of
//   - exceptDir: the directory to exclude from the deletion
//
// Returns:
//   - An error if there is an error deleting the directory contents
func DeleteDirContentsExcept(dir string, exceptDir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		// lets ignore the top level dir
		if dir != file && !strings.HasSuffix(file, exceptDir) {
			err = os.RemoveAll(file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// RecreateDirs recreates the given directories by deleting and recreating each recursively.
//
// Arguments:
//   - dirs: the directories to recreate
//
// Returns:
//   - An error if there is an error recreating the directories
func RecreateDirs(dirs ...string) error {
	for _, dir := range dirs {
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
		err = os.MkdirAll(dir, DefaultDirWritePermissions)
		if err != nil {
			return err
		}

	}
	return nil
}

// ListDirectory logs the directory at path.
//
// Arguments:
//   - root: the root directory to list
//   - recurse: whether to recurse into subdirectories
//
// Returns:
//   - An error if there is an error listing the directory
func ListDirectory(root string, recurse bool) ([]string, error) {
	var result []string
	if info, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("unable to list %s as does not exist: %w", root, err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", root)
		}
		return nil, fmt.Errorf("stat %s: %w", root, err)
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		dir, _ := filepath.Split(path)
		if !recurse && dir != root {
			// No recursion and we aren't in the root dir
			return nil
		}
		info, err = os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat %s: %w", path, err)
		}
		result = append(result, fmt.Sprintf("%v %d %s %s", info.Mode().String(), info.Size(), info.ModTime().Format(time.RFC822), info.Name()))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
