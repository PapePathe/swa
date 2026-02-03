package ast

import (
	"swahili/lang/lexer"
)

type TupleAssignmentExpression struct {
	Assignees *TupleExpression
	Operator  lexer.Token
	Value     Expression
	Tokens    []lexer.Token
}

var _ Expression = (*TupleAssignmentExpression)(nil)

func (e *TupleAssignmentExpression) Accept(g CodeGenerator) error {
	return g.VisitTupleAssignmentExpression(e)
}

func (e *TupleAssignmentExpression) TokenStream() []lexer.Token {
	return e.Tokens
}

func (e *TupleAssignmentExpression) VisitedSwaType() Type {
	panic("VisitedSwaType for TupleAssignmentExpression not implemented")
}
