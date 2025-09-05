package parser

import (
	"swahili/lang/ast"
	"swahili/lang/helpers"
	"swahili/lang/lexer"
)

func ParseStructInstantiationExpression(
	p *Parser,
	left ast.Expression,
	bp BindingPower,
) ast.Expression {
	structName := helpers.ExpectType[ast.SymbolExpression](left).Value
	properties := map[string]ast.Expression{}

	p.expect(lexer.OpenCurly)

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		propertyName := p.expect(lexer.Identifier).Value
		p.expect(lexer.Colon)
		expr := parseExpression(p, Logical)
		properties[propertyName] = expr

		if p.currentToken().Kind != lexer.CloseCurly {
			p.expect(lexer.Comma)
		}
	}

	p.expect(lexer.CloseCurly)

	return ast.StructInitializationExpression{
		Name:       structName,
		Properties: properties,
	}
}
