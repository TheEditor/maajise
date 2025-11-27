package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsutil-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Non-existent path
	nonExistentPath := filepath.Join(tmpDir, "non-existent.txt")
	if FileExists(nonExistentPath) {
		t.Error("FileExists should return false for non-existent file")
	}

	// Create a file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// File exists
	if !FileExists(testFile) {
		t.Error("FileExists should return true for existing file")
	}

	// Directory should return false (not a file)
	if FileExists(tmpDir) {
		t.Error("FileExists should return false for directories")
	}
}

func TestDirExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsutil-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Non-existent path
	nonExistentPath := filepath.Join(tmpDir, "non-existent-dir")
	if DirExists(nonExistentPath) {
		t.Error("DirExists should return false for non-existent directory")
	}

	// Directory exists
	if !DirExists(tmpDir) {
		t.Error("DirExists should return true for existing directory")
	}

	// File should return false (not a directory)
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	if DirExists(testFile) {
		t.Error("DirExists should return false for files")
	}
}

func TestPathExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsutil-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Non-existent path
	nonExistentPath := filepath.Join(tmpDir, "non-existent")
	if PathExists(nonExistentPath) {
		t.Error("PathExists should return false for non-existent path")
	}

	// Directory exists
	if !PathExists(tmpDir) {
		t.Error("PathExists should return true for existing directory")
	}

	// File exists
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	if !PathExists(testFile) {
		t.Error("PathExists should return true for existing file")
	}
}

func TestEnsureDir_NewDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsutil-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	newDir := filepath.Join(tmpDir, "new", "nested", "dir")
	if err := EnsureDir(newDir); err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}

	if !DirExists(newDir) {
		t.Error("EnsureDir should create the directory")
	}
}

func TestEnsureDir_ExistingDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsutil-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Should succeed when directory already exists
	if err := EnsureDir(tmpDir); err != nil {
		t.Fatalf("EnsureDir should not error for existing directory: %v", err)
	}

	if !DirExists(tmpDir) {
		t.Error("Directory should still exist")
	}
}

func TestEnsureDir_FileExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsutil-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file with the same name
	filePath := filepath.Join(tmpDir, "file-not-dir")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Should error when path exists as a file
	if err := EnsureDir(filePath); err == nil {
		t.Error("EnsureDir should error when path exists as a file")
	}
}

func TestEnsureParentDir_NewParents(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsutil-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "new", "nested", "file.txt")
	if err := EnsureParentDir(filePath); err != nil {
		t.Fatalf("EnsureParentDir failed: %v", err)
	}

	parentDir := filepath.Dir(filePath)
	if !DirExists(parentDir) {
		t.Error("EnsureParentDir should create parent directories")
	}
}

func TestEnsureParentDir_ExistingParents(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsutil-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "file.txt")
	if err := EnsureParentDir(filePath); err != nil {
		t.Fatalf("EnsureParentDir should not error for existing parents: %v", err)
	}

	if !DirExists(tmpDir) {
		t.Error("Parent directory should exist")
	}
}

func TestEnsureParentDir_RootPath(t *testing.T) {
	// Should handle root paths gracefully
	if err := EnsureParentDir("file.txt"); err != nil {
		t.Fatalf("EnsureParentDir should handle relative paths: %v", err)
	}
}

func TestEnsureParentDir_AbsolutePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsutil-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "file.txt")
	if err := EnsureParentDir(filePath); err != nil {
		t.Fatalf("EnsureParentDir failed: %v", err)
	}

	if !DirExists(tmpDir) {
		t.Error("Parent directory should exist")
	}
}

func TestFileExists_VsDirectory_Distinction(t *testing.T) {
	// This test specifically validates the bug fix:
	// FileExists should return false for directories
	tmpDir, err := os.MkdirTemp("", "fsutil-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create both a file and directory
	testFile := filepath.Join(tmpDir, "file.txt")
	testDir := filepath.Join(tmpDir, "testdir")

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Verify the distinction
	if !FileExists(testFile) {
		t.Error("FileExists should return true for files")
	}
	if FileExists(testDir) {
		t.Error("FileExists should return false for directories (bug fix)")
	}
	if !DirExists(testDir) {
		t.Error("DirExists should return true for directories")
	}
	if DirExists(testFile) {
		t.Error("DirExists should return false for files")
	}
}

func TestPathExists_Generic(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fsutil-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create both a file and directory
	testFile := filepath.Join(tmpDir, "file.txt")
	testDir := filepath.Join(tmpDir, "testdir")

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatal(err)
	}

	// PathExists should return true for both
	if !PathExists(testFile) {
		t.Error("PathExists should return true for files")
	}
	if !PathExists(testDir) {
		t.Error("PathExists should return true for directories")
	}

	// Should return false for non-existent paths
	if PathExists(filepath.Join(tmpDir, "non-existent")) {
		t.Error("PathExists should return false for non-existent paths")
	}
}
