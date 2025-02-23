package interpreter

import "swahili/lang/ast"

type InterpreterResult = any

type Interpreter interface {
	Interpret() (error, InterpreterResult)
}

// Run a program given an abstract syntax tree
func Run(program ast.BlockStatement) {
}
