package expression

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
)

func ExecuteStatement(statement string) error {
	parsedStatement, err := Parse(statement)
	if err != nil {
		return fmt.Errorf("could not parse statement: %w", err)
	}

	leftValue, err := executeExpression(parsedStatement.Left)
	if err != nil {
		return fmt.Errorf("could not parse left side expression: %w", err)
	}

	rightValue, err := executeExpression(parsedStatement.Right)
	if err != nil {
		return fmt.Errorf("could not parse left side expression: %w", err)
	}

	comparatorFunction, err := comparator.DefaultRegistry().Get(parsedStatement.Comparator)
	if err != nil {
		return fmt.Errorf("comparator not supported: %w", err)
	}

	return comparatorFunction.Compare(rightValue, leftValue)
}
