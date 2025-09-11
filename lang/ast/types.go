package ast

type SymbolType struct {
	Name string
}

var _ Type = (*SymbolType)(nil)

func (SymbolType) _type() DataType {
	return DataTypeSymbol
}

type ArrayType struct {
	Underlying Type
}

var _ Type = (*ArrayType)(nil)

func (ArrayType) _type() DataType {
	return DataTypeArray
}

type NumberType struct {
	Value int
}

var _ Type = (*NumberType)(nil)

func (NumberType) _type() DataType {
	return DataTypeNumber
}
