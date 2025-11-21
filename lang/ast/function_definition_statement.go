package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type FuncArg struct {
	Name    string
	ArgType Type
}

type FuncDeclStatement struct {
	Body       BlockStatement
	Name       string
	ReturnType Type
	Args       []FuncArg
	Tokens     []lexer.Token
}

var _ Statement = (*FuncDeclStatement)(nil)

func (fd FuncDeclStatement) extractType(t Type) llvm.Type {
	switch t.Value() {
	case DataTypeNumber:
		return llvm.GlobalContext().Int32Type()
	case DataTypeFloat:
		return llvm.GlobalContext().DoubleType()
	case DataTypeIntType:
		return llvm.GlobalContext().Int32Type()
	case DataTypeString:
		return llvm.PointerType(llvm.GlobalContext().Int8Type(), 0)
	default:
		panic(fmt.Errorf("argument type %v not supported", t))
	}
}

func (fd FuncDeclStatement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	params := []llvm.Type{}

	for _, arg := range fd.Args {
		params = append(params, fd.extractType(arg.ArgType))
	}

	returnType := fd.extractType(fd.ReturnType)
	newfuncType := llvm.FunctionType(returnType, params, false)
	newFunc := llvm.AddFunction(*ctx.Module, fd.Name, newfuncType)
	if err := ctx.AddFuncSymbol(fd.Name, &newfuncType); err != nil {
		return err, nil
	}

	for i, p := range newFunc.Params() {
		name := fd.Args[i].Name

		ctx.AddSymbol(name, &SymbolTableEntry{Value: p})
		p.SetName(name)
	}

	block := ctx.Context.AddBasicBlock(newFunc, "func-body")
	ctx.Builder.SetInsertPointAtEnd(block)

	if err, _ := fd.Body.CompileLLVM(ctx); err != nil {
		return err, nil
	}

	return nil, nil
}

func (expr FuncDeclStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (fd FuncDeclStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Args"] = fd.Args
	m["Body"] = fd.Body
	m["Name"] = fd.Name
	m["ReturnType"] = fd.ReturnType

	res := make(map[string]any)
	res["ast.FuncDeclStatement"] = m

	return json.Marshal(res)
}
