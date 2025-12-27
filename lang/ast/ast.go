package ast

import "swahili/lang/lexer"

// Node ...
type Node interface {
	Accept(g CodeGenerator) error
	TokenStream() []lexer.Token
}

// Statement ...
type Statement interface {
	Node
	CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult)
	MarshalJSON() ([]byte, error)
}

// Expression ...
type Expression interface {
	Node
	CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult)
	MarshalJSON() ([]byte, error)
}
