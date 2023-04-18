package mappings_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_OpenApiToModel_Outputs(t *testing.T) {
	in := openapi.Test{
		Outputs: []openapi.TestOutput{
			{
				Name: "OUTPUT",
				Selector: openapi.Selector{
					Query: `span[name="root"]`,
				},
				Value: "attr:tracetest.selected_spans.count",
			},
		},
	}

	expected := (maps.Ordered[string, model.Output]{}).MustAdd("OUTPUT", model.Output{
		Selector: `span[name="root"]`,
		Value:    "attr:tracetest.selected_spans.count",
	})

	m := mappings.New(traces.NewConversionConfig(), nil, nil)

	actual, err := m.In.Test(in)
	require.NoError(t, err)

	assert.Equal(t, expected, actual.Outputs)
}
