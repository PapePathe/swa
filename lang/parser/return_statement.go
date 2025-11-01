package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseReturnStatement(p *Parser) (ast.Statement, error) {
	p.expect(lexer.Return)

	rs := ast.ReturnStatement{}
	value, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}
	rs.Value = value

	p.expect(lexer.SemiColon)

	return rs, nil
}
