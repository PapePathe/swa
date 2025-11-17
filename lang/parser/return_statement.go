package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseReturnStatement(p *Parser) (ast.Statement, error) {
	stmt := ast.ReturnStatement{}
	p.currentStatement = &stmt
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.Return))

	value, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	stmt.Value = value
	stmt.Tokens = append(stmt.Tokens, value.TokenStream()...)
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.SemiColon))

	return stmt, nil
}
