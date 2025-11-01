package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseArrayAccess(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	expr := ast.ArrayAccessExpression{Name: left}

	p.expect(lexer.OpenBracket)

	index, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}
	expr.Index = index

	p.expect(lexer.CloseBracket)

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
	contents := []ast.Expression{}
	underlying := parseType(p, DefaultBindingPower)

	p.expect(lexer.OpenCurly)

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		expr, err := parseExpression(p, Logical)
		if err != nil {
			return nil, err
		}
		contents = append(contents, expr)

		if p.currentToken().Kind != lexer.CloseCurly {
			p.expect(lexer.Comma)
		}
	}

	p.expect(lexer.CloseCurly)

	return ast.ArrayInitializationExpression{
		Underlying: underlying,
		Contents:   contents,
	}, nil
}
