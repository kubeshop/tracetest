package filters

import (
	"fmt"
)

func toStringSlice(in []interface{}) []string {
	out := make([]string, 0, len(in))
	for _, obj := range in {
		out = append(out, fmt.Sprintf("%v", obj))
	}

	return out
}
