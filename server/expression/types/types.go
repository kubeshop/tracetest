package types

import "regexp"

type Type uint

const (
	TypeNil Type = iota
	TypeString
	TypeNumber
	TypeAttribute
	TypeDuration
	TypeVariable
	TypeArray
)

var typeNames = map[Type]string{
	TypeNil:       "nil",
	TypeString:    "string",
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

	if arrayRegex.Match([]byte(value)) {
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
