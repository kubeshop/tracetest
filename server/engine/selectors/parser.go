//nolint directives: structtag
package selectors

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
)

type Selector struct {
	SpanSelector []SpanSelector `parser:"( @@* ( \",\" @@ )*)"`
}

type SpanSelector struct {
	Filters       []Filter      `parser:"\"span\"\"[\"( @@* ( \",\" @@)*)\"]\""`
	PseudoClass   PseudoClass   `parser:"@@*"`
	ChildSelector *SpanSelector `parser:" @@*"`
}

type Filter struct {
	Property string `parser:"( @Ident ( @\".\" @Ident )*)"`
	Operator string `parser:"@(\"=\" | \"contains\" )"`
	Value    *Value `parser:"@@*"`
}

type Value struct {
	String  *string  `parser:" @String"`
	Int     *int64   `parser:" | @Int"`
	Float   *float64 `parser:" | @Float"`
	Boolean *bool    `parser:" | @(\"true\" | \"false\")"`
	Null    bool     `parser:" | @\"NULL\""`
}

type PseudoClass struct {
	Type  string `parser:"\":\" @(\"nth_child\")"`
	Value *Value `parser:"\"(\" @@* \")\""`
}

func CreateParser() (*participle.Parser, error) {
	parser, err := participle.Build(&Selector{})
	if err != nil {
		return nil, fmt.Errorf("could not create parser: %w", err)
	}

	return parser, nil
}
