package ast

import (
	"swahili/lang/lexer"
)

type ConditionalStatetement struct {
	Condition Expression
	Success   BlockStatement
	Failure   BlockStatement
	Tokens    []lexer.Token
}

var _ Statement = (*ConditionalStatetement)(nil)

func (cs ConditionalStatetement) Accept(g CodeGenerator) error {
	return g.VisitConditionalStatement(&cs)
}

func (expr ConditionalStatetement) TokenStream() []lexer.Token {
	return expr.Tokens
}
