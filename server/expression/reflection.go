package expression

import (
	"fmt"
)

type ReflectionToken struct {
	Identifier string
	Type       TermType
}

func GetTokens(statement string) ([]ReflectionToken, error) {
	parsedStatement, err := ParseStatement(statement)
	if err != nil {
		return []ReflectionToken{}, fmt.Errorf("could not parse statement: %w", err)
	}

	leftTokens := extractTokensFromExpression(parsedStatement.Left)
	rightTokens := []ReflectionToken{}
	if parsedStatement.Right != nil {
		rightTokens = extractTokensFromExpression(parsedStatement.Right)
	}

	allTokens := make([]ReflectionToken, 0, len(leftTokens)+len(rightTokens))
	allTokens = append(allTokens, leftTokens...)
	allTokens = append(allTokens, rightTokens...)

	return allTokens, nil
}

func GetTokensFromExpression(expression string) ([]ReflectionToken, error) {
	parsedExpression, err := Parse(expression)
	if err != nil {
		return []ReflectionToken{}, fmt.Errorf("could not parse statement: %w", err)
	}

	return extractTokensFromExpression(&parsedExpression), nil
}

func extractTokensFromExpression(expr *Expr) []ReflectionToken {
	tokens := make([]ReflectionToken, 0)
	tokens = append(tokens, extractTokenFromTerm(expr.Left)...)
	for _, opExpr := range expr.Right {
		rightTokens := extractTokenFromTerm(opExpr.Term)
		tokens = append(tokens, rightTokens...)
	}

	for _, filter := range expr.Filters {
		tokens = append(tokens, ReflectionToken{
			Identifier: filter.Name,
			Type:       FunctionCallType,
		})

		for _, arg := range filter.Args {
			tokens = append(tokens, extractTokenFromTerm(arg)...)
		}
	}

	return tokens
}

func extractTokenFromTerm(term *Term) []ReflectionToken {
	tokens := make([]ReflectionToken, 0)
	tokens = append(tokens, ReflectionToken{
		Identifier: extractIdentifierFromTerm(term),
		Type:       term.Type(),
	})

	termType := term.Type()
	if termType == StrType {
		for _, arg := range term.Str.Args {
			tokens = append(tokens, extractTokensFromExpression(&arg)...)
		}
	}

	return tokens
}

func extractIdentifierFromTerm(term *Term) string {
	if term.Type() == FunctionCallType {
		return term.FunctionCall.Name
	}

	if term.Type() == EnvironmentType {
		return term.Environment.name
	}

	// all other types don't have names, so return an empty string
	return ""
}
