package ast

import (
	"swahili/lang/lexer"
)

// TODO we should support simpler
// print statements like:
// print("Your int valus is {symbolName}");
type PrintStatetement struct {
	Values []Expression
	Tokens []lexer.Token
}

var _ Statement = (*PrintStatetement)(nil)

func (stmt *PrintStatetement) Accept(g CodeGenerator) error {
	return g.VisitPrintStatement(stmt)
}

func (stmt PrintStatetement) TokenStream() []lexer.Token {
	return stmt.Tokens
}
