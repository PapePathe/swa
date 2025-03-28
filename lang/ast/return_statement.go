package ast

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type ReturnStatement struct {
	Value Expression
}

func (rs ReturnStatement) Compile(ctx *Context) error {
	switch v := rs.Value.(type) {
	case NumberExpression:
		ctx.NewRet(constant.NewInt(types.I32, int64(v.Value)))
	default:
		err := fmt.Errorf("unknown expression in ReturnStatement <%s>", rs.Value)

		panic(err)
	}

	return nil
}
