package selectors_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/pkg/id"
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
			Name:            "EmptySelectorShouldSelectAllSpans",
			Expression:      ``,
			ExpectedSpanIds: []trace.SpanID{postImportSpanID, insertPokemonDatabaseSpanID, getPokemonFromExternalAPISpanID, updatePokemonDatabaseSpanID},
		},
		{
			Name:            "SelectorWithoutParametersShouldSelectAllSpans",
			Expression:      `span[]`,
			ExpectedSpanIds: []trace.SpanID{postImportSpanID, insertPokemonDatabaseSpanID, getPokemonFromExternalAPISpanID, updatePokemonDatabaseSpanID},
		},
		{
			Name:            "SelectorWithSpanName",
			Expression:      `span[name="Get pokemon from external API"]`,
			ExpectedSpanIds: []trace.SpanID{getPokemonFromExternalAPISpanID},
		},
		{
			Name:            "SelectorWithSimpleSingleAttributeQuerying",
			Expression:      `span[service.name="Pokeshop"]`,
			ExpectedSpanIds: []trace.SpanID{postImportSpanID, insertPokemonDatabaseSpanID},
		},
		{
			Name:            "MultipleSpanSelectors",
			Expression:      `span[service.name="Pokeshop"], span[service.name="Pokeshop-worker"]`,
			ExpectedSpanIds: []trace.SpanID{postImportSpanID, insertPokemonDatabaseSpanID, getPokemonFromExternalAPISpanID, updatePokemonDatabaseSpanID},
		},
		{
			Name:            "MultipleSpansUsingContains",
			Expression:      `span[service.name contains "Pokeshop"]`,
			ExpectedSpanIds: []trace.SpanID{postImportSpanID, insertPokemonDatabaseSpanID, getPokemonFromExternalAPISpanID, updatePokemonDatabaseSpanID},
		},
		{
			Name:            "SelectorWithMultipleAttributes",
			Expression:      `span[service.name="Pokeshop" tracetest.span.type="db"]`,
			ExpectedSpanIds: []trace.SpanID{insertPokemonDatabaseSpanID},
		},
		{
			Name:            "SelectorWithChildSelector",
			Expression:      `span[service.name="Pokeshop-worker"] span[tracetest.span.type="db"]`,
			ExpectedSpanIds: []trace.SpanID{updatePokemonDatabaseSpanID},
		},
		{
			Name:            "SelectorToSelectAllChildrenSpans",
			Expression:      `span[service.name="Pokeshop" tracetest.span.type="http"] span[]`,
			ExpectedSpanIds: []trace.SpanID{insertPokemonDatabaseSpanID, getPokemonFromExternalAPISpanID, updatePokemonDatabaseSpanID},
		},
		{
			Name:            "SelectorWithFirstPseudoClass",
			Expression:      `span[tracetest.span.type="db"]:first`,
			ExpectedSpanIds: []trace.SpanID{insertPokemonDatabaseSpanID},
		},
		{
			Name:            "SelectorWithFirstPseudoClass",
			Expression:      `span[tracetest.span.type="db"]:last`,
			ExpectedSpanIds: []trace.SpanID{updatePokemonDatabaseSpanID},
		},
		{
			Name:            "SelectorWithNthChildPseudoClass",
			Expression:      `span[tracetest.span.type="db"]:nth_child(2)`,
			ExpectedSpanIds: []trace.SpanID{updatePokemonDatabaseSpanID},
		},
		{
			Name:            "SelectorShouldNotMatchParentWhenChildrenAreSpecified",
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
