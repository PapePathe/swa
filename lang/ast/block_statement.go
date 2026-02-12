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

// FloatingBlockExpression
type FloatingBlockExpression struct {
	Stmt   Statement
	Tokens []lexer.Token
}

var _ Expression = (*FloatingBlockExpression)(nil)

func (f *FloatingBlockExpression) VisitedSwaType() Type {
	panic("unimplemented")
}

func (f *FloatingBlockExpression) Accept(g CodeGenerator) error {
	return g.VisitFloatingBlockExpression(f)
}

func (f *FloatingBlockExpression) TokenStream() []lexer.Token {
	return f.Tokens
}
