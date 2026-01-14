package compiler

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"swahili/lang/ast"

	"tinygo.org/x/go-llvm"
)

type CompilerResultType struct {
	Type    llvm.Type
	SubType llvm.Type
	Sentry  *StructSymbolTableEntry
	Aentry  *ArraySymbolTableEntry
	Entry   *SymbolTableEntry
}

type LLVMGenerator struct {
	Ctx            *CompilerCtx
	lastResult     *CompilerResult
	lastTypeResult *CompilerResultType
}

var _ ast.CodeGenerator = (*LLVMGenerator)(nil)

func (g *LLVMGenerator) VisitSymbolType(node *ast.SymbolType) error {
	err, entry := g.Ctx.FindStructSymbol(node.Name)
	if err != nil {
		return err
	}

	g.setLastTypeVisitResult(&CompilerResultType{
		Type:   entry.LLVMType,
		Sentry: entry,
	})

	return nil
}

func (g *LLVMGenerator) VisitNumberType(node *ast.NumberType) error {
	g.setLastTypeVisitResult(&CompilerResultType{
		Type: llvm.GlobalContext().Int32Type(),
	})

	return nil
}

func (g *LLVMGenerator) VisitNumber64Type(node *ast.Number64Type) error {
	g.setLastTypeVisitResult(&CompilerResultType{
		Type: llvm.GlobalContext().Int64Type(),
	})

	return nil
}

func (g *LLVMGenerator) VisitArrayType(node *ast.ArrayType) error {
	err := node.Underlying.Accept(g)
	if err != nil {
		return err
	}

	under := g.getLastTypeVisitResult()

	g.setLastTypeVisitResult(&CompilerResultType{
		Type:    llvm.ArrayType(under.Type, node.Size),
		SubType: under.Type,
	})

	return nil
}

func (g *LLVMGenerator) VisitFloatType(node *ast.FloatType) error {
	g.setLastTypeVisitResult(&CompilerResultType{
		Type: llvm.GlobalContext().DoubleType(),
	})

	return nil
}

func (g *LLVMGenerator) VisitPointerType(node *ast.PointerType) error {
	err := node.Underlying.Accept(g)
	if err != nil {
		return err
	}

	under := g.getLastTypeVisitResult()

	g.setLastTypeVisitResult(&CompilerResultType{
		Type:    llvm.PointerType(under.Type, 0),
		SubType: under.Type,
		Sentry:  under.Sentry,
	})

	return nil
}

func (g *LLVMGenerator) VisitStringType(node *ast.StringType) error {
	g.setLastTypeVisitResult(&CompilerResultType{
		Type:    llvm.PointerType(llvm.GlobalContext().Int8Type(), 0),
		SubType: llvm.GlobalContext().Int8Type(),
	})

	return nil
}

func (g *LLVMGenerator) VisitVoidType(node *ast.VoidType) error {
	g.setLastTypeVisitResult(&CompilerResultType{
		Type: llvm.GlobalContext().VoidType(),
	})

	return nil
}

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

	res := &CompilerResult{
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
	newCtx := NewCompilerContext(
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
	g.setLastResult(&CompilerResult{Value: &res})

	return nil
}

// VisitMainStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitMainStatement(node *ast.MainStatement) error {
	g.Ctx.InsideFunction = true

	defer func() { g.Ctx.InsideFunction = false }()

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
		result := &CompilerResult{
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
		result := &CompilerResult{
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

		if propIndex > len(result.SymbolTableEntry.Ref.PropertyTypes) ||
			len(result.SymbolTableEntry.Ref.PropertyTypes) == 0 {
			format := "Property named %s does not exist at index %d"

			return fmt.Errorf(format, prop.Value, propIndex)
		}

		propType := result.SymbolTableEntry.Ref.PropertyTypes[propIndex]

		finalresult := &CompilerResult{
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
	g.setLastResult(&CompilerResult{Value: &res})

	return nil
}

func (g *LLVMGenerator) NotImplemented(msg string) {
	fmt.Println(msg)
	os.Exit(1)
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
	if !g.Ctx.InsideFunction {
		value := llvm.ConstString(node.Value, true)
		g.setLastResult(&CompilerResult{Value: &value})

		return nil
	}

	valuePtr := g.Ctx.Builder.CreateGlobalStringPtr(node.Value, "")
	res := CompilerResult{Value: &valuePtr}

	g.setLastResult(&res)

	return nil
}

// VisitStructDeclaration implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitStructDeclaration(node *ast.StructDeclarationStatement) error {
	entry := &StructSymbolTableEntry{
		Metadata: *node,
		Embeds:   map[string]StructSymbolTableEntry{},
	}
	newtype := g.Ctx.Context.StructCreateNamed(node.Name)
	entry.LLVMType = newtype

	err := g.Ctx.AddStructSymbol(node.Name, entry)
	if err != nil {
		return err
	}

	for i, propertyName := range node.Properties {
		propertyType := node.Types[i]

		err := propertyType.Accept(g)
		if err != nil {
			return err
		}

		datatype := g.getLastTypeVisitResult()
		entry.PropertyTypes = append(entry.PropertyTypes, datatype.Type)

		if datatype.Type.TypeKind() == llvm.PointerTypeKind &&
			datatype.SubType.TypeKind() == llvm.StructTypeKind {
			if datatype.SubType.StructName() == node.Name {
				format := "struct with pointer reference to self not supported, property: %s"

				return fmt.Errorf(format, propertyName)
			}

			entry.Embeds[propertyName] = *datatype.Sentry
		}

		if datatype.Type.TypeKind() == llvm.StructTypeKind {
			if datatype.Type.StructName() == node.Name {
				format := "struct with reference to self not supported, property: %s"

				return fmt.Errorf(format, propertyName)
			}

			entry.Embeds[propertyName] = *datatype.Sentry
		}
	}

	newtype.StructSetBody(entry.PropertyTypes, false)

	err = g.Ctx.UpdateStructSymbol(node.Name, entry)
	if err != nil {
		return err
	}

	return nil
}

// VisitSymbolExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitSymbolExpression(node *ast.SymbolExpression) error {
	err, entry := g.Ctx.FindSymbol(node.Value)
	if err != nil {
		return err
	}

	if entry.Global {
		load := llvm.Value{}

		switch entry.DeclaredType.Value() {
		case ast.DataTypeFloat, ast.DataTypeNumber, ast.DataTypeNumber64:
			err := entry.DeclaredType.Accept(g)
			if err != nil {
				return err
			}

			datatype := g.getLastTypeVisitResult()

			load = g.Ctx.Builder.CreateLoad(datatype.Type, *entry.Address, "")
		case ast.DataTypeString:
			load = *entry.Address
		default:
			return fmt.Errorf("Unsupported datatype in global")
		}

		g.setLastResult(
			&CompilerResult{
				Value:            &load,
				SymbolTableEntry: entry,
			},
		)

		return nil
	}

	if entry.Address == nil {
		g.setLastResult(
			&CompilerResult{
				Value:            &entry.Value,
				SymbolTableEntry: entry,
			},
		)

		return nil
	}

	if entry.Address.IsAInstruction().IsNil() ||
		entry.Address.InstructionOpcode() != llvm.Alloca {
		g.setLastResult(
			&CompilerResult{
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
		&CompilerResult{
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
		err := node.ExplicitType.Accept(g)
		if err != nil {
			return err
		}

		typeresult := g.getLastTypeVisitResult()
		alloc := g.Ctx.Builder.CreateAlloca(typeresult.Type, fmt.Sprintf("alloc.%s", node.Name))
		g.Ctx.Builder.CreateStore(llvm.ConstNull(typeresult.Type), alloc)

		entry := &SymbolTableEntry{
			Value:        alloc,
			Address:      &alloc,
			DeclaredType: node.ExplicitType,
			Global:       !g.Ctx.InsideFunction,
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

func NewLLVMGenerator(ctx *CompilerCtx) *LLVMGenerator {
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

	g.setLastResult(&CompilerResult{
		Value:                 &itemPtr,
		ArraySymbolTableEntry: entry,
	})

	return nil
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
	g.setLastResult(&CompilerResult{Value: &val})

	return nil
}

func (g *LLVMGenerator) getProperty(expr *ast.MemberExpression) (string, error) {
	prop, ok := expr.Property.(ast.SymbolExpression)
	if !ok {
		return "", fmt.Errorf("struct property should be a symbol")
	}

	return prop.Value, nil
}

func (g *LLVMGenerator) prepareReturnValue(expr ast.Expression, res *CompilerResult) (llvm.Value, error) {
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
	err := node.Value.Accept(g)
	if err != nil {
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
		if !g.Ctx.InsideFunction {
			return fmt.Errorf("global var decl with %T not supported", node.Value)
		}

		finalAddr = res.Value

	case StyleLoadStore:
		if !g.Ctx.InsideFunction {
			return fmt.Errorf("global var decl with %T not supported", node.Value)
		}

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

	if _, ok := node.Value.(ast.ArrayInitializationExpression); ok {
		return g.Ctx.AddArraySymbol(node.Name, res.ArraySymbolTableEntry)
	}

	return nil
}

func (g *LLVMGenerator) findArrayOfStructsSymbolTableEntry(
	expr *ast.ArrayOfStructsAccessExpression,
) (error, *SymbolTableEntry, *ArraySymbolTableEntry, []llvm.Value) {
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

func (g *LLVMGenerator) resolveArrayOfStructsBase(nameNode ast.Node) (error, *SymbolTableEntry, *ArraySymbolTableEntry) {
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

		return nil, lastres.SymbolTableEntry, &ArraySymbolTableEntry{
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

func (g *LLVMGenerator) findArraySymbolTableEntry(
	expr *ast.ArrayAccessExpression,
) (error, *ArraySymbolTableEntry, *SymbolTableEntry, []llvm.Value) {
	var name string

	var array *SymbolTableEntry

	var entry *ArraySymbolTableEntry

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

		err = astType.Accept(g)
		if err != nil {
			return err, nil, nil, nil
		}

		etype := g.getLastTypeVisitResult()

		switch coltype := astType.(type) {
		case ast.PointerType:
			isPointerType = true
		case ast.ArrayType:
			elementsCount = coltype.Size
		default:
			err := fmt.Errorf("Property %s is not an array", propSym.Value)
			return err, nil, nil, nil
		}

		if isPointerType {
			pointerValue := g.Ctx.Builder.CreateLoad(*val.StuctPropertyValueType, *val.Value, "")
			array = &SymbolTableEntry{Value: pointerValue}
		} else {
			array = &SymbolTableEntry{Value: *val.Value}
		}

		entry = &ArraySymbolTableEntry{
			UnderlyingType: etype.SubType,
			Type:           etype.Type,
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
	structType *StructSymbolTableEntry,
	propName string,
) (int, error) {
	err, propIndex := structType.Metadata.PropertyIndex(propName)
	if err != nil {
		return 0, fmt.Errorf("struct %s has no field %s", structType.Metadata.Name, propName)
	}
	return propIndex, nil
}

func (g *LLVMGenerator) setLastResult(res *CompilerResult) {
	g.lastResult = res
}

func (g *LLVMGenerator) getLastResult() *CompilerResult {
	res := g.lastResult
	g.lastResult = nil

	return res
}

func (g *LLVMGenerator) setLastTypeVisitResult(res *CompilerResultType) {
	g.lastTypeResult = res
}

func (g *LLVMGenerator) getLastTypeVisitResult() *CompilerResultType {
	res := g.lastTypeResult
	g.lastTypeResult = nil

	return res
}
