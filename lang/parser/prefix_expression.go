package parser

import "swahili/lang/ast"

func ParsePrefixExpression(p *Parser) ast.Expression {
	operatorToken := p.advance()
	rightHandSide := parseExpression(p, DefaultBindingPower)

	return ast.PrefixExpression{
		Operator:        operatorToken,
		RightExpression: rightHandSide,
	}
}
