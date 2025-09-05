package parser

import (
	"swahili/lang/ast"
)

// ParseBinaryExpression ...
func ParseBinaryExpression(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	operatorToken := p.advance()
	right := parseExpression(p, bp)

	return ast.BinaryExpression{
		Left:     left,
		Right:    right,
		Operator: operatorToken,
	}
}
