package ast

import (
	"swahili/lang/lexer"
)

// StringExpression ...
type StringExpression struct {
	Value   string
	Tokens  []lexer.Token
	SwaType Type
}

var _ Expression = (*StringExpression)(nil)

func (expr StringExpression) String() string {
	return expr.Value
}

func (expr *StringExpression) Accept(g CodeGenerator) error {
	return g.VisitStringExpression(expr)
}

func (expr StringExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr StringExpression) VisitedSwaType() Type {
	return expr.SwaType
}

func (expr StringExpression) InstructionArg() string {
	return expr.Value
}
