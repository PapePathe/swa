package parser

import (
	"fmt"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseMemberCallExpression(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.expect(lexer.Dot)

	if p.currentToken().Kind == lexer.OpenParen {
		panic("ParseMemberCallExpression: function calls not yet supported")
	}

	_member := parseExpression(p, Member)

	switch left.(type) {
	case ast.ArrayAccessExpression:
		arr, _ := left.(ast.ArrayAccessExpression)
		return ast.ArrayOfStructsAccessExpression{
			Property: _member,
			Name:     arr.Name,
			Index:    arr.Index,
		}
	case ast.SymbolExpression:
		return ast.MemberExpression{
			Object:   left,
			Property: _member,
		}
	default:
		panic(fmt.Sprintf("Unsupported expression %s", left))
	}
}
