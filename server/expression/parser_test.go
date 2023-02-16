package expression_test

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/server/expression"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type parserTestCase struct {
	Name           string
	Query          string
	ExpectedOutput expression.Statement
}

func TestSimpleParsingRules(t *testing.T) {
	testCases := []parserTestCase{
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
			Name:  "should_parse_100_ms=100_ms",
			Query: `100 ms = 100 ms`,
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
				Left:       strExpr(`my name is \"john\"`),
				Comparator: "=",
				Right:      strExpr(`my name is \'john\'`),
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
	testCases := []parserTestCase{
		{
			Name:  "should_parse_expressions_with_addition_operator",
			Query: "2 + 1 = 1 + 2",
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left:  &expression.Term{Number: strp("2")},
					Right: []*expression.OpTerm{{Operator: "+", Term: &expression.Term{Number: strp("1")}}},
				},
				Comparator: "=",
				Right: &expression.Expr{
					Left:  &expression.Term{Number: strp("1")},
					Right: []*expression.OpTerm{{Operator: "+", Term: &expression.Term{Number: strp("2")}}},
				},
			},
		},
		{
			Name:  "should_parse_expressions_with_subtraction_operator",
			Query: "2 - 1 != 1 - 2",
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left:  &expression.Term{Number: strp("2")},
					Right: []*expression.OpTerm{{Operator: "-", Term: &expression.Term{Number: strp("1")}}},
				},
				Comparator: "!=",
				Right: &expression.Expr{
					Left:  &expression.Term{Number: strp("1")},
					Right: []*expression.OpTerm{{Operator: "-", Term: &expression.Term{Number: strp("2")}}},
				},
			},
		},
		{
			Name:  "should_parse_expressions_with_multiplication_operator",
			Query: "2 * 1 = 1 * 2",
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left:  &expression.Term{Number: strp("2")},
					Right: []*expression.OpTerm{{Operator: "*", Term: &expression.Term{Number: strp("1")}}},
				},
				Comparator: "=",
				Right: &expression.Expr{
					Left:  &expression.Term{Number: strp("1")},
					Right: []*expression.OpTerm{{Operator: "*", Term: &expression.Term{Number: strp("2")}}},
				},
			},
		},
		{
			Name:  "should_parse_expressions_with_division_operator",
			Query: "2 / 1 != 1 / 2",
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left:  &expression.Term{Number: strp("2")},
					Right: []*expression.OpTerm{{Operator: "/", Term: &expression.Term{Number: strp("1")}}},
				},
				Comparator: "!=",
				Right: &expression.Expr{
					Left:  &expression.Term{Number: strp("1")},
					Right: []*expression.OpTerm{{Operator: "/", Term: &expression.Term{Number: strp("2")}}},
				},
			},
		},
	}

	runTestCases(t, testCases)
}

func TestStringInterpolation(t *testing.T) {
	testCases := []parserTestCase{
		{
			Name:  "should_interpolate_simple_values",
			Query: `attr:text = "your age is ${18}"`,
			ExpectedOutput: expression.Statement{
				Left:       attrExpr("text"),
				Comparator: "=",
				Right: &expression.Expr{
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
				Right: &expression.Expr{
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
				Right: &expression.Expr{
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
				Right: &expression.Expr{
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

func TestFilters(t *testing.T) {
	testCases := []parserTestCase{
		{
			Name:  "should_allow_filter_on_left_hand_side",
			Query: `attr:my_json_attribute | json_path '.id' = 32`,
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left: &expression.Term{
						Attribute: attrp("my_json_attribute"),
					},
					Filters: []*expression.Filter{
						{
							Name: "json_path", Args: []*expression.Term{
								{
									Str: &expression.Str{
										Text: ".id",
										Args: []expression.Expr{},
									},
								},
							},
						},
					},
				},
				Comparator: "=",
				Right: &expression.Expr{
					Left: &expression.Term{
						Number: strp("32"),
					},
				},
			},
		},
		{
			Name:  "should_allow_filter_on_right_hand_side",
			Query: `attr:user_id = attr:tracetest.response.body | json_path '.id'`,
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left: &expression.Term{
						Attribute: attrp("user_id"),
					},
				},
				Comparator: "=",
				Right: &expression.Expr{
					Left: &expression.Term{
						Attribute: attrp("tracetest.response.body"),
					},
					Filters: []*expression.Filter{
						{
							Name: "json_path", Args: []*expression.Term{
								{
									Str: &expression.Str{
										Text: ".id",
										Args: []expression.Expr{},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name:  "should_allow_filters_with_multiple_arguments",
			Query: `attr:my_json_attribute | my_function 'arg1' 'arg2' 42 = "abc"`,
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left: &expression.Term{
						Attribute: attrp("my_json_attribute"),
					},
					Filters: []*expression.Filter{
						{
							Name: "my_function", Args: []*expression.Term{
								{
									Str: &expression.Str{
										Text: "arg1",
										Args: []expression.Expr{},
									},
								},
								{
									Str: &expression.Str{
										Text: "arg2",
										Args: []expression.Expr{},
									},
								},
								{
									Number: strp("42"),
								},
							},
						},
					},
				},
				Comparator: "=",
				Right: &expression.Expr{
					Left: &expression.Term{
						Str: &expression.Str{
							Text: "abc",
							Args: []expression.Expr{},
						},
					},
				},
			},
		},
		{
			Name:  "should_allow_chaining_filters",
			Query: `attr:my_json_attribute | json_path '.name' | lowercase = "john"`,
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left: &expression.Term{
						Attribute: attrp("my_json_attribute"),
					},
					Filters: []*expression.Filter{
						{
							Name: "json_path", Args: []*expression.Term{
								{
									Str: &expression.Str{
										Text: ".name",
										Args: []expression.Expr{},
									},
								},
							},
						},
						{
							Name: "lowercase",
						},
					},
				},
				Comparator: "=",
				Right: &expression.Expr{
					Left: &expression.Term{
						Str: &expression.Str{
							Text: "john",
							Args: []expression.Expr{},
						},
					},
				},
			},
		},
		{
			Name:  "should_allow_filters_inside_string_interpolation",
			Query: `attr:message = "welcome to tracetest, ${attr:tracetest.response.body | json_path '.name'}"`,
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left: &expression.Term{
						Attribute: attrp("message"),
					},
				},
				Comparator: "=",
				Right: &expression.Expr{
					Left: &expression.Term{
						Str: &expression.Str{
							Text: "welcome to tracetest, %s",
							Args: []expression.Expr{
								{
									Left: &expression.Term{
										Attribute: attrp("tracetest.response.body"),
									},
									Filters: []*expression.Filter{
										{
											Name: "json_path", Args: []*expression.Term{
												{
													Str: &expression.Str{
														Text: ".name",
														Args: []expression.Expr{},
													},
												},
											},
										},
									},
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

func TestFunctions(t *testing.T) {
	testCases := []parserTestCase{
		{
			Name:  "should_parse_no_arg_function",
			Query: `myFunction() = 3`,
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left: &expression.Term{
						FunctionCall: &expression.FunctionCall{
							Name: "myFunction",
							Args: nil,
						},
					},
				},
				Comparator: "=",
				Right:      numberExpr(3),
			},
		},
		{
			Name:  "should_parse_one_arg_function",
			Query: `getRandomString(5) = "abcde"`,
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left: &expression.Term{
						FunctionCall: &expression.FunctionCall{
							Name: "getRandomString",
							Args: []*expression.Term{
								{
									Number: strp("5"),
								},
							},
						},
					},
				},
				Comparator: "=",
				Right:      strExpr("abcde"),
			},
		},
		{
			Name:  "should_parse_one_arg_function",
			Query: `randomInt(5, 10) = 7`,
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left: &expression.Term{
						FunctionCall: &expression.FunctionCall{
							Name: "randomInt",
							Args: []*expression.Term{
								{
									Number: strp("5"),
								},
								{
									Number: strp("10"),
								},
							},
						},
					},
				},
				Comparator: "=",
				Right:      numberExpr(7),
			},
		},
	}

	runTestCases(t, testCases)
}

func TestArrays(t *testing.T) {
	testCases := []parserTestCase{
		{
			Name:  "should_parse_empty_array",
			Query: `[] = []`,
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left: &expression.Term{
						Array: &expression.Array{},
					},
				},
				Comparator: "=",
				Right: &expression.Expr{
					Left: &expression.Term{
						Array: &expression.Array{},
					},
				},
			},
		},
		{
			Name:  "should_parse_single_item_arrays",
			Query: `[2] = [3]`,
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left: &expression.Term{
						Array: &expression.Array{
							Items: []*expression.Term{
								{Number: strp("2")},
							},
						},
					},
				},
				Comparator: "=",
				Right: &expression.Expr{
					Left: &expression.Term{
						Array: &expression.Array{
							Items: []*expression.Term{
								{Number: strp("3")},
							},
						},
					},
				},
			},
		},
		{
			Name:  "should_parse_multiple_items_arrays",
			Query: `[1, 2s, "3"] = ["3", 2s, 1]`,
			ExpectedOutput: expression.Statement{
				Left: &expression.Expr{
					Left: &expression.Term{
						Array: &expression.Array{
							Items: []*expression.Term{
								{Number: strp("1")},
								{Duration: durationp("2s")},
								{Str: &expression.Str{Text: "3", Args: []expression.Expr{}}},
							},
						},
					},
				},
				Comparator: "=",
				Right: &expression.Expr{
					Left: &expression.Term{
						Array: &expression.Array{
							Items: []*expression.Term{
								{Str: &expression.Str{Text: "3", Args: []expression.Expr{}}},
								{Duration: durationp("2s")},
								{Number: strp("1")},
							},
						},
					},
				},
			},
		},
	}

	runTestCases(t, testCases)
}

func runTestCases(t *testing.T, testCases []parserTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output, err := expression.ParseStatement(testCase.Query)
			require.NoError(t, err)
			assert.Equal(t, testCase.ExpectedOutput, output)
		})
	}
}

func strp(in string) *string {
	return &in
}

func durationp(in string) *expression.Duration {
	duration := expression.Duration(in)
	return &duration
}

func numberExpr(number int) *expression.Expr {
	return &expression.Expr{
		Left: &expression.Term{
			Number: strp(fmt.Sprintf("%d", number)),
		},
	}
}

func durationExpr(duration string) *expression.Expr {
	return &expression.Expr{
		Left: &expression.Term{
			Duration: durationp(duration),
		},
	}
}

func attrExpr(attrName string) *expression.Expr {
	return &expression.Expr{
		Left: &expression.Term{
			Attribute: attrp(attrName),
		},
	}
}

func attrp(attrName string) *expression.Attribute {
	attr := expression.NewAttribute(attrName)
	return &attr
}

func strExpr(str string) *expression.Expr {
	s := expression.Str{
		Text: str,
		Args: []expression.Expr{},
	}
	return &expression.Expr{
		Left: &expression.Term{
			Str: &s,
		},
	}
}
