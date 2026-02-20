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
	typeNud(lexer.TypeInt64, parseInt64Type)
	typeNud(lexer.TypeFloat, parseFloatType)
	typeNud(lexer.TypeString, parseStringType)
	typeNud(lexer.TypeError, parseErrorType)
	typeNud(lexer.TypeBool, parseBoolType)
	typeNud(lexer.OpenBracket, parseArrayType)
}

func parsePointerType(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{p.advance()}

	underlying, toks := parseType(p, DefaultBindingPower)
	tokens = append(tokens, toks...)

	return ast.PointerType{Underlying: underlying}, tokens
}

func parseBoolType(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{p.advance()}

	return &ast.BoolType{}, tokens
}

func parseFloatType(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{p.advance()}

	return ast.FloatType{}, tokens
}

func parseStringType(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{p.advance()}

	return ast.StringType{}, tokens
}

func parseErrorType(p *Parser) (ast.Type, []lexer.Token) {
	tokens := []lexer.Token{p.advance()}

	return ast.ErrorType{}, tokens
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

	if p.currentToken().Kind == lexer.Number {
		tok := p.expect(lexer.Number)
		tokens = append(tokens, tok)

		number, err := strconv.ParseInt(tok.Value, 10, 64)
		if err != nil {
			// FIXME return error instead of panic
			panic(err)
		}

		typ.Size = int(number)
	} else if p.currentToken().Kind == lexer.Identifier {
		expr, err := parseExpression(p, DefaultBindingPower)
		if err != nil {
			// FIXME return error instead of panic
			panic(err)
		}

		typ.DynSizeIdentifier = expr
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
		err := fmt.Errorf(
			"type nud handler expected for token kind: %s, value: %s line %d\n",
			tokenKind,
			p.currentToken().Value,
			p.currentToken().Line)
		p.errors = append(p.errors, err)

		return ast.SymbolType{Name: "couldnotparsetype"}, tokens
	}

	left, toks := nudFn(p)
	tokens = append(tokens, toks...)

	for typeBindingPowerLookup[p.currentToken().Kind] > bp {
		ledFn, exists := typeLedLookup[p.currentToken().Kind]

		if !exists {
			// FIXME return error instead of panic
			err := fmt.Errorf(
				"type led handler expected for token (%s: value(%s))\n",
				tokenKind,
				p.currentToken().Value,
			)
			p.errors = append(p.errors, err)
		}

		left, toks = ledFn(p, left, typeBindingPowerLookup[p.currentToken().Kind])
		tokens = append(tokens, toks...)
	}

	return left, tokens
}
