package ast

import (
	"fmt"
	"math"

	"tinygo.org/x/go-llvm"

	"swahili/lang/lexer"
)

// NumberExpression ...
type NumberExpression struct {
	Value  int64
	Tokens []lexer.Token
}

var _ Expression = (*NumberExpression)(nil)

func (e NumberExpression) String() string {
	return fmt.Sprintf("%d", e.Value)
}

func (se NumberExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	if se.Value < math.MinInt32 {
		return fmt.Errorf("%d is smaller than min value for int32", se.Value), nil
	}

	if se.Value > math.MaxInt32 {
		return fmt.Errorf("%d is greater than max value for int32", se.Value), nil
	}

	var signed bool

	if se.Value < 0 {
		signed = true
	}

	res := llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(se.Value), signed)

	return nil, &CompilerResult{Value: &res}
}

func (expr NumberExpression) Accept(g CodeGenerator) error {
	return g.VisitNumberExpression(&expr)
}

func (expr NumberExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}
