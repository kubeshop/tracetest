package output

import (
	"fmt"
	"strconv"

	"github.com/xoscar/xk6-tracetest-tracing/modules/tracetest"
	"go.k6.io/k6/lib/netext/httpext"
	"go.k6.io/k6/metrics"
)

func (o *Output) handleSample(sample metrics.SampleContainer) {
	if httpSample, ok := sample.(*httpext.Trail); ok {
		o.handleHttpSample(httpSample)
	}
}

func (o *Output) handleHttpSample(trail *httpext.Trail) {
	traceID, hasTrace := trail.Metadata["trace_id"]
	testID, hasTestID := trail.Metadata["test_id"]
	testDefinition, hasTestDefinition := trail.Metadata["test_definition"]

	if !hasTrace || (!hasTestID && !hasTestDefinition) {
		return
	}

	totalDuration := trail.Blocked + trail.ConnDuration + trail.Duration
	startTime := trail.EndTime.Add(-totalDuration)

	getTag := func(name string) string {
		val, _ := trail.Tags.Get(name)
		return val
	}

	strStatus := getTag("status")
	status, err := strconv.ParseInt(strStatus, 10, 64)
	if err != nil {
		o.logger.Warnf("unexpected error parsing status '%s': %w", strStatus, err)
		return
	}

	metadata := tracetest.Metadata{
		"StartTimeUnixNano": fmt.Sprint(startTime.UnixNano()),
		"EndTimeUnixNano":   fmt.Sprint(trail.EndTime.UnixNano()),
		"Group":             getTag("group"),
		"Scenario":          getTag("scenario"),
		"TraceID":           traceID,
		"HTTPUrl":           getTag("url"),
		"HTTPMethod":        getTag("method"),
		"HTTPStatus":        fmt.Sprint(status),
	}

	// do something with the metadata
	o.logger.Println("METADATA: ", metadata)

	if hasTestID {
		o.tracetest.RunTest(testID, traceID)
	} else {
		o.tracetest.RunFromDefinition(testDefinition, traceID)
	}
}
