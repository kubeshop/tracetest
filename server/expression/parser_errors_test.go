package expression_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/server/expression"
	"github.com/stretchr/testify/assert"
)

func TestStatementParsingErrors(t *testing.T) {
	executor := expression.NewExecutor()
	_, _, err := executor.Statement(`1 1 + = 2`)

	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), `invalid syntax "1 1 + = 2": `))

	unwrappedErr := errors.Unwrap(err)
	assert.False(t, strings.HasPrefix(unwrappedErr.Error(), `invalid syntax "1 1 + = 2": `))
}

func TestExpressionParsingErrors(t *testing.T) {
	executor := expression.NewExecutor()
	_, err := executor.Expression(`attr:attribute env:number`)

	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), `invalid syntax "attr:attribute env:number": `))

	unwrappedErr := errors.Unwrap(err)
	assert.False(t, strings.HasPrefix(unwrappedErr.Error(), `invalid syntax "attr:attribute env:number": `))
}
