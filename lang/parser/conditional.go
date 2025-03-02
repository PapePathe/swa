package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

func ParseBlockStatement(p *Parser) ast.BlockStatement {
	body := []ast.Statement{}

	p.expect(lexer.OpenCurly)

	for p.hasTokens() && p.currentToken().Kind != lexer.CloseCurly {
		body = append(body, ParseStatement(p))
	}

	p.expect(lexer.CloseCurly)

	return ast.BlockStatement{
		Body: body,
	}
}

func ParseConditionalExpression(p *Parser) ast.Statement {
	failBlock := ast.BlockStatement{}

	p.expect(lexer.KeywordIf)
	p.expect(lexer.OpenParen)

	condition := parseExpression(p, DefaultBindingPower)

	p.expect(lexer.CloseParen)

	successBlock := ParseBlockStatement(p)

	if p.currentToken().Kind == lexer.KeywordElse {
		p.expect(lexer.KeywordElse)

		failBlock = ParseBlockStatement(p)
	}

	return ast.ConditionalStatetement{
		Condition: condition,
		Success:   successBlock,
		Failure:   failBlock,
	}
}
