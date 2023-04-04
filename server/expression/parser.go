package expression

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
)

var defaultParser *participle.Parser
var defaultExpressionParser *participle.Parser

func createParser() (*participle.Parser, error) {
	if defaultParser != nil {
		return defaultParser, nil
	}

	parser, err := participle.Build(&Statement{}, participle.Lexer(languageLexer), participle.UseLookahead(2))
	if err != nil {
		return nil, fmt.Errorf("could not create parser: %w", err)
	}

	defaultParser = parser

	return defaultParser, nil
}

func ParseStatement(statement string) (Statement, error) {
	var parsedStatement Statement
	parser, err := createParser()
	if err != nil {
		return Statement{}, fmt.Errorf("could not create parser: %w", err)
	}

	err = parser.ParseString("", statement, &parsedStatement)
	if err != nil {
		return Statement{}, invalidSyntaxError(err, statement)
	}

	return parsedStatement, nil
}

func createExpressionParser() (*participle.Parser, error) {
	if defaultExpressionParser != nil {
		return defaultExpressionParser, nil
	}

	parser, err := participle.Build(&Expr{}, participle.Lexer(languageLexer), participle.UseLookahead(2))
	if err != nil {
		return nil, fmt.Errorf("could not create parser: %w", err)
	}

	defaultExpressionParser = parser

	return defaultExpressionParser, nil
}

func Parse(expression string) (Expr, error) {
	var parsedExpression Expr
	parser, err := createExpressionParser()
	if err != nil {
		return Expr{}, fmt.Errorf("could not create parser: %w", err)
	}

	err = parser.ParseString("", expression, &parsedExpression)
	if err != nil {
		return Expr{}, invalidSyntaxError(err, expression)
	}

	return parsedExpression, nil
}
