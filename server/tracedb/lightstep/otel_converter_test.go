package lightstep_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/tracedb/lightstep"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func traceId(id string) trace.TraceID {
	traceId, _ := trace.TraceIDFromHex(id)
	return traceId
}

func spanId(id string) trace.SpanID {
	spanId, _ := trace.SpanIDFromHex(id)
	return spanId
}

func TestResponseConversion(t *testing.T) {
	outputBytes, err := ioutil.ReadFile("./lightstep_output.json")
	require.NoError(t, err)

	expectedRoot := traces.Span{
		ID:        spanId("7c74a3963449b028"),
		Name:      "GET /api/tests/{testId}",
		StartTime: time.UnixMicro(1657229019567381),
		EndTime:   time.UnixMicro(1657229019626741),
		Parent:    nil,
		Attributes: traces.Attributes{
			"http.method":             "GET",
			"http.request.headers":    "{\"Accept\":[\"*/*\"],\"Accept-Language\":[\"en-US,en;q=0.9\"],\"Connection\":[\"Keep-Alive\"],\"Cookie\":[\"_ga=GA1.1.1562992345.1652882058; _ga_WP4XXN1FYN=GS1.1.1657209966.121.0.1657210189.0; _ga_ZP277L2M37=GS1.1.1657214237.23.1.1657214237.0\"],\"Referer\":[\"https://beta.tracetest.io/test/a868238d-e567-48ac-9548-6e3f68c4bfd0/run/51357e0e-b413-4de6-b1cb-4fb9e8260628\"],\"Sec-Ch-Ua\":[\"\\\".Not/A)Brand\\\";v=\\\"99\\\", \\\"Google Chrome\\\";v=\\\"103\\\", \\\"Chromium\\\";v=\\\"103\\\"\"],\"Sec-Ch-Ua-Mobile\":[\"?0\"],\"Sec-Ch-Ua-Platform\":[\"\\\"macOS\\\"\"],\"Sec-Fetch-Dest\":[\"empty\"],\"Sec-Fetch-Mode\":[\"cors\"],\"Sec-Fetch-Site\":[\"same-origin\"],\"User-Agent\":[\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36\"],\"Via\":[\"1.1 google\"],\"X-Cloud-Trace-Context\":[\"55e44bf1f6f589f38f0a8b14c0094897/17675920434189115578\"],\"X-Forwarded-For\":[\"200.55.228.62, 34.111.76.37\"],\"X-Forwarded-Proto\":[\"https\"]}",
			"http.request.params":     "{\"testId\":\"a868238d-e567-48ac-9548-6e3f68c4bfd0\"}",
			"http.request.query":      "{}",
			"http.route":              "/api/tests/{testId}",
			"http.target":             "/api/tests/a868238d-e567-48ac-9548-6e3f68c4bfd0",
			"instrumentation.name":    "tracetest",
			"instrumentation.version": "",
			"otlp.trace_id":           "1a840fa833f51bcdefedad84c9bb0cdf",
			"span.kind":               "internal",
		},
	}

	expectedChildSpan := traces.Span{
		ID:        spanId("f3120667225dd2ca"),
		Name:      "query SELECT",
		StartTime: time.UnixMicro(1657229019625260),
		EndTime:   time.UnixMicro(1657229019626167),
		Attributes: traces.Attributes{
			"instrumentation.name":    "github.com/j2gg0s/otsql",
			"instrumentation.version": "",
			"otlp.trace_id":           "1a840fa833f51bcdefedad84c9bb0cdf",
			"parent_span_guid":        "7c74a3963449b028",
			"service.name":            "tracetest",
			"span.kind":               "client",
			"sql.arg.1":               "a868238d-e567-48ac-9548-6e3f68c4bfd0",
			"sql.database":            "",
			"sql.instance":            "",
			"sql.query":               "\nSELECT t.test, d.definition\n\tFROM tests t\n\tJOIN definitions d ON d.test_id = t.id AND d.test_version = t.version\n WHERE t.id = $1 ORDER BY t.version DESC LIMIT 1",
		},
		Parent: &expectedRoot,
	}

	expectedRoot.Children = []*traces.Span{&expectedChildSpan}

	expectedTrace := traces.Trace{
		ID:       traceId("efedad84c9bb0cdf"),
		RootSpan: expectedRoot,
	}

	var response lightstep.GetTraceResponse
	err = json.Unmarshal(outputBytes, &response)
	require.NoError(t, err)

	trace := lightstep.ConvertResponseToOtelFormat(response)
	assert.Equal(t, expectedTrace, trace)
}
