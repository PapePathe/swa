/*
* swahili/lang
* Copyright (C) 2025  Papa Pathe SENE
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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

	if psr.hasTokens() {
		psr.expect(lexer.DialectDeclaration)
		psr.expect(lexer.Colon)
		dlct := psr.expect(lexer.Identifier)
		fmt.Println(dlct)
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
			panic(err)
		}
	}

	return p.advance()
}

func (p *Parser) expect(kind lexer.TokenKind) lexer.Token {
	return p.expectError(kind, nil)
}
