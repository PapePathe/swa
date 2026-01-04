package compiler

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

type LLVMGenerator struct {
	Ctx        *ast.CompilerCtx
	lastResult *ast.CompilerResult
}

var _ ast.CodeGenerator = (*LLVMGenerator)(nil)

func (g *LLVMGenerator) VisitArrayOfStructsAccessExpression(node *ast.ArrayOfStructsAccessExpression) error {
	if g.Ctx.Debugging {
		fmt.Printf("VisitArrayOfStructsAccessExpression %s[%s].%s\n", node.Name, node.Index, node.Property)
	}

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

	proptype := array.Ref.PropertyTypes[propIndex]

	if proptype.TypeKind() == llvm.StructTypeKind {
		nestedstruct, _ := array.Ref.Embeds[propName.Value]
		array.Ref = &nestedstruct
	}

	res := &ast.CompilerResult{
		Value:                  &structPtr,
		SymbolTableEntry:       array,
		StuctPropertyValueType: &proptype,
		ArraySymbolTableEntry:  arrayEntry,
	}

	g.setLastResult(res)

	return nil
}

// VisitAssignmentExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitAssignmentExpression(node *ast.AssignmentExpression) error {
	if g.Ctx.Debugging {
		fmt.Printf("VisitAssignmentExpression %s\n", node)
	}

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

	var address *llvm.Value
	if compiledAssignee.SymbolTableEntry != nil &&
		compiledAssignee.SymbolTableEntry.Address != nil {
		// TODO: figure out why we are doing this
		if compiledAssignee.SymbolTableEntry.Ref != nil {
			address = compiledAssignee.Value
		} else {
			address = compiledAssignee.SymbolTableEntry.Address
		}

	} else {
		address = compiledAssignee.Value
	}

	g.Ctx.Builder.CreateStore(valueToBeAssigned, *address)

	if g.Ctx.Debugging {
		fmt.Printf("VisitAssignmentExpression assignee: %s\n", compiledAssignee.Value.String())
		fmt.Printf("VisitAssignmentExpression value: %s\n", valueToBeAssigned.String())
	}

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
	if g.Ctx.Debugging {
		fmt.Printf("VisitMemberExpression %s.%s\n", node.Object, node.Property)
	}

	switch objectType := node.Object.(type) {
	case ast.SymbolExpression:
		err, varDef := g.Ctx.FindSymbol(objectType.Value)
		if err != nil {
			return fmt.Errorf("variable %s is not defined", objectType.Value)
		}

		if varDef.Ref == nil {
			return fmt.Errorf("variable %s is not a struct instance", objectType.Value)
		}

		propName, err := g.getProperty(node)
		if err != nil {
			return err
		}

		propIndex, err := g.resolveStructAccess(varDef.Ref, propName)
		if err != nil {
			return err
		}

		var baseValue llvm.Value = varDef.Value
		if varDef.Address != nil {
			baseValue = *varDef.Address
		}

		addr := g.Ctx.Builder.CreateStructGEP(varDef.Ref.LLVMType, baseValue, propIndex, "")
		propType := varDef.Ref.PropertyTypes[propIndex]
		result := &ast.CompilerResult{
			Value:                  &addr,
			SymbolTableEntry:       varDef,
			StuctPropertyValueType: &propType,
		}

		if propType.TypeKind() == llvm.StructTypeKind {
			prop, _ := varDef.Ref.Embeds[propName]
			result.SymbolTableEntry.Ref = &prop
		}

		g.setLastResult(result)

		return nil
	case ast.ArrayOfStructsAccessExpression:
		err := node.Object.Accept(g)
		if err != nil {
			return err
		}

		lastresult := g.getLastResult()

		propName, err := g.getProperty(node)
		if err != nil {
			return err
		}

		propIndex, err := g.resolveStructAccess(lastresult.SymbolTableEntry.Ref, propName)
		if err != nil {
			return err
		}

		propType := lastresult.SymbolTableEntry.Ref.PropertyTypes[propIndex]
		addr := g.Ctx.Builder.CreateStructGEP(
			lastresult.SymbolTableEntry.Ref.LLVMType,
			*lastresult.Value,
			propIndex,
			"")
		result := &ast.CompilerResult{
			Value:                  &addr,
			SymbolTableEntry:       lastresult.SymbolTableEntry,
			StuctPropertyValueType: &propType,
		}

		if propType.TypeKind() == llvm.StructTypeKind {
			// TODO handle the case where propName is not in embeds
			prop, _ := lastresult.SymbolTableEntry.Ref.Embeds[propName]
			result.SymbolTableEntry.Ref = &prop
		}

		g.setLastResult(result)

		return nil
	case ast.MemberExpression:
		err := node.Object.Accept(g)
		if err != nil {
			return err
		}

		prop, _ := node.Property.(ast.SymbolExpression)
		result := g.getLastResult()

		propIndex, err := g.resolveStructAccess(result.SymbolTableEntry.Ref, prop.Value)
		if err != nil {
			return err
		}

		addr := g.Ctx.Builder.CreateStructGEP(*result.StuctPropertyValueType, *result.Value, propIndex, "")
		propType := result.SymbolTableEntry.Ref.PropertyTypes[propIndex]
		finalresult := &ast.CompilerResult{
			Value:                  &addr,
			SymbolTableEntry:       result.SymbolTableEntry,
			StuctPropertyValueType: &propType,
		}

		if propType.TypeKind() == llvm.StructTypeKind {
			sEntry, ok := result.SymbolTableEntry.Ref.Embeds[prop.Value]
			if ok {
				result.SymbolTableEntry.Ref = &sEntry
			}
		}

		g.setLastResult(finalresult)

		return nil

	default:
		g.NotImplemented(fmt.Sprintf("VisitMemberExpression not implemented for %T", objectType))
	}

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
	reflect.TypeFor[ast.FunctionCallExpression](): extractDirect,
	reflect.TypeFor[ast.StringExpression]():       extractDirect,
	reflect.TypeFor[ast.NumberExpression]():       extractDirect,
	reflect.TypeFor[ast.FloatExpression]():        extractDirect,
	reflect.TypeFor[ast.BinaryExpression]():       extractDirect,
	reflect.TypeFor[ast.SymbolExpression]():       extractDirect,
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
	entry := &ast.StructSymbolTableEntry{
		Metadata: *node,
		Embeds:   map[string]ast.StructSymbolTableEntry{},
	}

	for i, propertyName := range node.Properties {
		propertyType := node.Types[i]

		err, llvmType := propertyType.LLVMType(g.Ctx)
		if err != nil {
			return err
		}

		if llvmType.TypeKind() == llvm.StructTypeKind {
			structType, _ := propertyType.(ast.SymbolType)
			err, sEntry := g.Ctx.FindStructSymbol(structType.Name)
			if err != nil {
				return err
			}

			entry.Embeds[propertyName] = *sEntry
		}

		entry.PropertyTypes = append(entry.PropertyTypes, llvmType)
	}

	newtype := g.Ctx.Context.StructCreateNamed(node.Name)
	entry.LLVMType = newtype
	newtype.StructSetBody(entry.PropertyTypes, false)

	return g.Ctx.AddStructSymbol(node.Name, entry)
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

func (g *LLVMGenerator) VisitFunctionDefinition(node *ast.FuncDeclStatement) error {
	oldCtx := g.Ctx
	g.Ctx = ast.NewCompilerContext(
		g.Ctx.Context,
		g.Ctx.Builder,
		g.Ctx.Module,
		g.Ctx.Dialect,
		g.Ctx,
	)

	err, params := g.funcParams(g.Ctx, node)
	if err != nil {
		return err
	}

	err, returnType := g.extractType(g.Ctx, node.ReturnType)
	if err != nil {
		return err
	}

	newfuncType := llvm.FunctionType(returnType.typ, params, node.ArgsVariadic)
	newFunc := llvm.AddFunction(*g.Ctx.Module, node.Name, newfuncType)

	err = oldCtx.AddFuncSymbol(node.Name, &newfuncType)
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

			switch argType.(type) {
			case ast.ArrayType:
				param := p
				entry = ast.SymbolTableEntry{
					Value:        param,
					DeclaredType: argType,
					Address:      &param,
				}
			case ast.PointerType:
				param := p
				entry = ast.SymbolTableEntry{
					Value:        param,
					DeclaredType: argType,
					Address:      &param,
				}
			case ast.SymbolType:
				entry = ast.SymbolTableEntry{
					Value:        p,
					DeclaredType: argType,
					Address:      &p,
				}
			case ast.FloatType, ast.NumberType, ast.Number64Type:
				entry = ast.SymbolTableEntry{
					Value:        p,
					DeclaredType: argType,
					Address:      &p,
				}
			default:
				alloca := g.Ctx.Builder.CreateAlloca(p.Type(), name)
				g.Ctx.Builder.CreateStore(p, alloca)

				entry = ast.SymbolTableEntry{
					Value:        alloca,
					DeclaredType: argType,
					Address:      &alloca,
				}
			}

			err, eType := g.extractType(g.Ctx, argType)
			if err != nil {
				return err
			}

			if eType.sEntry != nil {
				entry.Ref = eType.sEntry
			}

			if eType.aEntry != nil {
				entry.Ref = eType.aEntry.UnderlyingTypeDef

				err := g.Ctx.AddArraySymbol(name, eType.aEntry)
				if err != nil {
					return fmt.Errorf("failed to add parameter %s to arrays symbol table: %w", name, err)
				}
			}

			err = g.Ctx.AddSymbol(name, &entry)
			if err != nil {
				return fmt.Errorf("failed to add parameter %s to symbol table: %w", name, err)
			}

		}

		err := node.Body.Accept(g)
		if err != nil {
			return err
		}
	} else {
		newFunc.SetLinkage(llvm.ExternalLinkage)
	}

	g.Ctx = oldCtx

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
	reflect.TypeFor[ast.MemberExpression]():               StyleLoadStore,
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
	style, ok := nodeVariableDeclarationStyles[reflect.TypeOf(node.Value)]
	if !ok {
		return fmt.Errorf("var decl with %s not supported", node.Value)
	}

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

	if res.SymbolTableEntry != nil && res.SymbolTableEntry.Ref != nil {
		entry.Ref = res.SymbolTableEntry.Ref
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
	err, baseSym, arrayMeta := g.resolveArrayOfStructsBase(expr.Name)
	if err != nil {
		return err, nil, nil, nil
	}

	err, indices := g.resolveGepIndices(expr.Index)
	if err != nil {
		return err, nil, nil, nil
	}

	return nil, baseSym, arrayMeta, indices
}

func (g *LLVMGenerator) resolveArrayOfStructsBase(nameNode ast.Node) (error, *ast.SymbolTableEntry, *ast.ArraySymbolTableEntry) {
	switch typednode := nameNode.(type) {
	case ast.MemberExpression:
		err := typednode.Object.Accept(g)
		if err != nil {
			return err, nil, nil
		}

		lastres := g.getLastResult()

		propName, err := g.getProperty(&typednode)
		if err != nil {
			return err, nil, nil
		}

		propIndex, err := g.resolveStructAccess(lastres.SymbolTableEntry.Ref, propName)
		if err != nil {
			return err, nil, nil
		}

		propType := lastres.SymbolTableEntry.Ref.PropertyTypes[propIndex]
		underlyingTypeDef := lastres.SymbolTableEntry.Ref

		if propType.TypeKind() == llvm.ArrayTypeKind {
			if propType.ElementType().TypeKind() == llvm.StructTypeKind {
				proptype := lastres.SymbolTableEntry.Ref.Metadata.Types[propIndex]
				arrtype, _ := proptype.(ast.ArrayType)
				stype, _ := arrtype.Underlying.(ast.SymbolType)

				err, sym := g.Ctx.FindStructSymbol(stype.Name)
				if err != nil {
					return err, nil, nil
				}

				underlyingTypeDef = sym
				lastres.SymbolTableEntry.Ref = sym
			}
		}

		return nil, lastres.SymbolTableEntry, &ast.ArraySymbolTableEntry{
			ElementsCount:     propType.ArrayLength(),
			UnderlyingType:    propType.ElementType(),
			Type:              propType,
			UnderlyingTypeDef: underlyingTypeDef,
		}
	case ast.SymbolExpression:
		name := typednode.Value

		err, symEntry := g.Ctx.FindSymbol(name)
		if err != nil {
			return g.Ctx.Dialect.Error("ArrayOfStructsAccessExpression.NotFoundInSymbolTable", name), nil, nil
		}

		err, arrEntry := g.Ctx.FindArraySymbol(name)
		if err != nil {
			return err, nil, nil
		}

		return nil, symEntry, arrEntry
	default:
		return fmt.Errorf("ArrayOfStructsAccessExpression not implemented for %T", typednode), nil, nil
	}
}

// resolveGepIndices prepares the indices for a CreateGEP call.
// It ensures the first index is 0 (dereference) and the second is the evaluated index.
func (g *LLVMGenerator) resolveGepIndices(indexNode ast.Node) (error, []llvm.Value) {
	if err := indexNode.Accept(g); err != nil {
		return err, nil
	}

	res := g.getLastResult()
	if res == nil {
		return fmt.Errorf("failed to evaluate index expression"), nil
	}

	idxVal := *res.Value
	if idxVal.Type().TypeKind() == llvm.PointerTypeKind {
		idxVal = g.Ctx.Builder.CreateLoad(idxVal.AllocatedType(), idxVal, "idx.load")
	}

	i32 := g.Ctx.Context.Int32Type()
	indices := []llvm.Value{
		llvm.ConstInt(i32, 0, false), // Step into the array pointer
		idxVal,                       // The actual offset
	}

	return nil, indices
}
