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
	Ctx                   *CompilerCtx
	lastResult            *CompilerResult
	lastTypeResult        *CompilerResultType
	zeroValues            map[reflect.Type]*CompilerResult
	logger                *Logger
	currentFuncReturnType ast.Type
}

func (g *LLVMGenerator) VisitBooleanExpression(node *ast.BooleanExpression) error {
	value := llvm.ConstInt(llvm.GlobalContext().Int1Type(), uint64(0), false)

	if node.Value {
		value = llvm.ConstInt(llvm.GlobalContext().Int1Type(), uint64(1), false)
	}

	g.setLastResult(&CompilerResult{
		Value: &value,
	})

	return nil
}

var _ ast.CodeGenerator = (*LLVMGenerator)(nil)

func NewLLVMGenerator(ctx *CompilerCtx) *LLVMGenerator {
	logger := NewLogger("LLVM")

	return &LLVMGenerator{
		Ctx:        ctx,
		logger:     logger,
		zeroValues: map[reflect.Type]*CompilerResult{},
	}
}

func (g *LLVMGenerator) Debugf(format string, args ...any) {
	if g.Ctx.Debugging {
		g.logger.Debug(fmt.Sprintf(format, args...))
	}
}

func (g *LLVMGenerator) VisitSymbolAdressExpression(node *ast.SymbolAdressExpression) error {
	return fmt.Errorf("TODO implement VisitSymbolAdressExpression")
}

func (g *LLVMGenerator) VisitSymbolValueExpression(node *ast.SymbolValueExpression) error {
	return fmt.Errorf("TODO implement VisitSymbolValueExpression")
}

func (g *LLVMGenerator) VisitFloatingBlockExpression(node *ast.FloatingBlockExpression) error {
	return node.Stmt.Accept(g)
}

func (g *LLVMGenerator) VisitZeroExpression(node *ast.ZeroExpression) error {
	err := node.T.AcceptZero(g)
	if err != nil {
		return err
	}

	res := g.getLastResult()
	res.SwaType = node.T
	g.setLastResult(res)

	return nil
}

// VisitAssignmentExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitAssignmentExpression(node *ast.AssignmentExpression) error {
	old := g.logger.Step("AssignmentExpr")

	defer g.logger.Restore(old)

	g.Debugf("%d", node.Tokens[0].Line)

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
	case *ast.ArrayAccessExpression:
		valueToBeAssigned = g.Ctx.Builder.CreateLoad(
			compiledValue.Value.AllocatedType(),
			*compiledValue.Value,
			"",
		)
	default:
		valueToBeAssigned = *compiledValue.Value
	}

	var address *llvm.Value

	g.Debugf("SymbolTableEntry of compiledAssignee %+v", compiledAssignee.SymbolTableEntry)

	if compiledAssignee.SymbolTableEntry != nil &&
		compiledAssignee.SymbolTableEntry.Address != nil {
		// TODO: figure out why we are doing this
		// TODO this triggers LLVM Verify Module Error
		// Store operand must be a pointer. store i32 %23, i32 %21, align 4
		// File: tests/./regression/rsa.swa
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
		fmt.Printf("VisitAssignmentExpression address to be assigned: %s\n", address)
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
	old := g.logger.Step("FloatExpr")

	defer g.logger.Restore(old)

	g.Debugf("%v", node.Value)

	res := llvm.ConstFloat(
		g.Ctx.Context.DoubleType(),
		node.Value,
	)
	node.SwaType = ast.FloatType{}

	g.setLastResult(&CompilerResult{Value: &res})

	return nil
}

// VisitMainStatement implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitMainStatement(node *ast.MainStatement) error {
	g.Ctx.InsideFunction = true
	g.Ctx.IncrementMainOccurrences()

	old := g.logger.Step("MainStmt")

	defer g.logger.Restore(old)

	defer func() { g.Ctx.InsideFunction = false }()

	fnType := llvm.FunctionType(
		g.Ctx.Context.Int32Type(),
		[]llvm.Type{},
		false,
	)
	fn := llvm.AddFunction(*g.Ctx.Module, "main", fnType)
	block := g.Ctx.Context.AddBasicBlock(fn, "entry")
	g.Ctx.Builder.SetInsertPointAtEnd(block)

	g.currentFuncReturnType = &ast.NumberType{}

	err := node.Body.Accept(g)
	if err != nil {
		return err
	}

	g.currentFuncReturnType = nil

	return nil
}

// VisitNumberExpression implements [ast.CodeGenerator].
func (g *LLVMGenerator) VisitNumberExpression(node *ast.NumberExpression) error {
	old := g.logger.Step("NumberExpr")

	defer g.logger.Restore(old)

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
	node.SwaType = ast.NumberType{}

	g.setLastResult(&CompilerResult{Value: &res, SwaType: node.SwaType})

	return nil
}

func (g *LLVMGenerator) NotImplemented(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func (g *LLVMGenerator) VisitReturnStatement(node *ast.ReturnStatement) error {
	old := g.logger.Step("ReturnStmt")

	defer g.logger.Restore(old)

	err := node.Value.Accept(g)
	if err != nil {
		return err
	}

	res := g.getLastResult()

	retVal, err := g.prepareReturnValue(node.Value, res)
	if err != nil {
		return err
	}

	retValType := node.Value.VisitedSwaType()

	if !retValType.Equals(g.currentFuncReturnType) {
		return fmt.Errorf(
			"expected return value of function to be %s, but got %s",
			g.currentFuncReturnType,
			retValType,
		)
	}

	g.Debugf(
		"Func return type %s, Func return type %s",
		g.currentFuncReturnType.Value().String(),
		retValType.Value().String())

	g.Ctx.Builder.CreateRet(retVal)

	return nil
}

func (g *LLVMGenerator) VisitStringExpression(node *ast.StringExpression) error {
	old := g.logger.Step("StringExpr")

	defer g.logger.Restore(old)

	if !g.Ctx.InsideFunction {
		value := llvm.ConstString(node.Value, true)
		g.setLastResult(&CompilerResult{Value: &value, SwaType: ast.StringType{}})

		return nil
	}

	valuePtr := g.Ctx.Builder.CreateGlobalStringPtr(node.Value, "")
	res := CompilerResult{Value: &valuePtr, SwaType: ast.StringType{}}

	node.SwaType = ast.StringType{}

	g.setLastResult(&res)

	return nil
}

func (g *LLVMGenerator) VisitStructDeclaration(node *ast.StructDeclarationStatement) error {
	old := g.logger.Step("StructDeclExpr")

	defer g.logger.Restore(old)

	g.Debugf(node.Name)

	if len(node.Properties) == 0 {
		key := "LLVMGenerator.VisitStructDeclaration.EmptyStruct"

		return g.Ctx.Dialect.Error(key, node.Name)
	}

	entry := &StructSymbolTableEntry{
		Metadata:    *node,
		Embeds:      map[string]StructSymbolTableEntry{},
		ArrayEmbeds: map[string]ArraySymbolTableEntry{},
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

		if datatype.Type.TypeKind() == llvm.ArrayTypeKind &&
			datatype.SubType.TypeKind() == llvm.StructTypeKind {
			g.Debugf("processing property %s as array of structs", propertyName)

			entry.ArrayEmbeds[propertyName] = *datatype.Aentry
		}

		if datatype.Type.TypeKind() == llvm.PointerTypeKind &&
			datatype.SubType.TypeKind() == llvm.StructTypeKind {
			if datatype.SubType.StructName() == node.Name {
				key := "VisitStructDeclaration.SelfPointerReferenceNotAllowed"

				return g.Ctx.Dialect.Error(key, propertyName)
			}

			g.Debugf("processing property %s as struct", propertyName)

			entry.Embeds[propertyName] = *datatype.Sentry
		}

		if datatype.Type.TypeKind() == llvm.StructTypeKind {
			if datatype.Type.StructName() == node.Name {
				key := "VisitStructDeclaration.SelfReferenceNotAllowed"

				return g.Ctx.Dialect.Error(key, propertyName)
			}

			g.Debugf("processing property %s as nested struct", propertyName)

			entry.Embeds[propertyName] = *datatype.Sentry
		}
	}

	newtype.StructSetBody(entry.PropertyTypes, false)

	err = g.Ctx.UpdateStructSymbol(node.Name, entry)
	if err != nil {
		return err
	}

	node.SwaType = ast.SymbolType{Name: node.Name}

	return nil
}

func (g *LLVMGenerator) VisitSymbolExpression(node *ast.SymbolExpression) error {
	old := g.logger.Step("SymbolExpr")

	defer g.logger.Restore(old)

	g.Debugf(node.Value)

	err, entry := g.Ctx.FindSymbol(node.Value)
	if err != nil {
		return err
	}

	if entry.Global {
		load := llvm.Value{}

		switch entry.DeclaredType.Value() {
		case ast.DataTypeFloat, ast.DataTypeNumber,
			ast.DataTypeNumber64, ast.DataTypeBool:
			err := entry.DeclaredType.Accept(g)
			if err != nil {
				return err
			}

			datatype := g.getLastTypeVisitResult()

			load = g.Ctx.Builder.CreateLoad(datatype.Type, *entry.Address, "")
		case ast.DataTypeString:
			load = *entry.Address
		default:
			key := "VisitSymbolExpression.UnsupportedTypeAsGlobal"

			return g.Ctx.Dialect.Error(key, entry.DeclaredType.Value().String())
		}

		g.setLastResult(
			&CompilerResult{
				Value:            &load,
				SymbolTableEntry: entry,
				SwaType:          entry.DeclaredType,
			},
		)

		node.SwaType = entry.DeclaredType

		return nil
	}

	if entry.Address == nil {
		g.setLastResult(
			&CompilerResult{
				Value:            &entry.Value,
				SymbolTableEntry: entry,
				SwaType:          entry.DeclaredType,
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
				SwaType:          entry.DeclaredType,
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

	node.SwaType = entry.DeclaredType

	g.setLastResult(
		&CompilerResult{
			Value:            &loadedValue,
			SymbolTableEntry: entry,
			SwaType:          entry.DeclaredType,
		},
	)

	return nil
}

func (g *LLVMGenerator) VisitPrefixExpression(node *ast.PrefixExpression) error {
	old := g.logger.Step("PrefixExpr")

	defer g.logger.Restore(old)

	err := node.RightExpression.Accept(g)
	if err != nil {
		return err
	}

	res := g.getLastResult()

	handler, ok := prefixOpHandlers[node.Operator.Kind]
	if !ok {
		key := "VisitPrefixExpression.OperatorNotSupported"

		return g.Ctx.Dialect.Error(key, node.Operator.Kind)
	}

	val := handler(g, *res.Value)
	node.SwaType = res.SwaType

	g.setLastResult(&CompilerResult{Value: &val})

	return nil
}

func (g *LLVMGenerator) getProperty(expr *ast.MemberExpression) (string, error) {
	prop, ok := expr.Property.(*ast.SymbolExpression)
	if !ok {
		key := "LLVMGenerator.getProperty.NotASymbol"

		return "", g.Ctx.Dialect.Error(key)
	}

	return prop.Value, nil
}

func (g *LLVMGenerator) prepareReturnValue(expr ast.Expression, res *CompilerResult) (llvm.Value, error) {
	switch expr.(type) {
	case *ast.ArrayAccessExpression:
		return g.Ctx.Builder.CreateLoad((*g.Ctx.Context).Int32Type(), *res.Value, ""), nil
	case *ast.MemberExpression:
		return g.Ctx.Builder.CreateLoad(*res.StuctPropertyValueType, *res.Value, ""), nil
	case *ast.StructInitializationExpression:
		return g.Ctx.Builder.CreateLoad(res.Value.AllocatedType(), *res.Value, ""), nil
	case *ast.StringExpression:
		alloc := g.Ctx.Builder.CreateAlloca(res.Value.Type(), "")
		g.Ctx.Builder.CreateStore(*res.Value, alloc)

		return alloc, nil
	case *ast.SymbolExpression, *ast.BinaryExpression, *ast.FunctionCallExpression,
		*ast.NumberExpression, *ast.FloatExpression, *ast.TupleExpression,
		*ast.ErrorExpression, *ast.ZeroExpression, *ast.BooleanExpression,
		*ast.PrefixExpression:
		return *res.Value, nil
	default:
		key := "VisitReturnStatement.UnsupportedExpression"

		return llvm.Value{}, g.Ctx.Dialect.Error(key, expr)
	}
}

// resolveGepIndices prepares the indices for a CreateGEP call.
// It ensures the first index is 0 (dereference) and the second is the evaluated index.
func (g *LLVMGenerator) resolveGepIndices(indexNode ast.Node) (error, []llvm.Value) {
	err := indexNode.Accept(g)
	if err != nil {
		return err, nil
	}

	res := g.getLastResult()
	if res == nil {
		key := "LLVMGenerator.resolveGepIndices.FailedToEvaluate"

		return g.Ctx.Dialect.Error(key), nil
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
