package parser

import (
	"swahili/lang/ast"
	"swahili/lang/helpers"
	"swahili/lang/lexer"
)

func ParseStructInstantiationExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	tok := helpers.ExpectType[ast.SymbolExpression](left)
	structName := tok.Value
	properties := []string{}
	values := []ast.Expression{}
	tokens := []lexer.Token{}

	tokens = append(tokens, left.TokenStream()...)

	tokens = append(tokens, p.expect(lexer.OpenCurly))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		tok := p.expect(lexer.Identifier)
		propertyName := tok.Value
		tokens = append(tokens, tok)
		tokens = append(tokens, p.expect(lexer.Colon))
		expr, err := parseExpression(p, Logical)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, expr.TokenStream()...)
		values = append(values, expr)
		properties = append(properties, propertyName)

		if p.currentToken().Kind != lexer.CloseCurly {
			tokens = append(tokens, p.expect(lexer.Comma))
		}
	}

	tokens = append(tokens, p.expect(lexer.CloseCurly))

	return ast.StructInitializationExpression{
		Name:       structName,
		Properties: properties,
		Values:     values,
		Tokens:     tokens,
	}, nil
}
