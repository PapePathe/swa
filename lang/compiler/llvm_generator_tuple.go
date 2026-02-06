package compiler

import (
	"fmt"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitTupleExpression(node *ast.TupleExpression) error {
	var values []llvm.Value
	var types []llvm.Type

	for i, expr := range node.Expressions {
		if err := expr.Accept(g); err != nil {
			return err
		}

		res := g.getLastResult()
		if res == nil {
			return fmt.Errorf("Tuple expr element %d returned nil result", i)
		}

		if res.Value == nil {
			return fmt.Errorf("Tuple expr element %d returned result with nil Value", i)
		}

		values = append(values, *res.Value)
		types = append(types, res.Value.Type())
	}

	structType := g.Ctx.Context.StructType(types, false)
	structVal := llvm.Undef(structType)

	for i, val := range values {
		structVal = g.Ctx.Builder.CreateInsertValue(structVal, val, i, "")
	}

	g.setLastResult(&CompilerResult{
		Value: &structVal,
	})
	return nil
}

func (g *LLVMGenerator) VisitTupleType(node *ast.TupleType) error {
	// A tuple type in LLVM is a struct containing all types.
	llvmTypes := []llvm.Type{}
	for _, typ := range node.Types {
		err := typ.Accept(g)
		if err != nil {
			return err
		}

		llvmTypes = append(llvmTypes, g.getLastTypeVisitResult().Type)
	}

	structType := g.Ctx.Context.StructType(llvmTypes, false)
	res := CompilerResultType{Type: structType}

	g.setLastTypeVisitResult(&res)

	return nil
}

func (g *LLVMGenerator) ZeroOfTupleType(node *ast.TupleType) error {
	// TODO this should not be nil
	return nil
}
