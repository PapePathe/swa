package compiler

import (
	"fmt"
	"reflect"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

// VisitVarDeclaration implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitVarDeclaration(node *ast.VarDeclarationStatement) error {
	old := g.logger.Step("VarDeclStmt")

	defer g.logger.Restore(old)

	g.Debugf(node.Name)

	if g.Ctx.SymbolExistsInCurrentScope(node.Name) {
		key := "LLVMGenerator.VisitVarDeclaration.AlreadyExisits"

		return g.Ctx.Dialect.Error(key, node.Name)
	}

	switch node.Value {
	case nil:
		g.Debugf("Initializing with Zero Value")

		err := node.ExplicitType.Accept(g)
		if err != nil {
			return err
		}

		typeresult := g.getLastTypeVisitResult()
		alloc := g.Ctx.Builder.CreateAlloca(typeresult.Type, fmt.Sprintf("alloc.%s", node.Name))

		err = node.ExplicitType.AcceptZero(g)
		if err != nil {
			return err
		}

		zero := g.getLastResult()

		g.Ctx.Builder.CreateStore(*zero.Value, alloc)

		entry := &SymbolTableEntry{
			Value:        alloc,
			Address:      &alloc,
			DeclaredType: node.ExplicitType,
			Global:       !g.Ctx.InsideFunction,
		}

		if zero.StructSymbolTableEntry != nil {
			entry.Ref = zero.StructSymbolTableEntry
		}

		g.Debugf("setting symbol table entry %v", entry)

		err = g.Ctx.AddSymbol(node.Name, entry)
		if err != nil {
			return err
		}

		if typeresult.Aentry != nil {
			g.Debugf("setting array symbol table entry %v", typeresult.Aentry)

			err := g.Ctx.AddArraySymbol(node.Name, typeresult.Aentry)
			if err != nil {
				return err
			}
		}

		g.Debugf("Initializing with Zero Value finished")

		return nil
	default:
		return g.declareVarWithInitializer(node)
	}
}

type InitializationStyle int

const (
	StyleDefault   InitializationStyle = iota // Alloca + Store
	StyleDirect                               // Use compiled value as address (Literals)
	StyleLoadStore                            // Load from pointer, then Alloca + Store (Accessors)
)

var nodeVariableDeclarationStyles = map[reflect.Type]InitializationStyle{
	reflect.TypeFor[*ast.ArrayAccessExpression]():          StyleLoadStore,
	reflect.TypeFor[*ast.ArrayOfStructsAccessExpression](): StyleLoadStore,
	reflect.TypeFor[*ast.MemberExpression]():               StyleLoadStore,
	reflect.TypeFor[*ast.StringExpression]():               StyleDefault,
	reflect.TypeFor[*ast.NumberExpression]():               StyleDefault,
	reflect.TypeFor[*ast.FloatExpression]():                StyleDefault,
	reflect.TypeFor[*ast.FunctionCallExpression]():         StyleDefault,
	reflect.TypeFor[*ast.SymbolExpression]():               StyleDefault,
	reflect.TypeFor[*ast.BinaryExpression]():               StyleDefault,
	reflect.TypeFor[*ast.StructInitializationExpression](): StyleDirect,
	reflect.TypeFor[*ast.ArrayInitializationExpression]():  StyleDirect,
	reflect.TypeFor[*ast.PrefixExpression]():               StyleDefault,
	reflect.TypeFor[*ast.ZeroExpression]():                 StyleDefault,
	reflect.TypeFor[*ast.ErrorExpression]():                StyleDefault,
}

func (g *LLVMGenerator) declareVarWithInitializer(node *ast.VarDeclarationStatement) error {
	err := node.Value.Accept(g)
	if err != nil {
		return err
	}

	res := g.getLastResult()

	style, ok := nodeVariableDeclarationStyles[reflect.TypeOf(node.Value)]
	if !ok {
		key := "LLVMGenerator.VisitVarDeclaration.UnsupportedInitializerType"

		return g.Ctx.Dialect.Error(key, node.Value)
	}

	var finalAddr *llvm.Value

	switch style {
	case StyleDirect:
		if !g.Ctx.InsideFunction {
			key := "LLVMGenerator.VisitVarDeclaration.UnsupportedTypeAsGlobal"

			return g.Ctx.Dialect.Error(key, node.Value)
		}

		finalAddr = res.Value

	case StyleLoadStore:
		if !g.Ctx.InsideFunction {
			key := "LLVMGenerator.VisitVarDeclaration.NotInsideFunction"

			return g.Ctx.Dialect.Error(key, node.Value)
		}

		// Logic for extracting value from an accessor
		typ := res.Value.AllocatedType()
		if _, ok := node.Value.(*ast.ArrayOfStructsAccessExpression); ok {
			typ = *res.StuctPropertyValueType
		}
		loadedVal := g.Ctx.Builder.CreateLoad(typ, *res.Value, "tmp.load")

		alloc := g.Ctx.Builder.CreateAlloca(typ, node.Name)
		g.Ctx.Builder.CreateStore(loadedVal, alloc)
		finalAddr = &alloc

	default: // StyleDefault
		if g.Ctx.InsideFunction {
			alloc := g.Ctx.Builder.CreateAlloca(res.Value.Type(), node.Name)
			g.Ctx.Builder.CreateStore(*res.Value, alloc)
			finalAddr = &alloc
		} else {
			glob := llvm.AddGlobal(*g.Ctx.Module, res.Value.Type(), node.Name)
			glob.SetInitializer(*res.Value)
			finalAddr = &glob
		}
	}

	// Unified Metadata Management
	return g.finalizeSymbol(node, finalAddr, res)
}

// Helper to keep metadata logic separate from IR generation logic
func (g *LLVMGenerator) finalizeSymbol(
	node *ast.VarDeclarationStatement,
	addr *llvm.Value,
	res *CompilerResult,
) error {
	entry := &SymbolTableEntry{
		Address:      addr,
		DeclaredType: node.ExplicitType,
		Global:       !g.Ctx.InsideFunction,
		Ref:          res.StructSymbolTableEntry,
	}

	if res.ArraySymbolTableEntry != nil && res.ArraySymbolTableEntry.UnderlyingTypeDef != nil {
		entry.Ref = res.ArraySymbolTableEntry.UnderlyingTypeDef
	}

	if res.SymbolTableEntry != nil && res.SymbolTableEntry.Ref != nil {
		entry.Ref = res.SymbolTableEntry.Ref
	}

	if entry.Ref == nil &&
		node.ExplicitType.Value() == ast.DataTypeSymbol {
		err := node.ExplicitType.Accept(g)
		if err != nil {
			return err
		}

		compiledres := g.getLastTypeVisitResult()

		if compiledres.Sentry != nil {
			entry.Ref = compiledres.Sentry
		}
	}

	err := g.Ctx.AddSymbol(node.Name, entry)
	if err != nil {
		return err
	}

	if _, ok := node.Value.(*ast.ArrayInitializationExpression); ok {
		return g.Ctx.AddArraySymbol(node.Name, res.ArraySymbolTableEntry)
	}

	return nil
}
