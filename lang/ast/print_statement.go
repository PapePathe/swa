package ast

import (
	"swahili/lang/lexer"
)

type PrintStatetement struct {
	Values []Expression
	Tokens []lexer.Token
}

var _ Statement = (*PrintStatetement)(nil)

func (ps PrintStatetement) Accept(g CodeGenerator) error {
	return g.VisitPrintStatement(&ps)
}

func (expr PrintStatetement) TokenStream() []lexer.Token {
	return expr.Tokens
}
