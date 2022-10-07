package expression

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/traces"
)

var ErrNoMatch error = errors.New("no match")

type Executor struct {
	Stores map[string]DataStore
}

func NewExecutor(attributeDataStore DataStore, metaAttributesDataStore DataStore) Executor {
	return Executor{
		Stores: map[string]DataStore{
			"attr":                     attributeDataStore,
			"tracetest.selected_spans": metaAttributesDataStore,
		},
	}
}

func (e Executor) ExecuteStatement(statement string) (string, string, error) {
	parsedStatement, err := Parse(statement)
	if err != nil {
		return "", "", fmt.Errorf("could not parse statement: %w", err)
	}

	leftValue, leftType, err := e.executeExpression(parsedStatement.Left)
	if err != nil {
		return "", "", fmt.Errorf("could not parse left side expression: %w", err)
	}

	rightValue, rightType, err := e.executeExpression(parsedStatement.Right)
	if err != nil {
		return "", "", fmt.Errorf("could not parse left side expression: %w", err)
	}

	// https://github.com/kubeshop/tracetest/issues/1203
	if leftType == TYPE_DURATION || rightType == TYPE_DURATION {
		leftValue = getRoundedDurationValue(leftValue)
		rightValue = getRoundedDurationValue(rightValue)
	}

	comparatorFunction, err := comparator.DefaultRegistry().Get(parsedStatement.Comparator)
	if err != nil {
		return "", "", fmt.Errorf("comparator not supported: %w", err)
	}

	err = comparatorFunction.Compare(rightValue, leftValue)
	if err == comparator.ErrNoMatch {
		err = ErrNoMatch
	}

	return leftValue, rightValue, err
}

type executionValue struct {
	Value string
	Type  Type
}

func (e Executor) executeExpression(expr Expr) (string, Type, error) {
	currentValue, currentType, err := e.resolveTerm(expr.Left)
	if err != nil {
		return "", TYPE_NIL, fmt.Errorf("could not resolve term: %w", err)
	}

	value := executionValue{currentValue, currentType}
	if expr.Right != nil {
		for _, opTerm := range expr.Right {
			currentValue, currentType, err = e.executeOperation(value, opTerm)
			if err != nil {
				return "", TYPE_NIL, fmt.Errorf("could not execute operation: %w", err)
			}

			value = executionValue{currentValue, currentType}
		}
	}

	if expr.Filters != nil {
		for _, filter := range expr.Filters {
			newValue, err := e.executeFilter(value, filter)
			if err != nil {
				return "", TYPE_NIL, fmt.Errorf("could not execute filter: %w", err)
			}

			value = newValue
			currentType = getType(value.Value)
		}
	}

	return value.Value, currentType, nil
}

func (e Executor) resolveTerm(term *Term) (string, Type, error) {
	if term.Attribute != nil {
		if term.Attribute.IsMeta() {
			selectedSpansDataStore := e.Stores["tracetest.selected_spans"]
			value, err := selectedSpansDataStore.Get(term.Attribute.Name())
			if err != nil {
				return "", TYPE_NIL, fmt.Errorf("could not resolve meta attribute: %w", err)
			}

			return value, getType(value), nil
		}

		attributeDataStore := e.Stores["attr"]
		value, err := attributeDataStore.Get(term.Attribute.Name())
		if err != nil {
			return "", TYPE_NIL, fmt.Errorf("could not resolve attribute %s: %w", *term.Attribute, err)
		}

		return value, getType(value), nil
	}

	if term.Duration != nil {
		nanoSeconds := traces.ConvertTimeFieldIntoNanoSeconds(*term.Duration)
		return fmt.Sprintf("%d", nanoSeconds), TYPE_DURATION, nil
	}

	if term.Number != nil {
		return *term.Number, TYPE_NUMBER, nil
	}

	if term.Str != nil {
		stringArgs := make([]any, 0, len(term.Str.Args))
		for _, arg := range term.Str.Args {
			stringArg, _, err := e.executeExpression(arg)
			if err != nil {
				return "", TYPE_NIL, fmt.Errorf("could not execute expression: %w", err)
			}

			stringArgs = append(stringArgs, stringArg)
		}

		value := term.Str.Text
		if len(stringArgs) > 0 {
			value = fmt.Sprintf(term.Str.Text, stringArgs...)
		}

		return value, TYPE_STRING, nil
	}

	return "", TYPE_NIL, fmt.Errorf("empty term")
}

func (e Executor) executeOperation(left executionValue, right *OpTerm) (string, Type, error) {
	rightValue, rightType, err := e.resolveTerm(right.Term)
	if err != nil {
		return "", TYPE_NIL, err
	}

	if left.Type != rightType {
		return "", TYPE_NIL, fmt.Errorf("types mismatch")
	}

	operatorFunction, err := getOperationRegistry().Get(right.Operator)
	if err != nil {
		return "", TYPE_NIL, err
	}

	newValue, err := operatorFunction(left, executionValue{rightValue, rightType})
	if err != nil {
		return "", TYPE_NIL, err
	}

	return newValue.Value, newValue.Type, nil
}

func getType(value string) Type {
	numberRegex := regexp.MustCompile(`^([0-9]+(\.[0-9]+)?)$`)
	durationRegex := regexp.MustCompile(`^([0-9]+(\.[0-9]+)?)(ns|us|ms|s|m|h)$`)

	if numberRegex.Match([]byte(value)) {
		return TYPE_NUMBER
	}

	if durationRegex.Match([]byte(value)) {
		return TYPE_DURATION
	}

	return TYPE_STRING
}

func (e Executor) executeFilter(value executionValue, filter *Filter) (executionValue, error) {
	args := make([]string, 0, len(filter.Args))
	for _, arg := range filter.Args {
		resolvedArg, _, err := e.resolveTerm(arg)
		if err != nil {
			return executionValue{}, err
		}

		args = append(args, resolvedArg)
	}

	newValue, err := executeFilter(value.Value, filter.FunctionName, args)
	if err != nil {
		return executionValue{}, err
	}

	return executionValue{
		Value: newValue,
		Type:  getType(newValue),
	}, nil
}

func getRoundedDurationValue(value string) string {
	numberValue, _ := strconv.Atoi(value)
	valueAsDuration := traces.ConvertNanoSecondsIntoProperTimeUnit(numberValue)
	roundedValue := traces.ConvertTimeFieldIntoNanoSeconds(valueAsDuration)

	return fmt.Sprintf("%d", roundedValue)
}
