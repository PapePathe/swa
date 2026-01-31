package parser

import (
	"swahili/lang/ast"
	"swahili/lang/helpers"
	"swahili/lang/lexer"
)

func ParseStructInstantiationExpression(p *Parser, left ast.Expression, bp BindingPower) (ast.Expression, error) {
	sInitExpr := ast.StructInitializationExpression{}
	p.currentExpression = &sInitExpr
	tok := helpers.ExpectType[*ast.SymbolExpression](left)
	sInitExpr.Name = tok.Value

	sInitExpr.Tokens = append(sInitExpr.Tokens, left.TokenStream()...)
	sInitExpr.Tokens = append(sInitExpr.Tokens, p.expect(lexer.OpenCurly))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		tok := p.expect(lexer.Identifier)
		propertyName := tok.Value
		sInitExpr.Tokens = append(sInitExpr.Tokens, tok)
		sInitExpr.Tokens = append(sInitExpr.Tokens, p.expect(lexer.Colon))

		expr, err := parseExpression(p, Logical)
		if err != nil {
			return nil, err
		}

		sInitExpr.Tokens = append(sInitExpr.Tokens, expr.TokenStream()...)
		sInitExpr.Values = append(sInitExpr.Values, expr)
		sInitExpr.Properties = append(sInitExpr.Properties, propertyName)

		if p.currentToken().Kind != lexer.CloseCurly {
			sInitExpr.Tokens = append(sInitExpr.Tokens, p.expect(lexer.Comma))
		}
	}

	sInitExpr.Tokens = append(sInitExpr.Tokens, p.expect(lexer.CloseCurly))

	p.currentExpression = nil

	return &sInitExpr, nil
}
