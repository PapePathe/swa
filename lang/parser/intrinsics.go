package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseMakeIntrinsic(p *Parser) (ast.Expression, error) {
	keywordTok := p.advance()

	expr := ast.FunctionCallExpression{}
	p.currentExpression = &expr

	expr.Name = &ast.SymbolExpression{
		Value:  keywordTok.Value,
		Tokens: []lexer.Token{keywordTok},
	}
	expr.Tokens = append(expr.Tokens, keywordTok)
	expr.Tokens = append(expr.Tokens, p.expect(lexer.OpenParen))

	if p.currentToken().Kind == lexer.CloseParen {
		expr.Tokens = append(expr.Tokens, p.expect(lexer.CloseParen))

		return &expr, nil
	}

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		var arg ast.Expression
		var err error

		if len(expr.Args) == 0 {
			typ, toks := parseType(p, DefaultBindingPower)
			arg = &ast.TypeExpression{
				Type:   typ,
				Tokens: toks,
			}
		} else {
			arg, err = parseExpression(p, DefaultBindingPower)
			if err != nil {
				return nil, err
			}
		}

		expr.Args = append(expr.Args, arg)
		expr.Tokens = append(expr.Tokens, arg.TokenStream()...)

		if p.currentToken().Kind == lexer.Comma {
			expr.Tokens = append(expr.Tokens, p.expect(lexer.Comma))
		}
	}

	expr.Tokens = append(expr.Tokens, p.expect(lexer.CloseParen))

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
