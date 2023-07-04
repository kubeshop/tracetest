package fileutil

import (
	"os"
	"path/filepath"
)

func ToAbsDir(filePath string) (string, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	return filepath.Dir(absPath), nil
}

func RelativeTo(path, relativeTo string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(relativeTo, path)
}

func IsFilePath(path string) bool {
	// for the current working dir, check if the file exists
	// by finding its absolute path and executing a stat command

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	info, err := os.Stat(absolutePath)
	if err != nil {
		return false
	}

	// if the string is empty the absolute path will the entire dir
	// otherwise the user also could send a directory by mistake
	return info != nil && !info.IsDir()
}
