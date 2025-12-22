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
		case ast.ArrayAccessExpression:
			load := g.Ctx.Builder.CreateLoad(
				lastResult.ArraySymbolTableEntry.UnderlyingType,
				*lastResult.Value,
				"",
			)
			printableValues = append(printableValues, load)
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
	properties := []llvm.Type{}
	for i := range node.Properties {
		propertyType := node.Types[i]
		err, llvmType := propertyType.LLVMType(g.Ctx)
		if err != nil {
			return err
		}
		properties = append(properties, llvmType)
	}

	newtype := g.Ctx.Context.StructCreateNamed(node.Name)
	newtype.StructSetBody(properties, false)

	entry := &ast.StructSymbolTableEntry{
		LLVMType:      newtype,
		Metadata:      *node,
		PropertyTypes: properties,
	}

	return g.Ctx.AddStructSymbol(node.Name, entry)
}

// VisitStructInitializationExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
	err, structType := g.Ctx.FindStructSymbol(node.Name)
	if err != nil {
		return err
	}

	structInstance := g.Ctx.Builder.CreateAlloca(structType.LLVMType, fmt.Sprintf("%s.instance", node.Name))

	for _, name := range node.Properties {
		err, propIndex := structType.Metadata.PropertyIndex(name)
		if err != nil {
			return fmt.Errorf("StructInitializationExpression: property %s not found", name)
		}

		expr := node.Values[propIndex]

		err = expr.Accept(g)
		if err != nil {
			return err
		}

		lastResult := g.getLastResult()
		pointerToField := g.Ctx.Builder.CreateStructGEP(
			structType.LLVMType,
			structInstance,
			propIndex,
			fmt.Sprintf("%s.instance.%s", node.Name, name),
		)

		switch expr.(type) {
		case ast.NumberExpression, ast.FloatExpression, ast.StringExpression:
			g.Ctx.Builder.CreateStore(*lastResult.Value, pointerToField)
		default:
			return fmt.Errorf("StructInitializationExpression expression: %v is not a known field type", expr)
		}
	}

	result := ast.CompilerResult{
		Value:                  &structInstance,
		StructSymbolTableEntry: structType,
	}

	g.setLastResult(&result)

	return nil
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
	name := node.Name

	switch node.Value.(type) {
	case ast.StringExpression:
		alloc := g.Ctx.Builder.CreateAlloca(compiledVal.Value.Type(), name)
		g.Ctx.Builder.CreateStore(*compiledVal.Value, alloc)

		entry := &ast.SymbolTableEntry{
			Address:      &alloc,
			DeclaredType: node.ExplicitType,
		}

		return g.Ctx.AddSymbol(node.Name, entry)
	case ast.NumberExpression, ast.FloatExpression, ast.StructInitializationExpression:
		alloc := g.Ctx.Builder.CreateAlloca(compiledVal.Value.Type(), name)
		g.Ctx.Builder.CreateStore(*compiledVal.Value, alloc)

		entry := &ast.SymbolTableEntry{
			Address:      &alloc,
			DeclaredType: node.ExplicitType,
		}
		if compiledVal.StructSymbolTableEntry != nil {
			entry.Ref = compiledVal.StructSymbolTableEntry
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
	err, entry, array, indices := g.findArraySymbolTableEntry(node)
	if err != nil {
		return err
	}

	itemPtr := g.Ctx.Builder.CreateInBoundsGEP(
		entry.UnderlyingType,
		*array.Address,
		indices,
		"",
	)

	g.setLastResult(&ast.CompilerResult{
		Value:                 &itemPtr,
		ArraySymbolTableEntry: entry,
	})

	return nil

}

func (g *LLVMGenerator) findArraySymbolTableEntry(expr *ast.ArrayAccessExpression) (error, *ast.ArraySymbolTableEntry, *ast.SymbolTableEntry, []llvm.Value) {
	var name string
	var array *ast.SymbolTableEntry
	var entry *ast.ArraySymbolTableEntry

	switch expr.Name.(type) {
	case ast.SymbolExpression:
		varName, _ := expr.Name.(ast.SymbolExpression)

		err, arrayEntry := g.Ctx.FindSymbol(varName.Value)
		if err != nil {
			key := "ArrayAccessExpression.NotFoundInSymbolTable"
			return g.Ctx.Dialect.Error(key, varName.Value), nil, nil, nil
		}
		array = arrayEntry

		err, arraySymEntry := g.Ctx.FindArraySymbol(varName.Value)
		if err != nil {
			key := "ArrayAccessExpression.NotFoundInArraySymbolTable"
			return g.Ctx.Dialect.Error(key, varName.Value), nil, nil, nil
		}
		entry = arraySymEntry

		name = varName.Value
	case ast.MemberExpression:
		if err := expr.Name.Accept(g); err != nil {
			return err, nil, nil, nil
		}
		val := g.getLastResult()

		var elementType llvm.Type
		var arrayType llvm.Type
		var elementsCount int = -1
		var isPointerType bool = false

		if val.SymbolTableEntry == nil && val.SymbolTableEntry.Ref == nil {
			return fmt.Errorf("ArrayAccessExpression Missing SymbolTableEntry"), nil, nil, nil
		}

		propExpr, _ := expr.Name.(ast.MemberExpression)
		propSym, _ := propExpr.Property.(ast.SymbolExpression)

		if val.SymbolTableEntry.Ref == nil {
			format := "ArrayAccessExpression property %s is not an array"
			return fmt.Errorf(format, propSym.Value), nil, nil, nil
		}

		// resolveStructAccess is now in generator
		propIndex, err := g.resolveStructAccess(val.SymbolTableEntry.Ref, propSym.Value)
		if err != nil {
			return err, nil, nil, nil
		}

		astType := val.SymbolTableEntry.Ref.Metadata.Types[propIndex]

		switch coltype := astType.(type) {
		case ast.PointerType:
			isPointerType = true
			err, elementType = coltype.Underlying.LLVMType(g.Ctx)
			if err != nil {
				return err, nil, nil, nil
			}
			arrayType = elementType
		case ast.ArrayType:
			err, elementType = coltype.Underlying.LLVMType(g.Ctx)
			if err != nil {
				return err, nil, nil, nil
			}
			err, arrayType = coltype.LLVMType(g.Ctx)
			if err != nil {
				return err, nil, nil, nil
			}
			elementsCount = coltype.Size
		default:
			err := fmt.Errorf("Property %s is not an array", propSym.Value)
			return err, nil, nil, nil
		}

		if isPointerType {
			pointerValue := g.Ctx.Builder.CreateLoad(*val.StuctPropertyValueType, *val.Value, "")
			array = &ast.SymbolTableEntry{Value: pointerValue}
		} else {
			array = &ast.SymbolTableEntry{Value: *val.Value}
		}

		entry = &ast.ArraySymbolTableEntry{
			UnderlyingType: elementType,
			Type:           arrayType,
			ElementsCount:  elementsCount,
		}
	default:
		err := fmt.Errorf("ArrayAccessExpression not implemented")
		return err, nil, nil, nil
	}

	var indices []llvm.Value

	switch expr.Index.(type) {
	case ast.NumberExpression:
		idx, _ := expr.Index.(ast.NumberExpression)

		if int(idx.Value) < 0 {
			key := "ArrayAccessExpression.AccessedIndexIsNotANumber"
			return g.Ctx.Dialect.Error(key, expr.Index), nil, nil, nil
		}

		if entry.ElementsCount >= 0 && int(idx.Value) > entry.ElementsCount-1 {
			key := "ArrayAccessExpression.IndexOutOfBounds"
			return g.Ctx.Dialect.Error(key, int(idx.Value), name), nil, nil, nil
		}

		indices = []llvm.Value{
			llvm.ConstInt((*g.Ctx.Context).Int32Type(), uint64(idx.Value), false),
		}
	case ast.SymbolExpression, ast.BinaryExpression, ast.MemberExpression:
		if err := expr.Index.Accept(g); err != nil {
			return err, nil, nil, nil
		}
		res := g.getLastResult()

		if res.Value.Type().TypeKind() == llvm.PointerTypeKind {
			load := g.Ctx.Builder.CreateLoad(res.Value.AllocatedType(), *res.Value, "")
			indices = []llvm.Value{load}
		} else {
			indices = []llvm.Value{*res.Value}
		}
	default:
		key := "ArrayAccessExpression.AccessedIndexIsNotANumber"
		return g.Ctx.Dialect.Error(key, expr.Index), nil, nil, nil
	}

	return nil, entry, array, indices
}

func (g *LLVMGenerator) resolveStructAccess(
	structType *ast.StructSymbolTableEntry,
	propName string,
) (int, error) {
	err, propIndex := structType.Metadata.PropertyIndex(propName)
	if err != nil {
		return 0, fmt.Errorf("struct %s has no field %s", structType.Metadata.Name, propName)
	}
	return propIndex, nil
}
