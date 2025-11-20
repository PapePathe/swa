package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseMainStatement(p *Parser) (ast.Statement, error) {
	ms := ast.MainStatement{}
	p.currentStatement = &ms
	ms.Tokens = append(ms.Tokens, p.expect(lexer.Main))
	ms.Tokens = append(ms.Tokens, p.expect(lexer.OpenParen))
	ms.Tokens = append(ms.Tokens, p.expect(lexer.CloseParen))
	ms.Tokens = append(ms.Tokens, p.expect(lexer.TypeInt))

	body, err := ParseBlockStatement(p)
	if err != nil {
		return nil, err
	}

	ms.Body = body
	ms.Tokens = append(ms.Tokens, body.TokenStream()...)

	return ms, nil
}
