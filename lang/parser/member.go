package parser

import (
	"fmt"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseMemberCallExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	switch left.(type) {
	case *ast.ArrayAccessExpression:
		return parseArrayOfStructsAccessExpression(p, left, bp)
	case *ast.SymbolExpression, *ast.MemberExpression,
		*ast.ArrayOfStructsAccessExpression:
		return parseMemberExpression(p, left, bp)
	default:
		return nil, fmt.Errorf("ParseMemberCallExpression expression %s not suppported", left)
	}
}

func parseMemberExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	expr := ast.MemberExpression{}
	p.currentExpression = &expr
	expr.Object = left
	expr.Tokens = append(expr.Tokens, left.TokenStream()...)
	expr.Tokens = append(expr.Tokens, p.expect(lexer.Dot))

	_member, err := parseExpression(p, Member)
	if err != nil {
		return nil, err
	}

	expr.Property = _member
	expr.Tokens = append(expr.Tokens, _member.TokenStream()...)

	return &expr, nil
}

func parseArrayOfStructsAccessExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	expr := ast.ArrayOfStructsAccessExpression{}
	p.currentExpression = &expr
	expr.Tokens = append(expr.Tokens, left.TokenStream()...)
	expr.Tokens = append(expr.Tokens, p.expect(lexer.Dot))
	arr, _ := left.(*ast.ArrayAccessExpression)
	expr.Name = arr.Name
	expr.Index = arr.Index

	_member, err := parseExpression(p, Member)
	if err != nil {
		return nil, err
	}

	expr.Property = _member
	expr.Tokens = append(expr.Tokens, _member.TokenStream()...)

	return &expr, nil
}
