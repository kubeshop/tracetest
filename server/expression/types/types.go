package types

import "regexp"

type Type uint

const (
	TypeNil Type = iota
	TypeString
	TypeNumber
	TypeAttribute
	TypeDuration
)

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
