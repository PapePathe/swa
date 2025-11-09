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

func (fd FuncDeclStatement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	params := []llvm.Type{}

	for _, arg := range fd.Args {
		var param llvm.Type

		switch arg.ArgType.Value() {
		case DataTypeNumber:
			param = llvm.GlobalContext().Int32Type()
		case DataTypeIntType:
			param = llvm.GlobalContext().Int32Type()
		case DataTypeString:
			param = llvm.PointerType(llvm.GlobalContext().Int8Type(), 0)
		default:
			panic(fmt.Errorf("argument type %v not supported", arg))
		}

		params = append(params, param)
	}

	var returnType llvm.Type
	switch fd.ReturnType.Value() {
	case DataTypeNumber:
		returnType = llvm.GlobalContext().Int32Type()
	case DataTypeIntType:
		returnType = llvm.GlobalContext().Int32Type()
	case DataTypeString:
		returnType = llvm.PointerType(llvm.GlobalContext().Int8Type(), 0)
	default:
		panic(fmt.Errorf("argument type %v not supported in return type", fd.ReturnType))
	}

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
