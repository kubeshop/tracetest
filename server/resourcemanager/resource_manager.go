package resourcemanager

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slices"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/validation"
)

type ResourceSpec interface {
	HasID() bool
	GetID() id.ID
	Validate() error
}

type ResourceList[T ResourceSpec] struct {
	Count int           `json:"count" yamlstream:"count"`
	Items []Resource[T] `json:"items" yamlstream:"items"`
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

func DisableDelete() managerOption {
	return func(c *config) {
		ops := []Operation{}
		for _, op := range availableOperations {
			if op == OperationDelete {
				continue
			}
			ops = append(ops, op)
		}

		c.enabledOperations = ops
	}
}

func WithTracer(tracer trace.Tracer) managerOption {
	return func(c *config) {
		c.tracer = tracer
	}
}

func CanBeAugmented() managerOption {
	return func(c *config) {
		c.enabledOperations = append(c.enabledOperations, augmentedOperations...)
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

	listHandler := m.methodNotAllowed
	if slices.Contains(enabledOps, OperationList) {
		listHandler = m.list
	}
	m.instrumentRoute(subrouter.HandleFunc("", listHandler).Methods(http.MethodGet).Name(fmt.Sprintf("%s.List", m.resourceTypePlural)))

	createHandler := m.methodNotAllowed
	if slices.Contains(enabledOps, OperationCreate) {
		createHandler = m.create
	}
	m.instrumentRoute(subrouter.HandleFunc("", createHandler).Methods(http.MethodPost).Name(fmt.Sprintf("%s.Create", m.resourceTypePlural)))

	upsertHandler := m.methodNotAllowed
	if slices.Contains(enabledOps, OperationCreate) && slices.Contains(enabledOps, OperationUpdate) {
		upsertHandler = m.upsert
	}
	m.instrumentRoute(subrouter.HandleFunc("", upsertHandler).Methods(http.MethodPut).Name(fmt.Sprintf("%s.Upsert", m.resourceTypePlural)))

	updateHandler := m.methodNotAllowed
	if slices.Contains(enabledOps, OperationUpdate) {
		updateHandler = m.update
	}
	m.instrumentRoute(subrouter.HandleFunc("/{id}", updateHandler).Methods(http.MethodPut).Name(fmt.Sprintf("%s.Update", m.resourceTypePlural)))

	getHandler := m.methodNotAllowed
	if slices.Contains(enabledOps, OperationGet) {
		getHandler = m.get
	}
	m.instrumentRoute(subrouter.HandleFunc("/{id}", getHandler).Methods(http.MethodGet).Name(fmt.Sprintf("%s.Get", m.resourceTypePlural)))

	deleteHandler := m.methodNotAllowed
	if slices.Contains(enabledOps, OperationDelete) {
		deleteHandler = m.delete
	}
	m.instrumentRoute(subrouter.HandleFunc("/{id}", deleteHandler).Methods(http.MethodDelete).Name(fmt.Sprintf("%s.Delete", m.resourceTypePlural)))

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
func (m *manager[T]) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	writeError(w, EncoderFromRequest(r), http.StatusMethodNotAllowed, fmt.Errorf("resource %s does not support the action", m.resourceTypeSingular))
}

func (m *manager[T]) create(w http.ResponseWriter, r *http.Request) {
	encoder := EncoderFromRequest(r)

	targetResource := Resource[T]{}
	err := encoder.DecodeRequestBody(&targetResource)
	if err != nil {
		writeError(w, encoder, http.StatusBadRequest, fmt.Errorf("cannot parse body: %w", err))
		return
	}

	// TODO: if resourceType != values.resourceType return error

	m.doCreate(w, r, encoder, targetResource.Spec)
}

func (m *manager[T]) doCreate(w http.ResponseWriter, r *http.Request, encoder Encoder, specs T) {
	if !specs.HasID() {
		specs = m.rh.SetID(specs, m.config.idgen())
	}

	created, err := m.rh.Create(r.Context(), specs)
	if err != nil {
		m.handleResourceHandlerError(w, "creating", err, encoder)
		return
	}

	newResource := Resource[T]{
		Type: m.resourceTypeSingular,
		Spec: created,
	}

	err = encoder.WriteEncodedResponse(w, http.StatusCreated, newResource)
	if err != nil {
		writeError(w, encoder, http.StatusInternalServerError, fmt.Errorf("cannot marshal entity: %w", err))
	}
}

func (m *manager[T]) upsert(w http.ResponseWriter, r *http.Request) {
	encoder := EncoderFromRequest(r)

	targetResource := Resource[T]{}
	err := encoder.DecodeRequestBody(&targetResource)
	if err != nil {
		writeError(w, encoder, http.StatusBadRequest, fmt.Errorf("cannot parse body: %w", err))
		return
	}

	// if there's no ID given, create the resource
	if !targetResource.Spec.HasID() {
		m.doCreate(w, r, encoder, targetResource.Spec)
		return
	}

	_, err = m.rh.Get(r.Context(), targetResource.Spec.GetID())
	if err != nil {
		// if the given ID is not found, create the resource
		if errors.Is(err, sql.ErrNoRows) {
			m.doCreate(w, r, encoder, targetResource.Spec)
			return
		} else {
			// some actual error, return it
			writeError(w, encoder, http.StatusInternalServerError, fmt.Errorf("could not get entity: %w", err))
			return
		}
	}

	// the resurce exists, update it
	m.doUpdate(w, r, encoder, targetResource.Spec)
}

func (m *manager[T]) update(w http.ResponseWriter, r *http.Request) {
	encoder := EncoderFromRequest(r)

	targetResource := Resource[T]{}
	err := encoder.DecodeRequestBody(&targetResource)
	if err != nil {
		writeError(w, encoder, http.StatusBadRequest, fmt.Errorf("cannot parse body: %w", err))
		return
	}

	// TODO: if resourceType != values.resourceType return error

	vars := mux.Vars(r)
	urlID := id.ID(vars["id"])
	if targetResource.Spec.HasID() && targetResource.Spec.GetID() != urlID {
		err := fmt.Errorf(
			"ID '%s' defined in resource spec does not match ID '%s' from URL",
			targetResource.Spec.GetID(),
			urlID,
		)
		writeError(w, encoder, http.StatusBadRequest, err)
		return
	}
	targetResource.Spec = m.rh.SetID(targetResource.Spec, urlID)

	m.doUpdate(w, r, encoder, targetResource.Spec)
}

func (m *manager[T]) doUpdate(w http.ResponseWriter, r *http.Request, encoder Encoder, specs T) {
	updated, err := m.rh.Update(r.Context(), specs)
	if err != nil {
		m.handleResourceHandlerError(w, "updating", err, encoder)
		return
	}

	newResource := Resource[T]{
		Type: m.resourceTypeSingular,
		Spec: updated,
	}

	err = encoder.WriteEncodedResponse(w, http.StatusOK, newResource)
	if err != nil {
		writeError(w, encoder, http.StatusInternalServerError, fmt.Errorf("cannot marshal entity: %w", err))
	}
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

	if take == 0 {
		take = 20
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
	encoder := EncoderFromRequest(r)

	ctx := r.Context()
	take, skip,
		query, sortBy,
		sortDirection, err := paginationParams(r, m.rh.SortingFields())
	if err != nil {
		writeError(w, encoder, http.StatusBadRequest, fmt.Errorf("cannot process request: %s", err.Error()))
		return
	}

	count, err := m.rh.Count(ctx, query)
	if err != nil {
		m.handleResourceHandlerError(w, "listing", err, encoder)
		return
	}

	listFn := m.rh.List
	if isRequestForAugmented(r) && m.rh.ListAugmented != nil {
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
		Items: []Resource[T]{},
	}

	for _, item := range items {
		resource := Resource[T]{
			Type: m.resourceTypeSingular,
			Spec: item,
		}

		resourceList.Items = append(resourceList.Items, resource)
	}

	err = encoder.WriteEncodedResponse(w, http.StatusOK, resourceList)

	if err != nil {
		writeError(w, encoder, http.StatusInternalServerError, fmt.Errorf("cannot marshal entity: %w", err))
	}
}

const HeaderAugmented = "X-Tracetest-Augmented"

func isRequestForAugmented(r *http.Request) bool {
	return r.Header.Get(HeaderAugmented) == "true"
}

func (m *manager[T]) get(w http.ResponseWriter, r *http.Request) {
	encoder := EncoderFromRequest(r)

	vars := mux.Vars(r)
	id := id.ID(vars["id"])

	getterFn := m.rh.Get
	if isRequestForAugmented(r) && m.rh.GetAugmented != nil {
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

	err = encoder.WriteEncodedResponse(w, http.StatusOK, newResource)
	if err != nil {
		writeError(w, encoder, http.StatusInternalServerError, fmt.Errorf("cannot marshal entity: %w", err))
	}
}

func (m *manager[T]) delete(w http.ResponseWriter, r *http.Request) {
	encoder := EncoderFromRequest(r)

	vars := mux.Vars(r)
	id := id.ID(vars["id"])

	err := m.rh.Delete(r.Context(), id)
	if err != nil {
		m.handleResourceHandlerError(w, "deleting", err, encoder)
		return
	}

	encoder.WriteEncodedResponse(w, http.StatusNoContent, nil)
}

func (m *manager[T]) handleResourceHandlerError(w http.ResponseWriter, verb string, err error, encoder Encoder) {
	// 404 - not found
	if errors.Is(err, sql.ErrNoRows) {
		encoder.WriteEncodedResponse(w, http.StatusNotFound, nil)
		return
	}

	if errors.Is(err, validation.ErrValidation) {
		writeError(w, encoder, http.StatusBadRequest, err)
	}

	// 500 - internal server error
	err = fmt.Errorf("error %s resource %s: %w", verb, m.resourceTypeSingular, err)
	writeError(w, encoder, http.StatusInternalServerError, err)
}

func writeError(w http.ResponseWriter, enc Encoder, code int, err error) {
	err = enc.WriteEncodedResponse(w, code, map[string]any{
		"code":  code,
		"error": err.Error(),
	})

	if err != nil {
		// this panic is intentional. Since we have a hardcoded map to encode
		// any errors means there's something very very wrong
		panic(fmt.Errorf("cannot marshal error: %w", err))
	}
}
