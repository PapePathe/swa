package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseArrayAccess(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	expr := ast.ArrayAccessExpression{Name: left}
	p.currentExpression = &expr

	expr.Tokens = append(expr.Tokens, left.TokenStream()...)
	expr.Tokens = append(expr.Tokens, p.expect(lexer.OpenBracket))

	index, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	expr.Index = index
	expr.Tokens = append(expr.Tokens, index.TokenStream()...)
	expr.Tokens = append(expr.Tokens, p.expect(lexer.CloseBracket))

	if p.currentToken().Kind == lexer.Dot {
		memberCall, err := ParseMemberCallExpression(p, expr, Member)
		if err != nil {
			return nil, err
		}

		return memberCall, nil
	}

	return expr, nil
}

func ParseArrayInitialization(p *Parser) (ast.Expression, error) {
	arrayExpr := ast.ArrayInitializationExpression{}
	p.currentExpression = &arrayExpr
	underlying, toks := parseType(p, DefaultBindingPower)
	arrayExpr.Underlying = underlying
	arrayExpr.Tokens = append(arrayExpr.Tokens, toks...)
	arrayExpr.Tokens = append(arrayExpr.Tokens, p.expect(lexer.OpenCurly))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		expr, err := parseExpression(p, Logical)
		if err != nil {
			return nil, err
		}

		arrayExpr.Tokens = append(arrayExpr.Tokens, expr.TokenStream()...)
		arrayExpr.Contents = append(arrayExpr.Contents, expr)
		arrayExpr.Tokens = append(arrayExpr.Tokens, expr.TokenStream()...)

		if p.currentToken().Kind != lexer.CloseCurly {
			arrayExpr.Tokens = append(arrayExpr.Tokens, p.expect(lexer.Comma))
		}
	}

	arrayExpr.Tokens = append(arrayExpr.Tokens, p.expect(lexer.CloseCurly))

	return arrayExpr, nil
}
