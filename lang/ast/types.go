package ast

type SymbolType struct {
	Name string
}

var _ Type = (*SymbolType)(nil)

func (SymbolType) _type() {}

type ArrayType struct {
	Underlying Type
}

var _ Type = (*ArrayType)(nil)

func (ArrayType) _type() {}

type NumberType struct {
	Value int
}

var _ Type = (*NumberType)(nil)

func (NumberType) _type() {}
