package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseSymbolValueExpression(p *Parser) (ast.Expression, error) {
	expr := ast.SymbolValueExpression{}
	expr.Tokens = append(expr.Tokens, p.expect(lexer.Star))

	inner, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	expr.Tokens = append(expr.Tokens, inner.TokenStream()...)
	expr.Exp = inner

	return &expr, nil
}

func ParseSymbolAddressExpression(p *Parser) (ast.Expression, error) {
	expr := ast.SymbolAdressExpression{}
	expr.Tokens = append(expr.Tokens, p.expect(lexer.Ampersand))

	inner, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	expr.Tokens = append(expr.Tokens, inner.TokenStream()...)
	expr.Exp = inner

	return &expr, nil
}

// ParseBinaryExpression ...
func ParseBinaryExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	expr := ast.BinaryExpression{}
	p.currentExpression = &expr
	expr.Left = left
	expr.Tokens = append(expr.Tokens, left.TokenStream()...)
	operatorToken := p.advance()
	expr.Operator = operatorToken
	expr.Tokens = append(expr.Tokens, operatorToken)

	right, err := parseExpression(p, bp)
	if err != nil {
		return nil, err
	}

	expr.Right = right
	expr.Tokens = append(expr.Tokens, right.TokenStream()...)

	return &expr, nil
}

func ParseZeroExpression(p *Parser) (ast.Expression, error) {
	expr := ast.ZeroExpression{}
	p.currentExpression = &expr

	expr.Tokens = append(expr.Tokens, p.expect(lexer.Zero))
	expr.Tokens = append(expr.Tokens, p.expect(lexer.OpenParen))

	typ, toks := parseType(p, DefaultBindingPower)

	expr.T = typ
	expr.Tokens = append(expr.Tokens, toks...)
	expr.Tokens = append(expr.Tokens, p.expect(lexer.CloseParen))

	return &expr, nil
}
