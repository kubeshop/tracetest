package selectors

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
)

type ParserSelector struct {
	SpanSelector []parserSpanSelector `( @@* ( "," @@ )*)`
}

type parserSpanSelector struct {
	Filters       []parserFilter      `"span""["( @@* ( "," @@)*)"]"`
	PseudoClass   parserPseudoClass   `@@*`
	ChildSelector *parserSpanSelector ` @@*`
}

type parserFilter struct {
	Property string       `( @Ident ( @"." @Ident )*)`
	Operator string       `@("=" | "contains" )`
	Value    *parserValue `@@*`
}

type parserValue struct {
	String  *string  ` @String`
	Int     *int64   ` | @Int`
	Float   *float64 ` | @Float`
	Boolean *bool    ` | @("true" | "false")`
	Null    bool     ` | @"NULL"`
}

type parserPseudoClass struct {
	Type  string       `":" @("nth_child")`
	Value *parserValue `"(" @@* ")"`
}

func CreateParser() (*participle.Parser, error) {
	parser, err := participle.Build(&ParserSelector{})
	if err != nil {
		return nil, fmt.Errorf("could not create parser: %w", err)
	}

	return parser, nil
}
