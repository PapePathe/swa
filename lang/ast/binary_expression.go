package ast

import (
	"encoding/json"
	"fmt"

	"tinygo.org/x/go-llvm"

	"swahili/lang/lexer"
)

// BinaryExpression ...
type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator lexer.Token
	Tokens   []lexer.Token
}

var _ Expression = (*BinaryExpression)(nil)

func (expr BinaryExpression) CompileLLVM(ctx *CompilerCtx) (error, *CompilerResult) {
	err, compiledLeftValue, compiledRightValue := expr.compileLeftAndRight(ctx)
	if err != nil {
		return err, nil
	}

	var finalLeftValue, finalRightValue llvm.Value

	// Load values from pointers first
	var leftValue, rightValue llvm.Value

	if compiledLeftValue.Type().TypeKind() == llvm.PointerTypeKind {
		// Get the element type of the pointer
		elementType := compiledLeftValue.Type().ElementType()
		leftValue = ctx.Builder.CreateLoad(elementType, *compiledLeftValue, "")
	} else {
		leftValue = *compiledLeftValue
	}

	if compiledRightValue.Type().TypeKind() == llvm.PointerTypeKind {
		// Get the element type of the pointer
		elementType := compiledRightValue.Type().ElementType()
		rightValue = ctx.Builder.CreateLoad(elementType, *compiledRightValue, "")
	} else {
		rightValue = *compiledRightValue
	}

	// Now determine the common type and cast both values
	ctype := commonType(leftValue, rightValue)
	finalLeftValue = expr.castToType(ctx, ctype, leftValue)
	finalRightValue = expr.castToType(ctx, ctype, rightValue)

	handler, ok := handlers[expr.Operator.Kind]
	if !ok {
		return fmt.Errorf("unsupported operator <%s>", expr.Operator.Kind), nil
	}
	return handler(ctx, finalLeftValue, finalRightValue)
}

func (expr BinaryExpression) TokenStream() []lexer.Token {
	return expr.Tokens
}

func (expr BinaryExpression) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	m["Left"] = expr.Left
	m["Right"] = expr.Right
	m["Operator"] = expr.Operator

	res := make(map[string]any)
	res["ast.BinaryExpression"] = m

	return json.Marshal(res)
}

func (expr BinaryExpression) compileLeftAndRight(ctx *CompilerCtx) (error, *llvm.Value, *llvm.Value) {
	err, compiledLeftValue := expr.Left.CompileLLVM(ctx)
	if err != nil {
		return err, nil, nil
	}

	if compiledLeftValue == nil {
		return fmt.Errorf("left side of expression is nil"), nil, nil
	}

	if compiledLeftValue.Value.Type().TypeKind() == llvm.VoidTypeKind {
		return fmt.Errorf("left side of expression is of type void"), nil, nil
	}

	err, compiledRightValue := expr.Right.CompileLLVM(ctx)
	if err != nil {
		return err, nil, nil
	}

	if compiledRightValue == nil {
		return fmt.Errorf("right side of expression is nil"), nil, nil
	}

	if compiledRightValue.Value.Type().TypeKind() == llvm.VoidTypeKind {
		return fmt.Errorf("right side of expression is of type void"), nil, nil
	}
	return nil, compiledLeftValue.Value, compiledRightValue.Value
}

func commonType(l, r llvm.Value) llvm.Type {
	// Both pointers
	if l.Type().TypeKind() == llvm.PointerTypeKind && r.Type().TypeKind() == llvm.PointerTypeKind {
		return l.GlobalValueType()
	}

	// Same type
	if l.Type() == r.Type() {
		return l.Type()
	}

	// Left is pointer
	if l.Type().TypeKind() == llvm.PointerTypeKind {
		return l.GlobalValueType()
	}

	// Right is pointer
	if r.Type().TypeKind() == llvm.PointerTypeKind {
		return r.GlobalValueType()
	}

	// Both integers - return the common integer type
	if l.Type().TypeKind() == llvm.IntegerTypeKind && r.Type().TypeKind() == llvm.IntegerTypeKind {
		return llvm.GlobalContext().Int32Type()
	}

	// Handle int vs float: promote to float
	if (l.Type().TypeKind() == llvm.IntegerTypeKind && (r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind)) ||
		((l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind) && r.Type().TypeKind() == llvm.IntegerTypeKind) {
		// Promote to double (float type)
		return llvm.GlobalContext().DoubleType()
	}

	// Both floats
	if (l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind) &&
		(r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind) {
		return llvm.GlobalContext().DoubleType()
	}

	panic(fmt.Errorf("Unhandled type combination: left=%s, right=%s", l.Type().TypeKind(), r.Type().TypeKind()))
}

func (expr BinaryExpression) castToType(ctx *CompilerCtx, t llvm.Type, v llvm.Value) llvm.Value {
	if t.TypeKind() == v.Type().TypeKind() {
		return v
	}

	switch v.Type().TypeKind() {
	case llvm.PointerTypeKind:
		switch t.TypeKind() {
		case llvm.IntegerTypeKind:
			return ctx.Builder.CreatePtrToInt(v, t, "")
		case llvm.DoubleTypeKind, llvm.FloatTypeKind:
			// This shouldn't happen in normal code, but return the value as-is
			return v
		default:
			panic(fmt.Errorf("Cannot cast pointer to %s", t.TypeKind()))
		}
	case llvm.IntegerTypeKind:
		switch t.TypeKind() {
		case llvm.DoubleTypeKind, llvm.FloatTypeKind:
			// Convert int to float
			return ctx.Builder.CreateSIToFP(v, t, "")
		case llvm.IntegerTypeKind:
			// Same type, just return
			return v
		default:
			panic(fmt.Errorf("Cannot cast integer to %s", t.TypeKind()))
		}
	case llvm.DoubleTypeKind, llvm.FloatTypeKind:
		switch t.TypeKind() {
		case llvm.IntegerTypeKind:
			// Convert float to int
			return ctx.Builder.CreateFPToSI(v, t, "")
		case llvm.DoubleTypeKind, llvm.FloatTypeKind:
			// Both floats, return as-is or cast between float/double
			return v
		default:
			panic(fmt.Errorf("Cannot cast float to %s", t.TypeKind()))
		}
	default:
		panic(fmt.Errorf("Unhandled type conversion from %s to %s", v.Type().TypeKind(), t.TypeKind()))
	}
}

type binaryHandlerFunc func(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult)

var handlers = map[lexer.TokenKind]binaryHandlerFunc{
	lexer.Plus:              add,
	lexer.Minus:             substract,
	lexer.Star:              multiply,
	lexer.Divide:            divide,
	lexer.GreaterThan:       greaterThan,
	lexer.GreaterThanEquals: greaterThanEquals,
	lexer.LessThan:          lessThan,
	lexer.LessThanEquals:    lessThanEquals,
	lexer.Equals:            equals,
}

func add(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFAdd(l, r, "")
	} else {
		res = ctx.Builder.CreateAdd(l, r, "")
	}

	return nil, &CompilerResult{Value: &res}
}

func divide(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	res := ctx.Builder.CreateSDiv(l, r, "")

	return nil, &CompilerResult{Value: &res}
}

func multiply(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	res := ctx.Builder.CreateMul(l, r, "")

	return nil, &CompilerResult{Value: &res}
}

func substract(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	res := ctx.Builder.CreateSub(l, r, "")

	return nil, &CompilerResult{Value: &res}
}

func greaterThan(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFCmp(llvm.FloatOGT, l, r, "")
	} else {
		res = ctx.Builder.CreateICmp(llvm.IntUGT, l, r, "")
	}

	return nil, &CompilerResult{Value: &res}
}

func greaterThanEquals(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFCmp(llvm.FloatOGE, l, r, "")
	} else {
		res = ctx.Builder.CreateICmp(llvm.IntUGE, l, r, "")
	}

	return nil, &CompilerResult{Value: &res}
}

func lessThan(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFCmp(llvm.FloatOLT, l, r, "")
	} else {
		res = ctx.Builder.CreateICmp(llvm.IntULT, l, r, "")
	}

	return nil, &CompilerResult{Value: &res}
}

func lessThanEquals(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFCmp(llvm.FloatOLE, l, r, "")
	} else {
		res = ctx.Builder.CreateICmp(llvm.IntULE, l, r, "")
	}

	return nil, &CompilerResult{Value: &res}
}

func equals(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	// Check if either side is a float type
	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFCmp(llvm.FloatOEQ, l, r, "")
	} else {
		res = ctx.Builder.CreateICmp(llvm.IntEQ, l, r, "")
	}

	return nil, &CompilerResult{Value: &res}
}
