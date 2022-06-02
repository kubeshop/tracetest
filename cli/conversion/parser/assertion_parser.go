package parser

import (
	"fmt"
	"strconv"

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
	Value     string `@Value`
}

var languageLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Operator", Pattern: `!=|<=|>=|=|<|>|contains`},
	{Name: "Attribute", Pattern: `[a-zA-Z_][a-zA-Z0-9_\.]*`},  // anything that starts with a letter
	{Name: "Value", Pattern: `([0-9]+(\.[0-9]+)?)|("[^"]*")`}, // number or quoted string
	{Name: "whitespace", Pattern: `[ |\t]+`},                  // spaces and tabs
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

	value := assertionParserObject.Value
	if value[0:1] == `"` {
		value, err = strconv.Unquote(value)
		if err != nil {
			return Assertion{}, fmt.Errorf("could not unquote value: %w", err)
		}
	}

	assertion := Assertion{
		Attribute: assertionParserObject.Attribute,
		Operator:  assertionParserObject.Operator,
		Value:     value,
	}

	return assertion, nil
}
