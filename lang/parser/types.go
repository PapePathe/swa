package parser

import (
	"fmt"

	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// TypeNudHandlerFunc...
type TypeNudHandlerFunc func(p *Parser) (ast.Type, []lexer.Token)

// TypeLedHandlerFunc...
type TypeLedHandlerFunc func(p *Parser, left ast.Type, bp BindingPower) (ast.Type, []lexer.Token)

// TypeNudLookup...
type TypeNudLookup map[lexer.TokenKind]TypeNudHandlerFunc

// TypeLedLookup...
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
	typeNud(lexer.Star, parsePointerType)
	typeNud(lexer.Identifier, parseSymbolType)
	typeNud(lexer.TypeInt, parseIntType)
	typeNud(lexer.TypeInt64, parseInt64Type)
	typeNud(lexer.TypeFloat, parseFloatType)
	typeNud(lexer.TypeString, parseStringType)
	typeNud(lexer.OpenBracket, parseArrayType)
}

func parsePointerType(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{p.advance()}

	underlying, toks := parseType(p, DefaultBindingPower)
	tokens = append(tokens, toks...)

	return ast.PointerType{Underlying: underlying}, tokens
}

func parseFloatType(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{p.advance()}

	return ast.FloatType{}, tokens
}

func parseStringType(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{p.advance()}

	return ast.StringType{}, tokens
}

func parseInt64Type(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{p.advance()}

	return ast.Number64Type{}, tokens
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
	typ := ast.ArrayType{}
	tokens := []lexer.Token{}
	tokens = append(tokens, p.advance())

	expr, err := parseExpression(p, DefaultBindingPower)
	if err != nil {
		panic(err)
	}

	typ.DynSizeIdentifier = expr

	if tpexpr, ok := expr.(ast.NumberExpression); ok {
		typ.Size = int(tpexpr.Value)
	}

	tokens = append(tokens, p.expect(lexer.CloseBracket))

	underlying, tokens := parseType(p, DefaultBindingPower)
	typ.Underlying = underlying

	return typ, tokens
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
