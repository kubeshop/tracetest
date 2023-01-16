package httpClient

import (
	"context"
	"net/http"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	k6HTTP "go.k6.io/k6/js/modules/k6/http"
)

type HttpClient struct {
	vu          modules.VU
	httpRequest HttpRequestFunc

	options Options
}

type HTTPResponse struct {
	*k6HTTP.Response `js:"-"`
	TraceID          string
}

type (
	HttpRequestFunc func(method string, url goja.Value, args ...goja.Value) (*k6HTTP.Response, error)
	HttpFunc        func(ctx context.Context, url goja.Value, args ...goja.Value) (*k6HTTP.Response, error)
	NewFunc         func(call goja.ConstructorCall) *goja.Object
)

func New(vu modules.VU) *HttpClient {
	r := k6HTTP.New().NewModuleInstance(vu).Exports().Default.(*goja.Object).Get("request")

	var httpRequest HttpRequestFunc
	err := vu.Runtime().ExportTo(r, &httpRequest)

	if err != nil {
		panic(err)
	}

	return &HttpClient{
		vu:          vu,
		httpRequest: httpRequest,
	}
}

func (h *HttpClient) Constructor(call goja.ConstructorCall) *goja.Object {
	rt := h.vu.Runtime()
	options, err := getOptions(h.vu, call.Argument(0))
	if err != nil {
		common.Throw(rt, err)
	}

	h.options = options

	return rt.ToValue(h).ToObject(rt)
}

func (c *HttpClient) Get(url goja.Value, args ...goja.Value) (*HTTPResponse, error) {
	args = append([]goja.Value{goja.Null()}, args...)
	return c.WithTrace(requestToHttpFunc(http.MethodGet, c.httpRequest), url, args...)
}

func (c *HttpClient) Post(url goja.Value, args ...goja.Value) (*HTTPResponse, error) {
	return c.WithTrace(requestToHttpFunc(http.MethodPost, c.httpRequest), url, args...)
}

func (c *HttpClient) Put(url goja.Value, args ...goja.Value) (*HTTPResponse, error) {
	return c.WithTrace(requestToHttpFunc(http.MethodPut, c.httpRequest), url, args...)
}

func (c *HttpClient) Del(url goja.Value, args ...goja.Value) (*HTTPResponse, error) {
	return c.WithTrace(requestToHttpFunc(http.MethodDelete, c.httpRequest), url, args...)
}

func (c *HttpClient) Head(url goja.Value, args ...goja.Value) (*HTTPResponse, error) {
	return c.WithTrace(requestToHttpFunc(http.MethodHead, c.httpRequest), url, args...)
}

func (c *HttpClient) Patch(url goja.Value, args ...goja.Value) (*HTTPResponse, error) {
	return c.WithTrace(requestToHttpFunc(http.MethodPatch, c.httpRequest), url, args...)
}

func (c *HttpClient) Options(url goja.Value, args ...goja.Value) (*HTTPResponse, error) {
	return c.WithTrace(requestToHttpFunc(http.MethodOptions, c.httpRequest), url, args...)
}
