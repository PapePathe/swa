package compiler

import (
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitWhileStatement(node *ast.WhileStatement) error {
	bodyBlock := g.Ctx.Builder.GetInsertBlock()
	parentFunc := bodyBlock.Parent()

	whileConditionBlock := g.Ctx.Context.AddBasicBlock(parentFunc, "while.cond")
	whileBodyBlock := g.Ctx.Context.AddBasicBlock(parentFunc, "while.body")
	whileMergeBlock := g.Ctx.Context.AddBasicBlock(parentFunc, "while.merge")

	g.Ctx.Builder.CreateBr(whileConditionBlock)

	g.Ctx.Builder.SetInsertPointAtEnd(whileConditionBlock)

	err := node.Condition.Accept(g)
	if err != nil {
		return err
	}

	condition := g.getLastResult()

	g.Ctx.Builder.CreateCondBr(*condition.Value, whileBodyBlock, whileMergeBlock)
	g.Ctx.Builder.SetInsertPointAtEnd(whileBodyBlock)

	err = node.Body.Accept(g)
	if err != nil {
		return err
	}

	lastBodyBlock := g.Ctx.Builder.GetInsertBlock()

	opcode := lastBodyBlock.LastInstruction().InstructionOpcode()
	if lastBodyBlock.LastInstruction().IsNil() || (opcode != llvm.Ret && opcode != llvm.Br) {
		g.Ctx.Builder.CreateBr(whileConditionBlock)
	}

	g.Ctx.Builder.SetInsertPointAtEnd(whileMergeBlock)

	return nil
}
