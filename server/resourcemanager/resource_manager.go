package resourcemanager

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/mitchellh/mapstructure"
)

type ResourceSpec interface {
	HasID() bool
	Validate() error
}

type Resource[T ResourceSpec] struct {
	Type string `mapstructure:"type"`
	Spec T      `mapstructure:"spec"`
}

type ResourceHandler[T ResourceSpec] interface {
	SetID(T, id.ID) T
	Create(context.Context, T) (T, error)
	Update(context.Context, T) (T, error)
	Get(context.Context, id.ID) (T, error)
	Delete(context.Context, id.ID) error
}

type manager[T ResourceSpec] struct {
	resourceType string
	handler      ResourceHandler[T]
	idgen        func() id.ID
}

func New[T ResourceSpec](resourceType string, handler ResourceHandler[T], idgenFn func() id.ID) *manager[T] {
	return &manager[T]{
		resourceType: resourceType,
		handler:      handler,
		idgen:        idgenFn,
	}
}

func (m *manager[T]) RegisterRoutes(r *mux.Router) *mux.Router {
	// prefix is /{resourceType | lowercase}/
	subrouter := r.PathPrefix("/" + strings.ToLower(m.resourceType)).Subrouter()

	subrouter.HandleFunc("/", m.create).Methods(http.MethodPost)
	subrouter.HandleFunc("/", m.update).Methods(http.MethodPut)

	subrouter.HandleFunc("/{id}", m.get).Methods(http.MethodGet)
	subrouter.HandleFunc("/{id}", m.delete).Methods(http.MethodDelete)

	return subrouter
}

func (m *manager[T]) create(w http.ResponseWriter, r *http.Request) {
	m.operationWithBody(w, r, http.StatusCreated, "creating", m.handler.Create)
}

func (m *manager[T]) update(w http.ResponseWriter, r *http.Request) {
	m.operationWithBody(w, r, http.StatusOK, "updating", m.handler.Update)
}

func (m *manager[T]) get(w http.ResponseWriter, r *http.Request) {
	encoder, err := encoderFromRequest(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf("cannot process request: %s", err.Error()))
		return
	}
	w.Header().Set("Content-Type", encoder.ResponseContentType())

	vars := mux.Vars(r)
	id := id.ID(vars["id"])

	item, err := m.handler.Get(r.Context(), id)
	if err != nil {
		m.handleResourceHandlerError(w, "getting", err, encoder)
		return
	}

	newResource := Resource[T]{
		Type: m.resourceType,
		Spec: item,
	}

	bytes, err := encodeValues(newResource, encoder)
	if err != nil {
		writeError(w, encoder, http.StatusInternalServerError, fmt.Errorf("cannot marshal entity: %w", err))
		return
	}

	writeResponse(w, 200, string(bytes))
}

func (m *manager[T]) delete(w http.ResponseWriter, r *http.Request) {
	encoder, err := encoderFromRequest(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf("cannot process request: %s", err.Error()))
		return
	}
	w.Header().Set("Content-Type", encoder.ResponseContentType())

	vars := mux.Vars(r)
	id := id.ID(vars["id"])

	err = m.handler.Delete(r.Context(), id)
	if err != nil {
		m.handleResourceHandlerError(w, "deleting", err, encoder)
		return
	}

	w.WriteHeader(204)
}

func (m *manager[T]) handleResourceHandlerError(w http.ResponseWriter, verb string, err error, encoder encoder) {
	// 404 - not found
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 500 - internal server error
	err = fmt.Errorf("error %s resource %s: %w", verb, m.resourceType, err)
	writeError(w, encoder, http.StatusInternalServerError, err)
}

func (m *manager[T]) operationWithBody(w http.ResponseWriter, r *http.Request, statusCode int, operationVerb string, fn func(context.Context, T) (T, error)) {
	encoder, err := encoderFromRequest(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf("cannot process request: %s", err.Error()))
		return
	}
	w.Header().Set("Content-Type", encoder.ResponseContentType())

	values, err := readValues(r, encoder)
	if err != nil {
		writeError(w, encoder, http.StatusBadRequest, fmt.Errorf("cannot parse body: %w", err))
		return
	}

	targetResource := Resource[T]{}
	err = mapstructure.Decode(values, &targetResource)
	if err != nil {
		writeError(w, encoder, http.StatusBadRequest, fmt.Errorf("cannot unmarshal body values: %w", err))
		return
	}

	// TODO: if resourceType != values.resourceType return error

	// TODO: check if this needs to be done per operation
	if !targetResource.Spec.HasID() {
		targetResource.Spec = m.handler.SetID(targetResource.Spec, m.idgen())
	}

	created, err := fn(r.Context(), targetResource.Spec)
	if err != nil {
		m.handleResourceHandlerError(w, operationVerb, err, encoder)
		return
	}

	newResource := Resource[T]{
		Type: m.resourceType,
		Spec: created,
	}

	bytes, err := encodeValues(newResource, encoder)
	if err != nil {
		writeError(w, encoder, http.StatusInternalServerError, fmt.Errorf("cannot marshal entity: %w", err))
		return
	}

	writeResponse(w, statusCode, string(bytes))
}

func writeResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}

func writeError(w http.ResponseWriter, enc encoder, code int, err error) {
	body, err := enc.Marshal(map[string]any{
		"code":  code,
		"error": err.Error(),
	})
	if err != nil {
		// this panic is intentional. Since we have a hardcoded map to encode
		// any errors means there's something very very wrong
		panic(fmt.Errorf("cannot marshal error: %w", err))
	}

	writeResponse(w, code, string(body))
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
