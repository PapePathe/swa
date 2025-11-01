package parser

import (
	"fmt"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseMemberCallExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	tokens := []lexer.Token{}
	tokens = append(tokens, p.expect(lexer.Dot))

	if p.currentToken().Kind == lexer.OpenParen {
		return nil, fmt.Errorf("ParseMemberCallExpression: function calls not yet supported")
	}

	_member, err := parseExpression(p, Member)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("ParseMemberCallExpression %s", left)
	}
}
