package comparator

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var (
	Basic = []Comparator{
		Eq,
		Neq,
		Gt,
		Gte,
		Lt,
		Lte,
		Contains,
		NotContains,
		StartsWith,
		EndsWith,
		JsonContains,
	}
)

// Eq
var Eq Comparator = eq{}

type eq struct{}

func (c eq) Compare(expected, actual string) error {
	if actual == expected {
		return nil
	}

	return ErrNoMatch
}

func (c eq) String() string {
	return "="
}

// Neq
var Neq Comparator = neq{}

type neq struct{}

func (c neq) Compare(expected, actual string) error {
	if actual != expected {
		return nil
	}

	return ErrNoMatch
}

func (c neq) String() string {
	return "!="
}

func parseNumber(s string) (int64, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(`cannot parse "%s" as integer`, s)
	}
	return n, nil
}

func parseNumbers(expected, actual string) (e, a int64, err error) {
	e, err = parseNumber(expected)
	if err != nil {
		return
	}

	a, err = parseNumber(actual)
	if err != nil {
		return
	}

	return
}

// Gt
var Gt Comparator = gt{}

type gt struct{}

func (c gt) Compare(expected, actual string) error {
	e, a, err := parseNumbers(expected, actual)
	if err != nil {
		return err
	}
	if a > e {
		return nil
	}

	return ErrNoMatch
}

func (c gt) String() string {
	return ">"
}

// Gte
var Gte Comparator = gte{}

type gte struct{}

func (c gte) Compare(expected, actual string) error {
	e, a, err := parseNumbers(expected, actual)
	if err != nil {
		return err
	}
	if a >= e {
		return nil
	}

	return ErrNoMatch
}

func (c gte) String() string {
	return ">="
}

// Lt
var Lt Comparator = lt{}

type lt struct{}

func (c lt) Compare(expected, actual string) error {
	e, a, err := parseNumbers(expected, actual)
	if err != nil {
		return err
	}
	if a < e {
		return nil
	}

	return ErrNoMatch
}

func (c lt) String() string {
	return "<"
}

// Lte
var Lte Comparator = lte{}

type lte struct{}

func (c lte) Compare(expected, actual string) error {
	e, a, err := parseNumbers(expected, actual)
	if err != nil {
		return err
	}
	if a <= e {
		return nil
	}

	return ErrNoMatch
}

func (c lte) String() string {
	return "<="
}

// Contains
var Contains Comparator = contains{}

type contains struct{}

func (c contains) Compare(expected, actual string) error {
	if strings.Contains(actual, expected) {
		return nil
	}

	return ErrNoMatch
}

func (c contains) String() string {
	return "contains"
}

// Not contains
var NotContains Comparator = notContains{}

type notContains struct{}

func (c notContains) Compare(expected, actual string) error {
	if strings.Contains(actual, expected) {
		return ErrNoMatch
	}

	return nil
}

func (c notContains) String() string {
	return "not-contains"
}

// StartsWith
var StartsWith Comparator = startsWith{}

type startsWith struct{}

func (c startsWith) Compare(expected, actual string) error {
	if strings.HasPrefix(actual, expected) {
		return nil
	}

	return ErrNoMatch
}

func (c startsWith) String() string {
	return "startsWith"
}

// EndsWith
var EndsWith Comparator = endsWith{}

type endsWith struct{}

func (c endsWith) Compare(expected, actual string) error {
	if strings.HasSuffix(actual, expected) {
		return nil
	}

	return ErrNoMatch
}

func (c endsWith) String() string {
	return "endsWith"
}

// JsonContains
var JsonContains Comparator = jsonContains{}

type jsonContains struct{}

func (c jsonContains) Compare(right, left string) error {
	supersetMap, err := c.parseJson(left)
	if err != nil {
		return fmt.Errorf("left side is not a JSON object")
	}

	subsetMap, err := c.parseJson(right)
	if err != nil {
		return fmt.Errorf("left side is not a JSON object")
	}

	return c.compare(supersetMap, subsetMap)
}

func (c jsonContains) parseJson(input string) ([]map[string]any, error) {
	if strings.HasPrefix(input, "[") && strings.HasSuffix(input, "]") {
		output := []map[string]any{}
		err := json.Unmarshal([]byte(input), &output)
		if err != nil {
			return []map[string]any{}, fmt.Errorf("invalid JSON array")
		}

		return output, nil
	}

	object := map[string]any{}
	err := json.Unmarshal([]byte(input), &object)
	if err != nil {
		return []map[string]any{}, fmt.Errorf("invalid JSON object")
	}

	return []map[string]any{object}, nil
}

func (c jsonContains) compare(left, right []map[string]any) error {
	for i, expected := range right {
		err := c.anyMatches(left, expected)
		if err != nil {
			return fmt.Errorf("left side array doesn't match item %d of right side: %w", i, ErrNoMatch)
		}
	}

	return nil
}

func (c jsonContains) anyMatches(array []map[string]any, expected map[string]any) error {
	for _, left := range array {
		err := c.compareObjects(left, expected)
		if err == nil {
			return nil
		}
	}

	return ErrNoMatch
}

func (c jsonContains) compareObjects(left, right map[string]any) error {
	for key, value := range right {
		leftValue, ok := left[key]
		if !ok {
			return fmt.Errorf(`field "%s" not found on left side: %w`, key, ErrNoMatch)
		}

		if leftValueMap, ok := leftValue.(map[string]any); ok {
			rightValueMap, ok := value.(map[string]any)
			if !ok {
				return fmt.Errorf(`field: "%s": %w`, key, ErrWrongType)
			}

			err := c.compareObjects(leftValueMap, rightValueMap)
			if err != nil {
				return fmt.Errorf(`field "%s": %w`, key, err)
			}

			return nil
		}

		// This minimizes the change of a panic due to comparing two []interface{} types
		leftValueString, _ := json.Marshal(leftValue)
		rightValueString, _ := json.Marshal(value)

		if string(leftValueString) != string(rightValueString) {
			return fmt.Errorf(`field "%s": %w`, key, ErrNoMatch)
		}
	}

	return nil
}

func (c jsonContains) String() string {
	return "json-contains"
}
