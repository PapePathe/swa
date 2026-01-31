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
