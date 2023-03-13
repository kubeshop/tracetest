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
	ResponseContentType() string
}

var encoders = []encoder{
	jsonEncoder,
	yamlEncoder,
}

var defaultEncoder = jsonEncoder

var jsonEncoder = basicEncoder{
	contentType: "application/json",
	marshalFn:   json.Marshal,
	unmarshalFn: json.Unmarshal,
}

var yamlEncoder = basicEncoder{
	contentType: "text/yaml",
	marshalFn:   yaml.Marshal,
	unmarshalFn: yaml.Unmarshal,
}

type basicEncoder struct {
	contentType string
	marshalFn   func(interface{}) ([]byte, error)
	unmarshalFn func([]byte, interface{}) error
}

func (e basicEncoder) ResponseContentType() string {
	return e.contentType
}

func (e basicEncoder) Accepts(contentType string) bool {
	return contentType == e.contentType
}

func (e basicEncoder) Marshal(in interface{}) (out []byte, err error) {
	return e.marshalFn(in)
}
func (e basicEncoder) Unmarshal(in []byte, out interface{}) (err error) {
	return e.unmarshalFn(in, out)
}

var errUnacceptableContentType = errors.New("unacceptable content type")

func encoderFromRequest(r *http.Request) (encoder, error) {
	contentType := r.Header.Get("Content-Type")
	accept := r.Header.Get("Accept")
	for _, enc := range encoders {
		if enc.Accepts(contentType) || enc.Accepts(accept) {
			return enc, nil
		}
	}

	if accept == "" && contentType == "" {
		return defaultEncoder, nil
	}

	return nil, fmt.Errorf("cannot handle content-type %s: %w", contentType, errUnacceptableContentType)
}
