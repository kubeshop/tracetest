package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kubeshop/tracetest/cli/variable"
	tracetestYaml "github.com/kubeshop/tracetest/server/model/yaml"
	"gopkg.in/yaml.v3"
)

type SpecWithID struct {
	ID string `yaml:"id"`
}

type File struct {
	path     string
	contents []byte
	file     tracetestYaml.File
}

func Read(path string) (*File, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read definition file %s: %w", path, err)
	}

	return New(path, b)
}

func ReadRaw(path string) (File, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return File{}, fmt.Errorf("could not read definition file %s: %w", path, err)
	}

	return NewFromRaw(path, b)
}

func New(path string, b []byte) (*File, error) {
	yf, err := tracetestYaml.Decode(b)
	if err != nil {
		return nil, fmt.Errorf("could not parse definition file: %w", err)
	}

	file := &File{
		contents: b,
		file:     yf,
		path:     path,
	}

	return file, nil
}

func NewFromRaw(path string, b []byte) (File, error) {
	var f tracetestYaml.File
	err := yaml.Unmarshal(b, &f)
	if err != nil {
		return File{}, fmt.Errorf("could not parse definition file: %w", err)
	}

	file := File{
		contents: b,
		path:     path,
		file:     f,
	}

	return file, nil
}

func (f *File) Path() string {
	return f.path
}

func (f File) AbsDir() string {
	abs, err := filepath.Abs(f.path)
	if err != nil {
		panic(fmt.Errorf(`cannot get absolute path from "%s": %w`, f.path, err))
	}

	return filepath.Dir(abs)
}

func (f *File) ResolveVariables() (*File, error) {
	variableInjector := variable.NewInjector(variable.WithVariableProvider(
		variable.EnvironmentVariableProvider{},
	))

	err := variableInjector.Inject(&f.file)
	if err != nil {
		return nil, err
	}

	bytes, err := tracetestYaml.Encode(f.file)
	if err != nil {
		return nil, err
	}

	f.contents = bytes

	return f, nil
}

func (f File) Definition() tracetestYaml.File {
	return f.file
}

func (f File) Contents() string {
	return string(f.contents)
}

func (f File) ContentType() string {
	return "text/yaml"
}

var (
	hasIDRegex      = regexp.MustCompile(`(?m:^\s+id:\s*[0-9a-zA-Z\-_]+$)`)
	indentSizeRegex = regexp.MustCompile(`(?m:^(\s+)\w+)`)
)

var ErrFileHasID = errors.New("file already has ID")

func (f File) HasID() bool {
	fileID := hasIDRegex.Find(f.contents)
	return fileID != nil
}

func (f *File) SetID(id string) (*File, error) {
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

	return New(f.path, newContents)
}

func (f File) Write() (*File, error) {
	err := os.WriteFile(f.path, f.contents, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not write file: %w", err)
	}

	return Read(f.path)
}

func (f File) WriteRaw() (File, error) {
	err := os.WriteFile(f.path, f.contents, 0644)
	if err != nil {
		return f, fmt.Errorf("could not write file: %w", err)
	}

	return ReadRaw(f.path)
}

func (f File) SaveChanges(changes string) File {
	f.contents = []byte(changes)

	return f
}
