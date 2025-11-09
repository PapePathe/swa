package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseWhileStatement(p *Parser) (ast.Statement, error) {
	stmt := ast.WhileStatement{}
	p.expect(lexer.KeywordWhile)
	p.expect(lexer.OpenParen)

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		expr, err := parseExpression(p, DefaultBindingPower)
		if err != nil {
			return nil, err
		}
		stmt.Condition = expr
	}
	p.expect(lexer.CloseParen)

	block, err := ParseBlockStatement(p)
	if err != nil {
		return nil, err
	}

	stmt.Body = block

	return stmt, nil
}

func ParseMainStatement(p *Parser) (ast.Statement, error) {
	ms := ast.MainStatement{}

	p.expect(lexer.Main)
	p.expect(lexer.OpenParen)
	p.expect(lexer.CloseParen)
	p.expect(lexer.TypeInt)

	body, err := ParseBlockStatement(p)
	if err != nil {
		return nil, err
	}

	ms.Body = body

	return ms, nil
}
