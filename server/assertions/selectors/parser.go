package selectors

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type ParserSelector struct {
	SpanSelectors []parserSpanSelector `( @@* ( "," @@ )*)*`
}

type parserSpanSelector struct {
	Filters       []parserFilter      `"span""["( @@* ( "," @@)*)*"]"`
	PseudoClass   parserPseudoClass   `@@*`
	ChildSelector *parserSpanSelector ` @@*`
}

type parserFilter struct {
	Property string       `( @Ident ( @"." @Ident )*)`
	Operator string       `@Comparator`
	Value    *parserValue `@@*`
}

type parserValue struct {
	String  *string  ` @String`
	Float   *float64 ` | @Float`
	Int     *int64   ` | @Int`
	Boolean *bool    ` | @("true" | "false")`
}

type parserPseudoClass struct {
	Type  string       `":" @("nth_child" | "first" | "last")`
	Value *parserValue `("(" @@* ")")*`
}

var languageLexer = lexer.MustStateful(lexer.Rules{
	"Root": {
		{Name: "whitespace", Pattern: `\s+`, Action: nil},
		{Name: "Punc", Pattern: `[(),|\[\]:]`, Action: nil},

		{Name: "Comparator", Pattern: `!=|<=|>=|=|<|>|contains|not-contains`},
		{Name: "Ident", Pattern: `[a-zA-Z][a-zA-Z0-9_\.]*`, Action: nil},
		{Name: "String", Pattern: `"(\\"|[^"])*"`, Action: nil},
		{Name: "Float", Pattern: `[0-9]+\.[0-9]+`},
		{Name: "Int", Pattern: `[0-9]+`},
	},
})

func CreateParser() (*participle.Parser, error) {
	parser, err := participle.Build(&ParserSelector{}, participle.Lexer(languageLexer), participle.UseLookahead(2))
	if err != nil {
		return nil, fmt.Errorf("could not create parser: %w", err)
	}

	return parser, nil
}
