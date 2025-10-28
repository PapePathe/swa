package ast

// Statement ...
type Statement interface {
	CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult)
}

// Expression ...
type Expression interface {
	CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult)
}
