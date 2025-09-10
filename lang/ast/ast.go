package ast

import (
	"tinygo.org/x/go-llvm"
)

// Statement ...
type Statement interface {
	//	Compile(ctx *Context) error
	CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value)
}

// Expression ...
type Expression interface {
	//	Compile(ctx *Context) (error, *CompileResult)
	CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value)
}

// Type
type Type interface {
	_type()
}

type CompilerCtx struct {
	Context     *llvm.Context
	Builder     *llvm.Builder
	Module      *llvm.Module
	SymbolTable map[string]llvm.Value
}
