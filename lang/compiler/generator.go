package compiler

import (
	"fmt"
	"math"
	"os"
	"swahili/lang/ast"
	"swahili/lang/lexer"

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

func (g *LLVMGenerator) VisitConditionalStatement(node *ast.ConditionalStatetement) error {
	err := node.Condition.Accept(g)
	if err != nil {
		return err
	}

	condition := g.getLastResult()
	bodyBlock := g.Ctx.Builder.GetInsertBlock()
	parentFunc := bodyBlock.Parent()
	mergeBlock := g.Ctx.Context.AddBasicBlock(parentFunc, "merge")
	thenBlock := g.Ctx.Context.AddBasicBlock(parentFunc, "if")
	elseBlock := g.Ctx.Context.AddBasicBlock(parentFunc, "else")

	g.Ctx.Builder.CreateCondBr(*condition.Value, thenBlock, elseBlock)
	g.Ctx.Builder.SetInsertPointAtEnd(thenBlock)

	if err := node.Success.Accept(g); err != nil {
		return err
	}

	successVal := g.getLastResult()
	if thenBlock.LastInstruction().InstructionOpcode() != llvm.Ret {
		g.Ctx.Builder.CreateBr(mergeBlock)
	}

	g.Ctx.Builder.SetInsertPointAtEnd(elseBlock)

	err = node.Failure.Accept(g)
	if err != nil {
		return err
	}

	failureVal := g.getLastResult()
	if elseBlock.LastInstruction().InstructionOpcode() != llvm.Ret {
		g.Ctx.Builder.CreateBr(mergeBlock)
	}

	g.Ctx.Builder.SetInsertPointAtEnd(mergeBlock)

	var phi llvm.Value
	if successVal != nil {
		phi = g.Ctx.Builder.CreatePHI(successVal.Value.Type(), "")
		phi.AddIncoming([]llvm.Value{*successVal.Value}, []llvm.BasicBlock{thenBlock})
	}

	if failureVal != nil {
		phi = g.Ctx.Builder.CreatePHI(successVal.Value.Type(), "")
		phi.AddIncoming([]llvm.Value{*successVal.Value}, []llvm.BasicBlock{thenBlock})
	}

	thenBlock.MoveAfter(bodyBlock)
	elseBlock.MoveAfter(thenBlock)
	mergeBlock.MoveAfter(thenBlock)

	if successVal != nil {
		g.setLastResult(&ast.CompilerResult{Value: &phi})
	}

	return nil
}

func (g *LLVMGenerator) VisitWhileStatement(node *ast.WhileStatement) error {
	bodyBlock := g.Ctx.Builder.GetInsertBlock()
	parentFunc := bodyBlock.Parent()

	whileConditionBlock := g.Ctx.Context.AddBasicBlock(parentFunc, "while.cond")
	whileBodyBlock := g.Ctx.Context.AddBasicBlock(parentFunc, "while.body")
	whileMergeBlock := g.Ctx.Context.AddBasicBlock(parentFunc, "while.merge")

	g.Ctx.Builder.CreateBr(whileConditionBlock)

	g.Ctx.Builder.SetInsertPointAtEnd(whileConditionBlock)

	err := node.Condition.Accept(g)
	if err != nil {
		return err
	}

	condition := g.getLastResult()

	g.Ctx.Builder.CreateCondBr(*condition.Value, whileBodyBlock, whileMergeBlock)
	g.Ctx.Builder.SetInsertPointAtEnd(whileBodyBlock)

	err = node.Body.Accept(g)
	if err != nil {
		return err
	}

	lastBodyBlock := g.Ctx.Builder.GetInsertBlock()

	opcode := lastBodyBlock.LastInstruction().InstructionOpcode()
	if lastBodyBlock.LastInstruction().IsNil() || (opcode != llvm.Ret && opcode != llvm.Br) {
		g.Ctx.Builder.CreateBr(whileConditionBlock)
	}

	g.Ctx.Builder.SetInsertPointAtEnd(whileMergeBlock)

	return nil
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
		case ast.SymbolExpression:
			err := v.Accept(g)
			if err != nil {
				return err
			}

			compiledVal := g.getLastResult()

			if compiledVal.SymbolTableEntry != nil && compiledVal.SymbolTableEntry.Ref != nil {
				load := g.Ctx.Builder.CreateLoad(compiledVal.SymbolTableEntry.Ref.LLVMType, *compiledVal.Value, "")
				g.Ctx.Builder.CreateStore(load, itemGep)
			} else {
				load := g.Ctx.Builder.CreateLoad(compiledVal.Value.Type(), *compiledVal.Value, "")
				g.Ctx.Builder.CreateStore(load, itemGep)
			}
		case ast.NumberExpression, ast.FloatExpression, *ast.StringExpression:
			err := v.Accept(g)
			if err != nil {
				return err
			}
			compiledVal := g.getLastResult()
			g.Ctx.Builder.CreateStore(*compiledVal.Value, itemGep)
		case ast.StructInitializationExpression:
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
	g.NotImplemented("VisitArrayOfStructsAccessExpression not implemented")

	return nil
}

// VisitAssignmentExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitAssignmentExpression(node *ast.AssignmentExpression) error {
	err := node.Assignee.Accept(g)
	if err != nil {
		return err
	}

	compiledAssignee := g.getLastResult()

	err = node.Value.Accept(g)
	if err != nil {
		return err
	}

	compiledValue := g.getLastResult()

	var valueToBeAssigned llvm.Value

	switch node.Value.(type) {
	case ast.ArrayAccessExpression:
		valueToBeAssigned = g.Ctx.Builder.CreateLoad(compiledValue.Value.AllocatedType(), *compiledValue.Value, "")
	default:
		valueToBeAssigned = *compiledValue.Value
	}

	address := compiledAssignee.Value
	if compiledAssignee.SymbolTableEntry.Address != nil {
		address = compiledAssignee.SymbolTableEntry.Address
	}

	g.Ctx.Builder.CreateStore(valueToBeAssigned, *address)

	return nil
}

// VisitBlockStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitBlockStatement(node *ast.BlockStatement) error {
	oldCtx := g.Ctx
	newCtx := ast.NewCompilerContext(
		oldCtx.Context,
		oldCtx.Builder,
		oldCtx.Module,
		oldCtx.Dialect,
		oldCtx,
	)
	g.Ctx = newCtx

	for _, v := range node.Body {
		err := v.Accept(g)
		if err != nil {
			return err
		}
	}

	g.Ctx = oldCtx

	return nil
}

// VisitCallExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitCallExpression(node *ast.CallExpression) error {
	g.NotImplemented("VisitCallExpression not implemented")

	return nil
}

// VisitExpressionStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitExpressionStatement(node *ast.ExpressionStatement) error {
	err := node.Exp.Accept(g)
	if err != nil {
		return err
	}

	return nil
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

func (g *LLVMGenerator) VisitFunctionCall(node *ast.FunctionCallExpression) error {
	name, ok := node.Name.(ast.SymbolExpression)
	if !ok {
		return fmt.Errorf("FunctionCallExpression: name is not a symbol")
	}
	err, funcType := g.Ctx.FindFuncSymbol(name.Value)
	if err != nil {
		return err
	}

	funcVal := g.Ctx.Module.NamedFunction(name.Value)
	if funcVal.IsNil() {
		return fmt.Errorf("function %s does not exist", name.Value)
	}

	args := []llvm.Value{}

	for _, arg := range node.Args {
		if err := arg.Accept(g); err != nil {
			return err
		}
		val := g.getLastResult()

		switch arg.(type) {
		case ast.SymbolExpression:
			if val.SymbolTableEntry.Ref != nil {
				alloca := g.Ctx.Builder.CreateAlloca(val.SymbolTableEntry.Ref.LLVMType, "")
				g.Ctx.Builder.CreateStore(*val.Value, alloca)
				args = append(args, alloca)

				break
			}

			if val.SymbolTableEntry != nil && val.SymbolTableEntry.Address != nil {
				args = append(args, *val.SymbolTableEntry.Address)
			} else {
				args = append(args, *val.Value)
			}

		default:
			args = append(args, *val.Value)
		}

	}

	val := g.Ctx.Builder.CreateCall(*funcType, funcVal, args, "")

	g.setLastResult(&ast.CompilerResult{Value: &val})

	return nil
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
	obj, ok := node.Object.(ast.SymbolExpression)
	if !ok {
		return fmt.Errorf("struct object should be a symbol")
	}

	err, varDef := g.Ctx.FindSymbol(obj.Value)
	if err != nil {
		return fmt.Errorf("variable %s is not defined", obj.Value)
	}

	if varDef.Ref == nil {
		return fmt.Errorf("variable %s is not a struct instance", obj.Value)
	}

	propName, err := g.getProperty(node)
	if err != nil {
		return err
	}

	propIndex, err := g.resolveStructAccess(varDef.Ref, propName)
	if err != nil {
		return err
	}

	var baseValue llvm.Value
	if varDef.Address != nil {
		baseValue = *varDef.Address
	} else {
		baseValue = varDef.Value
	}

	addr := g.Ctx.Builder.CreateStructGEP(varDef.Ref.LLVMType, baseValue, propIndex, "")
	propType := varDef.Ref.PropertyTypes[propIndex]

	g.setLastResult(&ast.CompilerResult{
		Value:                  &addr,
		SymbolTableEntry:       varDef,
		StuctPropertyValueType: &propType,
	})

	return nil
}

func (g *LLVMGenerator) getProperty(expr *ast.MemberExpression) (string, error) {
	prop, ok := expr.Property.(ast.SymbolExpression)
	if !ok {
		return "", fmt.Errorf("struct property should be a symbol")
	}

	return prop.Value, nil
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
		case ast.MemberExpression:
			loadedval := g.Ctx.Builder.CreateLoad(
				*lastResult.StuctPropertyValueType,
				*lastResult.Value,
				"",
			)
			printableValues = append(printableValues, loadedval)
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

func (g *LLVMGenerator) VisitReturnStatement(node *ast.ReturnStatement) error {
	if err := node.Value.Accept(g); err != nil {
		return err
	}
	res := g.getLastResult()

	retVal, err := g.prepareReturnValue(node.Value, res)
	if err != nil {
		return err
	}

	g.Ctx.Builder.CreateRet(retVal)
	return nil
}

func (g *LLVMGenerator) prepareReturnValue(expr ast.Expression, res *ast.CompilerResult) (llvm.Value, error) {
	switch expr.(type) {
	case ast.ArrayAccessExpression:
		return g.Ctx.Builder.CreateLoad((*g.Ctx.Context).Int32Type(), *res.Value, ""), nil
	case ast.MemberExpression:
		return g.Ctx.Builder.CreateLoad(*res.StuctPropertyValueType, *res.Value, ""), nil
	case ast.StringExpression:
		alloc := g.Ctx.Builder.CreateAlloca(res.Value.Type(), "")
		g.Ctx.Builder.CreateStore(*res.Value, alloc)
		return alloc, nil
	case ast.SymbolExpression, ast.BinaryExpression, ast.FunctionCallExpression, ast.NumberExpression, ast.FloatExpression:
		return *res.Value, nil
	default:
		return llvm.Value{}, fmt.Errorf("ReturnStatement unknown expression <%s>", expr)
	}
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

	if entry.Address == nil {
		g.setLastResult(
			&ast.CompilerResult{
				Value:            &entry.Value,
				SymbolTableEntry: entry,
			},
		)
		return nil
	}

	var loadedValue llvm.Value

	switch entry.DeclaredType.(type) {
	case ast.StringType, ast.ArrayType:
		loadedValue = g.Ctx.Builder.CreateLoad(entry.Address.Type(), *entry.Address, "")
	default:
		loadedValue = g.Ctx.Builder.CreateLoad(entry.Address.AllocatedType(), *entry.Address, "")
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
	if g.Ctx.SymbolExistsInCurrentScope(node.Name) {
		return fmt.Errorf("variable %s is already defined", node.Name)
	}

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
	case ast.StructInitializationExpression:

		entry := &ast.SymbolTableEntry{
			Address:      compiledVal.Value,
			DeclaredType: node.ExplicitType,
		}
		if compiledVal.StructSymbolTableEntry != nil {
			entry.Ref = compiledVal.StructSymbolTableEntry
		}

		return g.Ctx.AddSymbol(node.Name, entry)
	case ast.NumberExpression, ast.FunctionCallExpression, ast.FloatExpression,
		ast.BinaryExpression:
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

func (g *LLVMGenerator) VisitPrefixExpression(node *ast.PrefixExpression) error {
	if err := node.RightExpression.Accept(g); err != nil {
		return err
	}

	res := g.getLastResult()

	handler, ok := prefixOpHandlers[node.Operator.Kind]
	if !ok {
		return fmt.Errorf("PrefixExpression: operator %s not supported", node.Operator.Kind)
	}

	val := handler(g, *res.Value)
	g.setLastResult(&ast.CompilerResult{Value: &val})
	return nil
}

func (g *LLVMGenerator) VisitBinaryExpression(node *ast.BinaryExpression) error {
	err, leftResult, rightResult := g.compileLeftAndRightResult(node)
	if err != nil {
		return err
	}

	var finalLeftValue, finalRightValue llvm.Value

	var leftValue, rightValue llvm.Value

	switch leftResult.Value.Type().TypeKind() {
	case llvm.PointerTypeKind:
		if leftResult.ArraySymbolTableEntry != nil {
			leftValue = g.Ctx.Builder.CreateLoad(leftResult.ArraySymbolTableEntry.UnderlyingType, *leftResult.Value, "")

			break
		}

		if leftResult.StuctPropertyValueType != nil {
			elementType := leftResult.StuctPropertyValueType
			leftValue = g.Ctx.Builder.CreateLoad(*elementType, *leftResult.Value, "")

			break
		}

		if leftResult.Value.Type().ElementType().IsNil() {
			leftValue = *leftResult.Value

			break
		}

		elementType := leftResult.Value.Type().ElementType()
		leftValue = g.Ctx.Builder.CreateLoad(elementType, *leftResult.Value, "")
	case llvm.StructTypeKind:
		elementType := leftResult.StuctPropertyValueType
		leftValue = g.Ctx.Builder.CreateLoad(*elementType, *leftResult.Value, "")
	default:
		leftValue = *leftResult.Value
	}

	switch rightResult.Value.Type().TypeKind() {
	case llvm.PointerTypeKind:
		if rightResult.ArraySymbolTableEntry != nil {
			rightValue = g.Ctx.Builder.CreateLoad(rightResult.ArraySymbolTableEntry.UnderlyingType, *rightResult.Value, "")

			break
		}

		if rightResult.StuctPropertyValueType != nil {
			elementType := rightResult.StuctPropertyValueType
			rightValue = g.Ctx.Builder.CreateLoad(*elementType, *rightResult.Value, "")

			break
		}

		if rightResult.Value.Type().ElementType().IsNil() {
			rightValue = *rightResult.Value

			break
		}

		elementType := rightResult.Value.Type().ElementType()
		rightValue = g.Ctx.Builder.CreateLoad(elementType, *rightResult.Value, "")
	case llvm.StructTypeKind:
		elementType := rightResult.StuctPropertyValueType
		rightValue = g.Ctx.Builder.CreateLoad(*elementType, *rightResult.Value, "")

	default:
		rightValue = *rightResult.Value
	}

	// Determine the common type
	ctype, err := g.commonType(leftValue, rightValue)
	if err != nil {
		return fmt.Errorf("%w %v", err, node.TokenStream())
	}

	// Cast only if necessary
	if leftValue.Type() == ctype {
		finalLeftValue = leftValue
	} else {
		err, finalLeftValue = g.castToType(ctype, leftValue)
		if err != nil {
			return err
		}
	}

	if rightValue.Type() == ctype {
		finalRightValue = rightValue
	} else {
		err, finalRightValue = g.castToType(ctype, rightValue)
		if err != nil {
			return err
		}
	}

	err, res := g.handleBinaryOp(node.Operator.Kind, finalLeftValue, finalRightValue)
	if err != nil {
		return err
	}

	g.setLastResult(res)
	return nil
}

func (g *LLVMGenerator) compileLeftAndRightResult(node *ast.BinaryExpression) (error, *ast.CompilerResult, *ast.CompilerResult) {
	if err := node.Left.Accept(g); err != nil {
		return err, nil, nil
	}
	compiledLeftValue := g.getLastResult()

	if compiledLeftValue == nil {
		return fmt.Errorf("left side of expression is nil"), nil, nil
	}

	if compiledLeftValue.Value.Type().TypeKind() == llvm.VoidTypeKind {
		return fmt.Errorf("left side of expression is of type void"), nil, nil
	}

	switch node.Left.(type) {
	case ast.StringExpression:
		glob := llvm.AddGlobal(*g.Ctx.Module, compiledLeftValue.Value.Type(), "")
		glob.SetInitializer(*compiledLeftValue.Value)
		compiledLeftValue.Value = &glob
	default:
	}

	if err := node.Right.Accept(g); err != nil {
		return err, nil, nil
	}
	compiledRightValue := g.getLastResult()

	if compiledRightValue == nil {
		return fmt.Errorf("right side of expression is nil"), nil, nil
	}

	if compiledRightValue.Value.Type().TypeKind() == llvm.VoidTypeKind {
		return fmt.Errorf("right side of expression is of type void"), nil, nil
	}

	switch node.Right.(type) {
	case ast.StringExpression:
		glob := llvm.AddGlobal(*g.Ctx.Module, compiledRightValue.Value.Type(), "")
		glob.SetInitializer(*compiledRightValue.Value)
		compiledRightValue.Value = &glob
	default:
	}

	return nil, compiledLeftValue, compiledRightValue
}

func (g *LLVMGenerator) commonType(l, r llvm.Value) (llvm.Type, error) {
	lKind := l.Type().TypeKind()
	rKind := r.Type().TypeKind()

	// Both pointers - use int type as fallback (shouldn't happen in binary ops)
	if lKind == llvm.PointerTypeKind && rKind == llvm.PointerTypeKind {
		return (*g.Ctx.Context).Int32Type(), nil
	}

	// Same type
	if l.Type() == r.Type() {
		return l.Type(), nil
	}

	// Left is pointer, use right type
	if lKind == llvm.PointerTypeKind {
		return r.Type(), nil
	}

	// Right is pointer, use left type
	if rKind == llvm.PointerTypeKind {
		return l.Type(), nil
	}

	// Both integers - return the common integer type
	if lKind == llvm.IntegerTypeKind && rKind == llvm.IntegerTypeKind {
		return (*g.Ctx.Context).Int32Type(), nil
	}

	// Handle int vs float: promote to float
	if (lKind == llvm.IntegerTypeKind && (rKind == llvm.FloatTypeKind || rKind == llvm.DoubleTypeKind)) ||
		((lKind == llvm.FloatTypeKind || lKind == llvm.DoubleTypeKind) && rKind == llvm.IntegerTypeKind) {
		// Promote to double (float type)
		return (*g.Ctx.Context).DoubleType(), nil
	}

	// Both floats
	if (lKind == llvm.FloatTypeKind || lKind == llvm.DoubleTypeKind) &&
		(rKind == llvm.FloatTypeKind || rKind == llvm.DoubleTypeKind) {
		return (*g.Ctx.Context).DoubleType(), nil
	}

	format := "Unhandled type combination: left=%v, right=%s"

	return llvm.Type{}, fmt.Errorf(format, l, rKind)
}

func (g *LLVMGenerator) castToType(t llvm.Type, v llvm.Value) (error, llvm.Value) {
	vKind := v.Type().TypeKind()
	tKind := t.TypeKind()

	if vKind == tKind {
		return nil, v
	}

	switch vKind {
	case llvm.IntegerTypeKind:
		switch tKind {
		case llvm.DoubleTypeKind, llvm.FloatTypeKind:
			return nil, g.Ctx.Builder.CreateSIToFP(v, t, "")
		default:
			return fmt.Errorf("Cannot cast integer to %s", tKind), llvm.Value{}
		}
	case llvm.DoubleTypeKind, llvm.FloatTypeKind:
		switch tKind {
		case llvm.IntegerTypeKind:
			// Convert float to int
			return nil, g.Ctx.Builder.CreateFPToSI(v, t, "")
		default:
			return fmt.Errorf("Cannot cast %s to %s", vKind, tKind), llvm.Value{}
		}
	case llvm.PointerTypeKind:
		switch tKind {
		case llvm.IntegerTypeKind:
			return nil, g.Ctx.Builder.CreatePtrToInt(v, t, "")
		default:
			return nil, v
		}
	default:
		return fmt.Errorf("Unhandled type conversion from %s to %s", vKind, tKind), llvm.Value{}
	}
}

func (g *LLVMGenerator) handleBinaryOp(kind lexer.TokenKind, l, r llvm.Value) (error, *ast.CompilerResult) {
	handler, ok := binaryOpHandlers[kind]
	if !ok {
		return fmt.Errorf("Binary expressions : unsupported operator <%s>", kind), nil
	}
	res := handler(g, l, r)
	return nil, &ast.CompilerResult{Value: &res}
}

func (g *LLVMGenerator) handleAnd(l, r llvm.Value) llvm.Value {
	zero := llvm.ConstInt(l.Type(), 0, false)
	lBool := g.Ctx.Builder.CreateICmp(llvm.IntNE, l, zero, "")
	rBool := g.Ctx.Builder.CreateICmp(llvm.IntNE, r, zero, "")
	return g.Ctx.Builder.CreateAnd(lBool, rBool, "")
}

func (g *LLVMGenerator) handleOr(l, r llvm.Value) llvm.Value {
	zero := llvm.ConstInt(l.Type(), 0, false)
	lBool := g.Ctx.Builder.CreateICmp(llvm.IntNE, l, zero, "")
	rBool := g.Ctx.Builder.CreateICmp(llvm.IntNE, r, zero, "")
	return g.Ctx.Builder.CreateOr(lBool, rBool, "")
}

func (g *LLVMGenerator) handlePlus(l, r llvm.Value) llvm.Value {
	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind {
		return g.Ctx.Builder.CreateFAdd(l, r, "")
	}
	return g.Ctx.Builder.CreateAdd(l, r, "")
}

func (g *LLVMGenerator) handleMinus(l, r llvm.Value) llvm.Value {
	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind {
		return g.Ctx.Builder.CreateFSub(l, r, "")
	}
	return g.Ctx.Builder.CreateSub(l, r, "")
}

func (g *LLVMGenerator) handleStar(l, r llvm.Value) llvm.Value {
	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind {
		return g.Ctx.Builder.CreateFMul(l, r, "")
	}
	return g.Ctx.Builder.CreateMul(l, r, "")
}

func (g *LLVMGenerator) handleModulo(l, r llvm.Value) llvm.Value {
	return g.Ctx.Builder.CreateSRem(l, r, "")
}

func (g *LLVMGenerator) handleDivide(l, r llvm.Value) llvm.Value {
	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind {
		return g.Ctx.Builder.CreateFDiv(l, r, "")
	}
	return g.Ctx.Builder.CreateSDiv(l, r, "")
}

func (g *LLVMGenerator) handleGreaterThan(l, r llvm.Value) llvm.Value {
	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		return g.Ctx.Builder.CreateFCmp(llvm.FloatOGT, l, r, "")
	}
	return g.Ctx.Builder.CreateICmp(llvm.IntSGT, l, r, "")
}

func (g *LLVMGenerator) handleGreaterThanEquals(l, r llvm.Value) llvm.Value {
	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		return g.Ctx.Builder.CreateFCmp(llvm.FloatOGE, l, r, "")
	}
	return g.Ctx.Builder.CreateICmp(llvm.IntSGE, l, r, "")
}

func (g *LLVMGenerator) handleLessThan(l, r llvm.Value) llvm.Value {
	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		return g.Ctx.Builder.CreateFCmp(llvm.FloatOLT, l, r, "")
	}
	return g.Ctx.Builder.CreateICmp(llvm.IntSLT, l, r, "")
}

func (g *LLVMGenerator) handleLessThanEquals(l, r llvm.Value) llvm.Value {
	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		return g.Ctx.Builder.CreateFCmp(llvm.FloatOLE, l, r, "")
	}
	return g.Ctx.Builder.CreateICmp(llvm.IntSLE, l, r, "")
}

func (g *LLVMGenerator) handleEquals(l, r llvm.Value) llvm.Value {
	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		return g.Ctx.Builder.CreateFCmp(llvm.FloatOEQ, l, r, "")
	}
	return g.Ctx.Builder.CreateICmp(llvm.IntEQ, l, r, "")
}

func (g *LLVMGenerator) handleNotEquals(l, r llvm.Value) llvm.Value {
	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		return g.Ctx.Builder.CreateFCmp(llvm.FloatONE, l, r, "")
	}
	return g.Ctx.Builder.CreateICmp(llvm.IntNE, l, r, "")
}

func (g *LLVMGenerator) handlePrefixMinus(val llvm.Value) llvm.Value {
	if val.Type().TypeKind() == llvm.FloatTypeKind || val.Type().TypeKind() == llvm.DoubleTypeKind {
		return g.Ctx.Builder.CreateFNeg(val, "")
	}
	return g.Ctx.Builder.CreateNeg(val, "")
}

func (g *LLVMGenerator) handlePrefixNot(val llvm.Value) llvm.Value {
	zero := llvm.ConstInt(val.Type(), 0, false)
	boolVal := g.Ctx.Builder.CreateICmp(llvm.IntEQ, val, zero, "")
	return g.Ctx.Builder.CreateZExt(boolVal, val.Type(), "")
}

var binaryOpHandlers = map[lexer.TokenKind]func(g *LLVMGenerator, l, r llvm.Value) llvm.Value{
	lexer.And:               (*LLVMGenerator).handleAnd,
	lexer.Or:                (*LLVMGenerator).handleOr,
	lexer.Plus:              (*LLVMGenerator).handlePlus,
	lexer.Minus:             (*LLVMGenerator).handleMinus,
	lexer.Star:              (*LLVMGenerator).handleStar,
	lexer.Modulo:            (*LLVMGenerator).handleModulo,
	lexer.Divide:            (*LLVMGenerator).handleDivide,
	lexer.GreaterThan:       (*LLVMGenerator).handleGreaterThan,
	lexer.GreaterThanEquals: (*LLVMGenerator).handleGreaterThanEquals,
	lexer.LessThan:          (*LLVMGenerator).handleLessThan,
	lexer.LessThanEquals:    (*LLVMGenerator).handleLessThanEquals,
	lexer.Equals:            (*LLVMGenerator).handleEquals,
	lexer.NotEquals:         (*LLVMGenerator).handleNotEquals,
}

var prefixOpHandlers = map[lexer.TokenKind]func(g *LLVMGenerator, val llvm.Value) llvm.Value{
	lexer.Minus: (*LLVMGenerator).handlePrefixMinus,
	lexer.Not:   (*LLVMGenerator).handlePrefixNot,
}

func (g *LLVMGenerator) VisitFunctionDefinition(node *ast.FuncDeclStatement) error {
	newCtx := ast.NewCompilerContext(
		g.Ctx.Context,
		g.Ctx.Builder,
		g.Ctx.Module,
		g.Ctx.Dialect,
		g.Ctx,
	)

	err, params := g.funcParams(newCtx, node)
	if err != nil {
		return err
	}

	err, returnType := g.extractType(newCtx, node.ReturnType)
	if err != nil {
		return err
	}

	newfuncType := llvm.FunctionType(returnType.typ, params, node.ArgsVariadic)
	newFunc := llvm.AddFunction(*newCtx.Module, node.Name, newfuncType)

	err = g.Ctx.AddFuncSymbol(node.Name, &newfuncType)
	if err != nil {
		return err
	}

	for i, p := range newFunc.Params() {
		argType := node.Args[i].ArgType
		name := node.Args[i].Name
		p.SetName(name)

		entry := ast.SymbolTableEntry{Value: p, DeclaredType: argType}

		err, eType := g.extractType(newCtx, argType)
		if err != nil {
			return err
		}

		if eType.sEntry != nil {
			entry.Ref = eType.sEntry
		}

		err = newCtx.AddSymbol(name, &entry)
		if err != nil {
			return fmt.Errorf("failed to add parameter %s to symbol table: %w", name, err)
		}

		if eType.aEntry != nil {
			err := newCtx.AddArraySymbol(name, eType.aEntry)
			if err != nil {
				return fmt.Errorf("failed to add parameter %s to arrays symbol table: %w", name, err)
			}
		}
	}

	if len(node.Body.Body) > 0 {
		block := g.Ctx.Context.AddBasicBlock(newFunc, "body")
		g.Ctx.Builder.SetInsertPointAtEnd(block)

		newGenerator := &LLVMGenerator{Ctx: newCtx}
		if err := node.Body.Accept(newGenerator); err != nil {
			return err
		}

		return nil
	}

	// If we get here it means function has no Body
	// so it's just a declaration of an external function
	newFunc.SetLinkage(llvm.ExternalLinkage)

	return nil
}

type extractedType struct {
	typ    llvm.Type
	entry  *ast.SymbolTableEntry
	sEntry *ast.StructSymbolTableEntry
	aEntry *ast.ArraySymbolTableEntry
}

func (g *LLVMGenerator) extractType(ctx *ast.CompilerCtx, t ast.Type) (error, extractedType) {
	err, compiledType := t.LLVMType(ctx)
	if err != nil {
		return err, extractedType{}
	}

	switch typ := t.(type) {
	case ast.NumberType, ast.Number64Type, ast.FloatType, ast.StringType, ast.VoidType:
		return nil, extractedType{typ: compiledType}
	case ast.SymbolType:
		err, entry := ctx.FindStructSymbol(typ.Name)
		if err != nil {
			return err, extractedType{typ: llvm.Type{}}
		}

		etyp := extractedType{
			// TODO: need to dinstinguish between passing a struct as value and as a pointer
			typ:    llvm.PointerType(compiledType, 0),
			sEntry: entry,
		}

		return nil, etyp
	case ast.PointerType:
		var sEntry *ast.StructSymbolTableEntry

		switch undType := typ.Underlying.(type) {
		case ast.SymbolType:
			err, entry := ctx.FindStructSymbol(undType.Name)
			if err != nil {
				return err, extractedType{}
			}

			sEntry = entry
		default:
		}

		etype := extractedType{typ: compiledType, sEntry: sEntry}

		return nil, etype
	case ast.ArrayType:
		var sEntry *ast.StructSymbolTableEntry

		switch undType := typ.Underlying.(type) {
		case ast.SymbolType:
			err, entry := ctx.FindStructSymbol(undType.Name)
			if err != nil {
				return err, extractedType{}
			}

			sEntry = entry
		default:
		}

		etype := extractedType{
			typ: llvm.PointerType(compiledType.ElementType(), 0),
			aEntry: &ast.ArraySymbolTableEntry{
				UnderlyingType:    compiledType.ElementType(),
				UnderlyingTypeDef: sEntry,
				ElementsCount:     compiledType.ArrayLength(),
				Type:              llvm.ArrayType(compiledType.ElementType(), typ.Size),
			},
		}

		return nil, etype
	default:
		return fmt.Errorf("FuncDeclStatement argument type %v not supported", t), extractedType{}
	}
}

func (g *LLVMGenerator) funcParams(ctx *ast.CompilerCtx, node *ast.FuncDeclStatement) (error, []llvm.Type) {
	params := []llvm.Type{}

	for _, arg := range node.Args {
		err, typ := g.extractType(ctx, arg.ArgType)
		if err != nil {
			return err, nil
		}

		params = append(params, typ.typ)
	}

	return nil, params
}
