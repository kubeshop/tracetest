package selectors_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

var (
	gen                             = id.NewRandGenerator()
	postImportSpanID                = gen.SpanID()
	insertPokemonDatabaseSpanID     = gen.SpanID()
	getPokemonFromExternalAPISpanID = gen.SpanID()
	updatePokemonDatabaseSpanID     = gen.SpanID()
)
var pokeshopTrace = traces.Trace{
	ID: gen.TraceID(),
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
			Name:            "Empty_selector_should_select_all_spans",
			Expression:      ``,
			ExpectedSpanIds: []trace.SpanID{postImportSpanID, insertPokemonDatabaseSpanID, getPokemonFromExternalAPISpanID, updatePokemonDatabaseSpanID},
		},
		{
			Name:            "Selector_without_parameters_should_select_all_spans",
			Expression:      `span[]`,
			ExpectedSpanIds: []trace.SpanID{postImportSpanID, insertPokemonDatabaseSpanID, getPokemonFromExternalAPISpanID, updatePokemonDatabaseSpanID},
		},
		{
			Name:            "Selector_with_span_name",
			Expression:      `span[name="Get pokemon from external API"]`,
			ExpectedSpanIds: []trace.SpanID{getPokemonFromExternalAPISpanID},
		},
		{
			Name:            "Selector_with_simple_single_attribute_querying",
			Expression:      `span[service.name="Pokeshop"]`,
			ExpectedSpanIds: []trace.SpanID{postImportSpanID, insertPokemonDatabaseSpanID},
		},
		{
			Name:            "Multiple_span_selectors",
			Expression:      `span[service.name="Pokeshop"], span[service.name="Pokeshop-worker"]`,
			ExpectedSpanIds: []trace.SpanID{postImportSpanID, insertPokemonDatabaseSpanID, getPokemonFromExternalAPISpanID, updatePokemonDatabaseSpanID},
		},
		{
			Name:            "Multiple_spans_using_contains",
			Expression:      `span[service.name contains "Pokeshop"]`,
			ExpectedSpanIds: []trace.SpanID{postImportSpanID, insertPokemonDatabaseSpanID, getPokemonFromExternalAPISpanID, updatePokemonDatabaseSpanID},
		},
		{
			Name:            "Selector_with_multiple_attributes",
			Expression:      `span[service.name="Pokeshop" tracetest.span.type="db"]`,
			ExpectedSpanIds: []trace.SpanID{insertPokemonDatabaseSpanID},
		},
		{
			Name:            "Selector_with_child_selector",
			Expression:      `span[service.name="Pokeshop-worker"] span[tracetest.span.type="db"]`,
			ExpectedSpanIds: []trace.SpanID{updatePokemonDatabaseSpanID},
		},
		{
			Name:            "Selector_with_first_pseudo_class",
			Expression:      `span[tracetest.span.type="db"]:first`,
			ExpectedSpanIds: []trace.SpanID{insertPokemonDatabaseSpanID},
		},
		{
			Name:            "Selector_with_first_pseudo_class",
			Expression:      `span[tracetest.span.type="db"]:last`,
			ExpectedSpanIds: []trace.SpanID{updatePokemonDatabaseSpanID},
		},
		{
			Name:            "Selector_with_nth_child_pseudo_class",
			Expression:      `span[tracetest.span.type="db"]:nth_child(2)`,
			ExpectedSpanIds: []trace.SpanID{updatePokemonDatabaseSpanID},
		},
		{
			Name:            "Selector_should_not_match_parent_when_children_are_specified",
			Expression:      `span[service.name = "Pokeshop-worker"] span[service.name = "Pokeshop-worker"]`,
			ExpectedSpanIds: []trace.SpanID{updatePokemonDatabaseSpanID},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			selector, err := selectors.New(testCase.Expression)
			require.NoError(t, err)

			spans := selector.Filter(pokeshopTrace)
			ensureExpectedSpansWereReturned(t, testCase.ExpectedSpanIds, spans)
		})
	}
}

func ensureExpectedSpansWereReturned(t *testing.T, spanIDs []trace.SpanID, spans []traces.Span) {
	assert.Len(t, spans, len(spanIDs), "Should_return_the_same_number_of_spans_as_we_expected")
	for _, span := range spans {
		assert.Contains(t, spanIDs, span.ID, "span ID was returned but wasn't expected")
	}
}
