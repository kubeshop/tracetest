package parser

import (
	"fmt"
	"strconv"

	"github.com/alecthomas/participle/v2"
)

type Assertion struct {
	Attribute string
	Operator  string
	Value     string
}

type assertionParserObject struct {
	Attribute string       `( @Ident ( @"." @Ident )*)`
	Operator  string       `@( "contains" | "<" | ">" | "=" | "!" )+`
	Value     *parserValue `@@*`
}

type parserValue struct {
	String  *string  ` @String`
	Int     *int64   ` | @Int`
	Float   *float64 ` | @Float`
	Boolean *bool    ` | @("true" | "false")`
}

func (v *parserValue) GetString() string {
	if v.String != nil {
		return *v.String
	}

	if v.Int != nil {
		return fmt.Sprintf("%d", *v.Int)
	}

	if v.Float != nil {
		return fmt.Sprintf("%f", *v.Float)
	}

	if v.Boolean != nil {
		return fmt.Sprintf("%t", *v.Boolean)
	}

	return ""
}

var defaultParser *participle.Parser

func createParser() (*participle.Parser, error) {
	if defaultParser != nil {
		return defaultParser, nil
	}

	parser, err := participle.Build(&assertionParserObject{})
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

	value := assertionParserObject.Value.GetString()
	if value[0:1] == "\"" {
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
