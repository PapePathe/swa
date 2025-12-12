package ast

import (
	"encoding/json"
	"fmt"

	"tinygo.org/x/go-llvm"
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
	case DataTypePointer:
		return "DataTypePointer"
	default:
		panic(fmt.Sprintf("Unmatched data type %d", dt))
	}
}

const (
	DataTypeArray = iota
	DataTypeNumber
	DataTypeFloat
	DataTypeString
	DataTypeStruct
	DataTypeIntType
	DataTypeSymbol
	DataTypePointer
)

// Type
type Type interface {
	Value() DataType
	LLVMType() llvm.Type
}

type SymbolType struct {
	Name string
}

var _ Type = (*SymbolType)(nil)

func (SymbolType) Value() DataType {
	return DataTypeSymbol
}

func (typ SymbolType) LLVMTypeDyn(ctx *CompilerCtx) (error, llvm.Type) {
	err, sym := ctx.FindStructSymbol(typ.Name)
	if err != nil {
		return err, llvm.Type{}
	}

	return nil, sym.LLVMType
}

func (SymbolType) LLVMType() llvm.Type {
	panic("Please use SymbolType.LLVMTypeDyn")
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

func (a ArrayType) LLVMType() llvm.Type {
	return llvm.ArrayType(a.Underlying.LLVMType(), a.Size)
}

func (se ArrayType) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Value"] = se.Value().String()
	m["Underlying"] = se.Underlying

	res := make(map[string]any)
	res["ast.ArrayType"] = m

	return json.Marshal(res)
}

type NumberType struct{}

var _ Type = (*NumberType)(nil)

func (NumberType) Value() DataType {
	return DataTypeNumber
}

func (NumberType) LLVMType() llvm.Type {
	return llvm.GlobalContext().Int32Type()
}

func (se NumberType) MarshalJSON() ([]byte, error) {
	res := make(map[string]any)
	res["ast.NumberType"] = se.Value().String()

	return json.Marshal(res)
}

type StringType struct{}

var _ Type = (*StringType)(nil)

func (StringType) Value() DataType {
	return DataTypeString
}

func (se StringType) MarshalJSON() ([]byte, error) {
	res := make(map[string]any)
	res["ast.StringType"] = se.Value().String()

	return json.Marshal(res)
}

func (StringType) LLVMType() llvm.Type {
	return llvm.PointerType(
		llvm.GlobalContext().Int8Type(),
		0,
	)
}

type FloatType struct{}

var _ Type = (*FloatType)(nil)

func (FloatType) Value() DataType {
	return DataTypeFloat
}

func (FloatType) LLVMType() llvm.Type {
	return llvm.GlobalContext().DoubleType()
}

func (se FloatType) MarshalJSON() ([]byte, error) {
	res := make(map[string]any)
	res["ast.FloatType"] = se.Value().String()

	return json.Marshal(res)
}

type PointerType struct {
	Underlying Type
}

// LLVMType implements Type.
func (se PointerType) LLVMType() llvm.Type {
	return llvm.PointerType(se.Underlying.LLVMType(), 0)
}

var _ Type = (*PointerType)(nil)

func (PointerType) Value() DataType {
	return DataTypePointer
}

func (se PointerType) MarshalJSON() ([]byte, error) {
	res := make(map[string]any)
	res["ast.PointerType"] = se.Value().String()

	return json.Marshal(res)
}
