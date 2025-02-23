package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseMemberCallExpression(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.expect(lexer.Dot)

	if p.currentToken().Kind == lexer.OpenParen {
		panic("function calls not yet supported")
	}

	_member := parseExpression(p, Member)

	return ast.MemberExpression{
		Object:   left,
		Property: _member,
	}
}
