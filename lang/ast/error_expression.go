package ast

import (
	"swahili/lang/lexer"
)

type ErrorExpression struct {
	Tokens []lexer.Token
}

var _ Expression = (*ErrorExpression)(nil)

func (e *ErrorExpression) Accept(g CodeGenerator) error {
	return g.VisitErrorExpression(e)
}

func (e ErrorExpression) TokenStream() []lexer.Token {
	return e.Tokens
}

func (e ErrorExpression) VisitedSwaType() Type {
	return &ErrorType{}
}
