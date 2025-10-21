package ast

import "encoding/json"

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
	Value() DataType
}

type SymbolType struct {
	Name string
}

var _ Type = (*SymbolType)(nil)

func (SymbolType) Value() DataType {
	return DataTypeSymbol
}

type ArrayType struct {
	Underlying Type
}

var _ Type = (*ArrayType)(nil)

func (ArrayType) Value() DataType {
	return DataTypeArray
}

func (a ArrayType) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Type"] = a.Value()
	m["UnderlyingType"] = a.Underlying.Value()

	return json.Marshal(m)
}

type NumberType struct{}

var _ Type = (*NumberType)(nil)

func (NumberType) Value() DataType {
	return DataTypeNumber
}

type StringType struct{}

var _ Type = (*StringType)(nil)

func (StringType) Value() DataType {
	return DataTypeString
}
