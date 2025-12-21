package compiler

import (
	"fmt"
	"math"
	"os"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

type LLVMGenerator struct {
	Ctx        *ast.CompilerCtx
	lastResult *ast.CompilerResult
}

func (g *LLVMGenerator) setLastResult(res *ast.CompilerResult) {
	g.lastResult = res
}

func (g *LLVMGenerator) getLastResult() *ast.CompilerResult {
	res := g.lastResult
	g.lastResult = nil
	return res
}

// VisitArrayInitializationExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitArrayInitializationExpression(node *ast.ArrayInitializationExpression) error {
	err, llvmtyp := node.Underlying.LLVMType(g.Ctx)
	if err != nil {
		return err
	}

	arrayPointer := g.Ctx.Builder.CreateAlloca(llvmtyp, "")
	for i, v := range node.Contents {
		itemGep := g.Ctx.Builder.CreateGEP(
			llvmtyp,
			arrayPointer,
			[]llvm.Value{
				llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(0), false),
				llvm.ConstInt(llvm.GlobalContext().Int32Type(), uint64(i), false),
			},
			"",
		)
		switch v.(type) {
		case ast.StringExpression:
			err := v.Accept(g)
			if err != nil {
				return err
			}
			compiledVal := g.getLastResult()
			g.Ctx.Builder.CreateStore(*compiledVal.Value, itemGep)
		case ast.NumberExpression, ast.FloatExpression:
			err := v.Accept(g)
			if err != nil {
				return err
			}
			compiledVal := g.getLastResult()
			g.Ctx.Builder.CreateStore(*compiledVal.Value, itemGep)
		default:
			g.NotImplemented(fmt.Sprintf("VisitArrayInitializationExpression: Expression %s not supported", v))
		}
	}

	result := &ast.CompilerResult{
		Value: &arrayPointer,
		ArraySymbolTableEntry: &ast.ArraySymbolTableEntry{
			ElementsCount:  llvmtyp.ArrayLength(),
			UnderlyingType: llvmtyp.ElementType(),
			Type:           llvmtyp,
		},
	}

	g.setLastResult(result)

	return nil
}

// VisitArrayOfStructsAccessExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitArrayOfStructsAccessExpression(node *ast.ArrayOfStructsAccessExpression) error {
	panic("unimplemented")
}

// VisitAssignmentExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitAssignmentExpression(node *ast.AssignmentExpression) error {
	panic("unimplemented")
}

// VisitBinaryExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitBinaryExpression(node *ast.BinaryExpression) error {
	panic("unimplemented")
}

// VisitBlockStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitBlockStatement(node *ast.BlockStatement) error {
	for _, v := range node.Body {
		v.Accept(g)
	}

	return nil
}

// VisitCallExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitCallExpression(node *ast.CallExpression) error {
	panic("unimplemented")
}

// VisitConditionalStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitConditionalStatement(node *ast.ConditionalStatetement) error {
	panic("unimplemented")
}

// VisitExpressionStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitExpressionStatement(node *ast.ExpressionStatement) error {
	panic("unimplemented")
}

// VisitFloatExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitFloatExpression(node *ast.FloatExpression) error {
	res := llvm.ConstFloat(
		g.Ctx.Context.DoubleType(),
		node.Value,
	)
	g.setLastResult(&ast.CompilerResult{Value: &res})

	return nil
}

// VisitFunctionCall implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitFunctionCall(node *ast.FunctionCallExpression) error {
	panic("unimplemented")
}

// VisitFunctionDefinition implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitFunctionDefinition(node *ast.FuncDeclStatement) error {
	panic("unimplemented")
}

// VisitMainStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitMainStatement(node *ast.MainStatement) error {
	fnType := llvm.FunctionType(
		g.Ctx.Context.Int32Type(),
		[]llvm.Type{},
		false,
	)
	fn := llvm.AddFunction(*g.Ctx.Module, "main", fnType)
	block := g.Ctx.Context.AddBasicBlock(fn, "entry")
	g.Ctx.Builder.SetInsertPointAtEnd(block)

	return node.Body.Accept(g)
}

// VisitMemberExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitMemberExpression(node *ast.MemberExpression) error {
	panic("unimplemented")
}

// VisitNumberExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitNumberExpression(node *ast.NumberExpression) error {
	if node.Value < math.MinInt32 {
		return g.Ctx.Dialect.Error("NumberExpression.LessThanMinInt32", node.Value)
	}

	if node.Value > math.MaxInt32 {
		return g.Ctx.Dialect.Error("NumberExpression.GreaterThanMaxInt32", node.Value)
	}

	var signed bool

	if node.Value < 0 {
		signed = true
	}

	res := llvm.ConstInt(
		llvm.GlobalContext().Int32Type(),
		uint64(node.Value),
		signed,
	)
	g.setLastResult(&ast.CompilerResult{Value: &res})

	return nil
}

// VisitPrefixExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitPrefixExpression(node *ast.PrefixExpression) error {
	panic("unimplemented")
}

func (g *LLVMGenerator) NotImplemented(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

// VisitPrintStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitPrintStatement(node *ast.PrintStatetement) error {
	printableValues := []llvm.Value{}

	for _, v := range node.Values {
		err := v.Accept(g)
		if err != nil {
			return err
		}

		lastResult := g.getLastResult()

		switch v.(type) {
		case ast.StringExpression:
			printableValues = append(printableValues, *lastResult.Value)
		case ast.NumberExpression, ast.FloatExpression:
			printableValues = append(printableValues, *lastResult.Value)
		case ast.SymbolExpression:
			printableValues = append(printableValues, *lastResult.Value)
		default:
			format := "VisitPrintStatement unimplemented for %T"
			g.NotImplemented(fmt.Sprintf(format, v))
		}

	}

	g.Ctx.Builder.CreateCall(
		llvm.FunctionType(
			g.Ctx.Context.Int32Type(),
			[]llvm.Type{llvm.PointerType(g.Ctx.Context.Int8Type(), 0)},
			true,
		),
		g.Ctx.Module.NamedFunction("printf"),
		printableValues,
		"call.printf",
	)

	return nil
}

// VisitReturnStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitReturnStatement(node *ast.ReturnStatement) error {
	switch node.Value.(type) {
	case ast.NumberExpression:
		ret := llvm.ConstInt(g.Ctx.Context.Int32Type(), uint64(node.Value.(ast.NumberExpression).Value), false)
		g.Ctx.Builder.CreateRet(ret)
	default:
		panic("VisitReturnStatement unimplemented")
	}

	return nil
}

// VisitStringExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitStringExpression(node *ast.StringExpression) error {
	valuePtr := g.Ctx.Builder.CreateGlobalStringPtr(node.Value, "")
	res := ast.CompilerResult{Value: &valuePtr}

	g.setLastResult(&res)

	return nil
}

// VisitStructDeclaration implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitStructDeclaration(node *ast.StructDeclarationStatement) error {
	panic("unimplemented")
}

// VisitStructInitializationExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
	panic("unimplemented")
}

// VisitSymbolExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitSymbolExpression(node *ast.SymbolExpression) error {
	err, entry := g.Ctx.FindSymbol(node.Value)
	if err != nil {
		return err
	}
	var loadedValue llvm.Value

	switch entry.DeclaredType.(type) {
	case ast.StringType:
		loadedValue = g.Ctx.Builder.CreateLoad(entry.Address.Type(), *entry.Address, "")
	case ast.NumberType, ast.FloatType:
		loadedValue = g.Ctx.Builder.CreateLoad(entry.Address.AllocatedType(), *entry.Address, "")
	default:
		g.NotImplemented(fmt.Sprintf("VisitSymbolExpression unimplemented for %T", entry.DeclaredType))
	}

	g.setLastResult(
		&ast.CompilerResult{
			Value:            &loadedValue,
			SymbolTableEntry: entry,
		},
	)

	return nil
}

// VisitVarDeclaration implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitVarDeclaration(node *ast.VarDeclarationStatement) error {
	switch node.Value {
	case nil:
		return g.declareVarWithZeroValue(node)
	default:
		return g.declareVarWithInitializer(node)
	}
}

func (g *LLVMGenerator) declareVarWithInitializer(node *ast.VarDeclarationStatement) error {
	err := node.Value.Accept(g)
	if err != nil {
		return err
	}

	compiledVal := g.getLastResult()
	name := fmt.Sprintf("alloc.%s", node.Name)

	switch node.Value.(type) {
	case ast.StringExpression:
		alloc := g.Ctx.Builder.CreateAlloca(compiledVal.Value.Type(), name)
		g.Ctx.Builder.CreateStore(*compiledVal.Value, alloc)

		entry := &ast.SymbolTableEntry{
			Address:      &alloc,
			DeclaredType: node.ExplicitType,
		}

		return g.Ctx.AddSymbol(node.Name, entry)
	case ast.NumberExpression, ast.FloatExpression:
		alloc := g.Ctx.Builder.CreateAlloca(compiledVal.Value.Type(), name)
		g.Ctx.Builder.CreateStore(*compiledVal.Value, alloc)

		entry := &ast.SymbolTableEntry{
			Address:      &alloc,
			DeclaredType: node.ExplicitType,
		}
		if compiledVal.SymbolTableEntry != nil {
			entry.Ref = compiledVal.SymbolTableEntry.Ref
		}

		return g.Ctx.AddSymbol(node.Name, entry)
	case ast.ArrayInitializationExpression:
		entry := &ast.SymbolTableEntry{
			Address:      compiledVal.Value,
			DeclaredType: node.ExplicitType,
		}

		err := g.Ctx.AddSymbol(node.Name, entry)
		if err != nil {
			return err
		}

		return g.Ctx.AddArraySymbol(node.Name, compiledVal.ArraySymbolTableEntry)
	default:
		format := fmt.Sprintf("declareVarWithInitializer unimplemented for %T", node.Value)
		g.NotImplemented(format)
	}

	return nil

}

func (g *LLVMGenerator) declareVarWithZeroValue(node *ast.VarDeclarationStatement) error {
	err, llvmType := node.ExplicitType.LLVMType(g.Ctx)
	if err != nil {
		return err
	}

	alloc := g.Ctx.Builder.CreateAlloca(llvmType, fmt.Sprintf("alloc.%s", node.Name))
	g.Ctx.Builder.CreateStore(llvm.ConstNull(llvmType), alloc)

	entry := &ast.SymbolTableEntry{
		Value:        alloc,
		Address:      &alloc,
		DeclaredType: node.ExplicitType,
	}

	err = g.Ctx.AddSymbol(node.Name, entry)
	if err != nil {
		return err
	}

	return nil
}

// VisitWhileStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitWhileStatement(node *ast.WhileStatement) error {
	panic("unimplemented")
}

var _ ast.CodeGenerator = (*LLVMGenerator)(nil)

func NewLLVMGenerator(ctx *ast.CompilerCtx) *LLVMGenerator {
	return &LLVMGenerator{Ctx: ctx}
}

func (g *LLVMGenerator) VisitArrayAccessExpression(node *ast.ArrayAccessExpression) error {
	return nil
}
