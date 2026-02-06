package compiler

import (
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitSymbolType(node *ast.SymbolType) error {
	old := g.logger.Step("SymbolType")

	defer g.logger.Restore(old)

	err, entry := g.Ctx.FindStructSymbol(node.Name)
	if err != nil {
		return err
	}

	g.setLastTypeVisitResult(&CompilerResultType{
		Type:   entry.LLVMType,
		Sentry: entry,
	})

	return nil
}

func (g *LLVMGenerator) VisitNumberType(node *ast.NumberType) error {
	old := g.logger.Step("NumberType")

	defer g.logger.Restore(old)

	g.setLastTypeVisitResult(&CompilerResultType{
		Type: llvm.GlobalContext().Int32Type(),
	})

	return nil
}

func (g *LLVMGenerator) VisitNumber64Type(node *ast.Number64Type) error {
	old := g.logger.Step("Number64Type")

	defer g.logger.Restore(old)

	g.setLastTypeVisitResult(&CompilerResultType{
		Type: llvm.GlobalContext().Int64Type(),
	})

	return nil
}

func (g *LLVMGenerator) VisitArrayType(node *ast.ArrayType) error {
	old := g.logger.Step("ArrayType")

	defer g.logger.Restore(old)

	err := node.Underlying.Accept(g)
	if err != nil {
		return err
	}

	under := g.getLastTypeVisitResult()

	g.setLastTypeVisitResult(&CompilerResultType{
		Type:    llvm.ArrayType(under.Type, node.Size),
		SubType: under.Type,
		Aentry: &ArraySymbolTableEntry{
			Type:              llvm.ArrayType(under.Type, node.Size),
			UnderlyingTypeDef: under.Sentry,
			ElementsCount:     node.Size,
			UnderlyingType:    under.Type,
		},
	})

	return nil
}

func (g *LLVMGenerator) VisitFloatType(node *ast.FloatType) error {
	old := g.logger.Step("FloatType")

	defer g.logger.Restore(old)

	g.setLastTypeVisitResult(&CompilerResultType{
		Type: llvm.GlobalContext().DoubleType(),
	})

	return nil
}

func (g *LLVMGenerator) VisitPointerType(node *ast.PointerType) error {
	old := g.logger.Step("PointerType")

	defer g.logger.Restore(old)

	err := node.Underlying.Accept(g)
	if err != nil {
		return err
	}

	under := g.getLastTypeVisitResult()

	g.setLastTypeVisitResult(&CompilerResultType{
		Type:    llvm.PointerType(under.Type, 0),
		SubType: under.Type,
		Sentry:  under.Sentry,
	})

	return nil
}

func (g *LLVMGenerator) VisitStringType(node *ast.StringType) error {
	old := g.logger.Step("StringType")

	defer g.logger.Restore(old)

	g.setLastTypeVisitResult(&CompilerResultType{
		Type:    llvm.PointerType(llvm.GlobalContext().Int8Type(), 0),
		SubType: llvm.GlobalContext().Int8Type(),
	})

	return nil
}

func (g *LLVMGenerator) VisitVoidType(node *ast.VoidType) error {
	old := g.logger.Step("VoidType")

	defer g.logger.Restore(old)

	g.setLastTypeVisitResult(&CompilerResultType{
		Type: llvm.GlobalContext().VoidType(),
	})

	return nil
}

func (g *LLVMGenerator) VisitErrorType(node *ast.ErrorType) error {
	g.setLastTypeVisitResult(&CompilerResultType{
		Type:    llvm.PointerType(llvm.GlobalContext().Int8Type(), 0),
		SubType: llvm.GlobalContext().Int8Type(),
	})

	return nil
}
