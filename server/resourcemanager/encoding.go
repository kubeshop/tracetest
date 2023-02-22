package resourcemanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

type encoder interface {
	Marshal(in interface{}) (out []byte, err error)
	Unmarshal(in []byte, out interface{}) (err error)
	Accepts(contentType string) bool
}

var encoders = []encoder{
	jsonEncoder{},
	yamlEncoder{},
}

type yamlEncoder struct{}

func (e yamlEncoder) Accepts(contentType string) bool {
	return contentType == "text/yaml"
}
func (e yamlEncoder) Marshal(in interface{}) (out []byte, err error) {
	return yaml.Marshal(in)
}
func (e yamlEncoder) Unmarshal(in []byte, out interface{}) (err error) {
	return yaml.Unmarshal(in, out)
}

type jsonEncoder struct{}

func (e jsonEncoder) Accepts(contentType string) bool {
	return contentType == "application/json"
}
func (e jsonEncoder) Marshal(in interface{}) (out []byte, err error) {
	return json.Marshal(in)
}
func (e jsonEncoder) Unmarshal(in []byte, out interface{}) (err error) {
	return json.Unmarshal(in, out)
}

var errUnacceptableContentType = errors.New("unacceptable content type")

func encoderFromRequest(r *http.Request) (encoder, error) {
	contentType := r.Header.Get("Content-Type")
	for _, enc := range encoders {
		if enc.Accepts(contentType) {
			return enc, nil
		}
	}

	return nil, fmt.Errorf("cannot handle content-type %s: %w", contentType, errUnacceptableContentType)
}
