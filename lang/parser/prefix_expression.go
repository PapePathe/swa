package parser

import "swahili/lang/ast"

func ParsePrefixExpression(p *Parser) (ast.Expression, error) {
	operatorToken := p.advance()
	rightHandSide, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	return ast.PrefixExpression{
		Operator:        operatorToken,
		RightExpression: rightHandSide,
	}, nil
}
