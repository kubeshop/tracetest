package assertions_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/assertions"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLiteralValues(t *testing.T) {
	testCases := []struct {
		Name           string
		Expression     model.AssertionExpression
		Span           traces.Span
		ExpectedOutput string
	}{
		{
			Name: "should_return_number",
			Expression: model.AssertionExpression{
				LiteralValue: model.LiteralValue{
					Value: "200.3",
					Type:  "number",
				},
			},
			Span:           traces.Span{},
			ExpectedOutput: "200.3",
		},
		{
			Name: "should_return_string",
			Expression: model.AssertionExpression{
				LiteralValue: model.LiteralValue{
					Value: "this is a string",
					Type:  "string",
				},
			},
			Span:           traces.Span{},
			ExpectedOutput: "this is a string",
		},
		{
			Name: "should_return_duration",
			Expression: model.AssertionExpression{
				LiteralValue: model.LiteralValue{
					Value: "3s",
					Type:  "duration",
				},
			},
			Span:           traces.Span{},
			ExpectedOutput: "3s",
		},
		{
			Name: "should_return_span_attribute",
			Expression: model.AssertionExpression{
				LiteralValue: model.LiteralValue{
					Value: "service.name",
					Type:  "attribute",
				},
			},
			Span: traces.Span{
				Attributes: traces.Attributes{
					"service.name": "pokeshop",
				},
			},
			ExpectedOutput: "pokeshop",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := assertions.ExecuteExpression(testCase.Expression, testCase.Span)

			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedOutput, output)
		})
	}
}

func TestSimpleExpressions(t *testing.T) {
	testCases := []struct {
		Name           string
		Expression     model.AssertionExpression
		Span           traces.Span
		ExpectedOutput string
	}{
		{
			Name: "should_sum_numbers_two_values",
			Expression: model.AssertionExpression{
				LiteralValue: model.LiteralValue{
					Value: "8",
					Type:  "number",
				},
				Operation: "+",
				Expression: &model.AssertionExpression{
					LiteralValue: model.LiteralValue{
						Value: "16",
						Type:  "number",
					},
				},
			},
			Span:           traces.Span{},
			ExpectedOutput: "24",
		},
		{
			Name: "should_execute_multiple_operations",
			Expression: model.AssertionExpression{
				LiteralValue: model.LiteralValue{
					Value: "8",
					Type:  "number",
				},
				Operation: "+",
				Expression: &model.AssertionExpression{
					LiteralValue: model.LiteralValue{
						Value: "16",
						Type:  "number",
					},
					Operation: "-",
					Expression: &model.AssertionExpression{
						LiteralValue: model.LiteralValue{
							Value: "26",
							Type:  "number",
						},
					},
				},
			},
			Span:           traces.Span{},
			ExpectedOutput: "-2",
		},
		{
			Name: "should_multiply_numbers",
			Expression: model.AssertionExpression{
				LiteralValue: model.LiteralValue{
					Value: "30",
					Type:  "number",
				},
				Operation: "*",
				Expression: &model.AssertionExpression{
					LiteralValue: model.LiteralValue{
						Value: "50",
						Type:  "number",
					},
				},
			},
			Span:           traces.Span{},
			ExpectedOutput: "1500",
		},
		{
			Name: "should_multiply_numbers_float_numbers",
			Expression: model.AssertionExpression{
				LiteralValue: model.LiteralValue{
					Value: "2.5",
					Type:  "number",
				},
				Operation: "*",
				Expression: &model.AssertionExpression{
					LiteralValue: model.LiteralValue{
						Value: "50",
						Type:  "number",
					},
				},
			},
			Span:           traces.Span{},
			ExpectedOutput: "125.00",
		},
		{
			Name: "should_divide_numbers",
			Expression: model.AssertionExpression{
				LiteralValue: model.LiteralValue{
					Value: "80",
					Type:  "number",
				},
				Operation: "/",
				Expression: &model.AssertionExpression{
					LiteralValue: model.LiteralValue{
						Value: "20",
						Type:  "number",
					},
				},
			},
			Span:           traces.Span{},
			ExpectedOutput: "4",
		},
		{
			Name: "should_sum_durations",
			Expression: model.AssertionExpression{
				LiteralValue: model.LiteralValue{
					Value: "60s",
					Type:  "duration",
				},
				Operation: "+",
				Expression: &model.AssertionExpression{
					LiteralValue: model.LiteralValue{
						Value: "800ms",
						Type:  "duration",
					},
				},
			},
			Span:           traces.Span{},
			ExpectedOutput: "60800000000", // 60.8s in nanoseconds
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := assertions.ExecuteExpression(testCase.Expression, testCase.Span)

			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedOutput, output)
		})
	}
}

func TestExpressionsWithAttributes(t *testing.T) {
	testCases := []struct {
		Name           string
		Expression     model.AssertionExpression
		Span           traces.Span
		ExpectedOutput string
	}{
		{
			Name: "should_subtract_attribute_with_number",
			Expression: model.AssertionExpression{
				LiteralValue: model.LiteralValue{
					Value: "myapp.item_stock_count",
					Type:  "attribute",
				},
				Operation: "-",
				Expression: &model.AssertionExpression{
					LiteralValue: model.LiteralValue{
						Value: "8",
						Type:  "number",
					},
				},
			},
			Span: traces.Span{
				Attributes: traces.Attributes{
					"myapp.item_stock_count":  "684",
					"myapp.order_items_count": "2",
				},
			},
			ExpectedOutput: "676",
		},
		{
			Name: "should_subtract_attribute_with_attribute",
			Expression: model.AssertionExpression{
				LiteralValue: model.LiteralValue{
					Value: "myapp.item_stock_count",
					Type:  "attribute",
				},
				Operation: "-",
				Expression: &model.AssertionExpression{
					LiteralValue: model.LiteralValue{
						Value: "myapp.order_items_count",
						Type:  "attribute",
					},
				},
			},
			Span: traces.Span{
				Attributes: traces.Attributes{
					"myapp.item_stock_count":  "684",
					"myapp.order_items_count": "2",
				},
			},
			ExpectedOutput: "682",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := assertions.ExecuteExpression(testCase.Expression, testCase.Span)

			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedOutput, output)
		})
	}
}
