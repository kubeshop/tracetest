package expression

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/expression/functions"
	"github.com/kubeshop/tracetest/server/expression/types"
	"github.com/kubeshop/tracetest/server/expression/value"
	"github.com/kubeshop/tracetest/server/traces"
)

var ErrNoMatch error = errors.New("no match")

type Executor struct {
	Stores map[string]DataStore
}

func NewExecutor(dataStores ...DataStore) Executor {
	storesMap := make(map[string]DataStore, len(dataStores))
	for _, dataStore := range dataStores {
		// we can have nil dataStores from test cases
		if dataStore != nil {
			storesMap[dataStore.Source()] = dataStore
		}
	}

	return Executor{
		Stores: storesMap,
	}
}

func (e Executor) Statement(statement string) (string, string, error) {
	parsedStatement, err := ParseStatement(statement)
	if err != nil {
		return "", "", fmt.Errorf("could not parse statement: %w", err)
	}

	leftValue, err := e.Expression(parsedStatement.Left)
	if err != nil {
		return "", "", fmt.Errorf("could not parse left side expression: %w", err)
	}

	rightValue, err := e.Expression(parsedStatement.Right)
	if err != nil {
		return "", "", fmt.Errorf("could not parse left side expression: %w", err)
	}

	// https://github.com/kubeshop/tracetest/issues/1203
	if leftValue.Value().Type == types.TypeDuration || rightValue.Value().Type == types.TypeDuration {
		leftValue = value.New(types.TypedValue{
			Value: getRoundedDurationValue(leftValue.String()),
			Type:  types.TypeDuration,
		})
		rightValue = value.New(types.TypedValue{
			Value: getRoundedDurationValue(rightValue.String()),
			Type:  types.TypeDuration,
		})
	}

	comparatorFunction, err := comparator.DefaultRegistry().Get(parsedStatement.Comparator)
	if err != nil {
		return "", "", fmt.Errorf("comparator not supported: %w", err)
	}

	err = comparatorFunction.Compare(rightValue.String(), leftValue.String())
	if err == comparator.ErrNoMatch {
		err = ErrNoMatch
	}

	if leftValue.Value().Type == types.TypeDuration || rightValue.Value().Type == types.TypeDuration {
		// If any of the sides is a duration, there's a high change of the other side
		// to be a duration as well. So try to format both before returning it
		leftValue = value.NewFromString(maybeFormatDuration(leftValue))
		rightValue = value.NewFromString(maybeFormatDuration(rightValue))
	}

	return leftValue.String(), rightValue.String(), err
}

func (e Executor) Expression(expr Expr) (value.Value, error) {
	currentValue, err := e.resolveTerm(expr.Left)
	if err != nil {
		return value.Nil, fmt.Errorf("could not resolve term: %w", err)
	}

	if expr.Right != nil {
		for _, opTerm := range expr.Right {
			newValue, err := e.executeOperation(currentValue.Value(), opTerm)
			if err != nil {
				return value.Nil, fmt.Errorf("could not execute operation: %w", err)
			}

			currentValue = newValue
		}
	}

	if expr.Filters != nil {
		newValue, err := e.executeFilters(currentValue, expr.Filters)
		if err != nil {
			return value.Nil, err
		}

		currentValue = newValue
	}

	return currentValue, nil
}

func (e Executor) resolveTerm(term *Term) (value.Value, error) {
	if term.Attribute != nil {
		return e.resolveAttribute(term.Attribute)
	}

	if term.Variable != nil {
		return e.resolveVariable(term.Variable)
	}

	if term.FunctionCall != nil {
		return e.resolveFunctionCall(term.FunctionCall)
	}

	if term.Array != nil {
		return e.resolveArray(term.Array)
	}

	if term.Duration != nil {
		nanoSeconds := traces.ConvertTimeFieldIntoNanoSeconds(*term.Duration)
		typedValue := types.TypedValue{
			Value: fmt.Sprintf("%d", nanoSeconds),
			Type:  types.TypeDuration,
		}
		return value.New(typedValue), nil
	}

	if term.Number != nil {
		typedValue := types.TypedValue{
			Value: *term.Number,
			Type:  types.TypeNumber,
		}
		return value.New(typedValue), nil
	}

	if term.Str != nil {
		stringArgs := make([]any, 0, len(term.Str.Args))
		for _, arg := range term.Str.Args {
			newValue, err := e.Expression(arg)
			if err != nil {
				return value.Nil, fmt.Errorf("could not execute expression: %w", err)
			}

			stringArgs = append(stringArgs, newValue.String())
		}

		strValue := term.Str.Text
		if len(stringArgs) > 0 {
			strValue = fmt.Sprintf(term.Str.Text, stringArgs...)
		}

		return value.NewFromString(strValue), nil
	}

	return value.Nil, fmt.Errorf("empty term")
}

func (e Executor) resolveAttribute(attribute *Attribute) (value.Value, error) {
	if attribute.IsMeta() {
		selectedSpansDataStore := e.Stores[metaPrefix]
		attributeValue, err := selectedSpansDataStore.Get(attribute.Name())
		if err != nil {
			return value.Nil, fmt.Errorf("could not resolve meta attribute: %w", err)
		}

		return value.NewFromString(attributeValue), nil
	}

	attributeDataStore := e.Stores["attr"]
	attributeValue, err := attributeDataStore.Get(attribute.Name())
	if err != nil {
		return value.Nil, fmt.Errorf("could not resolve attribute %s: %w", attribute.Name(), err)
	}

	return value.NewFromString(attributeValue), nil
}

func (e Executor) resolveVariable(variable *Variable) (value.Value, error) {
	variableDataStore := e.Stores["var"]
	variableValue, err := variableDataStore.Get(variable.Name())
	if err != nil {
		return value.Nil, fmt.Errorf("could not resolve variable %s: %w", variable.Name(), err)
	}

	return value.NewFromString(variableValue), nil
}

func (e Executor) resolveFunctionCall(functionCall *FunctionCall) (value.Value, error) {
	args := make([]types.TypedValue, 0, len(functionCall.Args))
	for i, arg := range functionCall.Args {
		functionValue, err := e.resolveTerm(arg)
		if err != nil {
			return value.Nil, fmt.Errorf("could not execute function %s: invalid argument on index %d: %w", functionCall.Name, i, err)
		}

		args = append(args, functionValue.Value())
	}

	function, err := functions.DefaultRegistry().Get(functionCall.Name)
	if err != nil {
		return value.Nil, fmt.Errorf("function %s doesn't exist", functionCall.Name)
	}

	typedValue, err := function.Invoke(args...)
	return value.New(typedValue), err
}

func (e Executor) resolveArray(array *Array) (value.Value, error) {
	typedValues := make([]types.TypedValue, 0, len(array.Items))
	for index, item := range array.Items {
		termValue, err := e.resolveTerm(item)
		if err != nil {
			return value.Value{}, fmt.Errorf("could not resolve item at index %d: %w", index, err)
		}

		typedValues = append(typedValues, termValue.Value())
	}

	return value.NewArray(typedValues), nil
}

func (e Executor) executeOperation(left types.TypedValue, right *OpTerm) (value.Value, error) {
	rightValue, err := e.resolveTerm(right.Term)
	if err != nil {
		return value.Nil, err
	}

	if left.Type != rightValue.Value().Type {
		return value.Nil, fmt.Errorf("types mismatch")
	}

	operatorFunction, err := getOperationRegistry().Get(right.Operator)
	if err != nil {
		return value.Nil, err
	}

	newValue, err := operatorFunction(left, rightValue.Value())
	if err != nil {
		return value.Nil, err
	}

	return value.New(newValue), nil
}

func (e Executor) executeFilters(input value.Value, f []*Filter) (value.Value, error) {
	filterInput := input
	for _, filter := range f {
		output, err := e.executeFilter(filterInput, filter)
		if err != nil {
			return value.Nil, err
		}

		filterInput = output
	}

	return filterInput, nil
}

func (e Executor) executeFilter(input value.Value, filter *Filter) (value.Value, error) {
	args := make([]string, 0, len(filter.Args))
	for _, arg := range filter.Args {
		resolvedArg, err := e.resolveTerm(arg)
		if err != nil {
			return value.Value{}, err
		}

		args = append(args, resolvedArg.Value().Value)
	}

	newValue, err := executeFilter(input, filter.Name, args)
	if err != nil {
		return value.Value{}, err
	}

	return newValue, nil
}

func getRoundedDurationValue(value string) string {
	numberValue, _ := strconv.Atoi(value)
	valueAsDuration := traces.ConvertNanoSecondsIntoProperTimeUnit(numberValue)
	roundedValue := traces.ConvertTimeFieldIntoNanoSeconds(valueAsDuration)

	return fmt.Sprintf("%d", roundedValue)
}

func maybeFormatDuration(input value.Value) string {
	// Any type other than duration and number is certain to not be a duration field
	// We still try to convert types.TYPE_NUMBER because we store durations as long numbers,
	// so it's worth trying converting it.
	if input.Value().Type != types.TypeDuration && input.Value().Type != types.TypeNumber {
		return input.String()
	}

	intValue, _ := strconv.Atoi(input.String())
	return traces.ConvertNanoSecondsIntoProperTimeUnit(intValue)
}
