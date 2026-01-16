package compiler

import (
	"fmt"
	"reflect"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

type ElementInjector func(
	g *LLVMGenerator,
	expr ast.Expression,
	targetAddr llvm.Value,
) (error, *StructSymbolTableEntry)

var ArrayInitializationExpressionInjectors = map[reflect.Type]ElementInjector{
	reflect.TypeFor[ast.SymbolExpression]():               injectSymbol,
	reflect.TypeFor[ast.NumberExpression]():               injectLiteral,
	reflect.TypeFor[ast.FloatExpression]():                injectLiteral,
	reflect.TypeFor[ast.StringExpression]():               injectLiteral,
	reflect.TypeFor[ast.StructInitializationExpression](): injectStruct,
}

func (g *LLVMGenerator) VisitArrayInitializationExpression(node *ast.ArrayInitializationExpression) error {
	old := g.logger.Step("ArrayInitExpr")

	defer g.logger.Restore(old)

	if !g.Ctx.InsideFunction {
		return fmt.Errorf("array initialization should happen inside a function")
	}

	err := node.Underlying.Accept(g)
	if err != nil {
		return err
	}

	llvmtyp := g.getLastTypeVisitResult()
	arrayPointer := g.Ctx.Builder.CreateAlloca(llvmtyp.Type, "array_alloc")

	var discoveredEntry *StructSymbolTableEntry

	for i, expr := range node.Contents {
		itemGep := g.Ctx.Builder.CreateGEP(llvmtyp.Type, arrayPointer, []llvm.Value{
			llvm.ConstInt(llvm.GlobalContext().Int32Type(), 0, false),
			llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(i), false),
		}, "")

		injector, ok := ArrayInitializationExpressionInjectors[reflect.TypeOf(expr)]
		if !ok {
			return fmt.Errorf("unsupported array initialization element: %T", expr)
		}

		err, sEntry := injector(g, expr, itemGep)
		if err != nil {
			return err
		}

		if discoveredEntry == nil && sEntry != nil {
			discoveredEntry = sEntry
		}
	}

	g.setLastResult(&CompilerResult{
		Value: &arrayPointer,
		ArraySymbolTableEntry: &ArraySymbolTableEntry{
			ElementsCount:     llvmtyp.Type.ArrayLength(),
			UnderlyingTypeDef: discoveredEntry, // Correctly passed up!
			UnderlyingType:    llvmtyp.SubType,
			Type:              llvmtyp.Type,
		},
	})

	return nil
}

func injectLiteral(g *LLVMGenerator, expr ast.Expression, targetAddr llvm.Value) (error, *StructSymbolTableEntry) {
	err := expr.Accept(g)
	if err != nil {
		return err, nil
	}

	res := g.getLastResult()
	g.Ctx.Builder.CreateStore(*res.Value, targetAddr)

	return nil, nil // Literals don't define a struct subtype
}

func injectSymbol(g *LLVMGenerator, expr ast.Expression, targetAddr llvm.Value) (error, *StructSymbolTableEntry) {
	err := expr.Accept(g)
	if err != nil {
		return err, nil
	}

	res := g.getLastResult()

	// If it's a struct/complex type, we use the Ref (SymbolTableEntry) to find the type
	var loadType llvm.Type

	var structEntry *StructSymbolTableEntry

	if res.SymbolTableEntry != nil && res.SymbolTableEntry.Ref != nil {
		loadType = res.SymbolTableEntry.Ref.LLVMType
		structEntry = res.SymbolTableEntry.Ref
	} else {
		loadType = res.Value.Type()
	}

	var val llvm.Value
	if res.SymbolTableEntry.Address != nil {
		val = g.Ctx.Builder.CreateLoad(loadType, *res.SymbolTableEntry.Address, "arr.load.sym.from-address")
	} else {
		val = g.Ctx.Builder.CreateLoad(loadType, *res.Value, "arr.load.sym")
	}

	g.Ctx.Builder.CreateStore(val, targetAddr)

	return nil, structEntry
}

func injectStruct(g *LLVMGenerator, expr ast.Expression, targetAddr llvm.Value) (error, *StructSymbolTableEntry) {
	node, _ := expr.(ast.StructInitializationExpression)

	err, tblEntry := g.Ctx.FindStructSymbol(node.Name)
	if err != nil {
		return err, nil
	}

	for _, fieldName := range node.Properties {
		err, idx := tblEntry.Metadata.PropertyIndex(fieldName)
		if err != nil {
			return err, nil
		}

		fieldNode := node.Values[idx]

		err = fieldNode.Accept(g)
		if err != nil {
			return err, nil
		}

		fieldRes := g.getLastResult()

		fieldGep := g.Ctx.Builder.CreateGEP(
			tblEntry.LLVMType,
			targetAddr,
			[]llvm.Value{
				llvm.ConstInt(llvm.GlobalContext().Int32Type(), 0, false),
				llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(idx), false),
			},
			fmt.Sprintf("field.%s", fieldName),
		)

		g.Ctx.Builder.CreateStore(*fieldRes.Value, fieldGep)
	}

	return nil, tblEntry
}
