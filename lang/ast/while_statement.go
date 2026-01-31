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

func (stmt *WhileStatement) Accept(g CodeGenerator) error {
	return g.VisitWhileStatement(stmt)
}

func (stmt WhileStatement) TokenStream() []lexer.Token {
	return stmt.Tokens
}
