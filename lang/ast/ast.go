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

type StructSymbolTableEntry struct {
	LLVMType llvm.Type
	Metadata StructDeclarationStatement
}

type CompilerCtx struct {
	Context           *llvm.Context
	Builder           *llvm.Builder
	Module            *llvm.Module
	SymbolTable       map[string]llvm.Value
	StructSymbolTable map[string]StructSymbolTableEntry
	FuncSymbolTable   map[string]llvm.Type
}
