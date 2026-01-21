package ast

import (
	"swahili/lang/lexer"
)

type MainStatement struct {
	Body   BlockStatement
	Tokens []lexer.Token
}

func (ms MainStatement) Accept(g CodeGenerator) error {
	return g.VisitMainStatement(&ms)
}

func (expr MainStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}
