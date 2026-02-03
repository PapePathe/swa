package parser

import (
	"fmt"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseGroupingExpression(p *Parser) (ast.Expression, error) {
	tokens := []lexer.Token{}
	tokens = append(tokens, p.advance())

	var expressions []ast.Expression
	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		expr, err := parseExpression(p, DefaultBindingPower)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
		tokens = append(tokens, expr.TokenStream()...)

		if p.currentToken().Kind == lexer.Comma {
			tokens = append(tokens, p.expect(lexer.Comma))
		} else {
			break
		}
	}

	tokens = append(tokens, p.expect(lexer.CloseParen))

	if len(expressions) == 0 {
		return nil, fmt.Errorf("empty grouping expression")
	}
	if len(expressions) == 1 {
		return expressions[0], nil
	}

	return &ast.TupleExpression{
		Expressions: expressions,
		Tokens:      tokens,
	}, nil
}
