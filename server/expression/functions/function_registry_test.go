package functions_test

import (
	"strconv"
	"testing"

	"github.com/kubeshop/tracetest/server/expression/functions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFunctionWithoutArgs(t *testing.T) {
	registry := functions.GetFunctionRegistry()

	function, err := registry.Get("uuid")
	require.NoError(t, err)

	result, err := function.Invoke()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestFunctionWithArgs(t *testing.T) {
	registry := functions.GetFunctionRegistry()

	function, err := registry.Get("randomInt")
	require.NoError(t, err)

	args := []functions.Arg{
		{
			Value: "1",
			Type:  "number",
		},
		{
			Value: "100",
			Type:  "number",
		},
	}
	result, err := function.Invoke(args...)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	generatedNumber, err := strconv.Atoi(result)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, generatedNumber, 1)
	assert.LessOrEqual(t, generatedNumber, 100)
}

func TestFunctionWithWrongArgNumber(t *testing.T) {
	registry := functions.GetFunctionRegistry()

	function, err := registry.Get("randomInt")
	require.NoError(t, err)

	_, err = function.Invoke()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wrong number of arguments")
}

func TestFunctionWithWrongArgType(t *testing.T) {
	registry := functions.GetFunctionRegistry()

	function, err := registry.Get("randomInt")
	require.NoError(t, err)

	args := []functions.Arg{
		{
			Value: "1",
			Type:  "number",
		},
		{
			Value: "string",
			Type:  "string",
		},
	}

	_, err = function.Invoke(args...)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wrong argument type")
}

func TestInexistentFunction(t *testing.T) {
	registry := functions.GetFunctionRegistry()

	_, err := registry.Get("this should not exist!")
	assert.Error(t, err)
}
