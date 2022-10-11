package filters

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var arrayRegex = regexp.MustCompile(`^\[[^\]]*\]$`) // anything wrapped with []

func formatOutput(in interface{}) string {
	inType := reflect.TypeOf(in)
	switch inType.Kind() {
	case reflect.Array:
	case reflect.Slice:
		return "array"
	}

	return fmt.Sprintf("%v", in)
}

func formatArray(in []string) string {
	return fmt.Sprintf("[%s]", strings.Join(in, ","))
}

func toStringSlice(in []interface{}) []string {
	out := make([]string, 0, len(in))
	for _, obj := range in {
		out = append(out, fmt.Sprintf("%v", obj))
	}

	return out
}

func isArray(in string) bool {
	return arrayRegex.Match([]byte(in))
}
