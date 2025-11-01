package parser

import "swahili/lang/ast"

func ParseAssignmentExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	operatorToken := p.advance()
	rightHandSide, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	return ast.AssignmentExpression{
		Operator: operatorToken,
		Value:    rightHandSide,
		Assignee: left,
	}, nil
}
