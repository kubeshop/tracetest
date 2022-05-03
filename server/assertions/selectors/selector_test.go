package selectors_test

import (
	"testing"

	"github.com/kubeshop/tracetest/assertions/selectors"
	"github.com/kubeshop/tracetest/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

var (
	postImportSpanID                = createSpanID("00000001")
	insertPokemonDatabaseSpanID     = createSpanID("00000002")
	getPokemonFromExternalAPISpanID = createSpanID("00000003")
	updatePokemonDatabaseSpanID     = createSpanID("00000004")
)
var pokeshopTrace = traces.Trace{
	ID: trace.TraceID{},
	RootSpan: traces.Span{
		ID: postImportSpanID,
		Attributes: traces.Attributes{
			"service.name":        "Pokeshop",
			"tracetest.span.type": "http",
			"http.status_code":    "201",
		},
		Name: "POST /import",
		Children: []*traces.Span{
			{
				ID: insertPokemonDatabaseSpanID,
				Attributes: traces.Attributes{
					"service.name":        "Pokeshop",
					"tracetest.span.type": "db",
					"db.statement":        "INSERT INTO pokemon (id) values (?)",
				},
				Name: "Insert pokemon into database",
			},
			{
				ID: getPokemonFromExternalAPISpanID,
				Attributes: traces.Attributes{
					"service.name":        "Pokeshop-worker",
					"tracetest.span.type": "http",
					"http.status_code":    "200",
				},
				Name: "Get pokemon from external API",
				Children: []*traces.Span{
					{
						ID: updatePokemonDatabaseSpanID,
						Attributes: traces.Attributes{
							"service.name":        "Pokeshop-worker",
							"tracetest.span.type": "db",
							"db.statement":        "UPDATE pokemon (name = ?) WHERE id = ?",
						},
						Name: "Update pokemon on database",
					},
				},
			},
		},
	},
}

func TestSelector(t *testing.T) {
	testCases := []struct {
		Name            string
		Expression      string
		ExpectedSpanIds []trace.SpanID
	}{
		{
			Name:            "Selector with simple single attribute querying",
			Expression:      "span[service.name=\"Pokeshop\"]",
			ExpectedSpanIds: []trace.SpanID{postImportSpanID, insertPokemonDatabaseSpanID},
		},
		{
			Name:            "Selector with multiple attributes",
			Expression:      "span[service.name=\"Pokeshop\" tracetest.span.type=\"db\"]",
			ExpectedSpanIds: []trace.SpanID{insertPokemonDatabaseSpanID},
		},
		{
			Name:            "Selector with child selector",
			Expression:      "span[service.name=\"Pokeshop-worker\"] span[tracetest.span.type=\"db\"]",
			ExpectedSpanIds: []trace.SpanID{updatePokemonDatabaseSpanID},
		},
		{
			Name:            "Selector with pseudo class",
			Expression:      "span[service.name=\"tracetest.span.type=\"db\"]:nth_child(2)",
			ExpectedSpanIds: []trace.SpanID{updatePokemonDatabaseSpanID},
		},
	}

	builder, err := selectors.NewSelectorBuilder()
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			selector, err := builder.NewSelector(testCase.Expression)
			assert.NoError(t, err)

			spans := selector.Filter(pokeshopTrace)
			ensureExpectedSpansWereReturned(t, testCase.ExpectedSpanIds, spans)
		})
	}
}

func ensureExpectedSpansWereReturned(t *testing.T, spanIDs []trace.SpanID, spans []traces.Span) {
	assert.Len(t, spans, len(spanIDs), "Should return the same number of spans as we expected")
	for _, span := range spans {
		assert.Contains(t, spanIDs, span.ID, "span ID was returned but wasn't expected")
	}
}

func createSpanID(id string) trace.SpanID {
	stringBytes := []byte(id)
	bytes := [8]byte{}
	for i, b := range stringBytes {
		bytes[i] = b
	}
	return trace.SpanID(bytes)
}
