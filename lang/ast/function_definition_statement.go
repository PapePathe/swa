package ast

import (
	"encoding/json"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
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

func (fd FuncDeclStatement) Compile(ctx *Context) error {
	params := []*ir.Param{}

	for _, arg := range fd.Args {
		var param *ir.Param
		switch arg.ArgType {
		case "Entier_32":
			param = ir.NewParam(arg.Name, types.I32)
		case "Chaine":
			param = ir.NewParam(arg.Name, types.I8Ptr)
		default:
			panic(fmt.Errorf("argument type %s not supported", arg.ArgType))
		}
		params = append(params, param)
	}

	funcDef := ctx.mod.NewFunc(fd.Name, types.I32, params...)
	funcBlk := funcDef.NewBlock("")
	funcCtx := ctx.NewContext(funcBlk)

	for _, param := range params {
		alloc := funcCtx.NewAlloca(param.Type())

		funcCtx.NewStore(constant.NewInt(types.I32, 0), alloc)
		funcCtx.AddLocal(param.Name(), LocalVariable{Value: alloc})
	}
	fd.Body.Compile(funcCtx)

	return nil
}

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
		p.SetName(fd.Args[i].Name)
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
