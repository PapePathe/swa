package ast

import (
	"swahili/lang/lexer"
)

type ArrayInitializationExpression struct {
	Underlying Type
	Contents   []Expression
	Tokens     []lexer.Token
	SwaType    Type
}

var _ Expression = (*ArrayInitializationExpression)(nil)

func (expr ArrayInitializationExpression) Accept(g CodeGenerator) error {
	return g.VisitArrayInitializationExpression(&expr)
}

func (expr ArrayInitializationExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}
