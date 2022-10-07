package expression

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/traces"
)

var ErrNoMatch error = errors.New("no match")

type Executor struct {
	Stores map[string]DataStore
}

func NewExecutor(attributeDataStore DataStore) Executor {
	return Executor{
		Stores: map[string]DataStore{
			"attr": attributeDataStore,
		},
	}
}

func (e Executor) ExecuteStatement(statement string) (string, string, error) {
	parsedStatement, err := Parse(statement)
	if err != nil {
		return "", "", fmt.Errorf("could not parse statement: %w", err)
	}

	leftValue, err := e.executeExpression(parsedStatement.Left)
	if err != nil {
		return "", "", fmt.Errorf("could not parse left side expression: %w", err)
	}

	rightValue, err := e.executeExpression(parsedStatement.Right)
	if err != nil {
		return "", "", fmt.Errorf("could not parse left side expression: %w", err)
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

func (e Executor) executeExpression(expr Expr) (string, error) {
	currentValue, currentType, err := e.resolveTerm(expr.Left)
	if err != nil {
		return "", fmt.Errorf("could not resolve term: %w", err)
	}

	value := executionValue{currentValue, currentType}
	if expr.Right != nil {
		for _, opTerm := range expr.Right {
			currentValue, currentType, err = e.executeOperation(value, opTerm)
			if err != nil {
				return "", fmt.Errorf("could not execute operation: %w", err)
			}

			value = executionValue{currentValue, currentType}
		}
	}

	return currentValue, nil
}

func (e Executor) resolveTerm(term *Term) (string, Type, error) {
	if term.Attribute != nil {
		attributeDataStore := e.Stores["attr"]
		value, err := attributeDataStore.Get(string(*term.Attribute))
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
			stringArg, err := e.executeExpression(arg)
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
	numberRegex := regexp.MustCompile(`([0-9]+(\.[0-9]+)?)`)
	durationRegex := regexp.MustCompile(`([0-9]+(\.[0-9]+)?)(ns|us|ms|s|m|h)`)

	if numberRegex.Match([]byte(value)) {
		return TYPE_NUMBER
	}

	if durationRegex.Match([]byte(value)) {
		return TYPE_DURATION
	}

	return TYPE_STRING
}
