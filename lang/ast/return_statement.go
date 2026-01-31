package ast

import (
	"swahili/lang/lexer"
)

type ReturnStatement struct {
	Value  Expression
	Tokens []lexer.Token
}

func (stmt *ReturnStatement) Accept(g CodeGenerator) error {
	return g.VisitReturnStatement(stmt)
}

func (stmt ReturnStatement) TokenStream() []lexer.Token {
	return stmt.Tokens
}
