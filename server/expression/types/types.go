package types

import "regexp"

type Type uint

const (
	TYPE_NIL       Type = 0
	TYPE_STRING    Type = 1
	TYPE_NUMBER    Type = 2
	TYPE_ATTRIBUTE Type = 3
	TYPE_DURATION  Type = 4
)

func GetType(value string) Type {
	numberRegex := regexp.MustCompile(`^([0-9]+(\.[0-9]+)?)$`)
	durationRegex := regexp.MustCompile(`^([0-9]+(\.[0-9]+)?)(ns|us|ms|s|m|h)$`)

	if numberRegex.Match([]byte(value)) {
		return TYPE_NUMBER
	}

	if durationRegex.Match([]byte(value)) {
		return TYPE_DURATION
	}

	return TYPE_STRING
}

func GetTypedValue(value string) TypedValue {
	return TypedValue{Value: value, Type: GetType(value)}
}
