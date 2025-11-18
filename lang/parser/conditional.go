package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseBlockStatement(p *Parser) (ast.BlockStatement, error) {
	blockStatement := ast.BlockStatement{}
	body := []ast.Statement{}
	p.currentStatement = &blockStatement

	blockStatement.Tokens = append(blockStatement.Tokens, p.expect(lexer.OpenCurly))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		stmt, err := ParseStatement(p)
		if err != nil {
			return ast.BlockStatement{}, err
		}

		blockStatement.Tokens = append(blockStatement.Tokens, stmt.TokenStream()...)
		body = append(body, stmt)
	}

	blockStatement.Tokens = append(blockStatement.Tokens, p.expect(lexer.CloseCurly))
	blockStatement.Body = body

	return blockStatement, nil
}

func ParseConditionalExpression(p *Parser) (ast.Statement, error) {
	tokens := []lexer.Token{}
	failBlock := ast.BlockStatement{}

	tokens = append(tokens, p.expect(lexer.KeywordIf))
	tokens = append(tokens, p.expect(lexer.OpenParen))

	condition, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		return nil, err
	}

	tokens = append(tokens, condition.TokenStream()...)
	tokens = append(tokens, p.expect(lexer.CloseParen))

	successBlock, err := ParseBlockStatement(p)
	if err != nil {
		return nil, err
	}

	tokens = append(tokens, successBlock.TokenStream()...)

	if p.currentToken().Kind == lexer.KeywordElse {
		tokens = append(tokens, p.expect(lexer.KeywordElse))

		failBlock, err = ParseBlockStatement(p)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, failBlock.TokenStream()...)
	}

	return ast.ConditionalStatetement{
		Condition: condition,
		Success:   successBlock,
		Failure:   failBlock,
		Tokens:    tokens,
	}, nil
}
