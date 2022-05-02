package selectors

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
)

type Selector struct {
	SpanSelector []SpanSelector `( @@* ( "," @@ )*)`
}

type SpanSelector struct {
	Filters       []Filter      `"span""["( @@* ( "," @@)*)"]"`
	PseudoClass   PseudoClass   `@@*`
	ChildSelector *SpanSelector ` @@*`
}

type Filter struct {
	Property string `( @Ident ( @"." @Ident )*)`
	Operator string `@("=" | "contains" )`
	Value    *Value `@@*`
}

type Value struct {
	String  *string  ` @String`
	Int     *int64   ` | @Int`
	Float   *float64 ` | @Float`
	Boolean *bool    ` | @("true" | "false")`
	Null    bool     ` | @"NULL"`
}

type PseudoClass struct {
	Type  string `":" @("nth_child")`
	Value *Value `"(" @@* ")"`
}

func CreateParser() (*participle.Parser, error) {
	parser, err := participle.Build(&Selector{})
	if err != nil {
		return nil, fmt.Errorf("could not create parser: %w", err)
	}

	return parser, nil
}
