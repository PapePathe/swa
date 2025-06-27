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
	typeNud(lexer.OpenBracket, parseArrayType)
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
		panic(fmt.Sprintf("type nud handler expected for token %s\n", tokenKind))
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
