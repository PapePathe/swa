package ast

import (
	"encoding/json"
	"swahili/lang/lexer"
)

var _ Statement = (*WhileStatement)(nil)

type WhileStatement struct {
	Condition Expression
	Body      BlockStatement
	Tokens    []lexer.Token
}

func (expr WhileStatement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	bodyBlock := ctx.Builder.GetInsertBlock()
	parentFunc := bodyBlock.Parent()

	loopBlock := ctx.Context.AddBasicBlock(parentFunc, "loopBlock")
	conditionBlock := ctx.Context.AddBasicBlock(parentFunc, "conditionBlock")
	ctx.Builder.CreateBr(loopBlock)
	exitBlock := ctx.Context.AddBasicBlock(parentFunc, "exitBlock")
	ctx.Builder.SetInsertPointAtEnd(conditionBlock)

	err, cond := expr.Condition.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	ctx.Builder.CreateCondBr(*cond.Value, loopBlock, exitBlock)
	ctx.Builder.SetInsertPointAtEnd(loopBlock)

	err, _ = expr.Body.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	ctx.Builder.CreateBr(conditionBlock)
	conditionBlock.MoveAfter(bodyBlock)
	loopBlock.MoveAfter(bodyBlock)
	exitBlock.MoveAfter(loopBlock)
	ctx.Builder.SetInsertPointAtEnd(exitBlock)

	return nil, nil
}

func (expr WhileStatement) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr WhileStatement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Condition"] = expr.Condition
	m["Body"] = expr.Body

	res := make(map[string]any)
	res["ast.WhileStatement"] = m

	return json.Marshal(res)
}
