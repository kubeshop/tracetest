package parser

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Assertion struct {
	Attribute  string
	Operator   string
	Expression *Expression

	// Deprecated: this value is here due to retrocompatibility but it should be removed in our next release.
	// Reason: it was replaced by Expression as it is more flexible and powerful
	Value string
}

type Expression struct {
	LiteralValue ExprLiteral
	Operation    string
	Expression   *Expression
}

type assertionParserObject struct {
	Attribute string `@Attribute`
	Operator  string `@Operator`
	Value     Expr   `@@`
}

type Expr struct {
	Exp1     *ExprLiteral `@@`
	Operator string       `@(ExprOp)?`
	Exp2     *Expr        `@@?`
}

type ExprLiteral struct {
	Attribute    *string `( @Attribute`
	Duration     *string `| @Duration`
	Number       *string `| @Number`
	QuotedString *string `| @(QuotedString|SingleQuotedString) )`
}

func (e ExprLiteral) String() string {
	exprValue, _ := e.info()
	return exprValue
}

func (e ExprLiteral) Type() string {
	_, exprType := e.info()
	return exprType
}

func (e ExprLiteral) info() (string, string) {
	if e.Attribute != nil {
		return *e.Attribute, "attribute"
	}

	if e.Duration != nil {
		return *e.Duration, "duration"
	}

	if e.Number != nil {
		return *e.Number, "number"
	}

	if e.QuotedString != nil {
		return *e.QuotedString, "string"
	}

	return "", ""
}

func (e ExprLiteral) Unquote() ExprLiteral {
	return ExprLiteral{
		Attribute:    unquoteOrNil(e.Attribute),
		Duration:     unquoteOrNil(e.Duration),
		Number:       unquoteOrNil(e.Number),
		QuotedString: unquoteOrNil(e.QuotedString),
	}
}

func unquoteOrNil(in *string) *string {
	if in == nil {
		return nil
	}

	return strp(unquote(*in))
}

func strp(in string) *string {
	return &in
}

func (e Expr) String() string {
	if e.Operator == "" {
		return e.Exp1.String()
	}

	if e.Exp2 == nil {
		return e.Exp1.String()
	}

	return fmt.Sprintf("%s %s %s", e.Exp1.String(), e.Operator, e.Exp2.String())
}

var languageLexer = lexer.MustStateful(lexer.Rules{
	"Root": {
		{Name: "whitespace", Pattern: `\s+`, Action: nil},
		{Name: "Operator", Pattern: `!=|<=|>=|=|<|>|contains|not-contains`},
		{Name: "ExprOp", Pattern: `[\\+|\-|\\*|\/]`, Action: nil},
		{Name: "Attribute", Pattern: `[a-zA-Z_][a-zA-Z0-9_\.]*`},
		{Name: "Duration", Pattern: `([0-9]+(\.[0-9]+)?)(ns|us|ms|s|m|h)`},
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

	value := unquote(assertionParserObject.Value.String())

	assertion := Assertion{
		Attribute:  assertionParserObject.Attribute,
		Operator:   assertionParserObject.Operator,
		Value:      value,
		Expression: createExpression(&assertionParserObject.Value),
	}

	return assertion, nil
}

func createExpression(expr *Expr) *Expression {
	if expr == nil {
		return nil
	}

	return &Expression{
		LiteralValue: expr.Exp1.Unquote(),
		Operation:    expr.Operator,
		Expression:   createExpression(expr.Exp2),
	}
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
