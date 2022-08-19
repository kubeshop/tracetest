package comparator_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/stretchr/testify/assert"
)

func TestComparators(t *testing.T) {
	type compInput struct{ actual, expected string }
	type compErr struct{ actual, expected, err string }
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
			// actual = expected
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
			// actual != expected
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
			// actual > expected
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
			// actual >= expected
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
			// actual < expected
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
			// actual <= expected
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
			// actual CONTAINS expected
			expectSuccess: []compInput{
				{"hello", "he"},
				{"hello", "ll"},
				{"hello", "lo"},

				// https://github.com/kubeshop/tracetest/issues/617
				{`{"id":52}`, "52"},
			},
			expectNoMatch: []compInput{
				{"hello", "nop"},
			},
		},

		{
			name:       "Not contains",
			symbol:     "not-contains",
			comparator: comparator.NotContains,
			// actual NOT CONTAINS expected
			expectSuccess: []compInput{
				{"hello", "not"},
				{"hello", "ella"},
				{"hello", "helloo"},
				{`{"id":52}`, "56"},
			},
			expectNoMatch: []compInput{
				{"hello", "he"},
				{"hello", "hel"},
				{"hello", "ell"},
				{"hello", "ello"},
			},
		},

		// ***********
		{
			name:       "StartsWith",
			symbol:     "startsWith",
			comparator: comparator.StartsWith,
			expectSuccess: []compInput{
				{"hello", "he"},
			},
			expectNoMatch: []compInput{
				{"hello", "nop"},
				{"hello", "ll"},
				{"hello", "lo"},
			},
		},

		// ***********
		{
			name:       "EndsWith",
			symbol:     "endsWith",
			comparator: comparator.EndsWith,
			expectSuccess: []compInput{
				{"hello", "lo"},
			},
			expectNoMatch: []compInput{
				{"hello", "nop"},
				{"hello", "he"},
				{"hello", "ll"},
			},
		},
	}

	registry := comparator.DefaultRegistry()

	for _, c := range comps {
		t.Run(c.name, func(t *testing.T) {
			comp := c
			// t.Parallel()

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
