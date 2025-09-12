package ast

type DataType = int

const (
	DataTypeArray = iota
	DataTypeNumber
	DataTypeString
	DataTypeStruct
	DataTypeIntType
	DataTypeSymbol
)

// Type
type Type interface {
	_type() DataType
}

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

type NumberType struct{}

var _ Type = (*NumberType)(nil)

func (NumberType) _type() DataType {
	return DataTypeNumber
}

type StringType struct{}

var _ Type = (*StringType)(nil)

func (StringType) _type() DataType {
	return DataTypeString
}
