package ast

import (
	"swahili/lang/lexer"
)

type TupleExpression struct {
	Expressions []Expression
	Tokens      []lexer.Token
}

var _ Expression = (*TupleExpression)(nil)

func (e *TupleExpression) Accept(g CodeGenerator) error {
	return g.VisitTupleExpression(e)
}

func (e TupleExpression) TokenStream() []lexer.Token {
	return e.Tokens
}

func (e TupleExpression) VisitedSwaType() Type {
	types := make([]Type, len(e.Expressions))
	for i, expr := range e.Expressions {
		types[i] = expr.VisitedSwaType()
	}

	return &TupleType{Types: types}
}
