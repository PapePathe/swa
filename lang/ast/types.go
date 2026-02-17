package ast

import (
	"fmt"
	"os"
	"strings"
	"swahili/lang/lexer"
)

type DataType int

func (dt DataType) String() string {
	switch dt {
	case DataTypeString:
		return "String"
	case DataTypeIntType:
		return "IntType"
	case DataTypeNumber:
		return "Number"
	case DataTypeNumber64:
		return "Number64"
	case DataTypeFloat:
		return "Float"
	case DataTypeArray:
		return "Array"
	case DataTypeError:
		return "Error"
	case DataTypeTuple:
		return "Tuple"
	case DataTypeSymbol:
		return "Symbol"
	case DataTypePointer:
		return "Pointer"
	case DataTypeVoid:
		return "Void"
	case DataTypeBool:
		return "Bool"
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
	DataTypeError
	DataTypeTuple
	DataTypeBool
)

type SymbolType struct {
	Name string
}

func (typ SymbolType) Equals(other Type) bool {
	if typ.Value() != other.Value() {
		return false
	}

	oth, _ := other.(*SymbolType)

	if typ.Name != oth.Name {
		return false
	}

	return true
}

func (typ SymbolType) AcceptZero(g CodeGenerator) error {
	return g.ZeroOfSymbolType(&typ)
}

var _ Type = (*SymbolType)(nil)

func (SymbolType) Value() DataType {
	return DataTypeSymbol
}

func (typ SymbolType) Accept(g CodeGenerator) error {
	return g.VisitSymbolType(&typ)
}

func (t SymbolType) String() string {
	return fmt.Sprintf("%s(%s)", t.Value(), t.Name)
}

type ArrayType struct {
	Underlying        Type
	Size              int
	DynSizeIdentifier Expression
}

func (typ ArrayType) Equals(Type) bool {
	panic("unimplemented")
}

func (t ArrayType) String() string {
	return fmt.Sprintf("%s(%s)", t.Value(), t.Underlying.Value())
}

func (typ ArrayType) AcceptZero(g CodeGenerator) error {
	return g.ZeroOfArrayType(&typ)
}

var _ Type = (*ArrayType)(nil)

func (ArrayType) Value() DataType {
	return DataTypeArray
}

func (typ ArrayType) Accept(g CodeGenerator) error {
	return g.VisitArrayType(&typ)
}

type NumberType struct{}

func (t NumberType) String() string {
	return t.Value().String()
}

func (typ NumberType) Equals(other Type) bool {
	return typ.Value() == other.Value()
}

func (typ NumberType) AcceptZero(g CodeGenerator) error {
	return g.ZeroOfNumberType(&typ)
}

var _ Type = (*NumberType)(nil)

func (NumberType) Value() DataType {
	return DataTypeNumber
}

func (typ NumberType) Accept(g CodeGenerator) error {
	return g.VisitNumberType(&typ)
}

type StringType struct{}

// Equals implements [Type].
func (typ StringType) Equals(other Type) bool {
	if typ.Value() != other.Value() {
		return false
	}

	_, ok := other.(StringType)

	return ok
}

func (t StringType) String() string {
	return t.Value().String()
}

// AcceptZero implements [Type].
func (typ StringType) AcceptZero(g CodeGenerator) error {
	return g.ZeroOfStringType(&typ)
}

var _ Type = (*StringType)(nil)

func (StringType) Value() DataType {
	return DataTypeString
}

func (typ StringType) Accept(g CodeGenerator) error {
	return g.VisitStringType(&typ)
}

type FloatType struct{}

func (typ FloatType) Equals(other Type) bool {
	if typ.Value() != other.Value() {
		return false
	}

	_, ok := other.(FloatType)

	return ok
}

func (t FloatType) String() string {
	return t.Value().String()
}

func (typ FloatType) AcceptZero(g CodeGenerator) error {
	return g.ZeroOfFloatType(&typ)
}

var _ Type = (*FloatType)(nil)

func (FloatType) Value() DataType {
	return DataTypeFloat
}

func (typ FloatType) Accept(g CodeGenerator) error {
	return g.VisitFloatType(&typ)
}

type PointerType struct {
	Underlying Type
}

func (t PointerType) String() string {
	return fmt.Sprintf("%s(%s)", t.Value(), t.Underlying.Value())
}

func (typ PointerType) Equals(other Type) bool {
	if typ.Value() != other.Value() {
		return false
	}

	oth, ok := other.(PointerType)

	if !ok {
		return false
	}

	return typ.Underlying.Equals(oth.Underlying)
}

func (typ PointerType) AcceptZero(g CodeGenerator) error {
	return g.ZeroOfPointerType(&typ)
}

func (typ PointerType) Accept(g CodeGenerator) error {
	return g.VisitPointerType(&typ)
}

var _ Type = (*PointerType)(nil)

func (PointerType) Value() DataType {
	return DataTypePointer
}

type VoidType struct{}

func (t VoidType) String() string {
	return t.Value().String()
}

func (typ VoidType) Equals(other Type) bool {
	if typ.Value() != other.Value() {
		return false
	}

	_, ok := other.(VoidType)

	return ok
}

func (typ VoidType) AcceptZero(g CodeGenerator) error {
	return g.ZeroOfVoidType(&typ)
}

var _ Type = (*VoidType)(nil)

func (VoidType) Value() DataType {
	return DataTypeVoid
}

func (typ VoidType) Accept(g CodeGenerator) error {
	return g.VisitVoidType(&typ)
}

type Number64Type struct{}

func (typ Number64Type) Equals(other Type) bool {
	if typ.Value() != other.Value() {
		return false
	}

	_, ok := other.(VoidType)

	return ok
}

func (t Number64Type) String() string {
	return t.Value().String()
}

func (typ Number64Type) AcceptZero(g CodeGenerator) error {
	return g.ZeroOfNumber64Type(&typ)
}

var _ Type = (*Number64Type)(nil)

func (Number64Type) Value() DataType {
	return DataTypeNumber64
}

func (typ Number64Type) Accept(g CodeGenerator) error {
	return g.VisitNumber64Type(&typ)
}

type ErrorType struct{}

func (t ErrorType) String() string {
	return t.Value().String()
}

func (typ ErrorType) Equals(other Type) bool {
	if typ.Value() != other.Value() {
		return false
	}

	_, ok := other.(ErrorType)

	return ok
}

func (typ ErrorType) AcceptZero(g CodeGenerator) error {
	return g.ZeroOfErrorType(&typ)
}

var _ Type = (*ErrorType)(nil)

func (ErrorType) Value() DataType {
	return DataTypeError
}

func (typ ErrorType) Accept(g CodeGenerator) error {
	return g.VisitErrorType(&typ)
}

type TupleType struct {
	Types  []Type
	Tokens []lexer.Token
}

func (t TupleType) Equals(other Type) bool {
	if t.Value() != other.Value() {
		return false
	}

	oth, ok := other.(*TupleType)
	if !ok {
		return false
	}

	if len(t.Types) != len(oth.Types) {
		return false
	}

	for i, v := range t.Types {
		if !v.Equals(oth.Types[i]) {
			return false
		}
	}

	return true
}

var _ Type = (*TupleType)(nil)

func (t TupleType) Value() DataType {
	return DataTypeTuple
}

func (t *TupleType) Accept(g CodeGenerator) error {
	return g.VisitTupleType(t)
}

func (t *TupleType) AcceptZero(g CodeGenerator) error {
	return g.ZeroOfTupleType(t)
}

func (t *TupleType) String() string {
	sb := strings.Builder{}
	lminus := len(t.Types) - 1

	sb.WriteString(t.Value().String())
	sb.WriteString("(")

	for i, typ := range t.Types {
		sb.WriteString(typ.Value().String())

		if i < lminus {
			sb.WriteString(",")
		}
	}

	sb.WriteString(")")

	return sb.String()
}

type BoolType struct{}

func (b BoolType) Equals(other Type) bool {
	if b.Value() != other.Value() {
		return false
	}

	_, ok := other.(*BoolType)

	return ok
}

func (t *BoolType) String() string {
	return t.Value().String()
}

var _ Type = (*BoolType)(nil)

func (b *BoolType) Accept(g CodeGenerator) error {
	return g.VisitBoolType(b)
}

func (b *BoolType) AcceptZero(g CodeGenerator) error {
	return g.ZeroOfBoolType(b)
}

func (b *BoolType) Value() DataType {
	return DataTypeBool
}
