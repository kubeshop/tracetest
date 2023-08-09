package expression

type Token struct {
	Identifier string
	Type       TermType
}

func GetTokens(statement string) ([]Token, error) {
	parsedStatement, err := ParseStatement(statement)
	if err != nil {
		return []Token{}, invalidSyntaxError(err, statement)
	}

	leftTokens := extractTokensFromExpression(parsedStatement.Left)
	rightTokens := []Token{}
	if parsedStatement.Right != nil {
		rightTokens = extractTokensFromExpression(parsedStatement.Right)
	}

	allTokens := make([]Token, 0, len(leftTokens)+len(rightTokens))
	allTokens = append(allTokens, leftTokens...)
	allTokens = append(allTokens, rightTokens...)

	return allTokens, nil
}

func GetTokensFromExpression(expression string) ([]Token, error) {
	parsedExpression, err := Parse(expression)
	if err != nil {
		return []Token{}, invalidSyntaxError(err, expression)
	}

	return extractTokensFromExpression(&parsedExpression), nil
}

func extractTokensFromExpression(expr *Expr) []Token {
	tokens := make([]Token, 0)
	tokens = append(tokens, extractTokenFromTerm(expr.Left)...)
	for _, opExpr := range expr.Right {
		rightTokens := extractTokenFromTerm(opExpr.Term)
		tokens = append(tokens, rightTokens...)
	}

	for _, filter := range expr.Filters {
		tokens = append(tokens, Token{
			Identifier: filter.Name,
			Type:       FunctionCallType,
		})

		for _, arg := range filter.Args {
			tokens = append(tokens, extractTokenFromTerm(arg)...)
		}
	}

	return tokens
}

func extractTokenFromTerm(term *Term) []Token {
	tokens := make([]Token, 0)
	tokens = append(tokens, Token{
		Identifier: extractIdentifierFromTerm(term),
		Type:       term.Type(),
	})

	termType := term.Type()
	if termType == StrType {
		for _, arg := range term.Str.Args {
			tokens = append(tokens, extractTokensFromExpression(&arg)...)
		}
	}

	return tokens
}

func extractIdentifierFromTerm(term *Term) string {
	if term.Type() == FunctionCallType {
		return term.FunctionCall.Name
	}

	if term.Type() == EnvironmentType {
		return term.Environment.name
	}

	if term.Type() == VariableType {
		return term.Variable.name
	}

	// all other types don't have names, so return an empty string
	return ""
}
