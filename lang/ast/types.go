package ast

import (
	"fmt"
	"os"

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
	case DataTypeNumber64:
		return "DataTypeNumber64 bits"
	case DataTypeFloat:
		return "DataTypeFloat"
	case DataTypeArray:
		return "DataTypeArray"
	case DataTypeSymbol:
		return "DataTypeSymbol"
	case DataTypePointer:
		return "DataTypePointer"
	case DataTypeVoid:
		return "DataTypeVoid"
	default:
		fmt.Printf("Unmatched data type %d", dt)
		os.Exit(1)
	}

	return ""
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
	DataTypePointer
	DataTypeVoid
)

// Type
type Type interface {
	Accept(g CodeGenerator) error
	Value() DataType
	LLVMType(ctx *CompilerCtx) (error, llvm.Type)
}

type SymbolType struct {
	Name string
}

var _ Type = (*SymbolType)(nil)

func (SymbolType) Value() DataType {
	return DataTypeSymbol
}

func (typ SymbolType) Accept(g CodeGenerator) error {
	return g.VisitSymbolType(&typ)
}

func (typ SymbolType) LLVMType(ctx *CompilerCtx) (error, llvm.Type) {
	err, sym := ctx.FindStructSymbol(typ.Name)
	if err != nil {
		return err, llvm.Type{}
	}

	return nil, sym.LLVMType
}

type ArrayType struct {
	Underlying        Type
	Size              int
	DynSizeIdentifier Expression
}

var _ Type = (*ArrayType)(nil)

func (ArrayType) Value() DataType {
	return DataTypeArray
}

func (typ ArrayType) Accept(g CodeGenerator) error {
	return g.VisitArrayType(&typ)
}

func (a ArrayType) LLVMType(ctx *CompilerCtx) (error, llvm.Type) {
	err, under := a.Underlying.LLVMType(ctx)
	if err != nil {
		return err, llvm.Type{}
	}

	return nil, llvm.ArrayType(under, a.Size)
}

type NumberType struct{}

var _ Type = (*NumberType)(nil)

func (NumberType) Value() DataType {
	return DataTypeNumber
}

func (typ NumberType) Accept(g CodeGenerator) error {
	return g.VisitNumberType(&typ)
}

func (NumberType) LLVMType(*CompilerCtx) (error, llvm.Type) {
	return nil, llvm.GlobalContext().Int32Type()
}

type StringType struct{}

var _ Type = (*StringType)(nil)

func (StringType) Value() DataType {
	return DataTypeString
}

func (typ StringType) Accept(g CodeGenerator) error {
	return g.VisitStringType(&typ)
}

func (StringType) LLVMType(*CompilerCtx) (error, llvm.Type) {
	return nil, llvm.PointerType(
		llvm.GlobalContext().Int8Type(),
		0,
	)
}

type FloatType struct{}

var _ Type = (*FloatType)(nil)

func (FloatType) Value() DataType {
	return DataTypeFloat
}

func (FloatType) LLVMType(*CompilerCtx) (error, llvm.Type) {
	return nil, llvm.GlobalContext().DoubleType()
}

func (typ FloatType) Accept(g CodeGenerator) error {
	return g.VisitFloatType(&typ)
}

type PointerType struct {
	Underlying Type
}

func (typ PointerType) Accept(g CodeGenerator) error {
	return g.VisitPointerType(&typ)
}

func (se PointerType) LLVMType(ctx *CompilerCtx) (error, llvm.Type) {
	err, under := se.Underlying.LLVMType(ctx)
	if err != nil {
		return err, llvm.Type{}
	}

	return nil, llvm.PointerType(under, 0)
}

var _ Type = (*PointerType)(nil)

func (PointerType) Value() DataType {
	return DataTypePointer
}

type VoidType struct{}

var _ Type = (*VoidType)(nil)

func (VoidType) Value() DataType {
	return DataTypeVoid
}

func (VoidType) LLVMType(ctx *CompilerCtx) (error, llvm.Type) {
	return nil, llvm.GlobalContext().VoidType()
}

func (typ VoidType) Accept(g CodeGenerator) error {
	return g.VisitVoidType(&typ)
}

type Number64Type struct{}

var _ Type = (*Number64Type)(nil)

func (Number64Type) Value() DataType {
	return DataTypeNumber64
}

func (Number64Type) LLVMType(ctx *CompilerCtx) (error, llvm.Type) {
	return nil, llvm.GlobalContext().Int64Type()
}

func (typ Number64Type) Accept(g CodeGenerator) error {
	return g.VisitNumber64Type(&typ)
}
