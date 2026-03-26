package ast

import "swahili/lang/lexer"

type ListExpression struct {
	DataType Type
	Length   Expression
	Capacity Expression
	Tokens   []lexer.Token
}

var _ Expression = (*ListExpression)(nil)

func (m *ListExpression) Accept(g CodeGenerator) error {
	panic("unimplemented")
}

func (m *ListExpression) TokenStream() []lexer.Token {
	panic("unimplemented")
}

func (m *ListExpression) VisitedSwaType() Type {
	panic("unimplemented")
}

type ListCountExpression struct {
	Expr   Expression
	Tokens []lexer.Token
}

var _ Expression = (*ListCountExpression)(nil)

func (l *ListCountExpression) Accept(g CodeGenerator) error {
	return g.VisitListCountExpression(l)
}

func (l *ListCountExpression) TokenStream() []lexer.Token {
	return l.Tokens
}

func (l *ListCountExpression) VisitedSwaType() Type {
	panic("unimplemented")
}
