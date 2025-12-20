package compiler

import (
	"fmt"
	"swahili/lang/ast"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type LLVMGenerator struct {
	Ctx        *ast.CompilerCtx
	lastResult *ast.CompilerResult
}

func NewLLVMGenerator(ctx *ast.CompilerCtx) *LLVMGenerator {
	return &LLVMGenerator{Ctx: ctx}
}

func (g *LLVMGenerator) setLastResult(res *ast.CompilerResult) {
	g.lastResult = res
}

func (g *LLVMGenerator) getLastResult() *ast.CompilerResult {
	res := g.lastResult
	g.lastResult = nil
	return res
}

// Statements

func (g *LLVMGenerator) VisitBlockStatement(node *ast.BlockStatement) error {
	for _, stmt := range node.Body {
		if err := stmt.Accept(g); err != nil {
			return err
		}
	}
	return nil
}

func (g *LLVMGenerator) VisitVarDeclaration(node *ast.VarDeclarationStatement) error {
	if g.Ctx.SymbolExistsInCurrentScope(node.Name) {
		return fmt.Errorf("variable %s is aleady defined", node.Name)
	}

	var val *ast.CompilerResult
	if node.Value != nil {
		if err := node.Value.Accept(g); err != nil {
			return err
		}
		val = g.getLastResult()
	}

	if val == nil {
		if node.ExplicitType == nil {
			return fmt.Errorf("VarDeclarationStatement: return value is nil and no explicit type <%s>", node.Name)
		}

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

		if arrayType, ok := node.ExplicitType.(ast.ArrayType); ok {
			err, underlyingLLVMType := arrayType.Underlying.LLVMType(g.Ctx)
			if err != nil {
				return err
			}

			arrayEntry := &ast.ArraySymbolTableEntry{
				UnderlyingType: underlyingLLVMType,
				Type:           llvmType,
				ElementsCount:  arrayType.Size,
			}

			err = g.Ctx.AddArraySymbol(node.Name, arrayEntry)
			if err != nil {
				return err
			}
		} else if symType, ok := node.ExplicitType.(ast.SymbolType); ok {
			err, typeDef := g.Ctx.FindStructSymbol(symType.Name)
			if err != nil {
				return err
			}

			entry.Ref = typeDef
		}

		err = g.Ctx.AddSymbol(node.Name, entry)
		if err != nil {
			return err
		}

		return nil
	}

	if _, ok := node.Value.(ast.MemberExpression); ok {
		if val.SymbolTableEntry != nil && val.SymbolTableEntry.Ref != nil {
			memberExpr, _ := node.Value.(ast.MemberExpression)
			propExpr, _ := memberExpr.Property.(ast.SymbolExpression)

			err, propIndex := val.SymbolTableEntry.Ref.Metadata.PropertyIndex(propExpr.Value)
			if err != nil {
				return err
			}
			// Load the value from the address so TypeCheck can work properly
			loadedValue := g.Ctx.Builder.CreateLoad(val.SymbolTableEntry.Ref.PropertyTypes[propIndex], *val.Value, "")
			val.Value = &loadedValue
		} else {
			return fmt.Errorf("VarDeclarationStatement/MemberExpression: unable to determine type")
		}
	}

	err := g.typeCheck(node.ExplicitType.Value(), val.Value.Type())
	if err != nil {
		return err
	}

	switch node.Value.(type) {
	case ast.ArrayOfStructsAccessExpression:
		alloc := g.Ctx.Builder.CreateAlloca(val.Value.Type(), fmt.Sprintf("alloc.%s", node.Name))
		g.Ctx.Builder.CreateStore(*val.Value, alloc)

		err = g.Ctx.AddSymbol(node.Name, &ast.SymbolTableEntry{Value: *val.Value, Address: &alloc, DeclaredType: node.ExplicitType})
		if err != nil {
			return err
		}
	case ast.ArrayAccessExpression:
		load := g.Ctx.Builder.CreateLoad(val.Value.AllocatedType(), *val.Value, "load.from-array")
		alloc := g.Ctx.Builder.CreateAlloca(val.Value.AllocatedType(), fmt.Sprintf("alloc.%s", node.Name))
		g.Ctx.Builder.CreateStore(load, alloc)

		err = g.Ctx.AddSymbol(node.Name, &ast.SymbolTableEntry{Value: load, Address: &alloc, DeclaredType: node.ExplicitType})
		if err != nil {
			return err
		}
	case ast.StructInitializationExpression:
		return g.compileStructInitializationExpression(node, val)
	case ast.StringExpression:
		return g.compileStringExpression(node, val)
	case ast.SymbolExpression:
		alloc := g.Ctx.Builder.CreateAlloca(val.Value.Type(), fmt.Sprintf("alloc.%s", node.Name))
		g.Ctx.Builder.CreateStore(*val.Value, alloc)

		entry := &ast.SymbolTableEntry{
			Value:        *val.Value,
			Address:      &alloc,
			DeclaredType: node.ExplicitType,
		}

		switch node.ExplicitType.(type) {
		case ast.SymbolType:
			entry.Ref = val.SymbolTableEntry.Ref
		default:
		}

		err = g.Ctx.AddSymbol(node.Name, entry)
		if err != nil {
			return err
		}
	case ast.NumberExpression, ast.FloatExpression, ast.BinaryExpression, ast.FunctionCallExpression, ast.MemberExpression:
		alloc := g.Ctx.Builder.CreateAlloca(val.Value.Type(), fmt.Sprintf("alloc.%s", node.Name))
		g.Ctx.Builder.CreateStore(*val.Value, alloc)

		entry := &ast.SymbolTableEntry{
			Value:        *val.Value,
			Address:      &alloc,
			DeclaredType: node.ExplicitType,
		}
		if val.SymbolTableEntry != nil {
			entry.Ref = val.SymbolTableEntry.Ref
		}

		err = g.Ctx.AddSymbol(node.Name, entry)
		if err != nil {
			return err
		}
	case ast.ArrayInitializationExpression:
		return g.compileArrArrayInitializationExpression(node, val)
	default:
		return fmt.Errorf("VarDeclarationStatement: Unhandled expression type (%v)", node.Value)
	}

	return nil
}

func (g *LLVMGenerator) typeCheck(t ast.DataType, k llvm.Type) error {
	switch k.TypeKind() {
	case llvm.ArrayTypeKind:
		switch t {
		case ast.DataTypeSymbol:
		case ast.DataTypeString, ast.DataTypeArray:
			// we are ok
		default:
			return fmt.Errorf("expected %s got %s", t, k.TypeKind())
		}
	case llvm.DoubleTypeKind:
		switch t {
		case ast.DataTypeSymbol:
		case ast.DataTypeFloat:
			// we are ok
		default:
			return fmt.Errorf("expected %s got %s", t, k.TypeKind())
		}
	case llvm.IntegerTypeKind:
		switch t {
		case ast.DataTypeSymbol:
		case ast.DataTypeNumber, ast.DataTypeNumber64:
			// we are ok
		default:
			return fmt.Errorf("expected %s got %s", t, k.TypeKind())
		}
	case llvm.PointerTypeKind:
		switch t {
		case ast.DataTypeString:
			switch k.ElementType() {
			case (*g.Ctx.Context).Int8Type():
				// we good
			default:
				// TODO fix this
				//	return fmt.Errorf("expected %s got pointer of unknown value %v", t, k.IsNil())
			}
		case ast.DataTypeNumber:
			switch k.ElementType().TypeKind() {
			case llvm.IntegerTypeKind:
				// we good
			default:
				// return fmt.Errorf("expected %s got pointer of %v", t, k)
			}
		case ast.DataTypeArray:
		case ast.DataTypeSymbol:
		default:
			return fmt.Errorf("unsupported element %s", t)
		}
	}

	return nil
}

func (g *LLVMGenerator) compileStructInitializationExpression(
	node *ast.VarDeclarationStatement,
	val *ast.CompilerResult,
) error {
	explicitType, ok := node.ExplicitType.(ast.SymbolType)
	if !ok {
		return fmt.Errorf("explicit type is not a symbol %v", node.ExplicitType)
	}

	err, typeDef := g.Ctx.FindStructSymbol(explicitType.Name)
	if err != nil {
		return fmt.Errorf("Could not find typedef for %s in structs symbol table", explicitType.Name)
	}

	entry := &ast.SymbolTableEntry{
		Value:        *val.Value,
		Ref:          typeDef,
		DeclaredType: node.ExplicitType,
	}

	err = g.Ctx.AddSymbol(node.Name, entry)
	if err != nil {
		return err
	}

	return nil
}

func (g *LLVMGenerator) compileStringExpression(
	node *ast.VarDeclarationStatement,
	val *ast.CompilerResult,
) error {
	glob := llvm.AddGlobal(*g.Ctx.Module, val.Value.Type(), fmt.Sprintf("global.%s", node.Name))
	glob.SetInitializer(*val.Value)

	alloc := g.Ctx.Builder.CreateAlloca(llvm.PointerType((*g.Ctx.Context).Int8Type(), 0), "")

	g.Ctx.Builder.CreateStore(glob, alloc)

	entry := &ast.SymbolTableEntry{
		Value:        *val.Value,
		Address:      &alloc,
		DeclaredType: node.ExplicitType,
	}

	err := g.Ctx.AddSymbol(node.Name, entry)
	if err != nil {
		return err
	}
	return nil
}

func (g *LLVMGenerator) compileArrArrayInitializationExpression(
	node *ast.VarDeclarationStatement,
	val *ast.CompilerResult,
) error {
	entry := &ast.SymbolTableEntry{
		Value:        *val.Value,
		Address:      val.Value,
		DeclaredType: node.ExplicitType,
	}

	err := g.Ctx.AddSymbol(node.Name, entry)
	if err != nil {
		return err
	}

	err = g.Ctx.AddArraySymbol(node.Name, val.ArraySymbolTableEntry)
	if err != nil {
		return err
	}

	return nil
}

func (g *LLVMGenerator) VisitReturnStatement(node *ast.ReturnStatement) error {
	if err := node.Value.Accept(g); err != nil {
		return err
	}
	res := g.getLastResult()

	switch node.Value.(type) {
	case ast.ArrayAccessExpression:
		val := g.Ctx.Builder.CreateLoad((*g.Ctx.Context).Int32Type(), *res.Value, "")
		g.Ctx.Builder.CreateRet(val)
	case ast.MemberExpression:
		loadedval := g.Ctx.Builder.CreateLoad(*res.StuctPropertyValueType, *res.Value, "")
		g.Ctx.Builder.CreateRet(loadedval)
	case ast.SymbolExpression, ast.BinaryExpression, ast.FunctionCallExpression:
		g.Ctx.Builder.CreateRet(*res.Value)
	case ast.StringExpression:
		alloc := g.Ctx.Builder.CreateAlloca(res.Value.Type(), "")
		g.Ctx.Builder.CreateStore(*res.Value, alloc)
		g.Ctx.Builder.CreateRet(alloc)
	case ast.NumberExpression, ast.FloatExpression:
		g.Ctx.Builder.CreateRet(*res.Value)
	default:
		return fmt.Errorf("ReturnStatement unknown expression <%s>", node.Value)
	}

	return nil
}

func (g *LLVMGenerator) VisitExpressionStatement(node *ast.ExpressionStatement) error {
	return node.Exp.Accept(g)
}

func (g *LLVMGenerator) VisitConditionalStatement(node *ast.ConditionalStatetement) error {
	if err := node.Condition.Accept(g); err != nil {
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

	if err := node.Failure.Accept(g); err != nil {
		return err
	}
	failureVal := g.getLastResult()

	lastElseBlock := g.Ctx.Builder.GetInsertBlock()
	opcode := lastElseBlock.LastInstruction().InstructionOpcode()
	if lastElseBlock.LastInstruction().IsNil() || (opcode != llvm.Ret && opcode != llvm.Br) {
		g.Ctx.Builder.CreateBr(mergeBlock)
	}
	g.Ctx.Builder.SetInsertPointAtEnd(mergeBlock)

	var phi llvm.Value
	if successVal != nil && failureVal != nil {
		phi = g.Ctx.Builder.CreatePHI(successVal.Value.Type(), "")
		phi.AddIncoming([]llvm.Value{*successVal.Value}, []llvm.BasicBlock{thenBlock})
		phi.AddIncoming([]llvm.Value{*failureVal.Value}, []llvm.BasicBlock{elseBlock})
	} else if successVal != nil {
		// If only success branch has value, we can't easily create a PHI node
		// unless we have a default value for the else branch.
		// For now, we'll just return successVal if we are in the then block,
		// but this is tricky for a visitor.
		// Matching original behavior:
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

	if err := node.Condition.Accept(g); err != nil {
		return err
	}
	condition := g.getLastResult()

	g.Ctx.Builder.CreateCondBr(*condition.Value, whileBodyBlock, whileMergeBlock)

	g.Ctx.Builder.SetInsertPointAtEnd(whileBodyBlock)

	if err := node.Body.Accept(g); err != nil {
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

		entry := ast.SymbolTableEntry{Value: p}

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
	case ast.NumberType, ast.Number64Type, ast.Number8Type, ast.FloatType, ast.StringType, ast.VoidType:
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

func (g *LLVMGenerator) VisitStructDeclaration(node *ast.StructDeclarationStatement) error {
	propertyTypes := []llvm.Type{}

	for _, t := range node.Types {
		err, llvmType := t.LLVMType(g.Ctx)
		if err != nil {
			return err
		}

		propertyTypes = append(propertyTypes, llvmType)
	}

	newtype := g.Ctx.Context.StructType(propertyTypes, false)

	err := g.Ctx.AddStructSymbol(node.Name, &ast.StructSymbolTableEntry{
		LLVMType:      newtype,
		PropertyTypes: propertyTypes,
		Metadata:      *node,
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *LLVMGenerator) VisitPrintStatement(node *ast.PrintStatetement) error {
	printableValues := []llvm.Value{}

	for _, v := range node.Values {
		switch strExpr := v.(type) {
		case ast.StringExpression:
			str := g.Ctx.Builder.CreateGlobalStringPtr(strExpr.Value, "")
			printableValues = append(printableValues, str)
		default:
			if err := v.Accept(g); err != nil {
				return err
			}
			res := g.getLastResult()

			switch v.(type) {
			case ast.MemberExpression:
				loadedval := g.Ctx.Builder.CreateLoad(*res.StuctPropertyValueType, *res.Value, "")
				printableValues = append(printableValues, loadedval)
			case ast.ArrayAccessExpression:
				vv, _ := v.(ast.ArrayAccessExpression)

				err, valForPrint := vv.CompileLLVMForPrint(g.Ctx)
				if err != nil {
					return err
				}

				printableValues = append(printableValues, *valForPrint)
			case ast.ArrayOfStructsAccessExpression, ast.NumberExpression, ast.FloatExpression, ast.SymbolExpression, ast.FunctionCallExpression, ast.BinaryExpression:
				printableValues = append(printableValues, *res.Value)
			default:
				return fmt.Errorf("Expression %v not supported in print statement", v)
			}
		}
	}

	printfFunc := g.Ctx.Module.NamedFunction("printf")
	g.Ctx.Builder.CreateCall(
		printfFunc.GlobalValueType(),
		printfFunc,
		printableValues,
		"",
	)

	return nil
}

func (g *LLVMGenerator) VisitMainStatement(node *ast.MainStatement) error {
	mainFuncType := llvm.FunctionType((*g.Ctx.Context).Int32Type(), []llvm.Type{}, false)
	mainFunc := llvm.AddFunction(*g.Ctx.Module, "main", mainFuncType)

	block := g.Ctx.Context.AddBasicBlock(mainFunc, "entry")
	g.Ctx.Builder.SetInsertPointAtEnd(block)

	if err := node.Body.Accept(g); err != nil {
		return err
	}

	lastBlock := g.Ctx.Builder.GetInsertBlock()
	if lastBlock.LastInstruction().IsNil() || lastBlock.LastInstruction().InstructionOpcode() != llvm.Ret {
		g.Ctx.Builder.CreateRet(llvm.ConstInt((*g.Ctx.Context).Int32Type(), 0, false))
	}

	return nil
}

// Expressions

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
	switch kind {
	case lexer.And:
		zero := llvm.ConstInt(l.Type(), 0, false)
		lBool := g.Ctx.Builder.CreateICmp(llvm.IntNE, l, zero, "")
		rBool := g.Ctx.Builder.CreateICmp(llvm.IntNE, r, zero, "")
		resBool := g.Ctx.Builder.CreateAnd(lBool, rBool, "")
		return nil, &ast.CompilerResult{Value: &resBool}
	case lexer.Or:
		zero := llvm.ConstInt(l.Type(), 0, false)
		lBool := g.Ctx.Builder.CreateICmp(llvm.IntNE, l, zero, "")
		rBool := g.Ctx.Builder.CreateICmp(llvm.IntNE, r, zero, "")
		resBool := g.Ctx.Builder.CreateOr(lBool, rBool, "")
		return nil, &ast.CompilerResult{Value: &resBool}
	case lexer.Plus:
		var res llvm.Value
		if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind {
			res = g.Ctx.Builder.CreateFAdd(l, r, "")
		} else {
			res = g.Ctx.Builder.CreateAdd(l, r, "")
		}
		return nil, &ast.CompilerResult{Value: &res}
	case lexer.Minus:
		var res llvm.Value
		if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind {
			res = g.Ctx.Builder.CreateFSub(l, r, "")
		} else {
			res = g.Ctx.Builder.CreateSub(l, r, "")
		}
		return nil, &ast.CompilerResult{Value: &res}
	case lexer.Star:
		var res llvm.Value
		if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind {
			res = g.Ctx.Builder.CreateFMul(l, r, "")
		} else {
			res = g.Ctx.Builder.CreateMul(l, r, "")
		}
		return nil, &ast.CompilerResult{Value: &res}
	case lexer.Modulo:
		res := g.Ctx.Builder.CreateSRem(l, r, "")
		return nil, &ast.CompilerResult{Value: &res}
	case lexer.Divide:
		var res llvm.Value
		if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind {
			res = g.Ctx.Builder.CreateFDiv(l, r, "")
		} else {
			res = g.Ctx.Builder.CreateSDiv(l, r, "")
		}
		return nil, &ast.CompilerResult{Value: &res}
	case lexer.GreaterThan:
		var res llvm.Value
		if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
			r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
			res = g.Ctx.Builder.CreateFCmp(llvm.FloatOGT, l, r, "")
		} else {
			res = g.Ctx.Builder.CreateICmp(llvm.IntSGT, l, r, "")
		}
		return nil, &ast.CompilerResult{Value: &res}
	case lexer.GreaterThanEquals:
		var res llvm.Value
		if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
			r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
			res = g.Ctx.Builder.CreateFCmp(llvm.FloatOGE, l, r, "")
		} else {
			res = g.Ctx.Builder.CreateICmp(llvm.IntSGE, l, r, "")
		}
		return nil, &ast.CompilerResult{Value: &res}
	case lexer.LessThan:
		var res llvm.Value
		if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
			r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
			res = g.Ctx.Builder.CreateFCmp(llvm.FloatOLT, l, r, "")
		} else {
			res = g.Ctx.Builder.CreateICmp(llvm.IntSLT, l, r, "")
		}
		return nil, &ast.CompilerResult{Value: &res}
	case lexer.LessThanEquals:
		var res llvm.Value
		if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
			r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
			res = g.Ctx.Builder.CreateFCmp(llvm.FloatOLE, l, r, "")
		} else {
			res = g.Ctx.Builder.CreateICmp(llvm.IntSLE, l, r, "")
		}
		return nil, &ast.CompilerResult{Value: &res}
	case lexer.Equals:
		var res llvm.Value
		if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
			r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
			res = g.Ctx.Builder.CreateFCmp(llvm.FloatOEQ, l, r, "")
		} else {
			res = g.Ctx.Builder.CreateICmp(llvm.IntEQ, l, r, "")
		}
		return nil, &ast.CompilerResult{Value: &res}
	case lexer.NotEquals:
		var res llvm.Value
		if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
			r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
			res = g.Ctx.Builder.CreateFCmp(llvm.FloatONE, l, r, "")
		} else {
			res = g.Ctx.Builder.CreateICmp(llvm.IntNE, l, r, "")
		}
		return nil, &ast.CompilerResult{Value: &res}
	default:
		return fmt.Errorf("Binary expressions : unsupported operator <%s>", kind), nil
	}
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

		args = append(args, *val.Value)
	}

	val := g.Ctx.Builder.CreateCall(*funcType, funcVal, args, "")

	g.setLastResult(&ast.CompilerResult{Value: &val})
	return nil
}

func (g *LLVMGenerator) VisitCallExpression(node *ast.CallExpression) error {
	if err := node.Caller.Accept(g); err != nil {
		return err
	}
	res := g.getLastResult()

	args := []llvm.Value{}

	for _, arg := range node.Arguments {
		if err := arg.Accept(g); err != nil {
			return err
		}
		val := g.getLastResult()

		args = append(args, *val.Value)
	}

	val := g.Ctx.Builder.CreateCall(res.Value.Type().ElementType(), *res.Value, args, "")

	g.setLastResult(&ast.CompilerResult{Value: &val})
	return nil
}

func (g *LLVMGenerator) VisitStringExpression(node *ast.StringExpression) error {
	val := llvm.ConstString(node.Value, true)
	g.setLastResult(&ast.CompilerResult{Value: &val})
	return nil
}

func (g *LLVMGenerator) VisitNumberExpression(node *ast.NumberExpression) error {
	val := llvm.ConstInt((*g.Ctx.Context).Int32Type(), uint64(node.Value), false)
	g.setLastResult(&ast.CompilerResult{Value: &val})
	return nil
}

func (g *LLVMGenerator) VisitFloatExpression(node *ast.FloatExpression) error {
	val := llvm.ConstFloat((*g.Ctx.Context).DoubleType(), node.Value)
	g.setLastResult(&ast.CompilerResult{Value: &val})
	return nil
}

func (g *LLVMGenerator) VisitSymbolExpression(node *ast.SymbolExpression) error {
	err, entry := g.Ctx.FindSymbol(node.Value)
	if err != nil {
		return err
	}

	if entry.Address != nil {
		var load llvm.Value

		switch entry.DeclaredType.(type) {
		case ast.StringType:
			load = g.Ctx.Builder.CreateLoad(entry.Address.Type(), *entry.Address, "")
		case ast.PointerType:
			load = g.Ctx.Builder.CreateLoad(entry.Address.AllocatedType(), *entry.Address, "")
		case ast.ArrayType:
			load = g.Ctx.Builder.CreateLoad(entry.Address.Type(), *entry.Address, "")
		case ast.NumberType, ast.FloatType:
			load = g.Ctx.Builder.CreateLoad(entry.Address.AllocatedType(), *entry.Address, "")
		default:
			load = g.Ctx.Builder.CreateLoad(entry.Address.AllocatedType(), *entry.Address, "")
		}

		g.setLastResult(&ast.CompilerResult{Value: &load, SymbolTableEntry: entry})
		return nil
	}

	g.setLastResult(&ast.CompilerResult{Value: &entry.Value, SymbolTableEntry: entry})
	return nil
}

func (g *LLVMGenerator) VisitAssignmentExpression(node *ast.AssignmentExpression) error {
	if err := node.Assignee.Accept(g); err != nil {
		return err
	}
	assignee := g.getLastResult()

	var valueToBeAssigned llvm.Value

	switch node.Value.(type) {
	case ast.StringExpression:
		if err := node.Value.Accept(g); err != nil {
			return err
		}
		value := g.getLastResult()

		glob := llvm.AddGlobal(*g.Ctx.Module, value.Value.Type(), "")
		glob.SetInitializer(*value.Value)
		valueToBeAssigned = glob
	case ast.ArrayAccessExpression:
		if err := node.Value.Accept(g); err != nil {
			return err
		}
		value := g.getLastResult()

		valueToBeAssigned = g.Ctx.Builder.CreateLoad(
			value.Value.AllocatedType(),
			*value.Value,
			"load.from-array",
		)
	default:
		if err := node.Value.Accept(g); err != nil {
			return err
		}
		value := g.getLastResult()

		valueToBeAssigned = *value.Value
	}

	if assignee.SymbolTableEntry != nil && assignee.SymbolTableEntry.Ref != nil {
		str := g.Ctx.Builder.CreateStore(valueToBeAssigned, *assignee.Value)
		g.setLastResult(&ast.CompilerResult{Value: &str})

		return nil
	}

	if assignee.SymbolTableEntry != nil && assignee.SymbolTableEntry.Address != nil {
		str := g.Ctx.Builder.CreateStore(valueToBeAssigned, *assignee.SymbolTableEntry.Address)
		g.setLastResult(&ast.CompilerResult{Value: &str})

		return nil
	}

	str := g.Ctx.Builder.CreateStore(valueToBeAssigned, *assignee.Value)
	g.setLastResult(&ast.CompilerResult{Value: &str})

	return nil
}

func (g *LLVMGenerator) VisitArrayAccessExpression(node *ast.ArrayAccessExpression) error {
	err, entry, array, itemIndex := g.findArraySymbolTableEntry(node)
	if err != nil {
		return err
	}

	itemPtr := g.Ctx.Builder.CreateInBoundsGEP(
		entry.UnderlyingType,
		array.Value,
		itemIndex,
		"",
	)

	g.setLastResult(&ast.CompilerResult{Value: &itemPtr, ArraySymbolTableEntry: entry})
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

func (g *LLVMGenerator) VisitMemberExpression(node *ast.MemberExpression) error {
	if _, ok := node.Object.(ast.ArrayAccessExpression); ok {
		return g.compileArrArrayAccessExpression(node)
	}

	if _, ok := node.Object.(ast.ArrayOfStructsAccessExpression); ok {
		return g.compileArrayOfStructsAccessExpression(node)
	}

	if _, ok := node.Object.(ast.MemberExpression); ok {
		return g.compileNestedMMemberExpression(node)
	}

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

func (g *LLVMGenerator) getNestedMemberAddress(member *ast.MemberExpression) (error, llvm.Value) {
	propName, err := g.getProperty(member)
	if err != nil {
		return err, llvm.Value{}
	}

	if arrayAccess, ok := member.Object.(ast.ArrayAccessExpression); ok {
		if err := arrayAccess.Accept(g); err != nil {
			return err, llvm.Value{}
		}
		res := g.getLastResult()

		if res.ArraySymbolTableEntry == nil || res.ArraySymbolTableEntry.UnderlyingTypeDef == nil {
			return fmt.Errorf("array elements are not structs"), llvm.Value{}
		}

		structType := res.ArraySymbolTableEntry.UnderlyingTypeDef

		propIndex, err := g.resolveStructAccess(structType, propName)
		if err != nil {
			return err, llvm.Value{}
		}

		addr := g.Ctx.Builder.CreateStructGEP(structType.LLVMType, *res.Value, propIndex, "")
		return nil, addr
	}

	if arrayStructAccess, ok := member.Object.(ast.ArrayOfStructsAccessExpression); ok {
		if err := arrayStructAccess.Accept(g); err != nil {
			return err, llvm.Value{}
		}
		res := g.getLastResult()

		if res.SymbolTableEntry == nil || res.SymbolTableEntry.Address == nil || res.SymbolTableEntry.Ref == nil {
			return fmt.Errorf("array element property is not a struct"), llvm.Value{}
		}

		structType := res.SymbolTableEntry.Ref

		propIndex, err := g.resolveStructAccess(structType, propName)
		if err != nil {
			return err, llvm.Value{}
		}

		addr := g.Ctx.Builder.CreateStructGEP(structType.LLVMType, *res.SymbolTableEntry.Address, propIndex, "")
		return nil, addr
	}

	if nestedMember, ok := member.Object.(ast.MemberExpression); ok {
		err, nestedAddr := g.getNestedMemberAddress(&nestedMember)
		if err != nil {
			return err, llvm.Value{}
		}

		baseObj, err := g.findBaseSymbol(nestedMember.Object)
		if err != nil {
			return err, llvm.Value{}
		}

		var structType *ast.StructSymbolTableEntry

		err, varDef := g.Ctx.FindSymbol(baseObj.Value)
		if err != nil {
			// Try finding array
			errArr, arrDef := g.Ctx.FindArraySymbol(baseObj.Value)
			if errArr != nil {
				return err, llvm.Value{}
			}
			structType = arrDef.UnderlyingTypeDef
		} else {
			structType = varDef.Ref
		}

		nestedStructType, err := g.getNestedStructType(&nestedMember, structType)
		if err != nil {
			return err, llvm.Value{}
		}

		propIndex, err := g.resolveStructAccess(nestedStructType, propName)
		if err != nil {
			return err, llvm.Value{}
		}

		addr := g.Ctx.Builder.CreateStructGEP(nestedStructType.LLVMType, nestedAddr, propIndex, "")
		return nil, addr
	}

	obj, ok := member.Object.(ast.SymbolExpression)
	if !ok {
		return fmt.Errorf("struct object should be a symbol"), llvm.Value{}
	}

	err, varDef := g.Ctx.FindSymbol(obj.Value)
	if err != nil {
		return err, llvm.Value{}
	}

	propIndex, err := g.resolveStructAccess(varDef.Ref, propName)
	if err != nil {
		return err, llvm.Value{}
	}

	addr := g.Ctx.Builder.CreateStructGEP(varDef.Ref.LLVMType, varDef.Value, propIndex, "")
	return nil, addr
}

func (g *LLVMGenerator) findBaseSymbol(obj ast.Expression) (ast.SymbolExpression, error) {
	switch v := obj.(type) {
	case ast.SymbolExpression:
		return v, nil
	case ast.MemberExpression:
		return g.findBaseSymbol(v.Object)
	case ast.ArrayAccessExpression:
		return g.findBaseSymbol(v.Name)
	case ast.ArrayOfStructsAccessExpression:
		return g.findBaseSymbol(v.Name)
	default:
		return ast.SymbolExpression{}, fmt.Errorf("cannot resolve base object")
	}
}

func (g *LLVMGenerator) getProperty(expr *ast.MemberExpression) (string, error) {
	prop, ok := expr.Property.(ast.SymbolExpression)
	if !ok {
		return "", fmt.Errorf("struct property should be a symbol")
	}
	return prop.Value, nil
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

func (g *LLVMGenerator) getNestedStructType(
	member *ast.MemberExpression,
	baseStructType *ast.StructSymbolTableEntry,
) (*ast.StructSymbolTableEntry, error) {
	propName, err := g.getProperty(member)
	if err != nil {
		return nil, err
	}

	currentStructType := baseStructType
	if nestedMember, ok := member.Object.(ast.MemberExpression); ok {
		currentStructType, err = g.getNestedStructType(&nestedMember, baseStructType)
		if err != nil {
			return nil, err
		}
	}

	propIndex, err := g.resolveStructAccess(currentStructType, propName)
	if err != nil {
		return nil, err
	}

	propType := currentStructType.Metadata.Types[propIndex]
	if symbolType, ok := propType.(ast.SymbolType); ok {
		err, structDef := g.Ctx.FindStructSymbol(symbolType.Name)
		if err != nil {
			return nil, fmt.Errorf("cannot find struct type %s", symbolType.Name)
		}
		return structDef, nil
	}

	return nil, fmt.Errorf("property %s is not a struct type", propName)
}

func (g *LLVMGenerator) compileArrArrayAccessExpression(expr *ast.MemberExpression) error {
	arrayAccess, _ := expr.Object.(ast.ArrayAccessExpression)

	if err := arrayAccess.Accept(g); err != nil {
		return err
	}
	res := g.getLastResult()

	if res.ArraySymbolTableEntry == nil || res.ArraySymbolTableEntry.UnderlyingTypeDef == nil {
		return fmt.Errorf("array elements are not structs")
	}

	structType := res.ArraySymbolTableEntry.UnderlyingTypeDef

	propName, err := g.getProperty(expr)
	if err != nil {
		return err
	}

	propIndex, err := g.resolveStructAccess(structType, propName)
	if err != nil {
		return err
	}

	addr := g.Ctx.Builder.CreateStructGEP(structType.LLVMType, *res.Value, propIndex, "")
	propType := structType.PropertyTypes[propIndex]

	g.setLastResult(&ast.CompilerResult{
		Value:                  &addr,
		SymbolTableEntry:       &ast.SymbolTableEntry{Ref: structType},
		StuctPropertyValueType: &propType,
	})
	return nil
}

func (g *LLVMGenerator) compileArrayOfStructsAccessExpression(expr *ast.MemberExpression) error {
	arrayStructAccess, _ := expr.Object.(ast.ArrayOfStructsAccessExpression)

	if err := arrayStructAccess.Accept(g); err != nil {
		return err
	}
	res := g.getLastResult()

	if res.SymbolTableEntry == nil || res.SymbolTableEntry.Address == nil || res.SymbolTableEntry.Ref == nil {
		return fmt.Errorf("array element property is not a struct")
	}

	structType := res.SymbolTableEntry.Ref

	propName, err := g.getProperty(expr)
	if err != nil {
		return err
	}

	propIndex, err := g.resolveStructAccess(structType, propName)
	if err != nil {
		return err
	}

	addr := g.Ctx.Builder.CreateStructGEP(structType.LLVMType, *res.SymbolTableEntry.Address, propIndex, "")
	propType := structType.PropertyTypes[propIndex]

	g.setLastResult(&ast.CompilerResult{
		Value:                  &addr,
		SymbolTableEntry:       &ast.SymbolTableEntry{Ref: structType},
		StuctPropertyValueType: &propType,
	})
	return nil
}

func (g *LLVMGenerator) compileNestedMMemberExpression(expr *ast.MemberExpression) error {
	nestedMember, _ := expr.Object.(ast.MemberExpression)
	err, nestedAddr := g.getNestedMemberAddress(&nestedMember)
	if err != nil {
		return err
	}

	baseObj, err := g.findBaseSymbol(nestedMember.Object)
	if err != nil {
		return err
	}

	var structType *ast.StructSymbolTableEntry

	err, varDef := g.Ctx.FindSymbol(baseObj.Value)
	if err != nil {
		errArr, arrDef := g.Ctx.FindArraySymbol(baseObj.Value)
		if errArr != nil {
			return err
		}

		structType = arrDef.UnderlyingTypeDef
	} else {
		structType = varDef.Ref
	}

	nestedStructType, err := g.getNestedStructType(&nestedMember, structType)
	if err != nil {
		return err
	}

	propName, err := g.getProperty(expr)
	if err != nil {
		return err
	}

	propIndex, err := g.resolveStructAccess(nestedStructType, propName)
	if err != nil {
		return err
	}

	addr := g.Ctx.Builder.CreateStructGEP(nestedStructType.LLVMType, nestedAddr, propIndex, "")
	propType := nestedStructType.PropertyTypes[propIndex]

	g.setLastResult(&ast.CompilerResult{
		Value:                  &addr,
		SymbolTableEntry:       &ast.SymbolTableEntry{Ref: nestedStructType},
		StuctPropertyValueType: &propType,
	})
	return nil
}

func (g *LLVMGenerator) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
	err, structDef := g.Ctx.FindStructSymbol(node.Name)
	if err != nil {
		return fmt.Errorf("StructInitializationExpression: Undefined struct named %s", node.Name)
	}

	structInstance := g.Ctx.Builder.CreateAlloca(structDef.LLVMType, fmt.Sprintf("%s.instance", node.Name))

	for _, name := range node.Properties {
		err, propIndex := structDef.Metadata.PropertyIndex(name)
		if err != nil {
			return fmt.Errorf("StructInitializationExpression: property %s not found", name)
		}

		expr := node.Values[propIndex]
		if err := expr.Accept(g); err != nil {
			return err
		}
		val := g.getLastResult()

		field1Ptr := g.Ctx.Builder.CreateStructGEP(structDef.LLVMType, structInstance, propIndex, "")

		switch expr.(type) {
		case ast.StringExpression:
			glob := llvm.AddGlobal(*g.Ctx.Module, val.Value.Type(), "")
			glob.SetInitializer(*val.Value)
			g.Ctx.Builder.CreateStore(glob, field1Ptr)
		case ast.NumberExpression, ast.FloatExpression:
			g.Ctx.Builder.CreateStore(*val.Value, field1Ptr)
		case ast.SymbolExpression:
			if val.SymbolTableEntry.Address != nil {
				if val.Value.Type().TypeKind() == llvm.PointerTypeKind {
					g.Ctx.Builder.CreateStore(*val.SymbolTableEntry.Address, field1Ptr)

					break
				}

				g.Ctx.Builder.CreateStore(*val.Value, field1Ptr)

				break
			}

			if val.SymbolTableEntry.Ref != nil {
				load := g.Ctx.Builder.CreateLoad(val.SymbolTableEntry.Ref.LLVMType, *val.Value, "")
				g.Ctx.Builder.CreateStore(load, field1Ptr)

				break
			}

			g.Ctx.Builder.CreateStore(*val.Value, field1Ptr)
		case ast.StructInitializationExpression:
			nestedStructType := structDef.PropertyTypes[propIndex]
			loadedNestedStruct := g.Ctx.Builder.CreateLoad(nestedStructType, *val.Value, "")
			g.Ctx.Builder.CreateStore(loadedNestedStruct, field1Ptr)
		case ast.ArrayInitializationExpression:
			astType := structDef.Metadata.Types[propIndex]
			if _, isPointer := astType.(ast.PointerType); isPointer {
				g.Ctx.Builder.CreateStore(*val.Value, field1Ptr)
			} else {
				arrayType := structDef.PropertyTypes[propIndex]
				loadedArray := g.Ctx.Builder.CreateLoad(arrayType, *val.Value, "")
				g.Ctx.Builder.CreateStore(loadedArray, field1Ptr)
			}
		default:
			return fmt.Errorf("StructInitializationExpression: expression %v not implemented", expr)
		}
	}

	g.setLastResult(&ast.CompilerResult{Value: &structInstance})
	return nil
}

func (g *LLVMGenerator) VisitArrayInitializationExpression(node *ast.ArrayInitializationExpression) error {
	llvmType, structDef, err := g.extractArrayInitializationType(node)
	if err != nil {
		return err
	}

	err, elements := g.getArrayElements(node, *llvmType)
	if err != nil {
		return err
	}

	arrayType := llvm.ArrayType(*llvmType, len(elements))
	alloc := g.Ctx.Builder.CreateAlloca(arrayType, "")

	for i, element := range elements {
		ptr := g.Ctx.Builder.CreateInBoundsGEP(
			arrayType,
			alloc,
			[]llvm.Value{
				llvm.ConstInt((*g.Ctx.Context).Int32Type(), 0, false),
				llvm.ConstInt((*g.Ctx.Context).Int32Type(), uint64(i), false),
			},
			"",
		)
		g.Ctx.Builder.CreateStore(element, ptr)
	}

	g.setLastResult(&ast.CompilerResult{
		Value: &alloc,
		ArraySymbolTableEntry: &ast.ArraySymbolTableEntry{
			UnderlyingType:    *llvmType,
			UnderlyingTypeDef: structDef,
			Type:              arrayType,
			ElementsCount:     len(elements),
		},
	})

	return nil
}

func (g *LLVMGenerator) extractArrayInitializationType(expr *ast.ArrayInitializationExpression) (*llvm.Type, *ast.StructSymbolTableEntry, error) {
	arrayType, ok := expr.Underlying.(ast.ArrayType)
	if !ok {
		return nil, nil, fmt.Errorf("Type (%s) cannot be casted to array type", expr.Underlying)
	}

	switch arrayType.Underlying.Value() {
	case ast.DataTypeFloat:
		typ := (*g.Ctx.Context).DoubleType()
		return &typ, nil, nil
	case ast.DataTypeNumber, ast.DataTypeNumber8, ast.DataTypeNumber64:
		var typ llvm.Type
		if arrayType.Underlying.Value() == ast.DataTypeNumber64 {
			typ = (*g.Ctx.Context).Int64Type()
		} else if arrayType.Underlying.Value() == ast.DataTypeNumber8 {
			typ = (*g.Ctx.Context).Int8Type()
		} else {
			typ = (*g.Ctx.Context).Int32Type()
		}
		return &typ, nil, nil
	case ast.DataTypeSymbol:
		sym, _ := arrayType.Underlying.(ast.SymbolType)
		err, sdef := g.Ctx.FindStructSymbol(sym.Name)
		if err != nil {
			return nil, nil, fmt.Errorf("Type (%s) is not a valid struct", sym.Name)
		}
		return &sdef.LLVMType, sdef, nil
	case ast.DataTypeString:
		typ := llvm.PointerType((*g.Ctx.Context).Int32Type(), 0)
		return &typ, nil, nil
	default:
		return nil, nil, fmt.Errorf("Type (%v) not implemented in array expression", arrayType.Underlying.Value())
	}
}

func (g *LLVMGenerator) getArrayElements(expr *ast.ArrayInitializationExpression, elementType llvm.Type) (error, []llvm.Value) {
	elements := []llvm.Value{}

	for _, element := range expr.Contents {
		if err := element.Accept(g); err != nil {
			return err, nil
		}
		res := g.getLastResult()

		if res.Value.Type() != elementType {
			return fmt.Errorf("ArrayInitializationExpression: element type mismatch"), nil
		}

		elements = append(elements, *res.Value)
	}

	return nil, elements
}

func (g *LLVMGenerator) VisitArrayOfStructsAccessExpression(node *ast.ArrayOfStructsAccessExpression) error {
	err, array, itemIndex := g.findArrayOfStructsSymbolTableEntry(node)
	if err != nil {
		return err
	}

	itemPtr := g.Ctx.Builder.CreateInBoundsGEP(
		array.Ref.LLVMType,
		array.Value,
		itemIndex,
		"",
	)

	g.setLastResult(&ast.CompilerResult{Value: &itemPtr, SymbolTableEntry: array})
	return nil
}

func (g *LLVMGenerator) findArrayOfStructsSymbolTableEntry(expr *ast.ArrayOfStructsAccessExpression) (error, *ast.SymbolTableEntry, []llvm.Value) {
	var array *ast.SymbolTableEntry

	switch expr.Name.(type) {
	case ast.SymbolExpression:
		varName, _ := expr.Name.(ast.SymbolExpression)

		err, arrayEntry := g.Ctx.FindSymbol(varName.Value)
		if err != nil {
			key := "ArrayOfStructsAccessExpression.NotFoundInSymbolTable"
			return g.Ctx.Dialect.Error(key, varName.Value), nil, nil
		}
		array = arrayEntry
	default:
		return fmt.Errorf("ArrayOfStructsAccessExpression not implemented"), nil, nil
	}

	var indices []llvm.Value

	switch expr.Index.(type) {
	case ast.NumberExpression:
		idx, _ := expr.Index.(ast.NumberExpression)

		indices = []llvm.Value{
			llvm.ConstInt((*g.Ctx.Context).Int32Type(), uint64(idx.Value), false),
		}
	case ast.SymbolExpression, ast.BinaryExpression, ast.MemberExpression:
		if err := expr.Index.Accept(g); err != nil {
			return err, nil, nil
		}
		res := g.getLastResult()

		if res.Value.Type().TypeKind() == llvm.PointerTypeKind {
			load := g.Ctx.Builder.CreateLoad(res.Value.AllocatedType(), *res.Value, "")
			indices = []llvm.Value{load}
		} else {
			indices = []llvm.Value{*res.Value}
		}
	default:
		key := "ArrayOfStructsAccessExpression.AccessedIndexIsNotANumber"
		return g.Ctx.Dialect.Error(key, expr.Index), nil, nil
	}

	return nil, array, indices
}

func (g *LLVMGenerator) VisitPrefixExpression(node *ast.PrefixExpression) error {
	if err := node.RightExpression.Accept(g); err != nil {
		return err
	}
	res := g.getLastResult()

	var val llvm.Value
	switch node.Operator.Kind {
	case lexer.Minus:
		if res.Value.Type().TypeKind() == llvm.FloatTypeKind || res.Value.Type().TypeKind() == llvm.DoubleTypeKind {
			val = g.Ctx.Builder.CreateFNeg(*res.Value, "")
		} else {
			val = g.Ctx.Builder.CreateNeg(*res.Value, "")
		}
	case lexer.Not:
		zero := llvm.ConstInt(res.Value.Type(), 0, false)
		boolVal := g.Ctx.Builder.CreateICmp(llvm.IntEQ, *res.Value, zero, "")
		val = g.Ctx.Builder.CreateZExt(boolVal, res.Value.Type(), "")
	default:
		return fmt.Errorf("PrefixExpression: operator %s not supported", node.Operator.Kind)
	}

	g.setLastResult(&ast.CompilerResult{Value: &val})
	return nil
}
