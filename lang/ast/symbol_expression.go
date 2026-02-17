package ast

import (
	"swahili/lang/lexer"
)

// SymbolExpression ...
type SymbolExpression struct {
	Value   string
	Tokens  []lexer.Token
	SwaType Type
}

var _ Expression = (*SymbolExpression)(nil)

func (e SymbolExpression) String() string {
	return e.Value
}

func (expr *SymbolExpression) Accept(g CodeGenerator) error {
	return g.VisitSymbolExpression(expr)
}

func (expr SymbolExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr SymbolExpression) VisitedSwaType() Type {
	return expr.SwaType
}

func (expr SymbolExpression) InstructionArg() string {
	return expr.String()
}

type SymbolAdressExpression struct {
	Exp     Expression
	Tokens  []lexer.Token
	SwaType Type
}

var _ Expression = (*SymbolAdressExpression)(nil)

func (expr *SymbolAdressExpression) Accept(g CodeGenerator) error {
	return g.VisitSymbolAdressExpression(expr)
}

func (expr SymbolAdressExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr SymbolAdressExpression) VisitedSwaType() Type {
	return expr.SwaType
}

type SymbolValueExpression struct {
	Exp     Expression
	Tokens  []lexer.Token
	SwaType Type
}

var _ Expression = (*SymbolValueExpression)(nil)

func (s *SymbolValueExpression) Accept(g CodeGenerator) error {
	return g.VisitSymbolValueExpression(s)
}

func (s *SymbolValueExpression) TokenStream() []lexer.Token {
	return s.Tokens
}

func (s *SymbolValueExpression) VisitedSwaType() Type {
	return s.SwaType
}
