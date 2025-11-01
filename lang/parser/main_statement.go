package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseMainStatement(p *Parser) (ast.Statement, error) {
	ms := ast.MainStatement{}

	p.expect(lexer.Main)
	p.expect(lexer.OpenParen)
	p.expect(lexer.CloseParen)
	p.expect(lexer.TypeInt)

	body, err := ParseBlockStatement(p)
	if err != nil {
		return nil, err
	}

	ms.Body = body

	return ms, nil
}
