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
