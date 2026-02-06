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

func (g *LLVMGenerator) VisitTupleAssignmentExpression(node *ast.TupleAssignmentExpression) error {
	old := g.logger.Step("VisitTupleAssignmentExpression")
	defer g.logger.Restore(old)

	// 1. Evaluate RHS
	err := node.Value.Accept(g)
	if err != nil {
		return err
	}

	typ, ok := node.Value.VisitedSwaType().(*ast.TupleType)
	if ok {
		l1 := len(typ.Types)
		l2 := len(node.Assignees.Expressions)

		if l1 != l2 {
			return fmt.Errorf("length of tuples does not match expected %d got %d", l1, l2)
		}
	}

	rhsRes := g.getLastResult()
	if rhsRes == nil || rhsRes.Value == nil {
		return fmt.Errorf("RHS of tuple assignment evaluated to nil")
	}

	rhsVal := *rhsRes.Value

	// 2. Iterate Assignees
	for i, expr := range node.Assignees.Expressions {
		// Extract value at index i
		elemVal := g.Ctx.Builder.CreateExtractValue(rhsVal, i, "")

		// Store into assignee
		switch a := expr.(type) {
		case *ast.MemberExpression, *ast.ArrayAccessExpression,
			*ast.ArrayOfStructsAccessExpression:
			err := expr.Accept(g)
			if err != nil {
				return err
			}

			res := g.getLastResult()
			g.Ctx.Builder.CreateStore(elemVal, *res.Value)
		case *ast.SymbolExpression:
			err, entry := g.Ctx.FindSymbol(a.Value)
			if err != nil {
				return err
			}

			if entry.Address == nil {
				format := "cannot assign to symbol %s which has no address"

				return fmt.Errorf(format, a.Value)
			}

			a.SwaType = entry.DeclaredType

			g.Ctx.Builder.CreateStore(elemVal, *entry.Address)
		default:
			return fmt.Errorf("Unsupported assignee type in tuple assignment: %T", a)
		}
	}

	return nil
}
