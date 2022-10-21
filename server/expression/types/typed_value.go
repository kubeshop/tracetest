package types

import "fmt"

type TypedValue struct {
	Type  Type
	Value string
}

func (tv TypedValue) FormattedString() string {
	if tv.Type == TypeString {
		return fmt.Sprintf(`"%s"`, tv.Value)
	}

	return tv.Value
}
