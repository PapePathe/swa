package parser

import (
	"fmt"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseMemberCallExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	p.expect(lexer.Dot)

	if p.currentToken().Kind == lexer.OpenParen {
		return ast.MemberExpression{}, fmt.Errorf("ParseMemberCallExpression: function calls not yet supported")
	}

	_member, err := parseExpression(p, Member)
	if err != nil {
		return ast.MemberExpression{}, err
	}

	switch left.(type) {
	case ast.ArrayAccessExpression:
		arr, _ := left.(ast.ArrayAccessExpression)
		return ast.ArrayOfStructsAccessExpression{
			Property: _member,
			Name:     arr.Name,
			Index:    arr.Index,
		}, nil
	case ast.SymbolExpression:
		return ast.MemberExpression{
			Object:   left,
			Property: _member,
		}, nil
	default:
		return ast.MemberExpression{}, fmt.Errorf("ParseMemberCallExpression %s", left)
	}
}
