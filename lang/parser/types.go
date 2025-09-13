package parser

import (
	"fmt"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// NudHandlerFunc ...
type TypeNudHandlerFunc func(p *Parser) ast.Type

// LedHandlerFunc ...
type TypeLedHandlerFunc func(p *Parser, left ast.Type, bp BindingPower) ast.Type

// NudLookup ...
type TypeNudLookup map[lexer.TokenKind]TypeNudHandlerFunc

// LedLookup ...
type TypeLedLookup map[lexer.TokenKind]TypeLedHandlerFunc

var (
	typeBindingPowerLookup BpLookup      = make(BpLookup)
	typeNudLookup          TypeNudLookup = make(TypeNudLookup)
	typeLedLookup          TypeLedLookup = make(TypeLedLookup)
)

//func typeLed(kind lexer.TokenKind, bp BindingPower, ledFn TypeLedHandlerFunc) {
//	typeBindingPowerLookup[kind] = bp
//	typeLedLookup[kind] = ledFn
//}

func typeNud(kind lexer.TokenKind, nudFn TypeNudHandlerFunc) {
	typeNudLookup[kind] = nudFn
}

func createTokenTypeLookups() {
	typeNud(lexer.Identifier, parseSymbolType)
	typeNud(lexer.TypeInt, parseIntType)
	typeNud(lexer.TypeString, parseStringType)
	typeNud(lexer.OpenBracket, parseArrayType)
}

func parseStringType(p *Parser) ast.Type {
	p.advance()

	return ast.StringType{}
}

func parseIntType(p *Parser) ast.Type {
	p.advance()

	return ast.NumberType{}
}

func parseSymbolType(p *Parser) ast.Type {
	return ast.SymbolType{
		Name: p.expect(lexer.Identifier).Value,
	}
}

func parseArrayType(p *Parser) ast.Type {
	p.advance()
	p.expect(lexer.CloseBracket)
	underlying := parseType(p, DefaultBindingPower)

	return ast.ArrayType{
		Underlying: underlying,
	}
}

func parseType(p *Parser, bp BindingPower) ast.Type {
	tokenKind := p.currentToken().Kind
	nudFn, exists := typeNudLookup[tokenKind]

	if !exists {
		panic(fmt.Sprintf("type nud handler expected for token kind: %s, value: %s\n", tokenKind, p.currentToken().Value))
	}

	left := nudFn(p)

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

		left = ledFn(p, left, typeBindingPowerLookup[p.currentToken().Kind])
	}

	return left
}
