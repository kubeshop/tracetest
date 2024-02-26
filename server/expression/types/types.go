package types

import (
	"encoding/json"
	"regexp"
)

type Type uint

const (
	TypeNil Type = iota
	TypeString
	TypeNumber
	TypeAttribute
	TypeDuration
	TypeVariable
	TypeArray
	TypeJson
)

var typeNames = map[Type]string{
	TypeNil:       "nil",
	TypeString:    "string",
	TypeJson:      "json",
	TypeNumber:    "number",
	TypeAttribute: "attribute",
	TypeDuration:  "duration",
	TypeVariable:  "variable",
	TypeArray:     "array",
}

func GetType(value string) Type {
	numberRegex := regexp.MustCompile(`^([0-9]+(\.[0-9]+)?)$`)
	durationRegex := regexp.MustCompile(`^([0-9]+(\.[0-9]+)?)(ns|us|ms|s|m|h)$`)
	arrayRegex := regexp.MustCompile(`\[[^\,]*(,[^\],]+)*\]`)

	if numberRegex.Match([]byte(value)) {
		return TypeNumber
	}

	if durationRegex.Match([]byte(value)) {
		return TypeDuration
	}

	if err := json.Unmarshal([]byte(value), &map[string]any{}); err == nil {
		return TypeJson
	}

	if arrayRegex.Match([]byte(value)) {
		var jsonArray = []map[string]any{}
		if err := json.Unmarshal([]byte(value), &jsonArray); err == nil {
			return TypeJson
		}

		return TypeArray
	}

	return TypeString
}

func GetTypedValue(value string) TypedValue {
	return TypedValue{Value: value, Type: GetType(value)}
}

func (t Type) String() string {
	return typeNames[t]
}
