package compiler

import (
	"fmt"
	"swahili/lang/ast"
	"swahili/lang/lexer"

	"tinygo.org/x/go-llvm"
)

func (g *LLVMGenerator) VisitBinaryExpression(node *ast.BinaryExpression) error {
	if err := node.Left.Accept(g); err != nil {
		return err
	}

	format := "Strings are not supported in %s of binary expression"
	if _, ok := node.Left.(ast.StringExpression); ok {
		return fmt.Errorf(format, "left")
	}

	if _, ok := node.Right.(ast.StringExpression); ok {
		return fmt.Errorf(format, "right")
	}

	leftRes := g.getLastResult()
	if leftRes.SymbolTableEntry != nil &&
		leftRes.SymbolTableEntry.DeclaredType.Value() == ast.DataTypeString {
		return fmt.Errorf(format, "symbol left")
	}

	leftVal := g.extractRValue(leftRes)

	if err := node.Right.Accept(g); err != nil {
		return err
	}

	rightRes := g.getLastResult()
	if rightRes.SymbolTableEntry != nil &&
		rightRes.SymbolTableEntry.DeclaredType.Value() == ast.DataTypeString {
		return fmt.Errorf(format, "symbol right")
	}

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
func (g *LLVMGenerator) extractRValue(res *CompilerResult) llvm.Value {
	val := *res.Value

	if val.Type().TypeKind() != llvm.PointerTypeKind {
		return val
	}

	var loadType llvm.Type

	switch {
	case res.ArraySymbolTableEntry != nil:
		if res.ArraySymbolTableEntry.UnderlyingType.TypeKind() == llvm.StructTypeKind {
			loadType = *res.StuctPropertyValueType
		} else {
			loadType = res.ArraySymbolTableEntry.UnderlyingType
		}
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

func (g *LLVMGenerator) handleBinaryOp(kind lexer.TokenKind, l, r llvm.Value) (error, *CompilerResult) {
	handler, ok := binaryOpHandlers[kind]
	if !ok {
		return fmt.Errorf("Binary expressions : unsupported operator <%s>", kind), nil
	}

	res := handler(g, l, r)

	return nil, &CompilerResult{Value: &res}
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
