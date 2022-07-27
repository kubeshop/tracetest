package assertions

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
)

type ExpressionOperation interface {
	Execute(model.LiteralValue, model.LiteralValue) (model.LiteralValue, error)
}

type sumOperation struct{}

func (*sumOperation) Execute(value1 model.LiteralValue, value2 model.LiteralValue) (model.LiteralValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 + f2
	})
}

var _ ExpressionOperation = &sumOperation{}

type subtractOperation struct{}

func (*subtractOperation) Execute(value1 model.LiteralValue, value2 model.LiteralValue) (model.LiteralValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 - f2
	})
}

var _ ExpressionOperation = &subtractOperation{}

type multiplyOperation struct{}

func (*multiplyOperation) Execute(value1 model.LiteralValue, value2 model.LiteralValue) (model.LiteralValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 * f2
	})
}

var _ ExpressionOperation = &multiplyOperation{}

type divideOperation struct{}

func (*divideOperation) Execute(value1 model.LiteralValue, value2 model.LiteralValue) (model.LiteralValue, error) {
	return runMathOperationOnNumbers(value1, value2, func(f1, f2 float64) float64 {
		return f1 / f2
	})
}

var _ ExpressionOperation = &divideOperation{}

func runMathOperationOnNumbers(value1 model.LiteralValue, value2 model.LiteralValue, operation func(float64, float64) float64) (model.LiteralValue, error) {
	if value1.Type != "number" || value2.Type != "number" {
		return model.LiteralValue{}, fmt.Errorf("sum operation is only allowed on numbers")
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
		Type:  "number",
	}, nil
}

func getOperationExecutor(operation string) (ExpressionOperation, error) {
	switch operation {
	case "+":
		return &sumOperation{}, nil
	case "-":
		return &subtractOperation{}, nil
	case "*":
		return &multiplyOperation{}, nil
	case "/":
		return &divideOperation{}, nil
	}

	return nil, fmt.Errorf(`unsupported operation "%s"`, operation)
}

func resolveLiteralValue(literalValue model.LiteralValue, span traces.Span) model.LiteralValue {
	if literalValue.Type == "attribute" {
		value := span.Attributes.Get(literalValue.Value)

		return model.LiteralValue{
			Value: value,
			Type:  getValueType(value),
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
