package ast

import (
	"swahili/lang/values"

	"github.com/llir/llvm/ir"
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

func (rs ReturnStatement) Compile(m *ir.Module, b *ir.Block) error {
	b.NewRet(constant.NewInt(types.I32, 0))

	return nil
}
