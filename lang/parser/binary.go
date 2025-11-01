package parser

import (
	"swahili/lang/ast"
)

// ParseBinaryExpression ...
func ParseBinaryExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	operatorToken := p.advance()
	right, err := parseExpression(p, bp)
	if err != nil {
		return nil, err
	}

	return ast.BinaryExpression{
		Left:     left,
		Right:    right,
		Operator: operatorToken,
	}, nil
}
