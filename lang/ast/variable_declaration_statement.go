package ast

import (
	"swahili/lang/lexer"
)

// VarDeclarationStatement ...
type VarDeclarationStatement struct {
	// The name of the variable
	Name string
	// Wether or not the variable is a constant
	IsConstant bool
	// The value assigned to the variable
	Value Expression
	// The explicit type of the variable
	ExplicitType Type
	Tokens       []lexer.Token
}

var _ Statement = (*VarDeclarationStatement)(nil)

func (stmt *VarDeclarationStatement) Accept(g CodeGenerator) error {
	return g.VisitVarDeclaration(stmt)
}

func (stmt VarDeclarationStatement) TokenStream() []lexer.Token {
	return stmt.Tokens
}
