package parser

import (
	"fmt"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// Parser ...
type Parser struct {
	tokens []lexer.Token
	pos    int
}

// Parse ...
func Parse(tokens []lexer.Token) ast.BlockStatement {
	body := make([]ast.Statement, 0)
	createTokenLookups()
	createTokenTypeLookups()
	psr := &Parser{tokens: tokens}

	for psr.hasTokens() {
		body = append(body, ParseStatement(psr))
	}

	return ast.BlockStatement{
		Body: body,
	}
}

func (p *Parser) currentToken() lexer.Token {
	return p.tokens[p.pos]

}

func (p *Parser) advance() lexer.Token {
	tk := p.currentToken()
	p.pos++

	return tk
}

func (p *Parser) hasTokens() bool {
	return p.pos < len(p.tokens) && p.currentToken().Kind != lexer.EOF
}

func (p *Parser) expectError(kind lexer.TokenKind, err any) lexer.Token {
	token := p.currentToken()

	if kind != token.Kind {
		if err == nil {
			err := fmt.Errorf("expected %s, but got %s", kind, token.Kind)
			panic(err)
		}

	}

	return p.advance()
}

func (p *Parser) expect(kind lexer.TokenKind) lexer.Token {
	return p.expectError(kind, nil)
}
