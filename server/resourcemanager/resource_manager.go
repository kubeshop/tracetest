package resourcemanager

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/pkg/id"
)

type ResourceSpec interface {
	HasID() bool
	Validate() error
}

type ResourceList[T ResourceSpec] struct {
	Count int   `json:"count"`
	Items []any `json:"items"`
}

type Resource[T ResourceSpec] struct {
	Type string `json:"type"`
	Spec T      `json:"spec"`
}

type Manager interface {
	EnabledOperations() []Operation
	Handler() any
	RegisterRoutes(*mux.Router) *mux.Router
	Provisioner
}

type manager[T ResourceSpec] struct {
	resourceTypeSingular string
	resourceTypePlural   string
	handler              any
	rh                   resourceHandler[T]
	config               config
}

type config struct {
	enabledOperations []Operation
	idgen             func() id.ID
	tracer            trace.Tracer
}

type managerOption func(*config)

func WithIDGen(fn func() id.ID) managerOption {
	return func(c *config) {
		c.idgen = fn
	}
}

func WithOperations(ops ...Operation) managerOption {
	return func(c *config) {
		c.enabledOperations = ops
	}
}

func WithTracer(tracer trace.Tracer) managerOption {
	return func(c *config) {
		c.tracer = tracer
	}
}

func New[T ResourceSpec](resourceTypeSingular, resourceTypePlural string, handler any, opts ...managerOption) Manager {
	rh := &resourceHandler[T]{}

	cfg := config{
		enabledOperations: availableOperations,
		idgen:             func() id.ID { return id.GenerateID() },
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	err := rh.bindOperations(cfg.enabledOperations, handler)

	if err != nil {
		err := fmt.Errorf(
			"cannot create Resourcemanager '%s': %w",
			resourceTypeSingular,
			err,
		)
		panic(err)
	}

	return &manager[T]{
		resourceTypeSingular: resourceTypeSingular,
		resourceTypePlural:   resourceTypePlural,
		handler:              handler,
		rh:                   *rh,
		config:               cfg,
	}
}

func (m *manager[T]) EnabledOperations() []Operation {
	return m.config.enabledOperations
}

func (m *manager[T]) Handler() any {
	return m.handler
}

func (m *manager[T]) RegisterRoutes(r *mux.Router) *mux.Router {
	// prefix is /{resourceType | lowercase}/
	subrouter := r.PathPrefix("/" + strings.ToLower(m.resourceTypePlural)).Subrouter()

	enabledOps := m.EnabledOperations()

	if slices.Contains(enabledOps, OperationList) {
		m.instrumentRoute(subrouter.HandleFunc("", m.list).Methods(http.MethodGet).Name("list"))
	}

	if slices.Contains(enabledOps, OperationCreate) {
		m.instrumentRoute(subrouter.HandleFunc("", m.create).Methods(http.MethodPost).Name(fmt.Sprintf("%s.Create", m.resourceTypePlural)))
	}

	if slices.Contains(enabledOps, OperationUpdate) {
		m.instrumentRoute(subrouter.HandleFunc("/{id}", m.update).Methods(http.MethodPut).Name(fmt.Sprintf("%s.Update", m.resourceTypePlural)))
	}

	if slices.Contains(enabledOps, OperationGet) {
		m.instrumentRoute(subrouter.HandleFunc("/{id}", m.get).Methods(http.MethodGet).Name(fmt.Sprintf("%s.Get", m.resourceTypePlural)))
	}

	if slices.Contains(enabledOps, OperationDelete) {
		m.instrumentRoute(subrouter.HandleFunc("/{id}", m.delete).Methods(http.MethodDelete).Name(fmt.Sprintf("%s.Delete", m.resourceTypePlural)))
	}

	return subrouter
}

func (m *manager[T]) instrumentRoute(route *mux.Route) {
	originalHandler := route.GetHandler()
	pathTemplate, _ := route.GetPathTemplate()

	newHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.config.tracer == nil {
			originalHandler.ServeHTTP(w, r)
			return
		}

		method := r.Method

		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		ctx, span := m.config.tracer.Start(ctx, fmt.Sprintf("%s %s", method, pathTemplate))
		defer span.End()

		params := make(map[string]interface{}, 0)
		for key, value := range mux.Vars(r) {
			params[key] = value
		}

		paramsJson, _ := json.Marshal(params)

		queryString := make(map[string]interface{}, 0)
		for key, value := range r.URL.Query() {
			queryString[key] = value
		}
		queryStringJson, _ := json.Marshal(queryString)

		headers := make(map[string]interface{}, 0)
		for key, value := range r.Header {
			headers[key] = value
		}
		headersJson, _ := json.Marshal(headers)

		span.SetAttributes(
			attribute.String(string(semconv.HTTPMethodKey), r.Method),
			attribute.String(string(semconv.HTTPRouteKey), pathTemplate),
			attribute.String(string(semconv.HTTPTargetKey), r.URL.String()),
			attribute.String("http.request.params", string(paramsJson)),
			attribute.String("http.request.query", string(queryStringJson)),
			attribute.String("http.request.headers", string(headersJson)),
		)

		originalHandler.ServeHTTP(w, r.WithContext(ctx))
	})

	route.Handler(newHandler)
}

func (m *manager[T]) create(w http.ResponseWriter, r *http.Request) {
	m.operationWithBody(w, r, http.StatusCreated, "creating", m.rh.Create)
}

func (m *manager[T]) update(w http.ResponseWriter, r *http.Request) {
	m.operationWithBody(w, r, http.StatusOK, "updating", m.rh.Update)
}

func getIntFromQuery(r *http.Request, key string) (int, error) {
	str := r.URL.Query().Get(key)
	if str == "" {
		return 0, nil
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("'%s' is not a number", str)
	}

	return val, nil
}

func paginationParams(r *http.Request, sortingFields []string) (take, skip int, query, sortBy, sortDirection string, err error) {
	take, err = getIntFromQuery(r, "take")
	if err != nil {
		err = fmt.Errorf("error reading take param: %w", err)
		return
	}

	skip, err = getIntFromQuery(r, "skip")
	if err != nil {
		err = fmt.Errorf("error reading skip param: %w", err)
		return
	}

	sortBy = r.URL.Query().Get("sortBy")
	if sortBy != "" && !slices.Contains(sortingFields, sortBy) {
		err = fmt.Errorf("invalid sort field: %s", sortBy)
		return
	}

	sortDirection = r.URL.Query().Get("sortDirection")

	query = r.URL.Query().Get("query")

	return
}

func (m *manager[T]) list(w http.ResponseWriter, r *http.Request) {
	encoder, err := encoderFromRequest(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf("cannot process request: %s", err.Error()))
		return
	}
	w.Header().Set("Content-Type", encoder.ResponseContentType())

	ctx := r.Context()
	take, skip,
		query, sortBy,
		sortDirection, err := paginationParams(r, m.rh.SortingFields())
	if err != nil {
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf("cannot process request: %s", err.Error()))
		return
	}

	count, err := m.rh.Count(ctx, query)
	if err != nil {
		m.handleResourceHandlerError(w, "listing", err, encoder)
		return
	}

	listFn := m.rh.List
	if isRequestForAugmented(r) {
		listFn = m.rh.ListAugmented
	}

	items, err := listFn(
		ctx,
		take,
		skip,
		query,
		sortBy,
		sortDirection,
	)
	if err != nil {
		m.handleResourceHandlerError(w, "listing", err, encoder)
		return
	}

	// TODO: the name "count" can be misleading when using pagination.
	//       an user can paginate the request and see a different number
	//       of records inside "item"
	resourceList := ResourceList[T]{
		Count: count,
		Items: []any{},
	}

	for _, item := range items {
		resource := Resource[T]{
			Type: m.resourceTypeSingular,
			Spec: item,
		}

		resourceList.Items = append(resourceList.Items, resource)
	}

	bytes, err := encodeValues(resourceList, encoder)
	if err != nil {
		writeError(w, encoder, http.StatusInternalServerError, fmt.Errorf("cannot marshal entity: %w", err))
		return
	}

	writeResponse(w, http.StatusOK, string(bytes))
}

const HeaderAugmented = "X-Tracetest-Augmented"

func isRequestForAugmented(r *http.Request) bool {
	return r.Header.Get(HeaderAugmented) == "true"
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

	getterFn := m.rh.Get
	if isRequestForAugmented(r) {
		getterFn = m.rh.GetAugmented
	}

	item, err := getterFn(r.Context(), id)
	if err != nil {
		m.handleResourceHandlerError(w, "getting", err, encoder)
		return
	}

	newResource := Resource[T]{
		Type: m.resourceTypeSingular,
		Spec: item,
	}

	bytes, err := encodeValues(newResource, encoder)
	if err != nil {
		writeError(w, encoder, http.StatusInternalServerError, fmt.Errorf("cannot marshal entity: %w", err))
		return
	}

	writeResponse(w, http.StatusOK, string(bytes))
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

	err = m.rh.Delete(r.Context(), id)
	if err != nil {
		m.handleResourceHandlerError(w, "deleting", err, encoder)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (m *manager[T]) handleResourceHandlerError(w http.ResponseWriter, verb string, err error, encoder encoder) {
	// 404 - not found
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 500 - internal server error
	err = fmt.Errorf("error %s resource %s: %w", verb, m.resourceTypeSingular, err)
	writeError(w, encoder, http.StatusInternalServerError, err)
}

func (m *manager[T]) operationWithBody(w http.ResponseWriter, r *http.Request, statusCode int, operationVerb string, fn func(context.Context, T) (T, error)) {
	encoder, err := encoderFromRequest(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf("cannot process request: %s", err.Error()))
		return
	}
	w.Header().Set("Content-Type", encoder.ResponseContentType())

	targetResource := Resource[T]{}
	err = readValues(r, encoder, &targetResource)
	if err != nil {
		writeError(w, encoder, http.StatusBadRequest, fmt.Errorf("cannot parse body: %w", err))
		return
	}

	// TODO: if resourceType != values.resourceType return error

	// TODO: check if this needs to be done per operation
	if !targetResource.Spec.HasID() {
		targetResource.Spec = m.rh.SetID(targetResource.Spec, m.config.idgen())
	}

	created, err := fn(r.Context(), targetResource.Spec)
	if err != nil {
		m.handleResourceHandlerError(w, operationVerb, err, encoder)
		return
	}

	newResource := Resource[T]{
		Type: m.resourceTypeSingular,
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
	return enc.Marshal(resource)
}

func readValues(r *http.Request, enc encoder, target any) error {
	body, err := readBody(r)
	if err != nil {
		return fmt.Errorf("cannot read yaml body: %w", err)
	}

	err = enc.Unmarshal(body, target)
	if err != nil {
		return fmt.Errorf("cannot unmarshal request: %w", err)
	}

	return nil
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

func decode(input any, output any) error {
	yamlContent, err := yaml.Marshal(input)
	if err != nil {
		return fmt.Errorf("couldn't convert marshal input into YAML: %w", err)
	}

	err = yaml.Unmarshal(yamlContent, output)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal YAML into target: %w", err)
	}

	return nil
}
