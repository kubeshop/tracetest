package httpClient

import (
	"fmt"
	"time"

	"github.com/dop251/goja"
	"github.com/xoscar/xk6-tracetest-tracing/models"
	"github.com/xoscar/xk6-tracetest-tracing/utils"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/metrics"
)

func (c *HttpClient) WithTrace(fn HttpFunc, url goja.Value, args ...goja.Value) (*HTTPResponse, error) {
	state := c.vu.State()
	if state == nil {
		return nil, fmt.Errorf("HTTP requests can only be made in the VU context")
	}

	traceID, _, err := (&models.TraceID{
		Prefix:            models.K6Prefix,
		Code:              models.K6Code_Cloud,
		UnixTimestampNano: uint64(time.Now().UnixNano()) / uint64(time.Millisecond),
	}).Encode()
	if err != nil {
		return nil, err
	}

	tracingHeaders := c.options.Propagator.GenerateHeaders(traceID)

	rt := c.vu.Runtime()
	var params *goja.Object
	if len(args) < 2 {
		params = rt.NewObject()
		if len(args) == 0 {
			args = []goja.Value{goja.Null(), params}
		} else {
			args = append(args, params)
		}
	} else {
		jsParams := args[1]
		if utils.IsNilly(jsParams) {
			params = rt.NewObject()
			args[1] = params
		} else {
			params = jsParams.ToObject(rt)
		}
	}

	var headers *goja.Object
	if jsHeaders := params.Get("headers"); utils.IsNilly(jsHeaders) {
		headers = rt.NewObject()
		params.Set("headers", headers)
	} else {
		headers = jsHeaders.ToObject(rt)
	}
	for key, val := range tracingHeaders {
		headers.Set(key, val)
	}

	c.setTags(rt, state, traceID, params)
	defer c.deleteTags(state)

	res, err := fn(c.vu.Context(), url, args...)
	return &HTTPResponse{Response: res, TraceID: traceID}, err
}

func (c *HttpClient) setTags(rt *goja.Runtime, state *lib.State, traceID string, params *goja.Object) {
	tracetestOptions := parseTracetestOptions(rt, params)
	state.Tags.Modify(func(tagsAndMeta *metrics.TagsAndMeta) {
		tagsAndMeta.SetMetadata("trace_id", traceID)

		if tracetestOptions.testID != "" {
			tagsAndMeta.SetMetadata("test_id", tracetestOptions.testID)
		} else if c.options.Tracetest.testID != "" {
			tagsAndMeta.SetMetadata("test_id", c.options.Tracetest.testID)
		}

		if tracetestOptions.testDefinition != "" {
			tagsAndMeta.SetMetadata("test_definition", tracetestOptions.testDefinition)
		} else if c.options.Tracetest.testDefinition != "" {
			tagsAndMeta.SetMetadata("test_definition", c.options.Tracetest.testDefinition)
		}
	})
}

func (c *HttpClient) deleteTags(state *lib.State) {
	state.Tags.Modify(func(tagsAndMeta *metrics.TagsAndMeta) {
		tagsAndMeta.DeleteMetadata("trace_id")
		tagsAndMeta.DeleteMetadata("test_id")
		tagsAndMeta.DeleteMetadata("test_definition")
	})
}
