package parser

import (
	"fmt"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// NudHandlerFunc ...
type TypeNudHandlerFunc func(p *Parser) (ast.Type, []lexer.Token)

// LedHandlerFunc ...
type TypeLedHandlerFunc func(p *Parser, left ast.Type, bp BindingPower) (ast.Type, []lexer.Token)

// NudLookup ...
type TypeNudLookup map[lexer.TokenKind]TypeNudHandlerFunc

// LedLookup ...
type TypeLedLookup map[lexer.TokenKind]TypeLedHandlerFunc

var (
	typeBindingPowerLookup BpLookup      = make(BpLookup)
	typeNudLookup          TypeNudLookup = make(TypeNudLookup)
	typeLedLookup          TypeLedLookup = make(TypeLedLookup)
)

func typeNud(kind lexer.TokenKind, nudFn TypeNudHandlerFunc) {
	typeNudLookup[kind] = nudFn
}

func createTokenTypeLookups() {
	typeNud(lexer.Identifier, parseSymbolType)
	typeNud(lexer.TypeInt, parseIntType)
	typeNud(lexer.TypeFloat, parseFloatType)
	typeNud(lexer.TypeString, parseStringType)
	typeNud(lexer.OpenBracket, parseArrayType)
}

func parseFloatType(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{p.advance()}

	return ast.FloatType{}, tokens
}

func parseStringType(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{p.advance()}

	return ast.StringType{}, tokens
}

func parseIntType(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{p.advance()}

	return ast.NumberType{}, tokens
}

func parseSymbolType(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{}

	name := p.expect(lexer.Identifier)
	tokens = append(tokens, name)

	return ast.SymbolType{
		Name: name.Value,
	}, tokens
}

func parseArrayType(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{}

	tokens = append(tokens, p.advance())
	tokens = append(tokens, p.expect(lexer.CloseBracket))

	underlying, tokens := parseType(p, DefaultBindingPower)

	return ast.ArrayType{
		Underlying: underlying,
	}, tokens
}

func parseType(p *Parser, bp BindingPower) (ast.Type, []lexer.Token) {
	tokenKind := p.currentToken().Kind
	nudFn, exists := typeNudLookup[tokenKind]
	tokens := []lexer.Token{}

	if !exists {
		panic(fmt.Sprintf("type nud handler expected for token kind: %s, value: %s\n", tokenKind, p.currentToken().Value))
	}

	left, toks := nudFn(p)
	tokens = append(tokens, toks...)

	for typeBindingPowerLookup[p.currentToken().Kind] > bp {
		ledFn, exists := typeLedLookup[p.currentToken().Kind]

		if !exists {
			panic(
				fmt.Sprintf(
					"type led handler expected for token (%s: value(%s))\n",
					tokenKind,
					p.currentToken().Value,
				),
			)
		}

		left, toks = ledFn(p, left, typeBindingPowerLookup[p.currentToken().Kind])
		tokens = append(tokens, toks...)
	}

	return left, tokens
}
