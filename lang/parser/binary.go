package parser

import (
	"swahili/lang/ast"
)

// ParseBinaryExpression ...
func ParseBinaryExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	expr := ast.BinaryExpression{}
	p.currentExpression = &expr
	expr.Left = left
	expr.Tokens = append(expr.Tokens, left.TokenStream()...)
	operatorToken := p.advance()
	expr.Operator = operatorToken
	expr.Tokens = append(expr.Tokens, operatorToken)

	right, err := parseExpression(p, bp)
	if err != nil {
		return nil, err
	}

	expr.Right = right
	expr.Tokens = append(expr.Tokens, right.TokenStream()...)

	return expr, nil
}
