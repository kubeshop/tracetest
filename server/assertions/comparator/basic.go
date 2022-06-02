package comparator

import (
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
		StartsWith,
		EndsWith,
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
