package parser

import (
	"fmt"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParsePrintStatement(p *Parser) (ast.Statement, error) {
	old := p.logger.Step("PrintStmt")
	defer p.logger.Restore(old)

	stmt := ast.PrintStatetement{}
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.Print))
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.OpenParen))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseParen {
		p.trace("currentToken: %+v", p.currentToken())

		expr, err := parseExpression(p, Logical)
		if err != nil {
			return nil, fmt.Errorf("ParsePrintStatement: %w", err)
		}

		stmt.Tokens = append(stmt.Tokens, expr.TokenStream()...)
		stmt.Values = append(stmt.Values, expr)

		if p.currentToken().Kind == lexer.Comma {
			stmt.Tokens = append(stmt.Tokens, p.expect(lexer.Comma))
		}
	}

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.CloseParen))
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.SemiColon))

	return &stmt, nil
}
