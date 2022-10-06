package expression

import "github.com/alecthomas/participle/v2/lexer"

type Statement struct {
	Left       Expr   `@@`
	Comparator string `@Comparator`
	Right      Expr   `@@`
}

type Expr struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

type OpTerm struct {
	Operator string `@Operator`
	Term     *Term  `@@`
}

type Term struct {
	Duration  *string    `( @Duration `
	Number    *string    `| @Number `
	Attribute *Attribute `| @Attribute `
	Str       *Str       `| @(QuotedString|SingleQuotedString) )`
}

var languageLexer = lexer.MustStateful(lexer.Rules{
	"Root": {
		{Name: "whitespace", Pattern: `\s+`, Action: nil},
		{Name: "punc", Pattern: "\\{\\\"'}", Action: nil},

		{Name: "Comparator", Pattern: `!=|<=|>=|=|<|>|contains|not-contains`},
		{Name: "Operator", Pattern: `(\+|\-|\*|\/)`, Action: nil},

		{Name: "Duration", Pattern: `([0-9]+(\.[0-9]+)?)(ns|us|ms|s|m|h)`},
		{Name: "Number", Pattern: `([0-9]+(\.[0-9]+)?)`},
		{Name: "Attribute", Pattern: `attr:[a-zA-Z_0-9][a-zA-Z_0-9.]*`, Action: nil},
		{Name: "QuotedString", Pattern: `"(\\"|[^"])*"`, Action: nil},
		{Name: "SingleQuotedString", Pattern: `'(\\'|[^'])*'`, Action: nil},
	},
})
