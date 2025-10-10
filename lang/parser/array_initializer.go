package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseArrayAccess(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	expr := ast.ArrayAccessExpression{Name: left}

	p.expect(lexer.OpenBracket)

	expr.Index = parseExpression(p, DefaultBindingPower)

	p.expect(lexer.CloseBracket)

	return expr
}

func ParseArrayInitialization(p *Parser) ast.Expression {
	contents := []ast.Expression{}
	underlying := parseType(p, DefaultBindingPower)

	p.expect(lexer.OpenCurly)

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		expr := parseExpression(p, Logical)
		contents = append(contents, expr)

		if p.currentToken().Kind != lexer.CloseCurly {
			p.expect(lexer.Comma)
		}
	}

	p.expect(lexer.CloseCurly)

	return ast.ArrayInitializationExpression{
		Underlying: underlying,
		Contents:   contents,
	}
}
