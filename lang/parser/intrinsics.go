package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseMakeExpression(p *Parser) (ast.Expression, error) {
	expr := ast.ListExpression{}
	p.currentExpression = &expr
	expr.Tokens = append(expr.Tokens, p.advance())
	expr.Tokens = append(expr.Tokens, p.expect(lexer.OpenParen))

	dtype, tokens := parseType(p, DefaultBindingPower)

	expr.Tokens = append(expr.Tokens, tokens...)
	expr.DataType = dtype

	expr.Tokens = append(expr.Tokens, p.expect(lexer.CloseParen))
	expr.Length = &ast.NumberExpression{Value: 0}
	expr.Capacity = &ast.NumberExpression{Value: 5}

	return &expr, nil
}

func ParseLenIntrinsic(p *Parser) (ast.Expression, error) {
	expr := ast.ListCountExpression{}
	p.currentExpression = &expr

	expr.Tokens = append(expr.Tokens, p.advance())
	expr.Tokens = append(expr.Tokens, p.expect(lexer.OpenParen))

	arg, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	expr.Expr = arg
	expr.Tokens = append(expr.Tokens, arg.TokenStream()...)
	expr.Tokens = append(expr.Tokens, p.expect(lexer.CloseParen))

	return &expr, nil
}
