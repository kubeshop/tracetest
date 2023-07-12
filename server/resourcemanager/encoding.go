package resourcemanager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"

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
	yamlStreamEncoder{},
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

type yamlStreamEncoder struct{}

func (e yamlStreamEncoder) ContentType() string {
	return "text/yaml-stream"
}

func (e yamlStreamEncoder) Accepts(contentType string) bool {
	return contentType == e.ContentType()
}

func (e yamlStreamEncoder) Marshal(in interface{}) (out []byte, err error) {
	targetField, err := getYamlStreamField(in, "items", reflect.Slice)
	if errors.Is(err, errNotAStruct) {
		// if the target is not a struct, marshal it as a single document
		return yaml.Marshal(in)
	}
	if err != nil {
		return nil, err
	}

	// iterate over each element in the slice
	for i := 0; i < targetField.Len(); i++ {
		// get the element
		elem := targetField.Index(i)

		// marshal the element
		elemBytes, err := yaml.Marshal(elem.Interface())
		if err != nil {
			return nil, fmt.Errorf("cannot marshal yaml: %w", err)
		}

		// append the document separator
		out = append(out, []byte("---\n")...)

		// append the marshaled element
		out = append(out, elemBytes...)
	}

	return out, nil
}

func (e yamlStreamEncoder) Unmarshal(in []byte, out interface{}) (err error) {
	targetField, err := getYamlStreamField(out, "items", reflect.Slice)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(bytes.NewReader(in), yaml.Strict())
	// iterate over each document in the yaml stream
	for {
		// we need to create a new instance of the slice element for each document
		// 1. get the type of the slice element
		elemType := targetField.Type().Elem()
		// 2. create a new instance of the slice element
		elem := reflect.New(elemType).Elem()

		// decode the yaml into the slice element. it needs to be a pointer to an interface{}
		err := decoder.Decode(elem.Addr().Interface())
		if errors.Is(err, io.EOF) {
			// no more documents in the stream
			break
		}

		if err != nil {
			// the current document is invalid. return the error
			return fmt.Errorf("cannot unmarshal yaml: %w", err)
		}

		// append the slice element to the target slice
		targetField.Set(reflect.Append(targetField, elem))
	}

	// if there's an error, ignore the count.
	countField, _ := getYamlStreamField(out, "count", reflect.Int)
	if countField.IsValid() {
		countField.SetInt(int64(targetField.Len()))
	}

	return nil
}

var errNotAStruct = errors.New("target is not a struct")

// getYamlStreamField returns the field in `target` with the name as the value of its yamlstream tag.
// it returns an error if the field is not found or if the field is not of the specified kind.
func getYamlStreamField(target interface{}, name string, kind reflect.Kind) (reflect.Value, error) {
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return reflect.Value{}, errNotAStruct
	}

	t := v.Type()
	var yamlStreamField reflect.Value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("yamlstream") == name {
			yamlStreamField = v.Field(i)
			break
		}
	}

	if !yamlStreamField.IsValid() {
		return reflect.Value{}, fmt.Errorf("no field defined as yamlstream %s found", name)
	}

	if yamlStreamField.Kind() != kind {
		return reflect.Value{}, fmt.Errorf("field defined as yamlstream %s is not of kind %s", name, kind)
	}

	return yamlStreamField, nil
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

	if len(body) == 0 {
		return fmt.Errorf("request body is empty")
	}

	return e.input.Unmarshal(body, out)
}

func (e Encoder) WriteEncodedResponse(w http.ResponseWriter, code int, data any) error {
	encoded, err := e.output.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot encode response data: %w", err)
	}

	w.Header().Set("Content-Type", e.output.ContentType())
	w.WriteHeader(code)
	if data == nil {
		return nil
	}

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
