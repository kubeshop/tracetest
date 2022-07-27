package assertions

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
)

func ExecuteExpression(expression model.AssertionExpression, span traces.Span) (string, error) {
	if expression.Operation == "" || expression.Expression == nil {
		// No operation to execute, so we just return the literal value
		return getLiteralValue(expression.LiteralValue, span), nil
	}

	operationExecutor, err := getOperationExecutor(expression.Operation)
	if err != nil {
		return "", fmt.Errorf("could not execute expression: %w", err)
	}

	firstValue := resolveLiteralValue(expression.LiteralValue, span)
	secondValue := resolveLiteralValue(expression.Expression.LiteralValue, span)

	literalValue, err := operationExecutor.Execute(firstValue, secondValue)
	if err != nil {
		return "", err
	}

	newExpression := expression.Expression
	newExpression.LiteralValue = literalValue

	return ExecuteExpression(*newExpression, span)
}
