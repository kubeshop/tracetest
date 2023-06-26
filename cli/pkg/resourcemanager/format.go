package resourcemanager

import (
	"fmt"
	"net/http"
)

type Format interface {
	BuildRequest(req *http.Request, verb Verb) error
	Format(data string) (string, error)
	String() string
}

type formatRegistry []Format

func (f formatRegistry) Get(format string) (Format, error) {
	for _, fr := range f {
		if fr.String() == format {
			return fr, nil
		}
	}

	return nil, fmt.Errorf("format '%s' not supported", format)
}

var Formats = formatRegistry{
	jsonFormat{},
	yamlFormat{},
	prettyFormat{},
}

type jsonFormat struct{}

func (j jsonFormat) String() string {
	return "json"
}

func (j jsonFormat) BuildRequest(req *http.Request, _ Verb) error {
	req.Header.Set("Content-Type", "application/json")
	return nil
}

func (j jsonFormat) Format(data string) (string, error) {
	return data, nil
}

type yamlFormat struct{}

func (y yamlFormat) String() string {
	return "yaml"
}

func (y yamlFormat) BuildRequest(req *http.Request, verb Verb) error {
	req.Header.Set("Content-Type", "application/x-yaml")
	if verb == VerbList {
		req.Header.Set("Accept", "text/yaml-stream")
	}
	return nil
}

func (y yamlFormat) Format(data string) (string, error) {
	return data, nil
}

type prettyFormat struct{}

func (p prettyFormat) String() string {
	return "pretty"
}

func (p prettyFormat) BuildRequest(req *http.Request, _ Verb) error {
	req.Header.Set("Content-Type", "application/json")
	return nil
}

func (p prettyFormat) Format(data string) (string, error) {
	return data, nil
}
