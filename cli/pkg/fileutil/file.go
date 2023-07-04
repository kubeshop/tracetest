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

type file struct {
	path     string
	contents []byte
}

func Read(filePath string) (file, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return file{}, fmt.Errorf("could not read definition file %s: %w", filePath, err)
	}

	return New(filePath, b), nil
}

func New(path string, b []byte) file {
	file := file{
		contents: b,
		path:     path,
	}

	return file
}

func (f file) Reader() io.Reader {
	return bytes.NewReader(f.contents)
}

var (
	hasIDRegex      = regexp.MustCompile(`(?m:^\s+id:\s*)`)
	indentSizeRegex = regexp.MustCompile(`(?m:^(\s+)\w+)`)
)

var ErrFileHasID = errors.New("file already has ID")

func (f file) HasID() bool {
	fileID := hasIDRegex.Find(f.contents)
	return fileID != nil
}

func (f file) SetID(id string) (file, error) {
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

func (f file) AbsDir() string {
	abs, err := filepath.Abs(f.path)
	if err != nil {
		panic(fmt.Errorf(`cannot get absolute path from "%s": %w`, f.path, err))
	}

	return filepath.Dir(abs)
}

func (f file) Write() (file, error) {
	err := os.WriteFile(f.path, f.contents, 0644)
	if err != nil {
		return f, fmt.Errorf("could not write file %s: %w", f.path, err)
	}

	return Read(f.path)
}

func (f file) ReadAll() (string, error) {
	return string(f.contents), nil
}
