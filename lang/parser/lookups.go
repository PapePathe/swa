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
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

// StatementHandlerFunc ...
type StatementHandlerFunc func(p *Parser) ast.Statement

// NudHandlerFunc ...
type NudHandlerFunc func(p *Parser) ast.Expression

// LedHandlerFunc ...
type LedHandlerFunc func(p *Parser, left ast.Expression, bp BindingPower) ast.Expression

// StatementLookup ...
type StatementLookup map[lexer.TokenKind]StatementHandlerFunc

// NudLookup ...
type NudLookup map[lexer.TokenKind]NudHandlerFunc

// LedLookup ...
type LedLookup map[lexer.TokenKind]LedHandlerFunc

// BpLookup ...
type BpLookup map[lexer.TokenKind]BindingPower

var (
	bindingPowerLookup BpLookup        = make(BpLookup)
	nudLookup          NudLookup       = make(NudLookup)
	ledLookup          LedLookup       = make(LedLookup)
	statementLookup    StatementLookup = make(StatementLookup)
)

func led(kind lexer.TokenKind, bp BindingPower, ledFn LedHandlerFunc) {
	bindingPowerLookup[kind] = bp
	ledLookup[kind] = ledFn
}

func nud(kind lexer.TokenKind, nudFn NudHandlerFunc) {
	nudLookup[kind] = nudFn
}

func statement(kind lexer.TokenKind, statementFn StatementHandlerFunc) {
	bindingPowerLookup[kind] = DefaultBindingPower
	statementLookup[kind] = statementFn
}

func createTokenLookups() {
	led(lexer.Assignment, Assignment, ParseAssignmentExpression)
	led(lexer.PlusEquals, Assignment, ParseAssignmentExpression)
	led(lexer.MinusEquals, Assignment, ParseAssignmentExpression)
	led(lexer.StarEquals, Assignment, ParseAssignmentExpression)

	led(lexer.And, Logical, ParseBinaryExpression)
	led(lexer.Or, Logical, ParseBinaryExpression)

	led(lexer.Equals, Relational, ParseBinaryExpression)
	led(lexer.LessThan, Relational, ParseBinaryExpression)
	led(lexer.LessThanEquals, Relational, ParseBinaryExpression)
	led(lexer.GreaterThan, Relational, ParseBinaryExpression)
	led(lexer.GreaterThanEquals, Relational, ParseBinaryExpression)

	led(lexer.Plus, Additive, ParseBinaryExpression)
	led(lexer.Minus, Additive, ParseBinaryExpression)
	led(lexer.Star, Multiplicative, ParseBinaryExpression)
	led(lexer.Divide, Multiplicative, ParseBinaryExpression)

	nud(lexer.Number, ParsePrimaryExpression)
	nud(lexer.String, ParsePrimaryExpression)
	nud(lexer.Identifier, ParsePrimaryExpression)
	nud(lexer.Minus, ParsePrefixExpression)
	nud(lexer.OpenParen, ParseGroupingExpression)

	led(lexer.Dot, Member, ParseMemberCallExpression)
	led(lexer.OpenCurly, Call, ParseStructInstantiationExpression)
	nud(lexer.OpenBracket, ParseArrayInitialization)

	statement(lexer.Print, ParsePrintStatement)
	statement(lexer.KeywordIf, ParseConditionalExpression)
	statement(lexer.Let, ParseVarDeclarationStatement)
	statement(lexer.Const, ParseVarDeclarationStatement)
	statement(lexer.Struct, ParseStructDeclarationStatement)
}
