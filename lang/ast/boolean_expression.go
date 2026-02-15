package ast

import "swahili/lang/lexer"

type BooleanExpression struct {
	Value   bool
	SwaType Type
	Tokens  []lexer.Token
}

var _ Expression = (*BooleanExpression)(nil)

func (b *BooleanExpression) Accept(g CodeGenerator) error {
	return g.VisitBooleanExpression(b)
}

func (b *BooleanExpression) TokenStream() []lexer.Token {
	return b.Tokens
}

func (b *BooleanExpression) VisitedSwaType() Type {
	return &BoolType{}
}
