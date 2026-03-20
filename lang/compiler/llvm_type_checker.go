package compiler

import (
	"fmt"
	"swahili/lang/ast"
	"swahili/lang/lexer"
)

type LLVMTypeChecker struct {
	ctx                   *CompilerCtx
	returnStatementsCount int
}

func (l *LLVMTypeChecker) VisitByteType(node *ast.ByteType) error {
	panic("unimplemented")
}

func (l *LLVMTypeChecker) ZeroOfByteType(node *ast.ByteType) error {
	panic("unimplemented")
}

var _ ast.CodeGenerator = (*LLVMTypeChecker)(nil)

func NewLLVMTypeChecker(ctx *CompilerCtx) *LLVMTypeChecker {
	return &LLVMTypeChecker{ctx: ctx}
}

func (l *LLVMTypeChecker) VisitAssignmentExpression(node *ast.AssignmentExpression) error {
	err := node.Assignee.Accept(l)
	if err != nil {
		return err
	}

	err = node.Value.Accept(l)
	if err != nil {
		return err
	}

	asstype := node.Assignee.VisitedSwaType()
	valType := node.Value.VisitedSwaType()

	if asstype == nil {
		return fmt.Errorf("LLVMTypeChecker VisitAssignmentExpression type of assignee is nil")
	}

	if valType == nil {
		return fmt.Errorf("LLVMTypeChecker VisitAssignmentExpression type of value is nil")
	}

	if !l.areTypesEquivalent(asstype, valType) {
		key := "LLVMTypeChecker.VisitAssignmentExpression.UnexpectedValue"
		return l.ctx.Dialect.Error(
			key,
			asstype.Value().String(),
			valType.Value().String(),
		)
	}

	return nil
}

func (l *LLVMTypeChecker) VisitBlockStatement(node *ast.BlockStatement) error {
	oldCtx := l.ctx
	l.ctx = NewCompilerContext(
		l.ctx.Context,
		l.ctx.Builder,
		l.ctx.Module,
		l.ctx.Dialect,
		l.ctx,
	)
	defer func() {
		l.ctx = oldCtx
	}()

	for _, v := range node.Body {
		err := v.Accept(l)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *LLVMTypeChecker) VisitMainStatement(node *ast.MainStatement) error {
	return node.Body.Accept(l)
}

func (l *LLVMTypeChecker) VisitConditionalStatement(node *ast.ConditionalStatetement) error {
	return nil
}

func (l *LLVMTypeChecker) VisitExpressionStatement(node *ast.ExpressionStatement) error {
	return node.Exp.Accept(l)
}

func (l *LLVMTypeChecker) VisitVarDeclaration(node *ast.VarDeclarationStatement) error {
	// Register symbol for subsequent expressions (always do this for declarations)
	err := l.ctx.AddSymbol(node.Name, &SymbolTableEntry{
		DeclaredType: node.ExplicitType,
	})
	if err != nil {
		return err
	}

	if node.Value == nil {
		return nil
	}

	// Make sure the value is visited to populate its type
	if node.Value.VisitedSwaType() == nil {
		err := node.Value.Accept(l)
		if err != nil {
			return err
		}
	}

	// Safety check again - if still nil after visiting, we can't check types
	if node.Value.VisitedSwaType() == nil {
		return nil
	}

	if !l.areTypesEquivalent(node.ExplicitType, node.Value.VisitedSwaType()) {
		key := "LLVMTypeChecker.VisitVarDeclaration.UnexpectedValue"
		return l.ctx.Dialect.Error(key,
			node.ExplicitType.Value().String(),
			node.Value.VisitedSwaType().Value().String())
	}

	return nil
}

func (l *LLVMTypeChecker) VisitStructDeclaration(node *ast.StructDeclarationStatement) error {
	return nil
}
func (l *LLVMTypeChecker) VisitStructInitializationExpression(node *ast.StructInitializationExpression) error {
	return nil
}

func (l *LLVMTypeChecker) VisitTupleAssignmentExpression(node *ast.TupleAssignmentExpression) error {
	valuestyp, ok := node.Value.VisitedSwaType().(*ast.TupleType)
	if !ok {
		return fmt.Errorf("Value is not a tuple")
	}

	assigneestyp, ok := node.Assignees.VisitedSwaType().(*ast.TupleType)
	if !ok {
		return fmt.Errorf("Assignee is not a tuple")
	}

	for i, typ := range valuestyp.Types {
		itemtyp := assigneestyp.Types[i]

		if !l.areTypesEquivalent(itemtyp, typ) {
			return fmt.Errorf(
				"Values mismatch in tuple expected %s got %s at index %d",
				typ.Value().String(), itemtyp.Value().String(), i)
		}
	}

	return nil
}
func (l *LLVMTypeChecker) VisitArrayAccessExpression(node *ast.ArrayAccessExpression) error {
	return nil
}
func (l *LLVMTypeChecker) VisitArrayInitializationExpression(node *ast.ArrayInitializationExpression) error {
	typ, _ := node.Underlying.(ast.ArrayType)

	for _, v := range node.Contents {
		if v.VisitedSwaType() != typ.Underlying {
			return fmt.Errorf("cannot insert %s in array of %s",
				v.VisitedSwaType().Value().String(),
				typ.Underlying.Value().String())
		}
	}

	return nil
}
func (l *LLVMTypeChecker) VisitArrayOfStructsAccessExpression(node *ast.ArrayOfStructsAccessExpression) error {
	return nil
}
func (l *LLVMTypeChecker) VisitArrayType(node *ast.ArrayType) error             { return nil }
func (l *LLVMTypeChecker) VisitSliceType(node *ast.SliceType) error             { return nil }
func (l *LLVMTypeChecker) VisitCallExpression(node *ast.CallExpression) error   { return nil }
func (l *LLVMTypeChecker) VisitErrorExpression(node *ast.ErrorExpression) error { return nil }
func (l *LLVMTypeChecker) VisitFloatExpression(node *ast.FloatExpression) error { return nil }
func (l *LLVMTypeChecker) VisitFloatType(node *ast.FloatType) error             { return nil }
func (l *LLVMTypeChecker) VisitFunctionCall(node *ast.FunctionCallExpression) error {
	symbol, ok := node.Name.(*ast.SymbolExpression)
	if !ok {
		return nil
	}

	if len(symbol.Tokens) > 0 {
		switch symbol.Tokens[0].Kind {
		case lexer.Make:
			typeExpr, ok := node.Args[0].(*ast.TypeExpression)
			if ok {
				node.SwaType = typeExpr.Type
			}
			return nil
		case lexer.Append:
			// append(slice, item) -> returns slice type
			err := node.Args[0].Accept(l)
			if err != nil {
				return err
			}
			node.SwaType = node.Args[0].VisitedSwaType()
			return nil
		case lexer.Len, lexer.Cap:
			node.SwaType = &ast.NumberType{}
			return nil
		}
	}

	name := symbol.Value
	err, funcType := l.ctx.FindFuncSymbol(name)
	if err != nil {
		return err
	}

	node.SwaType = funcType.meta.ReturnType

	for _, arg := range node.Args {
		err := arg.Accept(l)
		if err != nil {
			return err
		}
	}

	return nil
}
func (l *LLVMTypeChecker) VisitFunctionDefinition(node *ast.FuncDeclStatement) error {
	if node.Declaration {
		return nil
	}

	oldCtx := l.ctx
	l.ctx = NewCompilerContext(
		l.ctx.Context,
		l.ctx.Builder,
		l.ctx.Module,
		l.ctx.Dialect,
		l.ctx,
	)
	defer func() {
		l.ctx = oldCtx
	}()

	// Register parameters
	for _, p := range node.Args {
		err := l.ctx.AddSymbol(p.Name, &SymbolTableEntry{
			DeclaredType: p.ArgType,
		})
		if err != nil {
			return err
		}
	}

	l.returnStatementsCount = 0

	err := node.Body.Accept(l)
	if err != nil {
		return err
	}

	if l.returnStatementsCount == 0 && node.ReturnType.Value() != ast.DataTypeVoid {
		return fmt.Errorf("function %s has no return statement and expected %s", node.Name, node.ReturnType.Value().String())
	}

	l.returnStatementsCount = 0

	return nil
}
func (l *LLVMTypeChecker) VisitMemberExpression(node *ast.MemberExpression) error { return nil }
func (l *LLVMTypeChecker) VisitNumber64Type(node *ast.Number64Type) error         { return nil }
func (l *LLVMTypeChecker) VisitNumberExpression(node *ast.NumberExpression) error { return nil }
func (l *LLVMTypeChecker) VisitNumberType(node *ast.NumberType) error             { return nil }
func (l *LLVMTypeChecker) VisitPointerType(node *ast.PointerType) error           { return nil }
func (l *LLVMTypeChecker) VisitPrefixExpression(node *ast.PrefixExpression) error { return nil }
func (l *LLVMTypeChecker) VisitPrintStatement(node *ast.PrintStatetement) error   { return nil }
func (l *LLVMTypeChecker) VisitReturnStatement(node *ast.ReturnStatement) error {
	l.returnStatementsCount++

	return nil
}
func (l *LLVMTypeChecker) VisitStringExpression(node *ast.StringExpression) error { return nil }
func (l *LLVMTypeChecker) VisitStringType(node *ast.StringType) error             { return nil }
func (l *LLVMTypeChecker) VisitErrorType(node *ast.ErrorType) error               { return nil }
func (l *LLVMTypeChecker) VisitTupleType(node *ast.TupleType) error               { return nil }
func (l *LLVMTypeChecker) VisitSymbolExpression(node *ast.SymbolExpression) error {
	err, entry := l.ctx.FindSymbol(node.Value)
	if err != nil {
		return err
	}

	node.SwaType = entry.DeclaredType
	return nil
}
func (l *LLVMTypeChecker) VisitSymbolType(node *ast.SymbolType) error               { return nil }
func (l *LLVMTypeChecker) VisitTupleExpression(node *ast.TupleExpression) error     { return nil }
func (l *LLVMTypeChecker) VisitTypeExpression(node *ast.TypeExpression) error       { return nil }
func (l *LLVMTypeChecker) VisitVoidType(node *ast.VoidType) error                   { return nil }
func (l *LLVMTypeChecker) VisitWhileStatement(node *ast.WhileStatement) error       { return nil }
func (l *LLVMTypeChecker) ZeroOfArrayType(node *ast.ArrayType) error                { return nil }
func (l *LLVMTypeChecker) ZeroOfErrorType(node *ast.ErrorType) error                { return nil }
func (l *LLVMTypeChecker) ZeroOfFloatType(node *ast.FloatType) error                { return nil }
func (l *LLVMTypeChecker) ZeroOfNumber64Type(node *ast.Number64Type) error          { return nil }
func (l *LLVMTypeChecker) ZeroOfNumberType(node *ast.NumberType) error              { return nil }
func (l *LLVMTypeChecker) ZeroOfPointerType(node *ast.PointerType) error            { return nil }
func (l *LLVMTypeChecker) ZeroOfStringType(node *ast.StringType) error              { return nil }
func (l *LLVMTypeChecker) ZeroOfSymbolType(node *ast.SymbolType) error              { return nil }
func (l *LLVMTypeChecker) ZeroOfVoidType(node *ast.VoidType) error                  { return nil }
func (l *LLVMTypeChecker) ZeroOfTupleType(node *ast.TupleType) error                { return nil }
func (l *LLVMTypeChecker) ZeroOfSliceType(node *ast.SliceType) error                { return nil }
func (l *LLVMTypeChecker) VisitZeroExpression(node *ast.ZeroExpression) error       { return nil }
func (l *LLVMTypeChecker) VisitBinaryExpression(node *ast.BinaryExpression) error   { return nil }
func (l *LLVMTypeChecker) ZeroOfBoolType(node *ast.BoolType) error                  { return nil }
func (l *LLVMTypeChecker) VisitBoolType(node *ast.BoolType) error                   { return nil }
func (l *LLVMTypeChecker) VisitBooleanExpression(node *ast.BooleanExpression) error { return nil }

func (l *LLVMTypeChecker) VisitSymbolAdressExpression(node *ast.SymbolAdressExpression) error {
	return nil
}

func (l *LLVMTypeChecker) VisitSymbolValueExpression(node *ast.SymbolValueExpression) error {
	return nil
}

func (l *LLVMTypeChecker) VisitFloatingBlockExpression(node *ast.FloatingBlockExpression) error {
	return nil
}

func (l *LLVMTypeChecker) areTypesEquivalent(a, b ast.Type) bool {
	if a == nil || b == nil {
		return a == b
	}

	if a.Value() != b.Value() {
		return false
	}

	if a.Value() == ast.DataTypeSymbol {
		var nameA, nameB string
		if asym, ok := a.(*ast.SymbolType); ok {
			nameA = asym.Name
		} else if asym, ok := a.(ast.SymbolType); ok {
			nameA = asym.Name
		}

		if bsym, ok := b.(*ast.SymbolType); ok {
			nameB = bsym.Name
		} else if bsym, ok := b.(ast.SymbolType); ok {
			nameB = bsym.Name
		}
		return nameA == nameB
	}

	if a.Value() == ast.DataTypeSlice {
		var aslice, bslice ast.SliceType
		if st, ok := a.(ast.SliceType); ok {
			aslice = st
		} else if st, ok := a.(*ast.SliceType); ok {
			aslice = *st
		}

		if st, ok := b.(ast.SliceType); ok {
			bslice = st
		} else if st, ok := b.(*ast.SliceType); ok {
			bslice = *st
		}
		return l.areTypesEquivalent(aslice.Underlying, bslice.Underlying)
	}

	// For other types like ErrorType, NumberType etc, Value() equality is enough
	return true
}
