package compiler

import (
	"fmt"
	"reflect"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
	if !g.Ctx.InsideFunction {
		return fmt.Errorf("struct initialization should happen inside a function")
	}

	err, structType := g.Ctx.FindStructSymbol(node.Name)
	if err != nil {
		return err
	}

	instance := g.Ctx.Builder.CreateAlloca(structType.LLVMType, node.Name+".instance")

	for _, name := range node.Properties {
		err, propIndex := structType.Metadata.PropertyIndex(name)
		if err != nil {
			return err
		}

		expr := node.Values[propIndex]

		err = expr.Accept(g)
		if err != nil {
			return err
		}

		res := g.getLastResult()
		fieldAddr := g.Ctx.Builder.CreateStructGEP(
			structType.LLVMType,
			instance,
			propIndex,
			fmt.Sprintf("%s.%s", node.Name, name),
		)

		injector, ok := structInjectors[reflect.TypeOf(expr)]
		if !ok {
			return fmt.Errorf("struct field initialization unimplemented for %T", expr)
		}

		targetType := structType.PropertyTypes[propIndex]
		injector(g, res, fieldAddr, targetType)
	}

	g.setLastResult(&ast.CompilerResult{
		Value:                  &instance,
		StructSymbolTableEntry: structType,
	})

	return nil
}

type StructInitializationExpressionPropertyInjector func(
	g *LLVMGenerator,
	res *ast.CompilerResult,
	fieldAddr llvm.Value,
	targetType llvm.Type,
)

var structInjectors = map[reflect.Type]StructInitializationExpressionPropertyInjector{
	reflect.TypeFor[ast.NumberExpression]():               injectDirectly,
	reflect.TypeFor[ast.FloatExpression]():                injectDirectly,
	reflect.TypeFor[ast.StringExpression]():               injectDirectly,
	reflect.TypeFor[ast.BinaryExpression]():               injectDirectly,
	reflect.TypeFor[ast.MemberExpression]():               injectDirectly,
	reflect.TypeFor[ast.SymbolExpression]():               injectWithArrayDecay,
	reflect.TypeFor[ast.ArrayInitializationExpression]():  injectArrayLiteral,
	reflect.TypeFor[ast.ArrayAccessExpression]():          injectArrayAccess,
	reflect.TypeFor[ast.StructInitializationExpression](): injectNestedStruct,
}

func injectNestedStruct(g *LLVMGenerator, res *ast.CompilerResult, fieldAddr llvm.Value, targetType llvm.Type) {
	load := g.Ctx.Builder.CreateLoad(res.Value.AllocatedType(), *res.Value, "nested-struct.load")
	g.Ctx.Builder.CreateStore(load, fieldAddr)
}

func injectDirectly(g *LLVMGenerator, res *ast.CompilerResult, fieldAddr llvm.Value, targetType llvm.Type) {
	g.Ctx.Builder.CreateStore(*res.Value, fieldAddr)
}

func injectWithArrayDecay(g *LLVMGenerator, res *ast.CompilerResult, fieldAddr llvm.Value, targetType llvm.Type) {
	_, isArray := res.SymbolTableEntry.DeclaredType.(ast.ArrayType)

	if isArray && targetType.TypeKind() == llvm.PointerTypeKind {
		ptr := g.Ctx.Builder.CreateBitCast(*res.SymbolTableEntry.Address, targetType, "decay")
		g.Ctx.Builder.CreateStore(ptr, fieldAddr)

		return
	}

	g.Ctx.Builder.CreateStore(*res.Value, fieldAddr)
}

func injectArrayLiteral(g *LLVMGenerator, res *ast.CompilerResult, fieldAddr llvm.Value, targetType llvm.Type) {
	if targetType.TypeKind() == llvm.PointerTypeKind {
		ptr := g.Ctx.Builder.CreateBitCast(*res.Value, targetType, "array.ptr.cast")
		g.Ctx.Builder.CreateStore(ptr, fieldAddr)
	} else {
		load := g.Ctx.Builder.CreateLoad(res.ArraySymbolTableEntry.Type, *res.Value, "array.load")
		g.Ctx.Builder.CreateStore(load, fieldAddr)
	}
}

func injectArrayAccess(g *LLVMGenerator, res *ast.CompilerResult, fieldAddr llvm.Value, targetType llvm.Type) {
	load := g.Ctx.Builder.CreateLoad(res.ArraySymbolTableEntry.UnderlyingType, *res.Value, "access.load")
	g.Ctx.Builder.CreateStore(load, fieldAddr)
}
