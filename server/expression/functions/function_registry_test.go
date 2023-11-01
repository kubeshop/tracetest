package functions_test

import (
	"strconv"
	"testing"

	"github.com/kubeshop/tracetest/server/expression/functions"
	"github.com/kubeshop/tracetest/server/expression/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFunctionValidation(t *testing.T) {
	registry := functions.DefaultRegistry()

	emptyStringFn := func(args ...types.TypedValue) string { return "" }

	assert.Panics(t, func() {
		registry.Add("faulty", emptyStringFn,
			functions.Param(types.TypeString),
			functions.OptionalParam(types.TypeString),
			functions.Param(types.TypeNumber),
		)
	})
}

func TestFunctionWithoutArgs(t *testing.T) {
	registry := functions.DefaultRegistry()

	function, err := registry.Get("uuid")
	require.NoError(t, err)

	result, err := function.Invoke()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestFunctionWithArgs(t *testing.T) {
	registry := functions.DefaultRegistry()

	function, err := registry.Get("randomInt")
	require.NoError(t, err)

	args := []types.TypedValue{
		{
			Value: "1",
			Type:  types.TypeNumber,
		},
		{
			Value: "100",
			Type:  types.TypeNumber,
		},
	}
	result, err := function.Invoke(args...)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	generatedNumber, err := strconv.Atoi(result.Value)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, generatedNumber, 1)
	assert.LessOrEqual(t, generatedNumber, 100)
}

func TestFunctionWithWrongArgNumber(t *testing.T) {
	registry := functions.DefaultRegistry()

	function, err := registry.Get("randomInt")
	require.NoError(t, err)

	_, err = function.Invoke()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required parameters")
}

func TestFunctionWithWrongArgType(t *testing.T) {
	registry := functions.DefaultRegistry()

	function, err := registry.Get("randomInt")
	require.NoError(t, err)

	args := []types.TypedValue{
		{
			Value: "1",
			Type:  types.TypeNumber,
		},
		{
			Value: "string",
			Type:  types.TypeString,
		},
	}

	_, err = function.Invoke(args...)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid argument type")
}

func TestInexistentFunction(t *testing.T) {
	registry := functions.DefaultRegistry()

	_, err := registry.Get("this should not exist!")
	assert.Error(t, err)
}
