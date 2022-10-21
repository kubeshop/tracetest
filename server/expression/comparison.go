package expression

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/expression/value"
)

func compare(comparatorName string, leftValue, rightValue value.Value) error {
	comparatorFunction, err := comparator.DefaultRegistry().Get(comparatorName)
	if err != nil {
		return fmt.Errorf("comparator not supported: %w", err)
	}

	if leftValue.IsArray() && comparatorName == "contains" {
		return compareArrayContains(leftValue, rightValue)
	}

	err = comparatorFunction.Compare(rightValue.String(), leftValue.String())
	if err == comparator.ErrNoMatch {
		return ErrNoMatch
	}

	return nil
}

func compareArrayContains(array, expected value.Value) error {
	equalComparator, err := comparator.DefaultRegistry().Get("=")
	if err != nil {
		return fmt.Errorf("equal operator is not supported: %w", err)
	}

	for _, item := range array.Items {
		if err = equalComparator.Compare(item.Value, expected.String()); err == nil {
			return nil
		}
	}

	return ErrNoMatch
}
