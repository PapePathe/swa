package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseReturnStatement(p *Parser) ast.Statement {
	p.expect(lexer.Return)

	rs := ast.ReturnStatement{}
	rs.Value = parseExpression(p, DefaultBindingPower)

	p.expect(lexer.SemiColon)

	return rs
}
