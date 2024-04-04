package fileutil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type File struct {
	path     string
	contents []byte
}

func IsDir(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		return false
	}

	return file.IsDir()
}

func ReadDirFileNames(path string) []string {
	path, err := filepath.Abs(path)
	if err != nil {
		return []string{}
	}

	if !IsDir(path) {
		return []string{path}
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return []string{}
	}

	var result []string
	for _, file := range files {
		// TODO: add validation for file extensions, tracetest runnable definitions (?)
		if file.IsDir() {
			result = append(result, ReadDirFileNames(filepath.Join(path, file.Name()))...)
		} else {
			result = append(result, filepath.Join(path, file.Name()))
		}
	}

	return result
}

func Read(filePath string) (File, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return File{}, fmt.Errorf("could not read definition file %s: %w", filePath, err)
	}

	return New(filePath, b), nil
}

func New(path string, b []byte) File {
	file := File{
		contents: b,
		path:     path,
	}

	return file
}

func (f File) Reader() io.Reader {
	return bytes.NewReader(f.contents)
}

var (
	fileTypeRegex = regexp.MustCompile(`(?m:^type:\s*(\w+))`)
)

func (f File) Type() string {
	match := fileTypeRegex.FindSubmatch(f.contents)
	if len(match) < 2 {
		return ""
	}
	return string(match[1])
}

var (
	hasIDRegex      = regexp.MustCompile(`(?m:^\s+id:\s*)`)
	indentSizeRegex = regexp.MustCompile(`(?m:^(\s+)\w+)`)
)

var ErrFileHasID = errors.New("file already has ID")

func (f File) HasID() bool {
	fileID := hasIDRegex.Find(f.contents)
	return fileID != nil
}

func (f File) SetID(id string) (File, error) {
	if f.HasID() {
		return f, ErrFileHasID
	}

	indent := indentSizeRegex.FindSubmatchIndex(f.contents)
	if len(indent) < 4 {
		return f, fmt.Errorf("cannot detect indentation size")
	}

	indentSize := indent[3] - indent[2]
	// indent[2] is the index of the first indentation.
	// we can assume that's the first line within the `specs` block
	// so we can use it as the place to inejct the ID

	var newContents []byte
	newContents = append(newContents, f.contents[0:indent[2]]...)

	newContents = append(newContents, []byte(strings.Repeat(" ", indentSize))...)
	newContents = append(newContents, []byte("id: "+id+"\n")...)

	newContents = append(newContents, f.contents[indent[2]:]...)

	return New(f.path, newContents), nil
}

func (f File) AbsPath() string {
	abs, err := filepath.Abs(f.path)
	if err != nil {
		panic(fmt.Errorf(`cannot get absolute path from "%s": %w`, f.path, err))
	}

	return abs
}

func (f File) AbsDir() string {
	return filepath.Dir(f.AbsPath())
}

func (f File) RelativeFile(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(f.AbsDir(), path)
}

func (f File) Write() (File, error) {
	err := os.WriteFile(f.path, f.contents, 0644)
	if err != nil {
		return f, fmt.Errorf("could not write file %s: %w", f.path, err)
	}

	return Read(f.path)
}

func (f File) Contents() []byte {
	return f.contents
}
