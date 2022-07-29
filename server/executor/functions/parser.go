package functions

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type functionParserObject struct {
	FunctionName string `@Ident`
	Args         []*arg `"(" ( @@ ("," @@ )* )? ")"`
}

type arg struct {
	Number *string ` @Number`
	Bool   *string `| @("true" | "false")`
	String *string `| @(QuotedString|SingleQuotedString)`
}

func (a arg) Value() string {
	value, _ := a.info()
	return value
}

func (a arg) Type() string {
	_, argType := a.info()
	return argType
}

func (a arg) info() (string, string) {
	if a.Number != nil {
		return *a.Number, "number"
	}

	if a.String != nil {
		return *a.String, "string"
	}

	if a.Bool != nil {
		return *a.Bool, "boolean"
	}

	return "", ""
}

var languageLexer = lexer.MustStateful(lexer.Rules{
	"Root": {
		{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z_0-9]*`, Action: nil},
		{Name: "whitespace", Pattern: `\s+`, Action: nil},
		{Name: "Number", Pattern: `([0-9]+(\.[0-9]+)?)`},
		{Name: "QuotedString", Pattern: `".*`, Action: lexer.Push("QuotedString")},
		{Name: "SingleQuotedString", Pattern: `'.*`, Action: lexer.Push("SingleQuotedString")},
		{Name: "Punct", Pattern: `[(),]`, Action: nil},
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
var defaultFunctionRegistry = GetFunctionRegistry()

func createParser() (*participle.Parser, error) {
	if defaultParser != nil {
		return defaultParser, nil
	}

	parser, err := participle.Build(&functionParserObject{}, participle.Lexer(languageLexer))
	if err != nil {
		return nil, fmt.Errorf("could not create parser: %w", err)
	}

	defaultParser = parser

	return defaultParser, nil
}

func ParseFunction(input string) (Function, []FunctionArg, error) {
	parser, err := createParser()
	if err != nil {
		return Function{}, []FunctionArg{}, fmt.Errorf("could not create parser: %w", err)
	}

	var parserObject functionParserObject
	err = parser.ParseString("", input, &parserObject)
	if err != nil {
		return Function{}, []FunctionArg{}, fmt.Errorf(`could not parse input: "%s": %w`, input, err)
	}

	function, err := defaultFunctionRegistry.Get(parserObject.FunctionName)
	if err != nil {
		return Function{}, []FunctionArg{}, fmt.Errorf(`could not find function: %w`, err)
	}

	args := make([]FunctionArg, 0, len(parserObject.Args))
	for _, arg := range parserObject.Args {
		args = append(args, FunctionArg{
			Value: arg.Value(),
			Type:  arg.Type(),
		})
	}

	return function, args, nil
}
