package ast

import (
	"encoding/json"
	"fmt"

	"tinygo.org/x/go-llvm"
)

type FuncArg struct {
	Name    string
	ArgType string
}

type FuncDeclStatement struct {
	Body       BlockStatement
	Name       string
	ReturnType string
	Args       []FuncArg
}

var _ Statement = (*FuncDeclStatement)(nil)

func (fd FuncDeclStatement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	params := []llvm.Type{}

	for _, arg := range fd.Args {
		var param llvm.Type
		switch arg.ArgType {
		case "Entier_32":
			param = llvm.GlobalContext().Int32Type()
		case "Chaine":
			param = llvm.GlobalContext().Int8Type()
		default:
			panic(fmt.Errorf("argument type %s not supported", arg.ArgType))
		}
		params = append(params, param)
	}

	newFunc := llvm.AddFunction(
		*ctx.Module,
		fd.Name,
		llvm.FunctionType(
			llvm.GlobalContext().Int32Type(),
			params,
			false,
		),
	)

	for i, p := range newFunc.Params() {
		name := fd.Args[i].Name

		ctx.SymbolTable[name] = SymbolTableEntry{Value: p}
		p.SetName(name)
	}

	block := ctx.Context.AddBasicBlock(newFunc, "func-body")
	ctx.Builder.SetInsertPointAtEnd(block)

	if err, _ := fd.Body.CompileLLVM(ctx); err != nil {
		return err, nil
	}

	return nil, nil
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
