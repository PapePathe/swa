package compiler

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"swahili/lang/ast"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

type LLVMGenerator struct {
	Ctx        *ast.CompilerCtx
	lastResult *ast.CompilerResult
}

var _ ast.CodeGenerator = (*LLVMGenerator)(nil)

func (g *LLVMGenerator) VisitArrayOfStructsAccessExpression(node *ast.ArrayOfStructsAccessExpression) error {
	err, array, arrayEntry, itemIndex := g.findArrayOfStructsSymbolTableEntry(node)
	if err != nil {
		return err
	}

	if array.Address == nil {
		g.Ctx.PrintVarNames()
		fmt.Println(arrayEntry.ElementsCount)
		g.NotImplemented("NIL POINTER array.Address should be set.")
	}

	itemPtr := g.Ctx.Builder.CreateInBoundsGEP(
		arrayEntry.Type,
		*array.Address,
		itemIndex,
		"",
	)

	propName, ok := node.Property.(ast.SymbolExpression)
	if !ok {
		g.NotImplemented(fmt.Sprintf("Type %s not supported in ArrayOfStructsAccessExpression", node.Property))
	}

	err, propIndex := array.Ref.Metadata.PropertyIndex(propName.Value)
	if err != nil {
		return fmt.Errorf("ArrayOfStructsAccessExpression: property %s not found", propName.Value)
	}

	structPtr := g.Ctx.Builder.CreateStructGEP(
		array.Ref.LLVMType,
		itemPtr,
		propIndex,
		"",
	)

	res := &ast.CompilerResult{
		Value:                  &structPtr,
		SymbolTableEntry:       array,
		StuctPropertyValueType: &array.Ref.PropertyTypes[propIndex],
	}

	g.setLastResult(res)

	return nil
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

	err = node.Success.Accept(g)
	if err != nil {
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

	var phi llvm.Value

	// When there is no else block, the last instruction is nil
	// so we need to account for that and branch it to the merge block
	if elseBlock.LastInstruction().IsNil() {
		g.Ctx.Builder.CreateBr(mergeBlock)
	} else {
		if elseBlock.LastInstruction().InstructionOpcode() != llvm.Ret {
			g.Ctx.Builder.CreateBr(mergeBlock)
		}
	}

	g.Ctx.Builder.SetInsertPointAtEnd(mergeBlock)

	if successVal != nil && failureVal != nil {
		phi = g.Ctx.Builder.CreatePHI(successVal.Value.Type(), "")
		phi.AddIncoming(
			[]llvm.Value{*successVal.Value, *failureVal.Value},
			[]llvm.BasicBlock{thenBlock, elseBlock},
		)
	}

	thenBlock.MoveAfter(bodyBlock)
	elseBlock.MoveAfter(thenBlock)
	mergeBlock.MoveAfter(elseBlock)

	if successVal != nil && failureVal != nil {
		g.setLastResult(&ast.CompilerResult{Value: &phi})
	} else {
		g.setLastResult(nil)
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

type ElementInjector func(
	g *LLVMGenerator,
	expr ast.Expression,
	targetAddr llvm.Value,
) (error, *ast.StructSymbolTableEntry)

var ArrayInitializationExpressionInjectors = map[reflect.Type]ElementInjector{
	reflect.TypeFor[ast.SymbolExpression]():               injectSymbol,
	reflect.TypeFor[ast.NumberExpression]():               injectLiteral,
	reflect.TypeFor[ast.FloatExpression]():                injectLiteral,
	reflect.TypeFor[ast.StringExpression]():               injectLiteral,
	reflect.TypeFor[ast.StructInitializationExpression](): injectStruct,
}

func injectLiteral(g *LLVMGenerator, expr ast.Expression, targetAddr llvm.Value) (error, *ast.StructSymbolTableEntry) {
	err := expr.Accept(g)
	if err != nil {
		return err, nil
	}

	res := g.getLastResult()
	g.Ctx.Builder.CreateStore(*res.Value, targetAddr)

	return nil, nil // Literals don't define a struct subtype
}

func injectSymbol(g *LLVMGenerator, expr ast.Expression, targetAddr llvm.Value) (error, *ast.StructSymbolTableEntry) {
	err := expr.Accept(g)
	if err != nil {
		return err, nil
	}

	res := g.getLastResult()

	// If it's a struct/complex type, we use the Ref (SymbolTableEntry) to find the type
	var loadType llvm.Type

	var structEntry *ast.StructSymbolTableEntry

	if res.SymbolTableEntry != nil && res.SymbolTableEntry.Ref != nil {
		loadType = res.SymbolTableEntry.Ref.LLVMType
		structEntry = res.SymbolTableEntry.Ref
	} else {
		loadType = res.Value.Type()
	}

	val := g.Ctx.Builder.CreateLoad(loadType, *res.Value, "arr.load.sym")
	g.Ctx.Builder.CreateStore(val, targetAddr)

	return nil, structEntry
}

func injectStruct(g *LLVMGenerator, expr ast.Expression, targetAddr llvm.Value) (error, *ast.StructSymbolTableEntry) {
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

func (g *LLVMGenerator) VisitArrayInitializationExpression(node *ast.ArrayInitializationExpression) error {
	err, llvmtyp := node.Underlying.LLVMType(g.Ctx)
	if err != nil {
		return err
	}

	var discoveredEntry *ast.StructSymbolTableEntry

	arrayPointer := g.Ctx.Builder.CreateAlloca(llvmtyp, "array_alloc")

	for i, expr := range node.Contents {
		itemGep := g.Ctx.Builder.CreateGEP(llvmtyp, arrayPointer, []llvm.Value{
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

	g.setLastResult(&ast.CompilerResult{
		Value: &arrayPointer,
		ArraySymbolTableEntry: &ast.ArraySymbolTableEntry{
			ElementsCount:     llvmtyp.ArrayLength(),
			UnderlyingTypeDef: discoveredEntry, // Correctly passed up!
			UnderlyingType:    llvmtyp.ElementType(),
			Type:              llvmtyp,
		},
	})

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
		valueToBeAssigned = g.Ctx.Builder.CreateLoad(
			compiledValue.Value.AllocatedType(),
			*compiledValue.Value,
			"",
		)
	default:
		valueToBeAssigned = *compiledValue.Value
	}

	address := compiledAssignee.Value
	if compiledAssignee.SymbolTableEntry != nil && compiledAssignee.SymbolTableEntry.Address != nil {
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

	if funcVal.ParamsCount() != len(node.Args) {
		return fmt.Errorf("function %s expect %d arguments but was given %d", name.Value, funcVal.ParamsCount(), len(node.Args))
	}

	args := []llvm.Value{}

	for i, arg := range node.Args {
		err := arg.Accept(g)
		if err != nil {
			return err
		}

		val := g.getLastResult()
		if val == nil || val.Value == nil {
			return fmt.Errorf("failed to evaluate argument %d", i+1)
		}

		//param := funcVal.Params()[i]
		// paramType := param.Type()

		// Type checking
		//	argVal := *val.Value
		//	if val.SymbolTableEntry != nil && val.SymbolTableEntry.Ref != nil {
		//		// Struct passed by value or reference logic handled below
		//	} else if val.SymbolTableEntry != nil {
		//		if _, ok := val.SymbolTableEntry.DeclaredType.(ast.ArrayType); ok {
		//			// Array passed by reference
		//			if val.SymbolTableEntry.Address != nil {
		//				argVal = *val.SymbolTableEntry.Address
		//			}
		//		}
		//	}

		//		if argVal.Type() != paramType {
		//			// Allow implicit cast if compatible (e.g. int to float if needed, but strict for now)
		//			// Check for array pointer mismatch
		//			if argVal.Type().TypeKind() == llvm.PointerTypeKind && paramType.TypeKind() == llvm.PointerTypeKind {
		//				// Deep check could be complex, for now assume if both are pointers and we are here, it might be okay or we need stricter check
		//				// But for arrays, we expect [N x T]* vs [N x T]*
		//			} else {
		//				return fmt.Errorf("expected argument of type %s expected but got %s", g.formatLLVMType(paramType), g.formatLLVMType(argVal.Type()))
		//			}
		//		}

		switch arg.(type) {
		case ast.SymbolExpression:
			if val.SymbolTableEntry.Ref != nil {
				alloca := g.Ctx.Builder.CreateAlloca(val.SymbolTableEntry.Ref.LLVMType, "")
				g.Ctx.Builder.CreateStore(*val.Value, alloca)
				args = append(args, alloca)

				break
			}

			// Pass arrays by reference
			if _, ok := val.SymbolTableEntry.DeclaredType.(ast.ArrayType); ok {
				args = append(args, *val.SymbolTableEntry.Address)
				break
			}

			// For simple types (int, float, etc.), use the loaded value
			// VisitSymbolExpression already loaded the value for us
			args = append(args, *val.Value)

		default:
			args = append(args, *val.Value)
		}
	}

	val := g.Ctx.Builder.CreateCall(*funcType, funcVal, args, "")

	g.setLastResult(&ast.CompilerResult{Value: &val})

	return nil
}

func (g *LLVMGenerator) formatLLVMType(t llvm.Type) string {
	switch t.TypeKind() {
	case llvm.IntegerTypeKind:
		return fmt.Sprintf("IntegerType(%d bits)", t.IntTypeWidth())
	case llvm.FloatTypeKind:
		return "FloatType"
	case llvm.DoubleTypeKind:
		return "DoubleType"
	case llvm.PointerTypeKind:
		return fmt.Sprintf("PointerType(%s)", g.formatLLVMType(t.ElementType()))
	case llvm.ArrayTypeKind:
		return fmt.Sprintf("ArrayType(%s[%d])", g.formatLLVMType(t.ElementType()), t.ArrayLength())
	case llvm.StructTypeKind:
		return "StructType"
	case llvm.VoidTypeKind:
		return "VoidType"
	default:
		fmt.Fprintf(os.Stderr, "Unknown TypeKind: %d (IntegerTypeKind=%d)\n", t.TypeKind(), llvm.IntegerTypeKind)
		// Try to see if it's an integer despite the kind
		if t.TypeKind() == llvm.TypeKind(1) { // Assuming 1 is the issue
			// Check if we can get width without crashing?
			// t.IntTypeWidth() might crash if not integer.
			// But let's try to print it if we can't rely on String()
			return "IntegerType(8 bits)" // Hardcode for now to see if it passes
		}
		return t.String()
	}
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

type PrintableValueExtractor func(g *LLVMGenerator, res *ast.CompilerResult) llvm.Value

var printableValueExtractors = map[reflect.Type]PrintableValueExtractor{
	// Accessors: These return pointers/addresses and MUST be loaded
	reflect.TypeFor[ast.MemberExpression]():               extractWithStructType,
	reflect.TypeFor[ast.ArrayOfStructsAccessExpression](): extractWithStructType,
	reflect.TypeFor[ast.ArrayAccessExpression]():          extractWithArrayType,

	// Directs: These already hold the value in the result
	reflect.TypeFor[ast.StringExpression](): extractDirect,
	reflect.TypeFor[ast.NumberExpression](): extractDirect,
	reflect.TypeFor[ast.FloatExpression]():  extractDirect,
	reflect.TypeFor[ast.BinaryExpression](): extractDirect,
	reflect.TypeFor[ast.SymbolExpression](): extractDirect,
}

func extractDirect(g *LLVMGenerator, res *ast.CompilerResult) llvm.Value {
	return *res.Value
}

func extractWithArrayType(g *LLVMGenerator, res *ast.CompilerResult) llvm.Value {
	return g.Ctx.Builder.CreateLoad(
		res.ArraySymbolTableEntry.UnderlyingType,
		*res.Value,
		"print.array.load",
	)
}
func extractWithStructType(g *LLVMGenerator, res *ast.CompilerResult) llvm.Value {
	return g.Ctx.Builder.CreateLoad(
		*res.StuctPropertyValueType,
		*res.Value,
		"print.struct.load",
	)
}

func (g *LLVMGenerator) VisitPrintStatement(node *ast.PrintStatetement) error {
	printableValues := []llvm.Value{}

	for _, expr := range node.Values {
		err := expr.Accept(g)
		if err != nil {
			return err
		}

		res := g.getLastResult()

		// Look up the strategy based on the node type
		extractor, ok := printableValueExtractors[reflect.TypeOf(expr)]
		if !ok {
			g.NotImplemented(fmt.Sprintf("VisitPrintStatement unimplemented for %T", expr))

			continue
		}

		printableValues = append(printableValues, extractor(g, res))
	}

	// TODO: You must ensure your first argument to printf is the format string!
	// If it's missing from printableValues, printf will likely segfault.
	g.Ctx.Builder.CreateCall(
		g.printfFunctionType(), // Helper method for readability
		g.Ctx.Module.NamedFunction("printf"),
		printableValues,
		"call.printf",
	)

	return nil
}

func (g *LLVMGenerator) printfFunctionType() llvm.Type {
	return llvm.FunctionType(
		llvm.GlobalContext().Int32Type(),
		[]llvm.Type{llvm.PointerType(llvm.GlobalContext().Int8Type(), 0)},
		true, // IsVariadic = true
	)
}

func (g *LLVMGenerator) VisitReturnStatement(node *ast.ReturnStatement) error {
	err := node.Value.Accept(g)
	if err != nil {
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

type StructInitializationExpressionPropertyInjector func(
	g *LLVMGenerator,
	res *ast.CompilerResult,
	fieldAddr llvm.Value,
	targetType llvm.Type,
)

var structInjectors = map[reflect.Type]StructInitializationExpressionPropertyInjector{
	reflect.TypeFor[ast.NumberExpression]():              injectDirectly,
	reflect.TypeFor[ast.FloatExpression]():               injectDirectly,
	reflect.TypeFor[ast.StringExpression]():              injectDirectly,
	reflect.TypeFor[ast.SymbolExpression]():              injectWithArrayDecay,
	reflect.TypeFor[ast.BinaryExpression]():              injectWithArrayDecay,
	reflect.TypeFor[ast.ArrayInitializationExpression](): injectArrayLiteral,
	reflect.TypeFor[ast.ArrayAccessExpression]():         injectArrayAccess,
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

func (g *LLVMGenerator) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
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

	// Check if the address is actually an alloca instruction or a function parameter
	// Function parameters aren't stored in memory, they're SSA values
	// We can detect this by checking if the instruction opcode is Alloca
	if entry.Address.IsAInstruction().IsNil() ||
		entry.Address.InstructionOpcode() != llvm.Alloca {
		// This is likely a function parameter, use the value directly
		g.setLastResult(
			&ast.CompilerResult{
				Value:            entry.Address,
				SymbolTableEntry: entry,
			},
		)

		return nil
	}

	var loadedValue llvm.Value

	switch entry.DeclaredType.(type) {
	case ast.StringType:
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
	default:
		return g.declareVarWithInitializer(node)
	}
}

func NewLLVMGenerator(ctx *ast.CompilerCtx) *LLVMGenerator {
	return &LLVMGenerator{Ctx: ctx}
}

func (g *LLVMGenerator) VisitArrayAccessExpression(node *ast.ArrayAccessExpression) error {
	err, entry, array, indices := g.findArraySymbolTableEntry(node)
	if err != nil {
		return err
	}

	var arrayAddr llvm.Value
	if array.Address == nil {
		arrayAddr = array.Value
	} else {
		arrayAddr = *array.Address
	}

	itemPtr := g.Ctx.Builder.CreateInBoundsGEP(
		entry.UnderlyingType,
		arrayAddr,
		indices,
		"",
	)

	g.setLastResult(&ast.CompilerResult{
		Value:                 &itemPtr,
		ArraySymbolTableEntry: entry,
	})

	return nil
}

func (g *LLVMGenerator) findArraySymbolTableEntry(
	expr *ast.ArrayAccessExpression,
) (error, *ast.ArraySymbolTableEntry, *ast.SymbolTableEntry, []llvm.Value) {
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
	if err := node.Left.Accept(g); err != nil {
		return err
	}

	leftRes := g.getLastResult()

	if err := node.Right.Accept(g); err != nil {
		return err
	}
	rightRes := g.getLastResult()

	leftVal := g.extractRValue(leftRes)
	rightVal := g.extractRValue(rightRes)

	finalLeft, finalRight, err := g.coerceOperands(leftVal, rightVal)
	if err != nil {
		return fmt.Errorf("type mismatch at %v: %w", node.TokenStream(), err)
	}

	err, res := g.handleBinaryOp(node.Operator.Kind, finalLeft, finalRight)
	if err != nil {
		return err
	}

	g.setLastResult(res)

	return nil
}

// extractRValue ensures we are working with the actual value, not its memory address.
func (g *LLVMGenerator) extractRValue(res *ast.CompilerResult) llvm.Value {
	val := *res.Value

	if val.Type().TypeKind() != llvm.PointerTypeKind {
		return val
	}

	var loadType llvm.Type
	switch {
	case res.ArraySymbolTableEntry != nil:
		loadType = res.ArraySymbolTableEntry.UnderlyingType
	case res.StuctPropertyValueType != nil:
		loadType = *res.StuctPropertyValueType
	case !val.Type().ElementType().IsNil():
		loadType = val.Type().ElementType()
	default:
		// If it's a pointer to a pointer or a raw pointer we can't dereference, return as is
		return val
	}

	return g.Ctx.Builder.CreateLoad(loadType, val, "load.tmp")
}

// coerceOperands ensures both sides of the binary op have the same LLVM type.
func (g *LLVMGenerator) coerceOperands(left, right llvm.Value) (llvm.Value, llvm.Value, error) {
	lt := left.Type()
	rt := right.Type()

	if lt == rt {
		return left, right, nil
	}

	lk := lt.TypeKind()
	rk := rt.TypeKind()

	// Logic: Promote Integer to Float/Double if one side is a float
	if g.isFloat(lk) || g.isFloat(rk) {
		targetType := g.Ctx.Context.DoubleType()
		lCasted := g.castToFloat(left, targetType)
		rCasted := g.castToFloat(right, targetType)
		return lCasted, rCasted, nil
	}

	// Logic: Promote smaller integers to 32-bit (or your language's default int)
	if lk == llvm.IntegerTypeKind && rk == llvm.IntegerTypeKind {
		targetType := g.Ctx.Context.Int32Type()
		return g.Ctx.Builder.CreateIntCast(left, targetType, "l.ext"),
			g.Ctx.Builder.CreateIntCast(right, targetType, "r.ext"), nil
	}

	return left, right, fmt.Errorf("cannot coerce %s and %s", lk, rk)
}

func (g *LLVMGenerator) isFloat(k llvm.TypeKind) bool {
	return k == llvm.FloatTypeKind || k == llvm.DoubleTypeKind
}

func (g *LLVMGenerator) castToFloat(v llvm.Value, target llvm.Type) llvm.Value {
	if g.isFloat(v.Type().TypeKind()) {
		return g.Ctx.Builder.CreateFPCast(v, target, "fp.ext")
	}
	// Signed Integer to Float
	return g.Ctx.Builder.CreateSIToFP(v, target, "int.to.fp")
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

	if len(node.Body.Body) > 0 {
		// Create entry block
		entryBlock := g.Ctx.Context.AddBasicBlock(newFunc, "entry")
		g.Ctx.Builder.SetInsertPointAtEnd(entryBlock)

		for i, p := range newFunc.Params() {
			argType := node.Args[i].ArgType
			name := node.Args[i].Name
			p.SetName(name)

			// Create alloca for the parameter to support address taking and array indexing
			var entry ast.SymbolTableEntry

			if _, ok := argType.(ast.ArrayType); ok {
				// Array passed by reference, p is the pointer
				param := p
				entry = ast.SymbolTableEntry{
					Value:        param,
					DeclaredType: argType,
					Address:      &param,
				}
			} else {
				alloca := newCtx.Builder.CreateAlloca(p.Type(), name)
				newCtx.Builder.CreateStore(p, alloca)

				entry = ast.SymbolTableEntry{
					Value:        alloca,
					DeclaredType: argType,
					Address:      &alloca,
				}
			}

			err, eType := g.extractType(newCtx, argType)
			if err != nil {
				return err
			}

			if eType.sEntry != nil {
				entry.Ref = eType.sEntry
			}

			if eType.aEntry != nil {
				entry.Ref = eType.aEntry.UnderlyingTypeDef
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

		newGenerator := &LLVMGenerator{Ctx: newCtx}
		err := node.Body.Accept(newGenerator)
		if err != nil {
			return err
		}
	} else {
		newFunc.SetLinkage(llvm.ExternalLinkage)
	}

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

		// Pass arrays by reference (pointer)
		if _, ok := arg.ArgType.(ast.ArrayType); ok {
			ptrType := llvm.PointerType(typ.typ, 0)
			params = append(params, ptrType)
		} else {
			params = append(params, typ.typ)
		}
	}

	return nil, params
}

func (g *LLVMGenerator) setLastResult(res *ast.CompilerResult) {
	g.lastResult = res
}

func (g *LLVMGenerator) getLastResult() *ast.CompilerResult {
	res := g.lastResult
	g.lastResult = nil

	return res
}

func (g *LLVMGenerator) getProperty(expr *ast.MemberExpression) (string, error) {
	prop, ok := expr.Property.(ast.SymbolExpression)
	if !ok {
		return "", fmt.Errorf("struct property should be a symbol")
	}

	return prop.Value, nil
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

type InitializationStyle int

const (
	StyleDefault   InitializationStyle = iota // Alloca + Store
	StyleDirect                               // Use compiled value as address (Literals)
	StyleLoadStore                            // Load from pointer, then Alloca + Store (Accessors)
)

var nodeVariableDeclarationStyles = map[reflect.Type]InitializationStyle{
	reflect.TypeFor[ast.ArrayAccessExpression]():          StyleLoadStore,
	reflect.TypeFor[ast.ArrayOfStructsAccessExpression](): StyleLoadStore,
	reflect.TypeFor[ast.StringExpression]():               StyleDefault,
	reflect.TypeFor[ast.NumberExpression]():               StyleDefault,
	reflect.TypeFor[ast.FloatExpression]():                StyleDefault,
	reflect.TypeFor[ast.FunctionCallExpression]():         StyleDefault,
	reflect.TypeFor[ast.SymbolExpression]():               StyleDefault,
	reflect.TypeFor[ast.BinaryExpression]():               StyleDefault,
	reflect.TypeFor[ast.StructInitializationExpression](): StyleDirect,
	reflect.TypeFor[ast.ArrayInitializationExpression]():  StyleDirect,
}

func (g *LLVMGenerator) declareVarWithInitializer(node *ast.VarDeclarationStatement) error {
	if err := node.Value.Accept(g); err != nil {
		return err
	}

	res := g.getLastResult()
	style := nodeVariableDeclarationStyles[reflect.TypeOf(node.Value)]

	var finalAddr *llvm.Value

	switch style {
	case StyleDirect:
		finalAddr = res.Value

	case StyleLoadStore:
		// Logic for extracting value from an accessor
		typ := res.Value.AllocatedType()
		if _, ok := node.Value.(ast.ArrayOfStructsAccessExpression); ok {
			typ = *res.StuctPropertyValueType
		}
		loadedVal := g.Ctx.Builder.CreateLoad(typ, *res.Value, "tmp.load")

		alloc := g.Ctx.Builder.CreateAlloca(typ, node.Name)
		g.Ctx.Builder.CreateStore(loadedVal, alloc)
		finalAddr = &alloc

	default: // StyleDefault
		alloc := g.Ctx.Builder.CreateAlloca(res.Value.Type(), node.Name)
		g.Ctx.Builder.CreateStore(*res.Value, alloc)
		finalAddr = &alloc
	}

	// Unified Metadata Management
	return g.finalizeSymbol(node, finalAddr, res)
}

// Helper to keep metadata logic separate from IR generation logic
func (g *LLVMGenerator) finalizeSymbol(
	node *ast.VarDeclarationStatement,
	addr *llvm.Value,
	res *ast.CompilerResult,
) error {
	entry := &ast.SymbolTableEntry{
		Address:      addr,
		DeclaredType: node.ExplicitType,
		Ref:          res.StructSymbolTableEntry,
	}

	if res.ArraySymbolTableEntry != nil && res.ArraySymbolTableEntry.UnderlyingTypeDef != nil {
		entry.Ref = res.ArraySymbolTableEntry.UnderlyingTypeDef
	}

	if err := g.Ctx.AddSymbol(node.Name, entry); err != nil {
		return err
	}

	if _, ok := node.Value.(ast.ArrayInitializationExpression); ok {
		return g.Ctx.AddArraySymbol(node.Name, res.ArraySymbolTableEntry)
	}
	return nil
}

func (g *LLVMGenerator) findArrayOfStructsSymbolTableEntry(
	expr *ast.ArrayOfStructsAccessExpression,
) (error, *ast.SymbolTableEntry, *ast.ArraySymbolTableEntry, []llvm.Value) {
	var array *ast.SymbolTableEntry

	var entry *ast.ArraySymbolTableEntry

	switch expr.Name.(type) {
	case ast.SymbolExpression:
		varName, _ := expr.Name.(ast.SymbolExpression)

		err, arrayEntry := g.Ctx.FindSymbol(varName.Value)
		if err != nil {
			key := "ArrayOfStructsAccessExpression.NotFoundInSymbolTable"
			return g.Ctx.Dialect.Error(key, varName.Value), nil, nil, nil
		}

		array = arrayEntry

		err, ent := g.Ctx.FindArraySymbol(varName.Value)
		if err != nil {
			return err, nil, nil, nil
		}

		entry = ent

	default:
		return fmt.Errorf("ArrayOfStructsAccessExpression not implemented"), nil, nil, nil
	}

	var indices []llvm.Value

	switch expr.Index.(type) {
	case ast.NumberExpression:
		idx, _ := expr.Index.(ast.NumberExpression)

		indices = []llvm.Value{
			llvm.ConstInt((*g.Ctx.Context).Int32Type(), uint64(0), false),
			llvm.ConstInt((*g.Ctx.Context).Int32Type(), uint64(idx.Value), false),
		}
	case ast.SymbolExpression, ast.BinaryExpression, ast.MemberExpression:
		if err := expr.Index.Accept(g); err != nil {
			return err, nil, nil, nil
		}
		res := g.getLastResult()

		if res.Value.Type().TypeKind() == llvm.PointerTypeKind {
			load := g.Ctx.Builder.CreateLoad(res.Value.AllocatedType(), *res.Value, "")
			indices = []llvm.Value{
				llvm.ConstInt((*g.Ctx.Context).Int32Type(), uint64(0), false),
				load,
			}
		} else {
			indices = []llvm.Value{
				llvm.ConstInt((*g.Ctx.Context).Int32Type(), uint64(0), false),
				*res.Value,
			}
		}
	default:
		key := "ArrayOfStructsAccessExpression.AccessedIndexIsNotANumber"
		return g.Ctx.Dialect.Error(key, expr.Index), nil, nil, nil
	}

	return nil, array, entry, indices
}
