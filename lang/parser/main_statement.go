package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseMainStatement(p *Parser) ast.Statement {
	ms := ast.MainStatement{}

	p.expect(lexer.Main)
	p.expect(lexer.OpenParen)
	p.expect(lexer.CloseParen)
	p.expect(lexer.TypeInt)

	ms.Body = ParseBlockStatement(p)

	return ms
}
