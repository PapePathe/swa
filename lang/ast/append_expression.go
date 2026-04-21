package ast

import "swahili/lang/lexer"

type AppendExpression struct {
	Appendable Expression
	Value      Expression
	Tokens     []lexer.Token
}

var _ Expression = (*AppendExpression)(nil)

func (a *AppendExpression) Accept(g CodeGenerator) error {
	return g.VisitAppendExpression(a)
}

func (a *AppendExpression) TokenStream() []lexer.Token {
	return a.Tokens
}

func (a *AppendExpression) VisitedSwaType() Type {
	panic("unimplemented")
}
