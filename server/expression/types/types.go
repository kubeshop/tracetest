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
)

var typeNames = map[Type]string{
	TypeNil:       "nil",
	TypeString:    "string",
	TypeNumber:    "number",
	TypeAttribute: "attribute",
	TypeDuration:  "duration",
	TypeVariable:  "variable",
}

func GetType(value string) Type {
	numberRegex := regexp.MustCompile(`^([0-9]+(\.[0-9]+)?)$`)
	durationRegex := regexp.MustCompile(`^([0-9]+(\.[0-9]+)?)(ns|us|ms|s|m|h)$`)

	if numberRegex.Match([]byte(value)) {
		return TypeNumber
	}

	if durationRegex.Match([]byte(value)) {
		return TypeDuration
	}

	return TypeString
}

func GetTypedValue(value string) TypedValue {
	return TypedValue{Value: value, Type: GetType(value)}
}

func (t Type) String() string {
	return typeNames[t]
}
