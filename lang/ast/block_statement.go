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

func (bs BlockStatement) Accept(g CodeGenerator) error {
	return g.VisitBlockStatement(&bs)
}

func (expr BlockStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}
