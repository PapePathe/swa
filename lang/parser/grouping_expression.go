package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseGroupingExpression(p *Parser) (ast.Expression, error) {
	p.advance() // move past start of grouping expression

	expr, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}
	p.expect(lexer.CloseParen) // move past end of grouping expression

	return expr, nil
}
