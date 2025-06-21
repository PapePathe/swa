package ast

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"tinygo.org/x/go-llvm"
)

type ReturnStatement struct {
	Value Expression
}

func (rs ReturnStatement) Compile(ctx *Context) error {
	switch v := rs.Value.(type) {
	case NumberExpression:
		ctx.NewRet(constant.NewInt(types.I32, int64(v.Value)))
	case BinaryExpression:
		err, res := rs.Value.Compile(ctx)
		if err != nil {
			return err
		}
		ctx.NewRet(res.v)
	default:
		err := fmt.Errorf("unknown expression in ReturnStatement <%s>", rs.Value)

		panic(err)
	}

	return nil
}

func (rs ReturnStatement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	switch v := rs.Value.(type) {
	case SymbolExpression:
		val, ok := ctx.SymbolTable[v.Value]
		if !ok {
			return fmt.Errorf("Undefined variable %s", v.Value), nil
		}
		switch val.Type() {
		case llvm.GlobalContext().Int32Type():
			ctx.Builder.CreateRet(val)
		case llvm.PointerType(llvm.GlobalContext().Int32Type(), 0):
			loadedval := ctx.Builder.CreateLoad(llvm.GlobalContext().Int32Type(), val, "")
			ctx.Builder.CreateRet(loadedval)
		}
	case NumberExpression:
		ctx.Builder.CreateRet(llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(v.Value), false))
	case BinaryExpression:
		err, res := rs.Value.CompileLLVM(ctx)
		if err != nil {
			return err, nil
		}
		ctx.Builder.CreateRet(*res)
	default:
		err := fmt.Errorf("unknown expression in ReturnStatement <%s>", rs.Value)

		panic(err)
	}
	return nil, nil
}
