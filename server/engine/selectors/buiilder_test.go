package selectors_test

import (
	"testing"

	"github.com/kubeshop/tracetest/engine/selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleSelectorBuilder(t *testing.T) {
	testCases := []struct {
		Name             string
		Expression       string
		ShouldSucceed    bool
		ExpectedSelector selectors.Selector
	}{
		{
			Name:          "Selector with single attribute",
			Expression:    "span[tracetest.span.type=\"http\"]",
			ShouldSucceed: true,
		},
		{
			Name:          "Selector with multiple attributes",
			Expression:    "span[service.name=\"Pokeshop\" tracetest.error=true]",
			ShouldSucceed: true,
		},
		{
			Name:          "Selector with child selector",
			Expression:    "span[tracetest.span.duration=0.5] span[tracetest.span.type contains \"http\"]",
			ShouldSucceed: true,
		},
		{
			Name:          "Selector with pseudo class filter",
			Expression:    "span[http.status_code=200]:nth_child(3)",
			ShouldSucceed: true,
		},
		{
			Name:          "Selector with invalid syntax",
			Expression:    "span.tracetest.span.type=\"http\"",
			ShouldSucceed: false,
		},
	}

	builder, err := selectors.NewSelectorBuilder()
	require.NoError(t, err, "builder must be created successfully")

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			_, err := builder.NewSelector(testCase.Expression)
			if testCase.ShouldSucceed {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
