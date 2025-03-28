package ast

import (
	"fmt"
	"swahili/lang/values"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type ReturnStatement struct {
	Value Expression
}

func (rs ReturnStatement) statement() {}

func (rs ReturnStatement) Evaluate(s *Scope) (error, values.Value) {
	return nil, nil
}

func (rs ReturnStatement) Compile(ctx *Context) error {
	switch v := rs.Value.(type) {
	case NumberExpression:
		ctx.Block.NewRet(constant.NewInt(types.I32, int64(v.Value)))
	default:
		err := fmt.Errorf("Unknown expression in ReturnStatement <%s>", rs.Value)

		panic(err)
	}

	return nil
}
