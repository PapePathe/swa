package parser

import (
	"fmt"
	"os"
	"strings"
	"swahili/lang/ast"
	"swahili/lang/lexer"
	"swahili/lang/trc"
)

// Parser ...
type Parser struct {
	tokens            []lexer.Token
	pos               int
	currentStatement  ast.Statement
	currentExpression ast.Expression
	logger            trc.Logger
	tracing           bool
}

// Parse ...
func Parse(tokens []lexer.Token) (ast.BlockStatement, error) {
	body := make([]ast.Statement, 0)

	createTokenLookups()
	createTokenTypeLookups()

	psr := &Parser{
		tokens: tokens,
		logger: *trc.NewLogger("PARSER"),
		//	tracing: true,
		tracing: false,
	}

	if psr.hasTokens() {
		psr.expect(lexer.DialectDeclaration)
		psr.expect(lexer.Colon)
		psr.expect(lexer.Identifier)
		psr.expect(lexer.SemiColon)
	}

	for psr.hasTokens() {
		stmt, err := ParseStatement(psr)
		if err != nil {
			return ast.BlockStatement{}, err
		}

		body = append(body, stmt)
	}

	return ast.BlockStatement{
		Body: body,
	}, nil
}

func (p *Parser) trace(format string, args ...any) {
	if p.tracing {
		p.logger.Debug(fmt.Sprintf(format, args...))
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
			fmt.Println(p.unexpectedTokenError(kind))
			os.Exit(1)
		}

		fmt.Println(err)
		os.Exit(1)
	}

	return p.advance()
}

func (p *Parser) expect(kind lexer.TokenKind) lexer.Token {
	if p.currentStatement != nil {
		err := p.sourceError(kind, p.currentToken(), p.currentStatement.TokenStream())

		return p.expectError(kind, err)
	}

	if p.currentExpression != nil {
		err := p.sourceError(kind, p.currentToken(), p.currentExpression.TokenStream())

		return p.expectError(kind, err)
	}

	return p.expectError(kind, nil)
}

func (p *Parser) unexpectedTokenError(kind lexer.TokenKind) error {
	return fmt.Errorf(
		"expected %s, but got %s at line %d",
		kind,
		p.currentToken().Kind,
		p.currentToken().Line,
	)
}

const (
	ColorReset  = "\x1b[0m"
	ColorRed    = "\x1b[31m"
	ColorGreen  = "\x1b[32m"
	ColorYellow = "\x1b[33m"
	ColorBlue   = "\x1b[34m"
	ColorPurple = "\x1b[35m"
	ColorCyan   = "\x1b[36m"
	ColorWhite  = "\x1b[37m"
)

func (p *Parser) sourceError(kind lexer.TokenKind, token lexer.Token, stream []lexer.Token) error {
	line := 0
	if len(stream) > 0 {
		line = stream[0].Line
	}

	sb := strings.Builder{}

	sb.WriteString(ColorYellow)
	sb.WriteString(p.unexpectedTokenError(kind).Error())
	sb.WriteString(ColorReset)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("\n%s%d%s ", ColorBlue, line, ColorReset))

	for _, tok := range stream {
		if tok.Line > line {
			line = tok.Line
			sb.WriteString(fmt.Sprintf("\n%s%d%s ", ColorBlue, tok.Line, ColorReset))
		}

		sb.WriteString(fmt.Sprintf("%s%s%s ", ColorGreen, tok.Value, ColorReset))
	}

	if token.Line > line {
		sb.WriteString(fmt.Sprintf("\n%s%d%s ", ColorBlue, token.Line, ColorReset))
	}

	sb.WriteString(fmt.Sprintf("%s%s%s ", ColorRed, token.Value, ColorReset))

	return fmt.Errorf(sb.String())
}
