package comparator

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	Basic = []Comparator{
		Eq,
		Gt,
		Lt,
		Contains,
		StartsWith,
		EndsWith,
	}
)

// Eq
var Eq Comparator = eq{}

type eq struct{}

func (c eq) Compare(expected, actual string) error {
	if expected == actual {
		return nil
	}

	return ErrNoMatch
}

func (c eq) String() string {
	return "="
}

func parseNumber(s string) (int64, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(`cannot parse "%s" as integer`, s)
	}
	return n, nil
}

func parseNumbers(expected, actual string) (a, b int64, err error) {
	a, err = parseNumber(expected)
	if err != nil {
		return
	}

	b, err = parseNumber(actual)
	if err != nil {
		return
	}

	return
}

// Gt
var Gt Comparator = gt{}

type gt struct{}

func (c gt) Compare(expected, actual string) error {
	a, b, err := parseNumbers(expected, actual)
	if err != nil {
		return err
	}
	if a > b {
		return nil
	}

	return ErrNoMatch
}

func (c gt) String() string {
	return ">"
}

// Lt
var Lt Comparator = lt{}

type lt struct{}

func (c lt) Compare(expected, actual string) error {
	expectedNumber, actualNumber, err := parseNumbers(expected, actual)
	if err != nil {
		return err
	}
	if actualNumber < expectedNumber {
		return nil
	}

	return ErrNoMatch
}

func (c lt) String() string {
	return "<"
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
