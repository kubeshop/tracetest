package expression

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/expression/filters"
	"github.com/kubeshop/tracetest/server/expression/types"
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
	if leftType == types.TYPE_DURATION || rightType == types.TYPE_DURATION {
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

	if leftType == types.TYPE_DURATION || rightType == types.TYPE_DURATION {
		// If any of the sides is a duration, there's a high change of the other side
		// to be a duration as well. So try to format both before returning it
		leftValue = maybeFormatDuration(leftValue, leftType)
		rightValue = maybeFormatDuration(rightValue, rightType)
	}

	return leftValue, rightValue, err
}

func (e Executor) executeExpression(expr Expr) (string, types.Type, error) {
	currentValue, currentType, err := e.resolveTerm(expr.Left)
	if err != nil {
		return "", types.TYPE_NIL, fmt.Errorf("could not resolve term: %w", err)
	}

	value := types.TypedValue{Value: currentValue, Type: currentType}
	if expr.Right != nil {
		for _, opTerm := range expr.Right {
			currentValue, currentType, err = e.executeOperation(value, opTerm)
			if err != nil {
				return "", types.TYPE_NIL, fmt.Errorf("could not execute operation: %w", err)
			}

			value = types.TypedValue{Value: currentValue, Type: currentType}
		}
	}

	if expr.Filters != nil {
		newValue, err := e.executeFilters(value, expr.Filters)
		if err != nil {
			return "", types.TYPE_NIL, err
		}

		value = newValue
		currentType = types.GetType(value.Value)
	}

	return value.Value, currentType, nil
}

func (e Executor) resolveTerm(term *Term) (string, types.Type, error) {
	if term.Attribute != nil {
		if term.Attribute.IsMeta() {
			selectedSpansDataStore := e.Stores["tracetest.selected_spans"]
			value, err := selectedSpansDataStore.Get(term.Attribute.Name())
			if err != nil {
				return "", types.TYPE_NIL, fmt.Errorf("could not resolve meta attribute: %w", err)
			}

			return value, types.GetType(value), nil
		}

		attributeDataStore := e.Stores["attr"]
		value, err := attributeDataStore.Get(term.Attribute.Name())
		if err != nil {
			return "", types.TYPE_NIL, fmt.Errorf("could not resolve attribute %s: %w", *term.Attribute, err)
		}

		return value, types.GetType(value), nil
	}

	if term.Duration != nil {
		nanoSeconds := traces.ConvertTimeFieldIntoNanoSeconds(*term.Duration)
		return fmt.Sprintf("%d", nanoSeconds), types.TYPE_DURATION, nil
	}

	if term.Number != nil {
		return *term.Number, types.TYPE_NUMBER, nil
	}

	if term.Str != nil {
		stringArgs := make([]any, 0, len(term.Str.Args))
		for _, arg := range term.Str.Args {
			stringArg, _, err := e.executeExpression(arg)
			if err != nil {
				return "", types.TYPE_NIL, fmt.Errorf("could not execute expression: %w", err)
			}

			stringArgs = append(stringArgs, stringArg)
		}

		value := term.Str.Text
		if len(stringArgs) > 0 {
			value = fmt.Sprintf(term.Str.Text, stringArgs...)
		}

		return value, types.TYPE_STRING, nil
	}

	return "", types.TYPE_NIL, fmt.Errorf("empty term")
}

func (e Executor) executeOperation(left types.TypedValue, right *OpTerm) (string, types.Type, error) {
	rightValue, rightType, err := e.resolveTerm(right.Term)
	if err != nil {
		return "", types.TYPE_NIL, err
	}

	if left.Type != rightType {
		return "", types.TYPE_NIL, fmt.Errorf("types mismatch")
	}

	operatorFunction, err := getOperationRegistry().Get(right.Operator)
	if err != nil {
		return "", types.TYPE_NIL, err
	}

	newValue, err := operatorFunction(left, types.TypedValue{Value: rightValue, Type: rightType})
	if err != nil {
		return "", types.TYPE_NIL, err
	}

	return newValue.Value, newValue.Type, nil
}

func (e Executor) executeFilters(value types.TypedValue, f []*Filter) (types.TypedValue, error) {
	filterInput := filters.NewValue(value)

	for _, filter := range f {
		output, err := e.executeFilter(filterInput, filter)
		if err != nil {
			return types.TypedValue{}, err
		}

		filterInput = output
	}

	if filterInput.IsArray() {
		// we don't have to deal with arrays when doing comparisons
		// transform it into a string instead for the simplicity's sake
		return types.GetTypedValue(filterInput.String()), nil
	}

	return filterInput.Value(), nil
}

func (e Executor) executeFilter(input filters.Value, filter *Filter) (filters.Value, error) {
	args := make([]string, 0, len(filter.Args))
	for _, arg := range filter.Args {
		resolvedArg, _, err := e.resolveTerm(arg)
		if err != nil {
			return filters.Value{}, err
		}

		args = append(args, resolvedArg)
	}

	newValue, err := executeFilter(input, filter.FunctionName, args)
	if err != nil {
		return filters.Value{}, err
	}

	return newValue, nil
}

func getRoundedDurationValue(value string) string {
	numberValue, _ := strconv.Atoi(value)
	valueAsDuration := traces.ConvertNanoSecondsIntoProperTimeUnit(numberValue)
	roundedValue := traces.ConvertTimeFieldIntoNanoSeconds(valueAsDuration)

	return fmt.Sprintf("%d", roundedValue)
}

func maybeFormatDuration(value string, vType types.Type) string {
	// Any type other than duration and number is certain to not be a duration field
	// We still try to convert types.TYPE_NUMBER because we store durations as long numbers,
	// so it's worth trying converting it.
	if vType != types.TYPE_DURATION && vType != types.TYPE_NUMBER {
		return value
	}

	intValue, _ := strconv.Atoi(value)
	return traces.ConvertNanoSecondsIntoProperTimeUnit(intValue)
}
