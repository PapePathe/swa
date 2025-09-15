package ast

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

type ReturnStatement struct {
	Value Expression
}

func (rs ReturnStatement) CompileLLVM(ctx *CompilerCtx) (error, *llvm.Value) {
	switch v := rs.Value.(type) {
	case SymbolExpression:
		val, ok := ctx.SymbolTable[v.Value]
		if !ok {
			return fmt.Errorf("Undefined variable %s", v.Value), nil
		}

		switch val.Value.Type() {
		case llvm.GlobalContext().Int32Type():
			ctx.Builder.CreateRet(val.Value)
		case llvm.PointerType(llvm.GlobalContext().Int32Type(), 0):
			loadedval := ctx.Builder.CreateLoad(llvm.GlobalContext().Int32Type(), val.Value, "")
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
