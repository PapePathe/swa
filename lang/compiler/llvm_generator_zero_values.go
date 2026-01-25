package compiler

import (
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

// ZeroOfArrayType implements [ast.CodeGenerator].
func (g *LLVMGenerator) ZeroOfArrayType(node *ast.ArrayType) error {
	err := node.Underlying.AcceptZero(g)
	if err != nil {
		return err
	}

	lastres := g.getLastResult()

	err = g.VisitArrayType(node)
	if err != nil {
		return err
	}

	llvmtyp := g.getLastTypeVisitResult()
	arrayPointer := g.Ctx.Builder.CreateAlloca(llvmtyp.Type, "array_alloc")

	for i := range node.Size {
		itemGep := g.Ctx.Builder.CreateGEP(llvmtyp.Type, arrayPointer, []llvm.Value{
			llvm.ConstInt(llvm.GlobalContext().Int32Type(), 0, false),
			llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(i), false),
		}, "")

		g.Ctx.Builder.CreateStore(*lastres.Value, itemGep)
	}

	load := g.Ctx.Builder.CreateLoad(llvmtyp.Type, arrayPointer, "")

	res := &CompilerResult{
		Value:                 &load,
		ArraySymbolTableEntry: llvmtyp.Aentry,
	}

	g.Debugf("type visit result %+v %+v", g.formatLLVMType(llvmtyp.Type), llvmtyp.Aentry)

	if llvmtyp.Aentry.UnderlyingTypeDef != nil {
		res.StructSymbolTableEntry = llvmtyp.Aentry.UnderlyingTypeDef
	}

	g.setLastResult(res)

	return nil
}

// ZeroOfFloatType implements [ast.CodeGenerator].
func (g *LLVMGenerator) ZeroOfFloatType(node *ast.FloatType) error {
	err := node.Accept(g)
	if err != nil {
		return err
	}

	lastres := g.getLastTypeVisitResult()
	zero := llvm.ConstNull(lastres.Type)
	res := &CompilerResult{Value: &zero}

	g.setLastResult(res)

	return nil
}

// ZeroOfNumber64Type implements [ast.CodeGenerator].
func (g *LLVMGenerator) ZeroOfNumber64Type(node *ast.Number64Type) error {
	err := node.Accept(g)
	if err != nil {
		return err
	}

	lastres := g.getLastTypeVisitResult()
	zero := llvm.ConstNull(lastres.Type)
	res := &CompilerResult{Value: &zero}

	g.setLastResult(res)

	return nil
}

// ZeroOfNumberType implements [ast.CodeGenerator].
func (g *LLVMGenerator) ZeroOfNumberType(node *ast.NumberType) error {
	err := node.Accept(g)
	if err != nil {
		return err
	}

	lastres := g.getLastTypeVisitResult()
	zero := llvm.ConstNull(lastres.Type)
	res := &CompilerResult{Value: &zero}

	g.setLastResult(res)

	return nil
}

// ZeroOfPointerType implements [ast.CodeGenerator].
func (g *LLVMGenerator) ZeroOfPointerType(node *ast.PointerType) error {
	err := node.Accept(g)
	if err != nil {
		return err
	}

	lastres := g.getLastTypeVisitResult()
	zero := llvm.ConstNull(lastres.Type)
	res := &CompilerResult{Value: &zero}

	g.setLastResult(res)

	return nil
}

// ZeroOfStringType implements [ast.CodeGenerator].
func (g *LLVMGenerator) ZeroOfStringType(node *ast.StringType) error {
	zero := g.Ctx.Builder.CreateGlobalStringPtr("", "")
	res := &CompilerResult{Value: &zero}

	g.setLastResult(res)

	return nil
}

// ZeroOfSymbolType implements [ast.CodeGenerator].
func (g *LLVMGenerator) ZeroOfSymbolType(node *ast.SymbolType) error {
	err := node.Accept(g)
	if err != nil {
		return err
	}

	typ := g.getLastTypeVisitResult()
	alloc := g.Ctx.Builder.CreateAlloca(typ.Type, "")

	for index, v := range typ.Sentry.Metadata.Types {
		err := v.AcceptZero(g)
		if err != nil {
			return err
		}

		lastres := g.getLastResult()
		fieldAddr := g.Ctx.Builder.CreateStructGEP(
			typ.Type,
			alloc,
			index,
			"",
		)

		g.Ctx.Builder.CreateStore(*lastres.Value, fieldAddr)
	}
	load := g.Ctx.Builder.CreateLoad(alloc.AllocatedType(), alloc, "")

	res := &CompilerResult{Value: &load, StructSymbolTableEntry: typ.Sentry}
	g.setLastResult(res)

	return nil
}

// ZeroOfVoidType implements [ast.CodeGenerator].
func (g *LLVMGenerator) ZeroOfVoidType(node *ast.VoidType) error {
	panic("unimplemented")
}
