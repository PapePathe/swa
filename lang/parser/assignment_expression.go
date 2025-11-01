package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseAssignmentExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	tokens := []lexer.Token{}
	operatorToken := p.advance()
	tokens = append(tokens, operatorToken)
	rightHandSide, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	return ast.AssignmentExpression{
		Operator: operatorToken,
		Value:    rightHandSide,
		Assignee: left,
		Tokens:   tokens,
	}, nil
}
