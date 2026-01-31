package ast

import (
	"swahili/lang/lexer"
)

type FuncArg struct {
	Name    string
	ArgType Type
}

type FuncDeclStatement struct {
	Body         BlockStatement
	Name         string
	ReturnType   Type
	Args         []FuncArg
	Tokens       []lexer.Token
	ArgsVariadic bool
}

var _ Statement = (*FuncDeclStatement)(nil)

func (stmt *FuncDeclStatement) Accept(g CodeGenerator) error {
	return g.VisitFunctionDefinition(stmt)
}

func (stmt FuncDeclStatement) TokenStream() []lexer.Token {
	return stmt.Tokens
}
