package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseReturnStatement(p *Parser) (ast.Statement, error) {
	tokens := []lexer.Token{}
	tokens = append(tokens, p.expect(lexer.Return))

	value, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	tokens = append(tokens, value.TokenStream()...)
	tokens = append(tokens, p.expect(lexer.SemiColon))

	return ast.ReturnStatement{Value: value, Tokens: tokens}, nil
}
