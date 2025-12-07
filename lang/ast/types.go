package ast

import (
	"encoding/json"
	"fmt"
)

type DataType int

func (dt DataType) String() string {
	switch dt {
	case DataTypeString:
		return "DataTypeString"
	case DataTypeIntType:
		return "DataTypeIntType"
	case DataTypeNumber:
		return "DataTypeNumber"
	case DataTypeFloat:
		return "DataTypeFloat"
	case DataTypeArray:
		return "DataTypeArray"
	case DataTypeSymbol:
		return "DataTypeSymbol"
	case DataTypeVoid:
		return "DataTypeVoid"
	default:
		panic(fmt.Sprintf("Unmatched data type %d", dt))
	}
}

const (
	DataTypeArray = iota
	DataTypeNumber
	DataTypeNumber64
	DataTypeFloat
	DataTypeString
	DataTypeStruct
	DataTypeIntType
	DataTypeSymbol
	DataTypeVoid
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

func (se SymbolType) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = se.Value().String()
	m["Name"] = se.Name

	res := make(map[string]any)
	res["ast.SymbolType"] = m

	return json.Marshal(res)
}

type ArrayType struct {
	Underlying Type
	Size       int
}

var _ Type = (*ArrayType)(nil)

func (ArrayType) Value() DataType {
	return DataTypeArray
}

func (se ArrayType) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = se.Value().String()
	m["Underlying"] = se.Underlying

	res := make(map[string]any)
	res["ast.ArrayType"] = m

	return json.Marshal(res)
}

type Number64Type struct{}

var _ Type = (*Number64Type)(nil)

func (Number64Type) Value() DataType {
	return DataTypeNumber64
}

func (se Number64Type) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = se.Value().String()

	res := make(map[string]any)
	res["ast.Number64Type"] = m

	return json.Marshal(res)
}

type NumberType struct{}

var _ Type = (*NumberType)(nil)

func (NumberType) Value() DataType {
	return DataTypeNumber
}

func (se NumberType) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = se.Value().String()

	res := make(map[string]any)
	res["ast.NumberType"] = m

	return json.Marshal(res)
}

type StringType struct{}

var _ Type = (*StringType)(nil)

func (StringType) Value() DataType {
	return DataTypeString
}

func (se StringType) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = se.Value().String()

	res := make(map[string]any)
	res["ast.StringType"] = m

	return json.Marshal(res)
}

type FloatType struct{}

var _ Type = (*FloatType)(nil)

func (FloatType) Value() DataType {
	return DataTypeFloat
}

func (se FloatType) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = se.Value().String()

	res := make(map[string]any)
	res["ast.FloatType"] = m

	return json.Marshal(res)
}

type VoidType struct{}

var _ Type = (*VoidType)(nil)

func (VoidType) Value() DataType {
	return DataTypeVoid
}

func (se VoidType) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["ast.FloatType"] = se.Value().String()

	return json.Marshal(m)
}
