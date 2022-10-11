package expression

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kubeshop/tracetest/server/expression/types"
	"github.com/kubeshop/tracetest/server/traces"
)

type ExpressionOperation func(types.TypedValue, types.TypedValue) (types.TypedValue, error)

func sum(value1 types.TypedValue, value2 types.TypedValue) (types.TypedValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 + f2
	})
}

func subtract(value1 types.TypedValue, value2 types.TypedValue) (types.TypedValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 - f2
	})
}

func multiply(value1 types.TypedValue, value2 types.TypedValue) (types.TypedValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 * f2
	})
}

func divide(value1 types.TypedValue, value2 types.TypedValue) (types.TypedValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 / f2
	})
}

func runMathOperationOnNumbers(value1 types.TypedValue, value2 types.TypedValue, operation func(float64, float64) float64) (types.TypedValue, error) {
	if err := validateFieldType(value1); err != nil {
		return types.TypedValue{}, err
	}

	if err := validateFieldType(value2); err != nil {
		return types.TypedValue{}, err
	}

	operationType := types.TypeNumber
	if value1.Type == types.TypeDuration {
		value1.Value = fmt.Sprintf("%d", traces.ConvertTimeFieldIntoNanoSeconds(value1.Value))
		operationType = types.TypeDuration
	}

	if value2.Type == types.TypeDuration {
		value2.Value = fmt.Sprintf("%d", traces.ConvertTimeFieldIntoNanoSeconds(value2.Value))
		operationType = types.TypeDuration
	}

	number1, _ := strconv.ParseFloat(value1.Value, 64)
	number2, _ := strconv.ParseFloat(value2.Value, 64)

	result := operation(number1, number2)
	resultStr := fmt.Sprintf("%d", int64(result))
	if strings.Contains(value1.Value, ".") || strings.Contains(value2.Value, ".") {
		// float number
		resultStr = fmt.Sprintf("%.2f", result)
	}

	return types.TypedValue{
		Value: resultStr,
		Type:  operationType,
	}, nil
}

func validateFieldType(field types.TypedValue) error {
	if field.Type != types.TypeNumber && field.Type != types.TypeDuration {
		return fmt.Errorf("operation is only allowed on numbers and duration fields")
	}

	return nil
}
