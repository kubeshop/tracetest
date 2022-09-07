package assertions

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
)

type ExpressionOperation func(model.LiteralValue, model.LiteralValue) (model.LiteralValue, error)

func sum(value1 model.LiteralValue, value2 model.LiteralValue) (model.LiteralValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 + f2
	})
}

func subtract(value1 model.LiteralValue, value2 model.LiteralValue) (model.LiteralValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 - f2
	})
}

func multiply(value1 model.LiteralValue, value2 model.LiteralValue) (model.LiteralValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 * f2
	})
}

func divide(value1 model.LiteralValue, value2 model.LiteralValue) (model.LiteralValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 / f2
	})
}

func runMathOperationOnNumbers(value1 model.LiteralValue, value2 model.LiteralValue, operation func(float64, float64) float64) (model.LiteralValue, error) {
	if err := validateFieldType(value1); err != nil {
		return model.LiteralValue{}, err
	}

	if err := validateFieldType(value2); err != nil {
		return model.LiteralValue{}, err
	}

	operationType := "number"
	if value1.Type == "duration" {
		value1.Value = fmt.Sprintf("%d", traces.ConvertTimeFieldIntoNanoSeconds(value1.Value))
		operationType = "duration"
	}

	if value2.Type == "duration" {
		value2.Value = fmt.Sprintf("%d", traces.ConvertTimeFieldIntoNanoSeconds(value2.Value))
		operationType = "duration"
	}

	number1, _ := strconv.ParseFloat(value1.Value, 64)
	number2, _ := strconv.ParseFloat(value2.Value, 64)

	result := operation(number1, number2)
	resultStr := fmt.Sprintf("%d", int64(result))
	if strings.Contains(value1.Value, ".") || strings.Contains(value2.Value, ".") {
		// float number
		resultStr = fmt.Sprintf("%.2f", result)
	}

	return model.LiteralValue{
		Value: resultStr,
		Type:  operationType,
	}, nil
}

func validateFieldType(field model.LiteralValue) error {
	if field.Type != "number" && field.Type != "duration" {
		return fmt.Errorf("operation is only allowed on numbers and duration fields")
	}

	return nil
}

func resolveLiteralValue(literalValue model.LiteralValue, span traces.Span) model.LiteralValue {
	if literalValue.Type == "attribute" {
		value := span.Attributes.Get(literalValue.Value)

		return model.LiteralValue{
			Value: value,
			Type:  getValueType(value),
		}
	}

	if literalValue.Type == "duration" {
		value := traces.ConvertTimeFieldIntoNanoSeconds(literalValue.Value)

		return model.LiteralValue{
			Value: fmt.Sprintf("%d", value),
			Type:  "number",
		}
	}

	return literalValue
}

func getValueType(value string) string {
	numberRegex := regexp.MustCompile(`([0-9]+(\.[0-9]+)?)`)
	durationRegex := regexp.MustCompile(`([0-9]+(\.[0-9]+)?)(ns|us|ms|s|m|h)`)

	if numberRegex.Match([]byte(value)) {
		return "number"
	}

	if durationRegex.Match([]byte(value)) {
		return "duration"
	}

	return "string"
}

func getLiteralValue(literalValue model.LiteralValue, span traces.Span) string {
	resolvedValue := resolveLiteralValue(literalValue, span)
	return resolvedValue.Value
}
