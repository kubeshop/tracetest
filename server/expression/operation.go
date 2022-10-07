package expression

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kubeshop/tracetest/server/traces"
)

type ExpressionOperation func(executionValue, executionValue) (executionValue, error)

func sum(value1 executionValue, value2 executionValue) (executionValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 + f2
	})
}

func subtract(value1 executionValue, value2 executionValue) (executionValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 - f2
	})
}

func multiply(value1 executionValue, value2 executionValue) (executionValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 * f2
	})
}

func divide(value1 executionValue, value2 executionValue) (executionValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 / f2
	})
}

func runMathOperationOnNumbers(value1 executionValue, value2 executionValue, operation func(float64, float64) float64) (executionValue, error) {
	if err := validateFieldType(value1); err != nil {
		return executionValue{}, err
	}

	if err := validateFieldType(value2); err != nil {
		return executionValue{}, err
	}

	operationType := TYPE_NUMBER
	if value1.Type == TYPE_DURATION {
		value1.Value = fmt.Sprintf("%d", traces.ConvertTimeFieldIntoNanoSeconds(value1.Value))
		operationType = TYPE_DURATION
	}

	if value2.Type == TYPE_DURATION {
		value2.Value = fmt.Sprintf("%d", traces.ConvertTimeFieldIntoNanoSeconds(value2.Value))
		operationType = TYPE_DURATION
	}

	number1, _ := strconv.ParseFloat(value1.Value, 64)
	number2, _ := strconv.ParseFloat(value2.Value, 64)

	result := operation(number1, number2)
	resultStr := fmt.Sprintf("%d", int64(result))
	if strings.Contains(value1.Value, ".") || strings.Contains(value2.Value, ".") {
		// float number
		resultStr = fmt.Sprintf("%.2f", result)
	}

	return executionValue{
		Value: resultStr,
		Type:  operationType,
	}, nil
}

func validateFieldType(field executionValue) error {
	if field.Type != TYPE_NUMBER && field.Type != TYPE_DURATION {
		return fmt.Errorf("operation is only allowed on numbers and duration fields")
	}

	return nil
}
