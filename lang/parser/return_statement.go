package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseReturnStatement(p *Parser) (ast.Statement, error) {
	tokens := []lexer.Token{}
	tokens = append(tokens, p.expect(lexer.Return))

	rs := ast.ReturnStatement{}
	value, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	tokens = append(tokens, p.expect(lexer.SemiColon))
	rs.Value = value
	rs.Tokens = tokens

	return rs, nil
}
