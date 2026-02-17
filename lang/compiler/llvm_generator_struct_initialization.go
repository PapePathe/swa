package compiler

import (
	"fmt"
	"reflect"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
	old := g.logger.Step(fmt.Sprintf("StructInitExpr %s", node.Name))

	defer g.logger.Restore(old)

	if !g.Ctx.InsideFunction {
		key := "LLVMGenerator.VisitStructInitializationExpression.NotInsideFunction"

		return g.Ctx.Dialect.Error(key)
	}

	err, structType := g.Ctx.FindStructSymbol(node.Name)
	if err != nil {
		return err
	}

	instance := g.Ctx.Builder.CreateAlloca(structType.LLVMType, node.Name+".instance")

	if len(node.Values) == 0 {
		typ := ast.SymbolType{Name: node.Name}

		err := typ.AcceptZero(g)
		if err != nil {
			return err
		}

		zero := g.getLastResult()
		g.Ctx.Builder.CreateStore(*zero.Value, instance)
		g.setLastResult(&CompilerResult{
			Value:                  &instance,
			StructSymbolTableEntry: structType,
		})

		return nil
	}

	// TODO stucts should be allocated on the heap
	//	instance := g.Ctx.Builder.CreateMalloc(structType.LLVMType, node.Name+".instance")
	for _, name := range node.Properties {
		g.Debugf("initializing property %s", name)

		err, propIndex := structType.Metadata.PropertyIndex(name)
		if err != nil {
			return err
		}

		g.Debugf("property index: %d", propIndex)

		if propIndex >= len(node.Values) {
			return g.Ctx.Dialect.Error("All struct properties must be initialized")
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
			// TODO this is a developer error and the user should be informed
			key := "LLVMGenerator.VisitStructInitializationExpression.Unimplemented"

			return g.Ctx.Dialect.Error(key, expr)
		}

		targetType := structType.PropertyTypes[propIndex]
		injector(g, res, fieldAddr, targetType)
	}

	node.SwaType = ast.SymbolType{Name: node.Name}

	g.setLastResult(&CompilerResult{
		Value:                  &instance,
		StructSymbolTableEntry: structType,
	})

	return nil
}

type StructInitializationExpressionPropertyInjector func(
	g *LLVMGenerator,
	res *CompilerResult,
	fieldAddr llvm.Value,
	targetType llvm.Type,
)

var structInjectors = map[reflect.Type]StructInitializationExpressionPropertyInjector{
	reflect.TypeFor[*ast.ZeroExpression]():                 injectDirectly,
	reflect.TypeFor[*ast.BooleanExpression]():              injectDirectly,
	reflect.TypeFor[*ast.NumberExpression]():               injectDirectly,
	reflect.TypeFor[*ast.FloatExpression]():                injectDirectly,
	reflect.TypeFor[*ast.FunctionCallExpression]():         injectDirectly,
	reflect.TypeFor[*ast.StringExpression]():               injectDirectly,
	reflect.TypeFor[*ast.BinaryExpression]():               injectDirectly,
	reflect.TypeFor[*ast.MemberExpression]():               injectDirectly,
	reflect.TypeFor[*ast.SymbolExpression]():               injectWithArrayDecay,
	reflect.TypeFor[*ast.ArrayInitializationExpression]():  injectArrayLiteral,
	reflect.TypeFor[*ast.ArrayAccessExpression]():          injectArrayAccess,
	reflect.TypeFor[*ast.StructInitializationExpression](): injectNestedStruct,
}

func injectNestedStruct(g *LLVMGenerator, res *CompilerResult, fieldAddr llvm.Value, targetType llvm.Type) {
	load := g.Ctx.Builder.CreateLoad(res.Value.AllocatedType(), *res.Value, "nested-struct.load")
	g.Ctx.Builder.CreateStore(load, fieldAddr)
}

func injectDirectly(g *LLVMGenerator, res *CompilerResult, fieldAddr llvm.Value, targetType llvm.Type) {
	g.Ctx.Builder.CreateStore(*res.Value, fieldAddr)
}

func injectWithArrayDecay(g *LLVMGenerator, res *CompilerResult, fieldAddr llvm.Value, targetType llvm.Type) {
	_, isArray := res.SymbolTableEntry.DeclaredType.(ast.ArrayType)

	if isArray && targetType.TypeKind() == llvm.PointerTypeKind {
		ptr := g.Ctx.Builder.CreateBitCast(*res.SymbolTableEntry.Address, targetType, "decay")
		g.Ctx.Builder.CreateStore(ptr, fieldAddr)

		return
	}

	g.Ctx.Builder.CreateStore(*res.Value, fieldAddr)
}

func injectArrayLiteral(g *LLVMGenerator, res *CompilerResult, fieldAddr llvm.Value, targetType llvm.Type) {
	g.Debugf("injectArrayLiteral %s %v", fieldAddr, res)

	if targetType.TypeKind() == llvm.PointerTypeKind {
		g.Debugf("injectArrayLiteral of pointer type")

		ptr := g.Ctx.Builder.CreateBitCast(*res.Value, targetType, "array.ptr.cast")
		g.Ctx.Builder.CreateStore(ptr, fieldAddr)
	} else if targetType.TypeKind() == llvm.ArrayTypeKind {
		g.Debugf("injectArrayLiteral of array type")

		load := g.Ctx.Builder.CreateLoad(res.Value.AllocatedType(), *res.Value, "array.load.v1")
		g.Ctx.Builder.CreateStore(load, fieldAddr)
	} else {
		g.Debugf("injectArrayLiteral of other type")

		load := g.Ctx.Builder.CreateLoad(res.ArraySymbolTableEntry.Type, *res.Value, "array.load")
		g.Ctx.Builder.CreateStore(load, fieldAddr)
	}
}

func injectArrayAccess(g *LLVMGenerator, res *CompilerResult, fieldAddr llvm.Value, targetType llvm.Type) {
	load := g.Ctx.Builder.CreateLoad(res.ArraySymbolTableEntry.UnderlyingType, *res.Value, "access.load")
	g.Ctx.Builder.CreateStore(load, fieldAddr)
}
