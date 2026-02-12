package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseReturnStatement(p *Parser) (ast.Statement, error) {
	old := p.logger.Step("RetStmt")
	defer p.logger.Restore(old)

	stmt := ast.ReturnStatement{}
	p.currentStatement = &stmt
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.Return))

	var expressions []ast.Expression

	for p.hasTokens() &&
		p.currentToken().Kind != lexer.SemiColon &&
		p.currentToken().Kind != lexer.KeywordIf {
		p.trace("currentToken %v", p.currentToken())

		value, err := parseExpression(p, DefaultBindingPower)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, value)
		stmt.Tokens = append(stmt.Tokens, value.TokenStream()...)

		if p.currentToken().Kind == lexer.Comma {
			stmt.Tokens = append(stmt.Tokens, p.expect(lexer.Comma))
		} else {
			break
		}
	}

	if len(expressions) == 0 {
		// return; (void)
		stmt.Value = nil
	} else if len(expressions) == 1 {
		stmt.Value = expressions[0]
	} else {
		stmt.Value = &ast.TupleExpression{Expressions: expressions}
	}

	if p.currentToken().Kind == lexer.KeywordIf {
		p.trace("parsing cond return statement")

		cond := ast.ConditionalStatetement{}
		p.currentStatement = &cond
		cond.Tokens = append(cond.Tokens, p.expect(lexer.KeywordIf))
		cond.Tokens = append(cond.Tokens, p.expect(lexer.OpenParen))

		value, err := parseExpression(p, DefaultBindingPower)
		if err != nil {
			return nil, err
		}

		cond.Condition = value
		cond.Success = ast.BlockStatement{
			Body: []ast.Statement{&stmt},
		}

		cond.Tokens = append(cond.Tokens, p.expect(lexer.CloseParen))
		cond.Tokens = append(cond.Tokens, p.expect(lexer.SemiColon))

		return &cond, nil
	}

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.SemiColon))

	return &stmt, nil
}
