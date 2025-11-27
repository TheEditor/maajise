package fsutil

import (
	"os"
	"path/filepath"
)

// FileExists returns true if the path exists and is a regular file (not a directory).
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// DirExists returns true if the path exists and is a directory.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// PathExists returns true if the path exists (file or directory).
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// EnsureDir creates a directory with all parent directories if they don't exist.
// Returns an error if the path exists but is not a directory.
func EnsureDir(path string) error {
	if DirExists(path) {
		return nil
	}

	if PathExists(path) {
		// Path exists but is not a directory
		return os.ErrExist
	}

	return os.MkdirAll(path, 0755)
}

// EnsureParentDir creates all parent directories for the given path if they don't exist.
// Returns an error if a parent path exists but is not a directory.
func EnsureParentDir(path string) error {
	parent := filepath.Dir(path)
	if parent == "." || parent == "" {
		return nil
	}
	return EnsureDir(parent)
}
