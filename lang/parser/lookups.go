package parser

import (
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

type BindingPower int

const (
	DefaultBindingPower BindingPower = iota
	Comma
	Assignment
	Logical
	Relational
	Additive
	Multiplicative
	Unary
	Call
	Member
	Primary
)

type StatementHandlerFunc func(p *Parser) ast.Statement
type NudHandlerFunc func(p *Parser) ast.Expression
type LedHandlerFunc func(p *Parser, left ast.Expression, bp BindingPower) ast.Expression

type StatementLookup map[lexer.TokenKind]StatementHandlerFunc
type NudLookup map[lexer.TokenKind]NudHandlerFunc
type LedLookup map[lexer.TokenKind]LedHandlerFunc
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

func nud(kind lexer.TokenKind, bp BindingPower, nudFn NudHandlerFunc) {
	bindingPowerLookup[kind] = bp
	nudLookup[kind] = nudFn
}
func statement(kind lexer.TokenKind, statementFn StatementHandlerFunc) {
	bindingPowerLookup[kind] = DefaultBindingPower
	statementLookup[kind] = statementFn
}

func createTokenLookups() {
	led(lexer.AND, Logical, ParseBinaryExpression)
	led(lexer.OR, Logical, ParseBinaryExpression)

	led(lexer.LESS_THAN, Relational, ParseBinaryExpression)
	led(lexer.LESS_THAN_EQUALS, Relational, ParseBinaryExpression)
	led(lexer.GREATER_THAN, Relational, ParseBinaryExpression)
	led(lexer.GREATER_THAN_EQUALS, Relational, ParseBinaryExpression)

	led(lexer.PLUS, Additive, ParseBinaryExpression)
	led(lexer.MINUS, Additive, ParseBinaryExpression)

	led(lexer.STAR, Multiplicative, ParseBinaryExpression)
	led(lexer.DIVIDE, Multiplicative, ParseBinaryExpression)

	nud(lexer.NUMBER, Primary, ParsePrimaryExpression)
	nud(lexer.STRING, Primary, ParsePrimaryExpression)
	nud(lexer.IDENTIFIER, Primary, ParsePrimaryExpression)

	statement(lexer.LET, ParseVarDeclarationStatement)
	statement(lexer.CONST, ParseVarDeclarationStatement)
}
