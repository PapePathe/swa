package ast

import "swahili/lang/lexer"

// Statement ...
type Statement interface {
	CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult)
	MarshalJSON() ([]byte, error)
}

// Expression ...
type Expression interface {
	CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult)
	MarshalJSON() ([]byte, error)
	TokenStream() []lexer.Token
}
