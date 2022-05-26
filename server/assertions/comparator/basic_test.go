package comparator_test

import (
	"testing"

	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/stretchr/testify/assert"
)

func TestComparators(t *testing.T) {
	type compInput struct{ expected, actual string }
	type compErr struct{ expected, actual, err string }
	comps := []struct {
		name          string
		symbol        string
		comparator    comparator.Comparator
		expectSuccess []compInput
		expectNoMatch []compInput
		expectError   []compErr
	}{
		// ***********
		{
			name:       "Equal",
			symbol:     "=",
			comparator: comparator.Eq,
			expectSuccess: []compInput{
				{"a", "a"},
				{"1", "1"},
			},
			expectNoMatch: []compInput{
				{"", "2"},
				{"a", "b"},
			},
		},

		// ***********
		{
			name:       "NotEqual",
			symbol:     "!=",
			comparator: comparator.Neq,
			expectSuccess: []compInput{
				{"", "2"},
				{"a", "b"},
			},
			expectNoMatch: []compInput{
				{"a", "a"},
				{"1", "1"},
			},
		},

		// ***********
		{
			name:       "Gt",
			symbol:     ">",
			comparator: comparator.Gt,
			expectSuccess: []compInput{
				{"2", "1"},
				{"10", "2"},
			},
			expectNoMatch: []compInput{
				{"1", "1"},
				{"1", "2"},
			},
			expectError: []compErr{
				{"a", "1", `cannot parse "a" as integer`},
			},
		},

		// ***********
		{
			name:       "Gte",
			symbol:     ">=",
			comparator: comparator.Gte,
			expectSuccess: []compInput{
				{"2", "1"},
				{"2", "2"},
			},
			expectNoMatch: []compInput{
				{"1", "2"},
			},
			expectError: []compErr{
				{"a", "1", `cannot parse "a" as integer`},
			},
		},

		// ***********
		{
			name:       "Lt",
			symbol:     "<",
			comparator: comparator.Lt,
			expectSuccess: []compInput{
				{"1", "2"},
				{"2", "10"},
			},
			expectNoMatch: []compInput{
				{"1", "1"},
				{"2", "1"},
			},
			expectError: []compErr{
				{"a", "1", `cannot parse "a" as integer`},
			},
		},

		// ***********
		{
			name:       "Lte",
			symbol:     "<=",
			comparator: comparator.Lte,
			expectSuccess: []compInput{
				{"1", "2"},
				{"1", "1"},
			},
			expectNoMatch: []compInput{
				{"2", "1"},
			},
			expectError: []compErr{
				{"a", "1", `cannot parse "a" as integer`},
			},
		},

		// ***********
		{
			name:       "Contains",
			symbol:     "contains",
			comparator: comparator.Contains,
			expectSuccess: []compInput{
				{"he", "hello"},
				{"ll", "hello"},
				{"lo", "hello"},
			},
			expectNoMatch: []compInput{
				{"nop", "hello"},
			},
		},

		// ***********
		{
			name:       "StartsWith",
			symbol:     "startsWith",
			comparator: comparator.StartsWith,
			expectSuccess: []compInput{
				{"he", "hello"},
			},
			expectNoMatch: []compInput{
				{"nop", "hello"},
				{"ll", "hello"},
				{"lo", "hello"},
			},
		},

		// ***********
		{
			name:       "EndsWith",
			symbol:     "endsWith",
			comparator: comparator.EndsWith,
			expectSuccess: []compInput{
				{"lo", "hello"},
			},
			expectNoMatch: []compInput{
				{"nop", "hello"},
				{"he", "hello"},
				{"ll", "hello"},
			},
		},
	}

	registry := comparator.DefaultRegistry()

	for _, c := range comps {
		t.Run(c.name, func(t *testing.T) {
			comp := c
			t.Parallel()

			assert.Equal(t, comp.symbol, comp.comparator.String())
			found, err := registry.Get(comp.comparator.String())
			assert.NoError(
				t, err,
				`comparator "%s" not included in default registry`, comp.comparator.String(),
			)
			assert.Equal(
				t, comp.comparator, found,
				`comparator "%s" is incorrectly registered`, comp.comparator.String(),
			)

			for _, input := range comp.expectSuccess {
				assert.NoError(
					t, comp.comparator.Compare(input.expected, input.actual),
					`expected success comparing "%s" with "%s"`, input.expected, input.actual,
				)
			}

			for _, input := range comp.expectNoMatch {
				assert.ErrorIs(
					t, comp.comparator.Compare(input.expected, input.actual), comparator.ErrNoMatch,
					`expected NoMatch comparing "%s" with "%s"`, input.expected, input.actual,
				)
			}

			for _, input := range comp.expectError {
				assert.EqualError(
					t, comp.comparator.Compare(input.expected, input.actual), input.err,
					`expected "%s" err comparing "%s" with "%s"`, input.err, input.expected, input.actual,
				)
			}

		})
	}
}
