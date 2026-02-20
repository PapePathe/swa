package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseAssignmentExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	if tuple, ok := left.(*ast.TupleExpression); ok {
		expr := ast.TupleAssignmentExpression{}
		p.currentExpression = &expr
		expr.Assignees = tuple
		operatorToken := p.advance()
		expr.Operator = operatorToken
		expr.Tokens = append(expr.Tokens, operatorToken)
		expr.Tokens = append(expr.Tokens, left.TokenStream()...)

		rightHandSide, err := parseExpression(p, DefaultBindingPower)
		if err != nil {
			return nil, err
		}

		expr.Value = rightHandSide
		expr.Tokens = append(expr.Tokens, rightHandSide.TokenStream()...)

		return &expr, nil
	}

	expr := ast.AssignmentExpression{}
	p.currentExpression = &expr
	expr.Assignee = left
	operatorToken := p.advance()
	expr.Operator = operatorToken
	expr.Tokens = append(expr.Tokens, operatorToken)
	expr.Tokens = append(expr.Tokens, left.TokenStream()...)

	rightHandSide, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	expr.Value = rightHandSide
	expr.Tokens = append(expr.Tokens, rightHandSide.TokenStream()...)
	expr.Tokens = append(expr.Tokens, p.expect(lexer.SemiColon))

	return &expr, nil
}
