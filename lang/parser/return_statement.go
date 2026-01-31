package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseReturnStatement(p *Parser) (ast.Statement, error) {
	stmt := ast.ReturnStatement{}
	p.currentStatement = &stmt
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.Return))

	value, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	stmt.Value = value
	stmt.Tokens = append(stmt.Tokens, value.TokenStream()...)

	if p.currentToken().Kind == lexer.KeywordIf {
		cond := ast.ConditionalStatetement{}
		p.expect(lexer.KeywordIf)
		p.expect(lexer.OpenParen)
		value, err := parseExpression(p, DefaultBindingPower)
		if err != nil {
			return nil, err
		}

		cond.Condition = value
		cond.Success = ast.BlockStatement{
			Body: []ast.Statement{&stmt},
		}
		p.expect(lexer.CloseParen)
		p.expect(lexer.SemiColon)

		return &cond, nil
	}

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.SemiColon))

	return &stmt, nil
}
