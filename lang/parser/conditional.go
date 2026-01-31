package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseBlockStatement(p *Parser) (ast.BlockStatement, error) {
	blockStatement := ast.BlockStatement{}
	blockStatement.Tokens = append(blockStatement.Tokens, p.expect(lexer.OpenCurly))
	p.currentStatement = &blockStatement

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		stmt, err := ParseStatement(p)
		if err != nil {
			return ast.BlockStatement{}, err
		}

		blockStatement.Tokens = append(blockStatement.Tokens, stmt.TokenStream()...)
		blockStatement.Body = append(blockStatement.Body, stmt)
	}

	blockStatement.Tokens = append(blockStatement.Tokens, p.expect(lexer.CloseCurly))

	return blockStatement, nil
}

func ParseConditionalExpression(p *Parser) (ast.Statement, error) {
	stmt := ast.ConditionalStatetement{}
	p.currentStatement = &stmt

	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.KeywordIf))
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.OpenParen))

	condition, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	stmt.Condition = condition
	stmt.Tokens = append(stmt.Tokens, condition.TokenStream()...)
	stmt.Tokens = append(stmt.Tokens, p.expect(lexer.CloseParen))

	successBlock, err := ParseBlockStatement(p)
	if err != nil {
		return nil, err
	}

	stmt.Success = successBlock
	stmt.Tokens = append(stmt.Tokens, successBlock.TokenStream()...)

	if p.currentToken().Kind == lexer.KeywordElse {
		stmt.Tokens = append(stmt.Tokens, p.expect(lexer.KeywordElse))

		failBlock, err := ParseBlockStatement(p)
		if err != nil {
			return nil, err
		}

		stmt.Failure = failBlock
		stmt.Tokens = append(stmt.Tokens, failBlock.TokenStream()...)
	}

	return &stmt, nil
}
