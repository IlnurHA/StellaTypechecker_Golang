package utils

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil // Check if it is a directory
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // Path does not exist
	}
	return false, err // Other error (e.g., permission denied)
}

func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil // Check if it is a directory
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // Path does not exist
	}
	return false, err // Other error (e.g., permission denied)
}

func getFiles(dirPath string) ([]string, error) {
	testPaths := make([]string, 0)
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			testPaths = append(testPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking the path %v: %v\n", dirPath, err)
	}
	return testPaths, nil
}

func GetTestPaths(dirPath string) ([]string, error) {
	exists, err := DirExists(dirPath)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("Directory not found")
	}
	return getFiles(dirPath)
}

func GetProgramFromStdin() (string, error) {
	stdin, err := io.ReadAll(os.Stdin)
	return string(stdin), err
}
