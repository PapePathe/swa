package parser

import (
	"fmt"
	"strconv"

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

	if p.currentToken().Kind == lexer.Number {
		tok := p.expect(lexer.Number)

		number, err := strconv.ParseInt(tok.Value, 10, 64)
		if err != nil {
			p.Errors = append(p.Errors, err)
		}

		typ.Size = int(number)
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
		format := "token kind: %s with value: (%s) is not a valid start token for a type at line %d"
		err := fmt.Errorf(format, tokenKind, p.currentToken().Value, p.currentToken().Line)
		p.Errors = append(p.Errors, err)

		return nil, tokens
	}

	left, toks := nudFn(p)
	tokens = append(tokens, toks...)

	for typeBindingPowerLookup[p.currentToken().Kind] > bp {
		ledFn, exists := typeLedLookup[p.currentToken().Kind]

		if !exists {
			format := "token kind: %s with value: (%s) is not a valid start token for a led type at line %d"
			err := fmt.Errorf(format, tokenKind, p.currentToken().Value, p.currentToken().Line)
			p.Errors = append(p.Errors, err)

			return left, tokens
		}

		left, toks = ledFn(p, left, typeBindingPowerLookup[p.currentToken().Kind])
		tokens = append(tokens, toks...)
	}

	return left, tokens
}
