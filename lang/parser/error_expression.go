package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseErrorExpression(p *Parser) (ast.Expression, error) {
	expr := ast.ErrorExpression{}
	p.currentExpression = &expr
	expr.Tokens = append(expr.Tokens, p.expect(lexer.TypeError))
	expr.Tokens = append(expr.Tokens, p.expect(lexer.OpenParen))

	vexpr, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	expr.Exp = vexpr
	expr.Tokens = append(expr.Tokens, vexpr.TokenStream()...)
	expr.Tokens = append(expr.Tokens, p.expect(lexer.CloseParen))

	return &expr, nil
}
