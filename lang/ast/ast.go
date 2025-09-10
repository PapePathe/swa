package ast

import (
	"tinygo.org/x/go-llvm"
)

// Statement ...
type Statement interface {
	CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value)
}

// Expression ...
type Expression interface {
	CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value)
}

type DataType = int

const (
	DataTypeArray = iota
	DataTypeNumber
	DataTypeString
	DataTypeStruct
	DataTypeSymbol
)

// Type
type Type interface {
	_type() DataType
}

type CompilerCtx struct {
	Context           *llvm.Context
	Builder           *llvm.Builder
	Module            *llvm.Module
	SymbolTable       map[string]llvm.Value
	StructSymbolTable map[string]llvm.Type
	FuncSymbolTable   map[string]llvm.Type
}
