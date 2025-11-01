package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseArrayAccess(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	tokens := []lexer.Token{}
	expr := ast.ArrayAccessExpression{Name: left}

	tokens = append(tokens, left.TokenStream()...)
	tokens = append(tokens, p.expect(lexer.OpenBracket))

	index, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	tokens = append(tokens, index.TokenStream()...)

	expr.Index = index
	tokens = append(tokens, p.expect(lexer.CloseBracket))

	if p.currentToken().Kind == lexer.Dot {
		memberCall, err := ParseMemberCallExpression(p, expr, Member)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, memberCall.TokenStream()...)
		expr.Tokens = tokens

		return memberCall, nil
	}

	expr.Tokens = tokens

	return expr, nil
}

func ParseArrayInitialization(p *Parser) (ast.Expression, error) {
	contents := []ast.Expression{}
	tokens := []lexer.Token{}
	underlying, toks := parseType(p, DefaultBindingPower)

	tokens = append(tokens, toks...)
	tokens = append(tokens, p.expect(lexer.OpenCurly))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		expr, err := parseExpression(p, Logical)
		if err != nil {
			return nil, err
		}
		contents = append(contents, expr)
		tokens = append(tokens, expr.TokenStream()...)

		if p.currentToken().Kind != lexer.CloseCurly {
			tokens = append(tokens, p.expect(lexer.Comma))
		}
	}

	tokens = append(tokens, p.expect(lexer.CloseCurly))

	return ast.ArrayInitializationExpression{
		Underlying: underlying,
		Contents:   contents,
		Tokens:     tokens,
	}, nil
}
