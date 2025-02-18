package ast

type SymbolType struct {
	Name string
}

func (st SymbolType) _type() {}

type ArrayType struct {
	Underlying Type
}

func (st ArrayType) _type() {}
