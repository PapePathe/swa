package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParsePrefixExpression(p *Parser) (ast.Expression, error) {
	tokens := []lexer.Token{}
	operatorToken := p.advance()
	tokens = append(tokens, operatorToken)
	rightHandSide, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	tokens = append(tokens, rightHandSide.TokenStream()...)

	return ast.PrefixExpression{
		Operator:        operatorToken,
		RightExpression: rightHandSide,
		Tokens:          tokens,
	}, nil
}
