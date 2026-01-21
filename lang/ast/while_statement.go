package ast

import (
	"swahili/lang/lexer"
)

var _ Statement = (*WhileStatement)(nil)

type WhileStatement struct {
	Condition Expression
	Body      BlockStatement
	Tokens    []lexer.Token
}

func (ws WhileStatement) Accept(g CodeGenerator) error {
	return g.VisitWhileStatement(&ws)
}

func (expr WhileStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}
