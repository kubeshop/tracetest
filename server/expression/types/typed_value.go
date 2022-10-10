package types

import "fmt"

type TypedValue struct {
	Type  Type
	Value string
}

func (tv TypedValue) FormattedString() string {
	if tv.Type == TYPE_STRING {
		return fmt.Sprintf(`"%s"`, tv.Value)
	}

	return tv.Value
}
