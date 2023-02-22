package resourcemanager

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

type validator interface {
	Validate() error
}

type Resource[T validator] struct {
	Type string `mapstructure:"type"`
	Spec T      `mapstructure:"spec"`
}

type ResourceHandler[T validator] interface {
	Create(T) (T, error)
}

type manager[T validator] struct {
	resourceType string
	handler      ResourceHandler[T]
}

func New[T validator](resourceType string, handler ResourceHandler[T]) *manager[T] {
	return &manager[T]{
		resourceType: resourceType,
		handler:      handler,
	}
}

func (m *manager[T]) RegisterRoutes(r *mux.Router) *mux.Router {
	// prefix is /{resourceType | lowercase}/
	subrouter := r.PathPrefix("/" + strings.ToLower(m.resourceType)).Subrouter()

	subrouter.HandleFunc("/", m.create).Methods(http.MethodPost)

	return subrouter
}

func (m *manager[T]) create(w http.ResponseWriter, r *http.Request) {
	encoder, err := encoderFromRequest(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf("cannot process request: %s", err.Error()))
		return
	}

	values, err := readValues(r, encoder)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf("cannot parse body: %s", err.Error()))
		return
	}

	targetResource := Resource[T]{}
	err = mapstructure.Decode(values, &targetResource)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf("cannot unmarshal body values: %s", err.Error()))
		return
	}

	// TODO: if resoruceType != values.resourceType return error

	created, err := m.handler.Create(targetResource.Spec)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("cannot persist entity: %s", err.Error()))
		return
	}

	newResource := Resource[T]{
		Type: m.resourceType,
		Spec: created,
	}

	bytes, err := encodeValues(newResource, encoder)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("cannot marshal entity: %s", err.Error()))
		return
	}

	writeResponse(w, http.StatusCreated, string(bytes))

}

func writeResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func encodeValues(resource any, enc encoder) ([]byte, error) {
	// mapstructure doesn't have a `Decode`, but encoding with reversed provides this func.
	// See https://github.com/mitchellh/mapstructure/issues/53#issuecomment-273342420
	var values map[string]any
	err := mapstructure.Decode(resource, &values)
	if err != nil {
		return nil, fmt.Errorf("cannot encode values: %w", err)
	}

	return enc.Marshal(values)
}

func readValues(r *http.Request, enc encoder) (map[string]any, error) {
	body, err := readBody(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read yaml body: %w", err)
	}

	var out map[string]any
	err = enc.Unmarshal(body, &out)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal request: %w", err)
	}

	return out, err
}

func readBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("cannot read nil request body")
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read request body: %w", err)
	}

	return body, nil
}
