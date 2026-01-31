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

func (stmt *ConditionalStatetement) Accept(g CodeGenerator) error {
	return g.VisitConditionalStatement(stmt)
}

func (stmt ConditionalStatetement) TokenStream() []lexer.Token {
	return stmt.Tokens
}
