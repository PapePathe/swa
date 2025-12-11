package ast

import (
	"encoding/json"
	"fmt"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
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
	err, leftResult, rightResult := expr.compileLeftAndRightResult(ctx)
	if err != nil {
		return err, nil
	}

	var finalLeftValue, finalRightValue llvm.Value

	// Load values from pointers (e.g., from array access)
	var leftValue, rightValue llvm.Value

	// Left side
	if leftResult.Value.Type().TypeKind() == llvm.PointerTypeKind {
		if leftResult.ArraySymbolTableEntry != nil {
			leftValue = ctx.Builder.CreateLoad(leftResult.ArraySymbolTableEntry.UnderlyingType, *leftResult.Value, "")
		} else {
			if leftResult.StuctPropertyValueType != nil {
				elementType := leftResult.StuctPropertyValueType
				leftValue = ctx.Builder.CreateLoad(*elementType, *leftResult.Value, "")
			} else {
				if leftResult.Value.Type().ElementType().IsNil() {
					leftValue = *leftResult.Value
				} else {
					elementType := leftResult.Value.Type().ElementType()
					leftValue = ctx.Builder.CreateLoad(elementType, *leftResult.Value, "")
				}
			}
		}
	} else {
		leftValue = *leftResult.Value
	}

	// Right side
	if rightResult.Value.Type().TypeKind() == llvm.PointerTypeKind {
		if rightResult.ArraySymbolTableEntry != nil {
			rightValue = ctx.Builder.CreateLoad(rightResult.ArraySymbolTableEntry.UnderlyingType, *rightResult.Value, "")
		} else {
			if rightResult.StuctPropertyValueType != nil {
				elementType := rightResult.StuctPropertyValueType
				rightValue = ctx.Builder.CreateLoad(*elementType, *rightResult.Value, "")
			} else {
				if rightResult.Value.Type().ElementType().IsNil() {
					rightValue = *rightResult.Value
				} else {
					elementType := rightResult.Value.Type().ElementType()
					rightValue = ctx.Builder.CreateLoad(elementType, *rightResult.Value, "")
				}
			}
		}
	} else {
		rightValue = *rightResult.Value
	}

	// Determine the common type
	ctype := commonType(leftValue, rightValue)

	// Cast only if necessary
	if leftValue.Type() == ctype {
		finalLeftValue = leftValue
	} else {
		err, finalLeftValue = expr.castToType(ctx, ctype, leftValue)
		if err != nil {
			return err, nil
		}
	}

	if rightValue.Type() == ctype {
		finalRightValue = rightValue
	} else {
		err, finalRightValue = expr.castToType(ctx, ctype, rightValue)
		if err != nil {
			return err, nil
		}
	}

	handler, ok := handlers[expr.Operator.Kind]
	if !ok {
		return fmt.Errorf("Binary expressions : unsupported operator <%s>", expr.Operator.Kind), nil
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

func (expr BinaryExpression) compileLeftAndRightResult(ctx *CompilerCtx) (error, *CompilerResult, *CompilerResult) {
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

	switch expr.Left.(type) {
	case StringExpression:
		glob := llvm.AddGlobal(*ctx.Module, compiledLeftValue.Value.Type(), "")
		glob.SetInitializer(*compiledLeftValue.Value)
		compiledLeftValue.Value = &glob
	default:
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

	switch expr.Right.(type) {
	case StringExpression:
		glob := llvm.AddGlobal(*ctx.Module, compiledRightValue.Value.Type(), "")
		glob.SetInitializer(*compiledRightValue.Value)
		compiledRightValue.Value = &glob
	default:
	}

	return nil, compiledLeftValue, compiledRightValue
}

func commonType(l, r llvm.Value) llvm.Type {
	lKind := l.Type().TypeKind()
	rKind := r.Type().TypeKind()

	// Both pointers - use int type as fallback (shouldn't happen in binary ops)
	if lKind == llvm.PointerTypeKind && rKind == llvm.PointerTypeKind {
		return llvm.GlobalContext().Int32Type()
	}

	// Same type
	if l.Type() == r.Type() {
		return l.Type()
	}

	// Left is pointer, use right type
	if lKind == llvm.PointerTypeKind {
		return r.Type()
	}

	// Right is pointer, use left type
	if rKind == llvm.PointerTypeKind {
		return l.Type()
	}

	// Both integers - return the common integer type
	if lKind == llvm.IntegerTypeKind && rKind == llvm.IntegerTypeKind {
		return llvm.GlobalContext().Int32Type()
	}

	// Handle int vs float: promote to float
	if (lKind == llvm.IntegerTypeKind && (rKind == llvm.FloatTypeKind || rKind == llvm.DoubleTypeKind)) ||
		((lKind == llvm.FloatTypeKind || lKind == llvm.DoubleTypeKind) && rKind == llvm.IntegerTypeKind) {
		// Promote to double (float type)
		return llvm.GlobalContext().DoubleType()
	}

	// Both floats
	if (lKind == llvm.FloatTypeKind || lKind == llvm.DoubleTypeKind) &&
		(rKind == llvm.FloatTypeKind || rKind == llvm.DoubleTypeKind) {
		return llvm.GlobalContext().DoubleType()
	}

	panic(fmt.Errorf("Unhandled type combination: left=%s, right=%s", lKind, rKind))
}

func (expr BinaryExpression) castToType(ctx *CompilerCtx, t llvm.Type, v llvm.Value) (error, llvm.Value) {
	vKind := v.Type().TypeKind()
	tKind := t.TypeKind()

	if vKind == tKind {
		return nil, v
	}

	switch vKind {
	case llvm.IntegerTypeKind:
		switch tKind {
		case llvm.DoubleTypeKind, llvm.FloatTypeKind:
			return nil, ctx.Builder.CreateSIToFP(v, t, "")
		default:
			panic(fmt.Errorf("Cannot cast integer to %s", tKind))
		}
	case llvm.DoubleTypeKind, llvm.FloatTypeKind:
		switch tKind {
		case llvm.IntegerTypeKind:
			// Convert float to int
			return nil, ctx.Builder.CreateFPToSI(v, t, "")
		default:
			return fmt.Errorf("Cannot cast %s to %s", vKind, tKind), llvm.Value{}
		}
	case llvm.PointerTypeKind:
		switch tKind {
		case llvm.IntegerTypeKind:
			return nil, ctx.Builder.CreatePtrToInt(v, t, "")
		default:
			return nil, v
		}
	default:
		return fmt.Errorf("Unhandled type conversion from %s to %s", vKind, tKind), llvm.Value{}
	}
}

type binaryHandlerFunc func(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult)

var handlers = map[lexer.TokenKind]binaryHandlerFunc{
	lexer.And:               and,
	lexer.Or:                or,
	lexer.Plus:              add,
	lexer.Minus:             substract,
	lexer.Star:              multiply,
	lexer.Modulo:            modulo,
	lexer.Divide:            divide,
	lexer.GreaterThan:       greaterThan,
	lexer.GreaterThanEquals: greaterThanEquals,
	lexer.LessThan:          lessThan,
	lexer.LessThanEquals:    lessThanEquals,
	lexer.Equals:            equals,
	lexer.NotEquals:         notEquals,
}

func and(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	zero := llvm.ConstInt(l.Type(), 0, false)
	lBool := ctx.Builder.CreateICmp(llvm.IntNE, l, zero, "")
	rBool := ctx.Builder.CreateICmp(llvm.IntNE, r, zero, "")
	resBool := ctx.Builder.CreateAnd(lBool, rBool, "")

	return nil, &CompilerResult{Value: &resBool}
}

func or(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	zero := llvm.ConstInt(l.Type(), 0, false)
	lBool := ctx.Builder.CreateICmp(llvm.IntNE, l, zero, "")
	rBool := ctx.Builder.CreateICmp(llvm.IntNE, r, zero, "")
	resBool := ctx.Builder.CreateOr(lBool, rBool, "")

	return nil, &CompilerResult{Value: &resBool}
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

func modulo(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	res := ctx.Builder.CreateSRem(l, r, "")

	return nil, &CompilerResult{Value: &res}
}

func divide(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFDiv(l, r, "")
	} else {
		res = ctx.Builder.CreateSDiv(l, r, "")
	}

	return nil, &CompilerResult{Value: &res}
}

func multiply(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFMul(l, r, "")
	} else {
		res = ctx.Builder.CreateMul(l, r, "")
	}

	return nil, &CompilerResult{Value: &res}
}

func substract(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFSub(l, r, "")
	} else {
		res = ctx.Builder.CreateSub(l, r, "")
	}

	return nil, &CompilerResult{Value: &res}
}

func greaterThan(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFCmp(llvm.FloatOGT, l, r, "")
	} else {
		res = ctx.Builder.CreateICmp(llvm.IntSGT, l, r, "")
	}

	return nil, &CompilerResult{Value: &res}
}

func greaterThanEquals(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFCmp(llvm.FloatOGE, l, r, "")
	} else {
		res = ctx.Builder.CreateICmp(llvm.IntSGE, l, r, "")
	}

	return nil, &CompilerResult{Value: &res}
}

func lessThan(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFCmp(llvm.FloatOLT, l, r, "")
	} else {
		res = ctx.Builder.CreateICmp(llvm.IntSLT, l, r, "")
	}

	return nil, &CompilerResult{Value: &res}
}

func lessThanEquals(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFCmp(llvm.FloatOLE, l, r, "")
	} else {
		res = ctx.Builder.CreateICmp(llvm.IntSLE, l, r, "")
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

func notEquals(ctx *CompilerCtx, l, r llvm.Value) (error, *CompilerResult) {
	var res llvm.Value

	// Check if either side is a float type
	if l.Type().TypeKind() == llvm.FloatTypeKind || l.Type().TypeKind() == llvm.DoubleTypeKind ||
		r.Type().TypeKind() == llvm.FloatTypeKind || r.Type().TypeKind() == llvm.DoubleTypeKind {
		res = ctx.Builder.CreateFCmp(llvm.FloatONE, l, r, "")
	} else {
		res = ctx.Builder.CreateICmp(llvm.IntNE, l, r, "")
	}

	return nil, &CompilerResult{Value: &res}
}
