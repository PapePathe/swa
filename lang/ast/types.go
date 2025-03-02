package ast

type SymbolType struct {
	Name string
}

var _ Type = (*SymbolType)(nil)

func (st SymbolType) _type() {}

type ArrayType struct {
	Underlying Type
}

var _ Type = (*ArrayType)(nil)

func (st ArrayType) _type() {}
