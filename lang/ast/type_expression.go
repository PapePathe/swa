package ast

import "swahili/lang/lexer"

type TypeExpression struct {
	Type   Type
	Tokens []lexer.Token
}

func (t *TypeExpression) Accept(g CodeGenerator) error {
	return g.VisitTypeExpression(t)
}

func (t *TypeExpression) TokenStream() []lexer.Token {
	return t.Tokens
}

func (t *TypeExpression) VisitedSwaType() Type {
	return t.Type
}
