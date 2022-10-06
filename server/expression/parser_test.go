package expression_test

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/server/expression"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	Name           string
	Query          string
	ExpectedOutput expression.Statement
}

func TestSimpleParsingRules(t *testing.T) {
	testCases := []testCase{
		{
			Name:  "should_parse_1=1",
			Query: `1 = 1`,
			ExpectedOutput: expression.Statement{
				Left:       numberExpr(1),
				Comparator: "=",
				Right:      numberExpr(1),
			},
		},
		{
			Name:  "should_parse_100ms=100ms",
			Query: `100ms = 100ms`,
			ExpectedOutput: expression.Statement{
				Left:       durationExpr("100ms"),
				Comparator: "=",
				Right:      durationExpr("100ms"),
			},
		},
		{
			Name:  "should_parse_attr:abc=attr:abc",
			Query: `attr:abc = attr:abc`,
			ExpectedOutput: expression.Statement{
				Left:       attrExpr("abc"),
				Comparator: "=",
				Right:      attrExpr("abc"),
			},
		},
		{
			Name:  "should_parse_abc=abc",
			Query: `"abc" = "abc"`,
			ExpectedOutput: expression.Statement{
				Left:       strExpr("abc"),
				Comparator: "=",
				Right:      strExpr("abc"),
			},
		},
		{
			Name:  "should_parse_escaped_strings",
			Query: `"my name is \"john\"" = 'my name is \'john\''`,
			ExpectedOutput: expression.Statement{
				Left:       strExpr(`my name is "john"`),
				Comparator: "=",
				Right:      strExpr(`my name is 'john'`),
			},
		},
		{
			Name:  "should_parse_double_quotes_inside_single_quotes",
			Query: `attr:text = 'this is a "nice" test'`,
			ExpectedOutput: expression.Statement{
				Left:       attrExpr("text"),
				Comparator: "=",
				Right:      strExpr(`this is a "nice" test`),
			},
		},
	}

	runTestCases(t, testCases)
}

func TestExpressions(t *testing.T) {
	testCases := []testCase{
		{
			Name:  "should_parse_expressions_with_addition_operator",
			Query: "2 + 1 = 1 + 2",
			ExpectedOutput: expression.Statement{
				Left: expression.Expr{
					Left:  &expression.Term{Number: strp("2")},
					Right: []*expression.OpTerm{{Operator: "+", Term: &expression.Term{Number: strp("1")}}},
				},
				Comparator: "=",
				Right: expression.Expr{
					Left:  &expression.Term{Number: strp("1")},
					Right: []*expression.OpTerm{{Operator: "+", Term: &expression.Term{Number: strp("2")}}},
				},
			},
		},
		{
			Name:  "should_parse_expressions_with_subtraction_operator",
			Query: "2 - 1 != 1 - 2",
			ExpectedOutput: expression.Statement{
				Left: expression.Expr{
					Left:  &expression.Term{Number: strp("2")},
					Right: []*expression.OpTerm{{Operator: "-", Term: &expression.Term{Number: strp("1")}}},
				},
				Comparator: "!=",
				Right: expression.Expr{
					Left:  &expression.Term{Number: strp("1")},
					Right: []*expression.OpTerm{{Operator: "-", Term: &expression.Term{Number: strp("2")}}},
				},
			},
		},
		{
			Name:  "should_parse_expressions_with_multiplication_operator",
			Query: "2 * 1 = 1 * 2",
			ExpectedOutput: expression.Statement{
				Left: expression.Expr{
					Left:  &expression.Term{Number: strp("2")},
					Right: []*expression.OpTerm{{Operator: "*", Term: &expression.Term{Number: strp("1")}}},
				},
				Comparator: "=",
				Right: expression.Expr{
					Left:  &expression.Term{Number: strp("1")},
					Right: []*expression.OpTerm{{Operator: "*", Term: &expression.Term{Number: strp("2")}}},
				},
			},
		},
		{
			Name:  "should_parse_expressions_with_division_operator",
			Query: "2 / 1 != 1 / 2",
			ExpectedOutput: expression.Statement{
				Left: expression.Expr{
					Left:  &expression.Term{Number: strp("2")},
					Right: []*expression.OpTerm{{Operator: "/", Term: &expression.Term{Number: strp("1")}}},
				},
				Comparator: "!=",
				Right: expression.Expr{
					Left:  &expression.Term{Number: strp("1")},
					Right: []*expression.OpTerm{{Operator: "/", Term: &expression.Term{Number: strp("2")}}},
				},
			},
		},
	}

	runTestCases(t, testCases)
}

func TestStringInterpolation(t *testing.T) {
	testCases := []testCase{
		{
			Name:  "should_interpolate_simple_values",
			Query: `attr:text = "your age is ${18}"`,
			ExpectedOutput: expression.Statement{
				Left:       attrExpr("text"),
				Comparator: "=",
				Right: expression.Expr{
					Left: &expression.Term{
						Str: &expression.Str{
							Text: "your age is %s",
							Args: []expression.Expr{
								{Left: &expression.Term{Number: strp("18")}},
							},
						},
					},
				},
			},
		},
		{
			Name:  "should_interpolate_multiple_values",
			Query: `attr:text = "your age is ${18} but you must be ${21} to start drinking"`,
			ExpectedOutput: expression.Statement{
				Left:       attrExpr("text"),
				Comparator: "=",
				Right: expression.Expr{
					Left: &expression.Term{
						Str: &expression.Str{
							Text: "your age is %s but you must be %s to start drinking",
							Args: []expression.Expr{
								{Left: &expression.Term{Number: strp("18")}},
								{Left: &expression.Term{Number: strp("21")}},
							},
						},
					},
				},
			},
		},
		{
			Name:  "should_interpolate_complex_expressions",
			Query: `attr:text = "your age will be ${attr:user_age + 2} in two years"`,
			ExpectedOutput: expression.Statement{
				Left:       attrExpr("text"),
				Comparator: "=",
				Right: expression.Expr{
					Left: &expression.Term{
						Str: &expression.Str{
							Text: "your age will be %s in two years",
							Args: []expression.Expr{
								{
									Left: &expression.Term{Attribute: attrp("user_age")},
									Right: []*expression.OpTerm{
										{Operator: "+", Term: &expression.Term{Number: strp("2")}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name:  "should_support_json_formatting",
			Query: `attr:tracetest.response.body = '{"userID":"${attr:myapp.myspan.user_id}"}'`,
			ExpectedOutput: expression.Statement{
				Left:       attrExpr("tracetest.response.body"),
				Comparator: "=",
				Right: expression.Expr{
					Left: &expression.Term{
						Str: &expression.Str{
							Text: `{"userID":"%s"}`,
							Args: []expression.Expr{
								{
									Left: &expression.Term{Attribute: attrp("myapp.myspan.user_id")},
								},
							},
						},
					},
				},
			},
		},
	}

	runTestCases(t, testCases)
}

func runTestCases(t *testing.T, testCases []testCase) {
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := expression.Parse(testCase.Query)
			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedOutput, output)
		})
	}
}

func strp(in string) *string {
	return &in
}

func numberExpr(number int) expression.Expr {
	return expression.Expr{
		Left: &expression.Term{
			Number: strp(fmt.Sprintf("%d", number)),
		},
	}
}

func durationExpr(duration string) expression.Expr {
	return expression.Expr{
		Left: &expression.Term{
			Duration: &duration,
		},
	}
}

func attrExpr(attrName string) expression.Expr {
	return expression.Expr{
		Left: &expression.Term{
			Attribute: attrp(attrName),
		},
	}
}

func attrp(attrName string) *expression.Attribute {
	attr := expression.Attribute(attrName)
	return &attr
}

func strExpr(str string) expression.Expr {
	s := expression.Str{
		Text: str,
		Args: []expression.Expr{},
	}
	return expression.Expr{
		Left: &expression.Term{
			Str: &s,
		},
	}
}
