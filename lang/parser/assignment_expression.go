package parser

import "swahili/lang/ast"

func ParseAssignmentExpression(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	operatorToken := p.advance()
	rightHandSide := parseExpression(p, DefaultBindingPower)

	return ast.AssignmentExpression{
		Operator: operatorToken,
		Value:    rightHandSide,
		Assignee: left,
	}
}
