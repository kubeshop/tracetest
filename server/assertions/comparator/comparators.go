package comparator

import "fmt"

type Comparator interface {
	Compare(expected, actual string) error
	fmt.Stringer
}

var ErrNoMatch = fmt.Errorf("no match")
