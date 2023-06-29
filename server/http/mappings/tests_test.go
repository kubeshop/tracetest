package mappings_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_OpenApiToModel_Outputs(t *testing.T) {
	in := openapi.Test{
		Outputs: []openapi.TestOutput{
			{
				Name: "OUTPUT",
				SelectorParsed: openapi.Selector{
					Query: `span[name="root"]`,
				},
				Value: "attr:tracetest.selected_spans.count",
			},
		},
	}

	expected := (maps.Ordered[string, test.Output]{}).MustAdd("OUTPUT", test.Output{
		Selector: `span[name="root"]`,
		Value:    "attr:tracetest.selected_spans.count",
	})

	m := mappings.New(traces.NewConversionConfig(), nil)

	actual, err := m.In.Test(in)
	require.NoError(t, err)

	assert.Equal(t, expected, actual.Outputs)
}
