package ast

import (
	"swahili/lang/lexer"
)

// BlockStatement ...
type BlockStatement struct {
	// The body of the block statement
	Body   []Statement
	Tokens []lexer.Token
}

var _ Statement = (*BlockStatement)(nil)

func (stmt *BlockStatement) Accept(g CodeGenerator) error {
	return g.VisitBlockStatement(stmt)
}

func (stmt BlockStatement) TokenStream() []lexer.Token {
	return stmt.Tokens
}
