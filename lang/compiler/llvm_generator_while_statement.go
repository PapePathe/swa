package compiler

import (
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitWhileStatement(node *ast.WhileStatement) error {
	old := g.logger.Step("WhileStmt")

	defer g.logger.Restore(old)

	g.Debugf("%s", node.Condition)

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

	// TODO why on earth are we doing this
	lastBodyBlock := g.Ctx.Builder.GetInsertBlock()
	opcode := lastBodyBlock.LastInstruction()

	if lastBodyBlock.LastInstruction().IsNil() ||
		(!opcode.IsNil() && opcode.InstructionOpcode() != llvm.Ret && opcode.InstructionOpcode() != llvm.Br) {
		g.Ctx.Builder.CreateBr(whileConditionBlock)
	}

	g.Ctx.Builder.SetInsertPointAtEnd(whileMergeBlock)

	return nil
}
