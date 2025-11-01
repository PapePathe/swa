package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseBlockStatement(p *Parser) (ast.BlockStatement, error) {
	tokens := []lexer.Token{}
	body := []ast.Statement{}

	tokens = append(tokens, p.expect(lexer.OpenCurly))

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		stmt, err := ParseStatement(p)
		if err != nil {
			return ast.BlockStatement{}, err
		}

		tokens = append(tokens, stmt.TokenStream()...)
		body = append(body, stmt)
	}

	tokens = append(tokens, p.expect(lexer.CloseCurly))

	return ast.BlockStatement{
		Body:   body,
		Tokens: tokens,
	}, nil
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
