package ast

import (
	"swahili/lang/lexer"
)

type TupleType struct {
	Types  []Type
	Tokens []lexer.Token
}

var _ Type = (*TupleType)(nil)

func (t TupleType) Value() DataType {
	// We might need a DataTypeTuple if we want to distinguish it at the DataType level,
	// but for now let's see if we can just use the underlying types.
	// Actually, adding DataTypeTuple might be cleaner.
	return DataTypeSymbol // Placeholder, will update if needed
}

func (t *TupleType) Accept(g CodeGenerator) error {
	return g.VisitTupleType(t)
}

func (t *TupleType) AcceptZero(g CodeGenerator) error {
	return g.ZeroOfTupleType(t)
}
