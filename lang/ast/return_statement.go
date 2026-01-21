package ast

import (
	"swahili/lang/lexer"
)

type ReturnStatement struct {
	Value  Expression
	Tokens []lexer.Token
}

func (rs ReturnStatement) Accept(g CodeGenerator) error {
	return g.VisitReturnStatement(&rs)
}

func (expr ReturnStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}
