package compiler

import (
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitConditionalStatement(node *ast.ConditionalStatetement) error {
	err := node.Condition.Accept(g)
	if err != nil {
		return err
	}

	condition := g.getLastResult()
	bodyBlock := g.Ctx.Builder.GetInsertBlock()
	parentFunc := bodyBlock.Parent()
	mergeBlock := g.Ctx.Context.AddBasicBlock(parentFunc, "merge")
	thenBlock := g.Ctx.Context.AddBasicBlock(parentFunc, "if")
	elseBlock := g.Ctx.Context.AddBasicBlock(parentFunc, "else")

	g.Ctx.Builder.CreateCondBr(*condition.Value, thenBlock, elseBlock)
	g.Ctx.Builder.SetInsertPointAtEnd(thenBlock)

	err = node.Success.Accept(g)
	if err != nil {
		return err
	}

	successVal := g.getLastResult()

	if thenBlock.LastInstruction().InstructionOpcode() != llvm.Ret {
		g.Ctx.Builder.CreateBr(mergeBlock)
	}

	g.Ctx.Builder.SetInsertPointAtEnd(elseBlock)

	err = node.Failure.Accept(g)
	if err != nil {
		return err
	}

	failureVal := g.getLastResult()

	var phi llvm.Value

	// When there is no else block, the last instruction is nil
	// so we need to account for that and branch it to the merge block
	if elseBlock.LastInstruction().IsNil() {
		g.Ctx.Builder.CreateBr(mergeBlock)
	} else {
		if elseBlock.LastInstruction().InstructionOpcode() != llvm.Ret {
			g.Ctx.Builder.CreateBr(mergeBlock)
		}
	}

	g.Ctx.Builder.SetInsertPointAtEnd(mergeBlock)

	if successVal != nil && failureVal != nil {
		phi = g.Ctx.Builder.CreatePHI(successVal.Value.Type(), "")
		phi.AddIncoming(
			[]llvm.Value{*successVal.Value, *failureVal.Value},
			[]llvm.BasicBlock{thenBlock, elseBlock},
		)
	}

	thenBlock.MoveAfter(bodyBlock)
	elseBlock.MoveAfter(thenBlock)
	mergeBlock.MoveAfter(elseBlock)

	if successVal != nil && failureVal != nil {
		g.setLastResult(&ast.CompilerResult{Value: &phi})
	} else {
		g.setLastResult(nil)
	}

	return nil
}
