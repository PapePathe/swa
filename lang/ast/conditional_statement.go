package ast

import (
	"encoding/json"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type ConditionalStatetement struct {
	Condition Expression
	Success   BlockStatement
	Failure   BlockStatement
	Tokens    []lexer.Token
}

var _ Statement = (*ConditionalStatetement)(nil)

func (cs ConditionalStatetement) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	err, condition := cs.Condition.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	bodyBlock := ctx.Builder.GetInsertBlock()
	parentFunc := bodyBlock.Parent()

	mergeBlock := ctx.Context.AddBasicBlock(parentFunc, "merge")
	thenBlock := ctx.Context.AddBasicBlock(parentFunc, "if")
	elseBlock := ctx.Context.AddBasicBlock(parentFunc, "else")

	ctx.Builder.CreateCondBr(*condition.Value, thenBlock, elseBlock)

	ctx.Builder.SetInsertPointAtEnd(thenBlock)

	err, successVal := cs.Success.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	if thenBlock.LastInstruction().InstructionOpcode() != llvm.Ret {
		ctx.Builder.CreateBr(mergeBlock)
	}

	ctx.Builder.SetInsertPointAtEnd(elseBlock)

	err, failureVal := cs.Failure.CompileLLVM(ctx)
	if err != nil {
		return err, nil
	}

	ctx.Builder.CreateBr(mergeBlock)
	ctx.Builder.SetInsertPointAtEnd(mergeBlock)

	var phi llvm.Value
	if successVal != nil {
		phi := ctx.Builder.CreatePHI(successVal.Value.Type(), "")
		phi.AddIncoming([]llvm.Value{*successVal.Value}, []llvm.BasicBlock{thenBlock})
	}

	if failureVal != nil {
		phi := ctx.Builder.CreatePHI(successVal.Value.Type(), "")
		phi.AddIncoming([]llvm.Value{*successVal.Value}, []llvm.BasicBlock{thenBlock})
	}

	thenBlock.MoveAfter(bodyBlock)
	elseBlock.MoveAfter(thenBlock)
	mergeBlock.MoveAfter(thenBlock)

	return nil, &CompilerResult{Value: &phi}
}

func (cs ConditionalStatetement) Accept(g CodeGenerator) error {
	return g.VisitConditionalStatement(&cs)
}

func (expr ConditionalStatetement) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (cs ConditionalStatetement) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["success"] = cs.Success
	m["condition"] = cs.Condition
	m["failure"] = cs.Failure

	res := make(map[string]any)
	res["ast.ConditionalStatetement"] = m

	return json.Marshal(res)
}
