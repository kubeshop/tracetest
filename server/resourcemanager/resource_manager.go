package resourcemanager

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
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
	values, err := readValues(r)
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

	bytes, err := encodeValues(newResource, r)
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

func encodeValues(resource any, r *http.Request) ([]byte, error) {

	// mapstructure doesn't have a `Decode`, but encoding with reversed provides this func.
	// See https://github.com/mitchellh/mapstructure/issues/53#issuecomment-273342420
	var values map[string]any
	err := mapstructure.Decode(resource, &values)
	if err != nil {
		return nil, fmt.Errorf("cannot encode values: %w", err)
	}

	var (
		bytes []byte
	)

	switch r.Header.Get("Content-Type") {
	case "text/yaml":
		bytes, err = yaml.Marshal(values)
	case "application/json":
		bytes, err = json.Marshal(values)
	}

	return bytes, err
}

func readValues(r *http.Request) (map[string]any, error) {
	var (
		values map[string]any
		err    error
	)

	switch r.Header.Get("Content-Type") {
	case "text/yaml":
		values, err = readYaml(r)
	case "application/json":
		values, err = readJSON(r)
	}

	return values, err
}

func readYaml(r *http.Request) (map[string]any, error) {
	body, err := readBody(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read yaml body: %w", err)
	}

	out := make(map[string]any)
	err = yaml.Unmarshal(body, &out)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal yaml: %w", err)
	}

	return out, nil
}

func readJSON(r *http.Request) (map[string]any, error) {
	body, err := readBody(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read json body: %w", err)
	}

	out := make(map[string]any)
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal json: %w", err)
	}

	return out, nil
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

// // -xPOST http://localhost/api/[resource][/[id]]
// // -xDELETE http://localhost/api/[resource][/[id]] delete
// func (rm manager) CreateResource(res Resource) {
// 	rm.resourceHandler[res.Type].Create(res)
// }

// type controller struct{}

// //POST /api/tests/
// func (m *manager[T]) Create(openapiJSON string) {
// 	parsedOpenAPI := opeanpi.parse(json)
// 	resource, _ := DecodeResource(parsedOpenAPI)
// 	rm.CreateResource(resource)
// }

// //POST /api/tests/definition.yaml
// func (m *manager[T]) CreateFromYaml(yamlFile string) {
// 	parsedYaml := map[string]any{}
// 	_ = yaml.Unmarshal([]byte(yamlFile), &parsedYaml)
// 	resource, _ := DecodeResource(parsedYaml)
// 	rm.CreateResource(resource)
// }

// // package config

// // type: Config
// // name: TracetestConfig
// // id:  qewoifhesdaigseh
// // spec:
// //   analyticsEnabled: true

// type ConfigResource struct {
// 	AnalyticsEnabled bool `mapstructure:"analyticsEnabled`
// }

// type configResourceHandler struct {
// 	repo model.Repository
// }

// func (crh configResourceHandler) Create(res Resource) {
// 	env := res.Spec.(model.Environment)
// 	crh.repo.CreateEnvironment(context.Context, env)(Environment, error)
// }

// type environmentRepository struct {
// 	db sql.DB
// }

// // bun
