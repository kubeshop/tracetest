package expression

import (
	"fmt"
	"regexp"
	"strings"
)

type Str struct {
	Text string
	Args []Expr
}

func (s *Str) Capture(in []string) error {
	// Removes the quotes from the string field, so "abc" becomes abc instead.
	input := in[0]
	if len(input) >= 2 {
		input = input[1 : len(input)-1]
	}

	newInput, expressions, err := extractInterpolationArguments(input)
	if err != nil {
		return fmt.Errorf("could not parse interpolated expression: %w", err)
	}

	s.Text = newInput
	s.Args = expressions

	return nil
}

func extractInterpolationArguments(input string) (string, []Expr, error) {
	interpolationRegex := regexp.MustCompile(`\$\{[^}]+\}`)
	interpolationTokens := interpolationRegex.FindAllString(input, -1)

	if interpolationTokens == nil {
		// This string has no interpolation tokens
		return input, []Expr{}, nil
	}

	expressions := make([]Expr, 0, len(interpolationTokens))

	for _, interpolationToken := range interpolationTokens {
		innerExpression := interpolationToken[2 : len(interpolationToken)-1] // removes ${ and }
		expr, err := Parse(innerExpression)
		if err != nil {
			return "", []Expr{}, err
		}

		expressions = append(expressions, expr)
		input = strings.Replace(input, interpolationToken, "%s", 1)
	}

	return input, expressions, nil
}
