package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseWhileStatement(p *Parser) (ast.Statement, error) {
	stmt := ast.WhileStatement{}
	p.currentStatement = &stmt
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.KeywordWhile))
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.OpenParen))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		expr, err := parseExpression(p, DefaultBindingPower)
		if err != nil {
			return nil, err
		}
		stmt.Condition = expr
		stmt.Tokens = append(stmt.Tokens, expr.TokenStream()...)
	}

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.CloseParen))

	block, err := ParseBlockStatement(p)
	if err != nil {
		return nil, err
	}

	stmt.Body = block
	stmt.Tokens = append(stmt.Tokens, block.TokenStream()...)

	return stmt, nil
}
