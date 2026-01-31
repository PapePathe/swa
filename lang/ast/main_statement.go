package ast

import (
	"swahili/lang/lexer"
)

type MainStatement struct {
	Body   BlockStatement
	Tokens []lexer.Token
}

func (stmt *MainStatement) Accept(g CodeGenerator) error {
	return g.VisitMainStatement(stmt)
}

func (stmt MainStatement) TokenStream() []lexer.Token {
	return stmt.Tokens
}
