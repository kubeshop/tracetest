package resourcemanager

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/goccy/go-yaml"
)

type encoder interface {
	Marshal(in interface{}) (out []byte, err error)
	Unmarshal(in []byte, out interface{}) (err error)
	Accepts(contentType string) bool
	ContentType() string
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

func (e basicEncoder) ContentType() string {
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

type Encoder struct {
	req    *http.Request
	input  encoder
	output encoder
}

func (e Encoder) DecodeRequestBody(out interface{}) (err error) {
	body, err := io.ReadAll(e.req.Body)
	defer e.req.Body.Close()
	if err != nil {
		return fmt.Errorf("cannot read request body: %w", err)
	}

	return e.input.Unmarshal(body, out)
}

func (e Encoder) WriteEncodedResponse(w http.ResponseWriter, data any) error {
	encoded, err := e.output.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot encode response data: %w", err)
	}

	w.Header().Set("Content-Type", e.output.ContentType())
	_, err = w.Write(encoded)

	return err
}

func (e Encoder) RequestContentType() string {
	return e.input.ContentType()
}

func (e Encoder) ResponseContentType() string {
	return e.output.ContentType()
}

func EncoderFromRequest(r *http.Request) Encoder {
	return Encoder{
		req:    r,
		input:  getInputEncoder(r),
		output: getOutputEncoder(r),
	}
}

func getInputEncoder(r *http.Request) encoder {
	return getEncoderForHeader(r, "Content-Type", "Accept")
}

func getOutputEncoder(r *http.Request) encoder {
	return getEncoderForHeader(r, "Accept", "Content-Type")
}

func getEncoderForHeader(r *http.Request, headerWaterfall ...string) encoder {
	for _, header := range headerWaterfall {
		headerValue := r.Header.Get(header)
		if headerValue == "" {
			continue
		}

		for _, enc := range encoders {
			if enc.Accepts(headerValue) {
				return enc
			}
		}
	}

	return defaultEncoder
}
