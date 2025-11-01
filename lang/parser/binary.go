package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// ParseBinaryExpression ...
func ParseBinaryExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	tokens := []lexer.Token{}
	operatorToken := p.advance()
	tokens = append(tokens, operatorToken)
	right, err := parseExpression(p, bp)
	if err != nil {
		return nil, err
	}

	return ast.BinaryExpression{
		Left:     left,
		Right:    right,
		Operator: operatorToken,
		Tokens:   tokens,
	}, nil
}
