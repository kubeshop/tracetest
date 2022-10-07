package expression

import "fmt"

type executionValue struct {
	Value string
	Type  Type
}

func executeExpression(expr Expr) (string, error) {
	currentValue, _, err := resolveTerm(expr.Left)
	if err != nil {
		return "", fmt.Errorf("could not resolve term: %w", err)
	}

	// executionValue := executionValue{currentValue, currentType}
	// if expr.Right != nil {
	// 	for _, opTerm := range expr.Right {
	// 		currentValue, currentType, err = executeOperation(executionValue, opTerm)
	// 	}
	// }

	return currentValue, nil
}

func resolveTerm(term *Term) (string, Type, error) {
	if term.Attribute != nil {
		// get value from span
		return "value from span", TYPE_ATTRIBUTE, nil
	}

	if term.Duration != nil {
		return *term.Duration, TYPE_DURATION, nil
	}

	if term.Number != nil {
		return *term.Number, TYPE_NUMBER, nil
	}

	if term.Str != nil {
		stringArgs := make([]string, 0, len(term.Str.Args))
		for _, arg := range term.Str.Args {
			stringArg, err := executeExpression(arg)
			if err != nil {
				return "", TYPE_NIL, fmt.Errorf("could not execute expression: %w", err)
			}

			stringArgs = append(stringArgs, stringArg)
		}
		return fmt.Sprintf(term.Str.Text, stringArgs), TYPE_STRING, nil
	}

	return "", TYPE_NIL, fmt.Errorf("empty term")
}
