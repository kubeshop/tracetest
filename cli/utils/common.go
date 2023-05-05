package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func StringToIOReader(s string) io.Reader {
	return strings.NewReader(s)
}

func IOReadCloserToString(r io.ReadCloser) string {
	b, _ := io.ReadAll(r)
	return string(b)
}

func StringReferencesFile(path string) bool {
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

func Capitalize(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
