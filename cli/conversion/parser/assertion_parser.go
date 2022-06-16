package parser

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Assertion struct {
	Attribute string
	Operator  string
	Value     string
}

type assertionParserObject struct {
	Attribute string `@Attribute`
	Operator  string `@Operator`
	Value     string `@(Number|QuotedString|SingleQuotedString)`
}

var languageLexer = lexer.MustStateful(lexer.Rules{
	"Root": {
		{Name: "Operator", Pattern: `!=|<=|>=|=|<|>|contains`},
		{Name: "Attribute", Pattern: `[a-zA-Z_][a-zA-Z0-9_\.]*`},
		{Name: "whitespace", Pattern: `\s+`, Action: nil},
		{Name: "Number", Pattern: `([0-9]+(\.[0-9]+)?)`},
		{Name: "QuotedString", Pattern: `".*`, Action: lexer.Push("QuotedString")},
		{Name: "SingleQuotedString", Pattern: `'.*`, Action: lexer.Push("SingleQuotedString")},
	},

	"QuotedString": {
		{Name: "EscapedQuote", Pattern: `\\"`, Action: nil},
		{Name: "QuoteStringEnd", Pattern: `"`, Action: lexer.Pop()},
		{Name: "QuotedContent", Pattern: `[^"]*`, Action: nil},
	},

	"SingleQuotedString": {
		{Name: "EscapedSingleQuote", Pattern: `\\'`, Action: nil},
		{Name: "SingleQuoteStringEnd", Pattern: `'`, Action: lexer.Pop()},
		{Name: "SingleQuotedContent", Pattern: `[^']*`, Action: nil},
	},
})

var defaultParser *participle.Parser

func createParser() (*participle.Parser, error) {
	if defaultParser != nil {
		return defaultParser, nil
	}

	parser, err := participle.Build(&assertionParserObject{}, participle.Lexer(languageLexer))
	if err != nil {
		return nil, fmt.Errorf("could not create parser: %w", err)
	}

	defaultParser = parser

	return defaultParser, nil
}

func ParseAssertion(assertionQuery string) (Assertion, error) {
	parser, err := createParser()
	if err != nil {
		return Assertion{}, fmt.Errorf("could not create assertion parser: %w", err)
	}

	var assertionParserObject assertionParserObject
	err = parser.ParseString("", assertionQuery, &assertionParserObject)
	if err != nil {
		return Assertion{}, fmt.Errorf("could not parse assertion (%s): %w", assertionQuery, err)
	}

	value := unquote(assertionParserObject.Value)

	assertion := Assertion{
		Attribute: assertionParserObject.Attribute,
		Operator:  assertionParserObject.Operator,
		Value:     value,
	}

	return assertion, nil
}

func unquote(input string) string {
	isQuoted := input[0:1] == `"` || input[0:1] == `'`
	if !isQuoted {
		return input
	}

	unquotedValue := input[1 : len(input)-1]

	return unescapeQuotes(unquotedValue)
}

func unescapeQuotes(input string) string {
	input = strings.Replace(input, "\\\"", "\"", -1)
	return strings.Replace(input, "\\'", "'", -1)
}
