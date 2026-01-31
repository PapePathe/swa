package ast

import (
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type StructInitializationExpression struct {
	Name       string
	Properties []string
	Values     []Expression
	Tokens     []lexer.Token
	SwaType    Type
}

var _ Expression = (*StructInitializationExpression)(nil)

type StructItemValue struct {
	Position int
	Value    *llvm.Value
}

func (expr *StructInitializationExpression) Accept(g CodeGenerator) error {
	return g.VisitStructInitializationExpression(expr)
}

func (expr StructInitializationExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr StructInitializationExpression) VisitedSwaType() Type {
	return expr.SwaType
}
