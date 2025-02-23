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

	statement(lexer.Let, ParseVarDeclarationStatement)
	statement(lexer.Const, ParseVarDeclarationStatement)
	statement(lexer.Struct, ParseStructDeclarationStatement)
}
