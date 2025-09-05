package parser

import (
	"fmt"
	"os"
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

	if psr.hasTokens() {
		psr.expect(lexer.DialectDeclaration)
		psr.expect(lexer.Colon)
		psr.expect(lexer.Identifier)
		psr.expect(lexer.SemiColon)
	}

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
			err := fmt.Errorf(
				"expected %s, but got %s  current: %s",
				kind,
				token.Kind,
				token,
			)
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return p.advance()
}

func (p *Parser) expect(kind lexer.TokenKind) lexer.Token {
	return p.expectError(kind, nil)
}
