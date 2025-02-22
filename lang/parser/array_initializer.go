package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseArrayInitialization(p *Parser) ast.Expression {
	var underlying ast.Type
	var contents = []ast.Expression{}

	p.expect(lexer.OpenBracket)
	p.expect(lexer.CloseBracket)

	underlying = parseType(p, DefaultBindingPower)
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
